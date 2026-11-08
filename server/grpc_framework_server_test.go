/*
Copyright 2026 The grpc-framework Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package server

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"

	appapi "github.com/grpc-framework/grpc-framework/v2/example/apis/hello/apiv1"
	grpcclient "github.com/grpc-framework/grpc-framework/v2/pkg/grpc/client"
)

// blockingGreeter is a HelloGreeter whose SayHello sends on entered and
// then blocks. If release is set, SayHello ignores its context and
// returns when release is closed. Otherwise it returns linger after its
// context is done, as a handler does that cleans up after cancellation.
type blockingGreeter struct {
	appapi.UnimplementedHelloGreeterServer
	entered chan struct{}
	release chan struct{}
	linger  time.Duration
}

func (g *blockingGreeter) SayHello(
	ctx context.Context,
	req *appapi.SayHelloRequest,
) (*appapi.SayHelloResponse, error) {
	g.entered <- struct{}{}
	if g.release != nil {
		<-g.release
		return &appapi.SayHelloResponse{Message: "Hello, " + req.GetName()}, nil
	}
	<-ctx.Done()
	time.Sleep(g.linger)
	return nil, status.FromContextError(ctx.Err()).Err()
}

// rpcEndHandler is a stats.Handler that sends on ended when the server
// has finished an RPC: its handler and all interceptors have returned.
type rpcEndHandler struct {
	ended chan struct{}
}

func (h rpcEndHandler) TagRPC(ctx context.Context, _ *stats.RPCTagInfo) context.Context {
	return ctx
}

func (h rpcEndHandler) HandleRPC(_ context.Context, s stats.RPCStats) {
	if _, ok := s.(*stats.End); ok {
		h.ended <- struct{}{}
	}
}

func (h rpcEndHandler) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context {
	return ctx
}

func (h rpcEndHandler) HandleConn(context.Context, stats.ConnStats) {}

// startGreeterServer starts a GrpcFrameworkServer serving greeter on a
// free local port, with the gRPC server options opts, and returns it with
// a client connection to it.
func startGreeterServer(
	t *testing.T,
	greeter appapi.HelloGreeterServer,
	opts ...grpc.ServerOption,
) (*GrpcFrameworkServer, *grpc.ClientConn) {
	t.Helper()
	config := &ServerConfig{
		Name:         "testServer",
		Net:          "tcp",
		Address:      "127.0.0.1:0",
		Security:     &SecurityConfig{},
		AuditOutput:  io.Discard,
		AccessOutput: io.Discard,
	}
	config.RegisterGrpcServers(func(gs *grpc.Server) {
		appapi.RegisterHelloGreeterServer(gs, greeter)
	})
	config.WithServerOptions(opts...)

	s, err := NewGrpcFrameworkServer(config)
	require.NoError(t, err)
	require.NoError(t, s.Start())

	conn, err := grpcclient.Connect(s.Address(),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	require.NoError(t, err)
	return s, conn
}

// sayHello calls SayHello on conn in a new goroutine and returns a
// channel that receives the error of the call.
func sayHello(conn *grpc.ClientConn) <-chan error {
	done := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := appapi.NewHelloGreeterClient(conn).
			SayHello(ctx, &appapi.SayHelloRequest{Name: "test"})
		done <- err
	}()
	return done
}

// waitForHandler waits until the handler of the RPC started by sayHello
// has been entered.
func waitForHandler(t *testing.T, greeter *blockingGreeter, rpc <-chan error) {
	t.Helper()
	select {
	case <-greeter.entered:
	case err := <-rpc:
		t.Fatalf("SayHello returned before its handler was entered: %v", err)
	}
}

// TestServerRestartWithRpcInProgress stops a server while an RPC is in its
// handler and starts the next server, a few times, as a test suite does
// that starts a server per test. The handler returns a while after Stop
// cancels it, and the grpc_prometheus interceptor around it then reads
// grpc_prometheus.DefaultServerMetrics. Under -race the test fails if that
// handler still runs after Stop returned, while the next Start writes
// DefaultServerMetrics.
func TestServerRestartWithRpcInProgress(t *testing.T) {
	end := rpcEndHandler{ended: make(chan struct{}, 1)}
	waitForRpcEnd := func() {
		t.Helper()
		select {
		case <-end.ended:
		case <-time.After(10 * time.Second):
			t.Fatal("the RPC of the stopped server did not end")
		}
	}

	for i := 0; i < 5; i++ {
		greeter := &blockingGreeter{
			entered: make(chan struct{}, 1),
			linger:  100 * time.Millisecond,
		}
		s, conn := startGreeterServer(t, greeter, grpc.StatsHandler(end))

		// Let the RPC of the stopped server end before this server gets one,
		// so that no other RPC reads the metrics between this Start and it
		if i > 0 {
			waitForRpcEnd()
		}

		rpc := sayHello(conn)
		waitForHandler(t, greeter, rpc)

		s.Stop()
		assert.False(t, s.IsRunning())
		assert.Error(t, <-rpc)
		conn.Close()
	}
	waitForRpcEnd()
}

// TestServerStopWaitsForHandlers checks that Stop returns only after the
// handler of an RPC in progress has returned, even when the handler does
// not watch its context.
func TestServerStopWaitsForHandlers(t *testing.T) {
	greeter := &blockingGreeter{
		entered: make(chan struct{}, 1),
		release: make(chan struct{}),
	}
	s, conn := startGreeterServer(t, greeter)
	defer conn.Close()

	rpc := sayHello(conn)
	waitForHandler(t, greeter, rpc)

	stopped := make(chan struct{})
	go func() {
		s.Stop()
		close(stopped)
	}()

	select {
	case <-stopped:
		close(greeter.release)
		t.Fatal("Stop returned while a handler was running")
	case <-time.After(200 * time.Millisecond):
	}

	close(greeter.release)
	select {
	case <-stopped:
	case <-time.After(10 * time.Second):
		t.Fatal("Stop did not return after the handler returned")
	}
	assert.False(t, s.IsRunning())
	assert.Error(t, <-rpc)
}
