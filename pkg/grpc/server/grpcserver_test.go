/*
Package server is a generic gRPC server manager
Copyright 2018 Portworx

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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// testServer is a simple struct used abstract
// the creation and setup of the gRPC service
type testServer struct {
	conn   *grpc.ClientConn
	server *GrpcServer
}

func newTestServer(t *testing.T) *testServer {
	tester := &testServer{}

	var err error
	// Setup simple driver
	tester.server, err = New(&GrpcServerConfig{
		Name:    "unit-test",
		Net:     "tcp",
		Address: "127.0.0.1:0",
	})
	assert.Nil(t, err)

	// Test Start
	var testGrpcServer *grpc.Server
	startCalled := false
	err = tester.server.Start(func(grpcServer *grpc.Server) {
		testGrpcServer = grpcServer
		startCalled = true
	})
	assert.Nil(t, err)
	assert.NotNil(t, testGrpcServer)
	assert.True(t, startCalled)

	// Setup a connection to the driver
	tester.conn, err = grpc.Dial(tester.server.Address(), grpc.WithInsecure())
	assert.Nil(t, err)

	return tester
}

func (s *testServer) Stop() {
	// Shutdown servers
	s.conn.Close()
	s.server.Stop()
}

func (s *testServer) Conn() *grpc.ClientConn {
	return s.conn
}

func (s *testServer) Server() *GrpcServer {
	return s.server
}

func TestGrpcServerStartConfig(t *testing.T) {
	_, err := New(nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Configuration")

	_, err = New(&GrpcServerConfig{})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "must be provided")

	_, err = New(&GrpcServerConfig{
		Name: "name",
		Net:  "net",
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Address")

	_, err = New(&GrpcServerConfig{
		Address: "address",
		Net:     "net",
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Name")

	_, err = New(&GrpcServerConfig{
		Name:    "name",
		Address: "address",
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Net")

}

func TestGrpcServerStart(t *testing.T) {
	s := newTestServer(t)
	assert.True(t, s.Server().IsRunning())
	defer s.Stop()

	// Check if we can still talk to the server
	// after starting multiple times.
	called := false
	err := s.Server().Start(func(grpcserver *grpc.Server) {
		called = true
	})
	assert.False(t, called)
	assert.True(t, s.Server().IsRunning())
	assert.NotNil(t, err)
	err = s.Server().Start(func(grpcserver *grpc.Server) {
		called = true
	})
	assert.False(t, called)
	assert.True(t, s.Server().IsRunning())
	assert.NotNil(t, err)
	err = s.Server().Start(func(grpcserver *grpc.Server) {
		called = true
	})
	assert.False(t, called)
	assert.True(t, s.Server().IsRunning())
	assert.NotNil(t, err)

}

func TestServerStop(t *testing.T) {
	s := newTestServer(t)
	assert.True(t, s.Server().IsRunning())
	s.Stop()
	assert.False(t, s.Server().IsRunning())

	assert.NotPanics(t, s.Stop)
	assert.False(t, s.Server().IsRunning())
	assert.NotPanics(t, s.Stop)
	assert.False(t, s.Server().IsRunning())
	assert.NotPanics(t, s.Stop)
	assert.False(t, s.Server().IsRunning())
	assert.NotPanics(t, s.Stop)
	assert.False(t, s.Server().IsRunning())
}

func TestServerStartAfterStop(t *testing.T) {
	s := newTestServer(t)
	s.Stop()

	called := false
	err := s.Server().Start(func(grpcserver *grpc.Server) {
		called = true
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be started again")
	assert.False(t, called)
	assert.False(t, s.Server().IsRunning())

	err = s.Server().StartWithServer(func() *grpc.Server {
		called = true
		return grpc.NewServer()
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be started again")
	assert.False(t, called)
	assert.False(t, s.Server().IsRunning())
}

func TestServerStopBeforeServe(t *testing.T) {
	s, err := New(&GrpcServerConfig{
		Name:    "unit-test",
		Net:     "tcp",
		Address: "127.0.0.1:0",
	})
	assert.NoError(t, err)

	// Serve returns grpc.ErrServerStopped when Stop runs before it, which
	// must not exit the process
	err = s.StartWithServer(func() *grpc.Server {
		gs := grpc.NewServer()
		gs.Stop()
		return gs
	})
	assert.NoError(t, err)

	// Stop waits for the goroutine calling Serve to return
	s.Stop()
	assert.False(t, s.IsRunning())
}

// registerBlockingService registers the service test.Blocking with one
// streaming method, Block. Its handler sends on entered and returns when
// release is closed, ignoring the cancellation of its context.
func registerBlockingService(
	gs *grpc.Server,
	entered chan<- struct{},
	release <-chan struct{},
) {
	gs.RegisterService(&grpc.ServiceDesc{
		ServiceName: "test.Blocking",
		HandlerType: (*any)(nil),
		Streams: []grpc.StreamDesc{{
			StreamName:    "Block",
			ServerStreams: true,
			Handler: func(any, grpc.ServerStream) error {
				entered <- struct{}{}
				<-release
				return nil
			},
		}},
	}, struct{}{})
}

func TestServerStopWaitsForHandlers(t *testing.T) {
	s, err := New(&GrpcServerConfig{
		Name:    "unit-test",
		Net:     "tcp",
		Address: "127.0.0.1:0",
	})
	assert.NoError(t, err)

	entered := make(chan struct{}, 1)
	release := make(chan struct{})
	err = s.Start(func(gs *grpc.Server) {
		registerBlockingService(gs, entered, release)
	})
	assert.NoError(t, err)

	conn, err := grpc.NewClient(s.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = conn.NewStream(ctx, &grpc.StreamDesc{ServerStreams: true},
		"/test.Blocking/Block")
	assert.NoError(t, err)
	select {
	case <-entered:
	case <-ctx.Done():
		t.Fatal("the handler of the RPC was not entered")
	}

	stopped := make(chan struct{})
	go func() {
		s.Stop()
		close(stopped)
	}()

	// The handler ignores the cancellation of its context, so Stop must
	// wait until it is released.
	select {
	case <-stopped:
		close(release)
		t.Fatal("Stop returned while a handler was running")
	case <-time.After(200 * time.Millisecond):
	}

	close(release)
	select {
	case <-stopped:
	case <-time.After(10 * time.Second):
		t.Fatal("Stop did not return after the handler returned")
	}
	assert.False(t, s.IsRunning())
}
