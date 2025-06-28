/*
Copyright 2022 Pure Storage

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
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	appapi "github.com/grpc-framework/grpc-framework/v2/example/apis/hello/apiv1"
	appserver "github.com/grpc-framework/grpc-framework/v2/example/pkg/server"
	grpcclient "github.com/grpc-framework/grpc-framework/v2/pkg/grpc/client"
	"github.com/lpabon/lputils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	grpcSocket = "/tmp/grpc-framework-testServer.sock"
)

type testServer struct {
	config *ServerConfig
	server *Server
	conn   *grpc.ClientConn
}

func (s *testServer) Stop() {
	s.server.Stop()
	s.conn.Close()
}

func (s *testServer) Address() string {
	return s.server.Address()
}

func (s *testServer) Conn() *grpc.ClientConn {
	return s.conn
}

func newDefaultConfig(t *testing.T) *ServerConfig {
	config := &ServerConfig{
		Name:    "testServer",
		Net:     "tcp",
		Address: "127.0.0.1:0",
		Socket:  grpcSocket,
	}
	config.WithDefaultRestServer("9001").
		RegisterGrpcServers(func(gs *grpc.Server) {
			appapi.RegisterHelloGreeterServer(gs, &appserver.HelloGreeter{})
			appapi.RegisterHelloIdentityServer(gs, &appserver.HelloGreeter{})
		}).
		RegisterRestHandlers(
			appapi.RegisterHelloGreeterHandler,
			appapi.RegisterHelloIdentityHandler,
		)

	return config
}

func newTestServer(t *testing.T, config *ServerConfig) *testServer {

	if config.Socket != "" {
		os.Remove(config.Socket)
	}

	s, err := New(config)
	assert.NoError(t, err)

	err = s.Start()
	assert.NoError(t, err)

	// Setup connection to server
	conn, err := grpcclient.Connect(
		s.Address(),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	assert.NoError(t, err)

	return &testServer{
		config: config,
		server: s,
		conn:   conn,
	}
}

func newDefaultTestServer(t *testing.T) *testServer {
	ts := newTestServer(t, newDefaultConfig(t))

	// test connection
	identity := appapi.NewHelloIdentityClient(ts.Conn())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	assert.NotNil(t, ctx)
	r, err := identity.ServerVersion(ctx, &appapi.ServerVersionRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.NotNil(t, r.GetServerVersion())
	assert.Equal(t, int32(appapi.HelloVersion_MAJOR), r.GetServerVersion().GetMajor())
	assert.Equal(t, int32(appapi.HelloVersion_MINOR), r.GetServerVersion().GetMinor())
	assert.Equal(t, int32(appapi.HelloVersion_PATCH), r.GetServerVersion().GetPatch())

	return ts
}

func TestSimpleServer(t *testing.T) {
	s := newDefaultTestServer(t)
	defer s.Stop()
}

func TestServerWithoutRest(t *testing.T) {
	config := &ServerConfig{
		Name:    "testServer",
		Net:     "tcp",
		Address: "127.0.0.1:0",
		Socket:  grpcSocket,
	}

	s := newTestServer(t, config)
	assert.Nil(t, s.server.restGateway)
	assert.NotNil(t, s.server.udsServer)
	assert.NotNil(t, s.server.netServer)
	defer s.Stop()
}

func TestServerWithoutUdsAndRest(t *testing.T) {
	config := &ServerConfig{
		Name:    "testServer",
		Net:     "tcp",
		Address: "127.0.0.1:0",
	}

	s := newTestServer(t, config)
	assert.Nil(t, s.server.restGateway)
	assert.Nil(t, s.server.udsServer)
	assert.NotNil(t, s.server.netServer)
	defer s.Stop()
}

func TestServerErrorWhenRestWithoutUds(t *testing.T) {
	config := &ServerConfig{
		Name:    "testServer",
		Net:     "tcp",
		Address: "127.0.0.1:0",
	}
	config.WithDefaultRestServer("9001").
		RegisterGrpcServers(func(gs *grpc.Server) {
			appapi.RegisterHelloGreeterServer(gs, &appserver.HelloGreeter{})
		})
	_, err := New(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must provide unix domain socket for REST")
}

func TestSimpleServerLockTest(t *testing.T) {
	s := newDefaultTestServer(t)
	defer s.Stop()

	value := 0
	err := s.server.Transaction(func() error {
		value = 1
		return nil
	})
	assert.Equal(t, value, 1)
	assert.NoError(t, err)

	err = s.server.Transaction(func() error {
		value = 2
		return fmt.Errorf("ERROR")
	})
	assert.Equal(t, value, 2)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "ERROR")
}

func runminigw(t *testing.T, address string) (*http.Server, error) {
	ctx := context.Background()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	c, err := grpcclient.Connect(address, opts)
	assert.NoError(t, err)
	err = appapi.RegisterHelloIdentityHandler(ctx, mux, c)
	assert.NoError(t, err)

	//err := appapi.RegisterHelloIdentityHandlerFromEndpoint(ctx, mux, address, opts)

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	server := &http.Server{Addr: ":9002", Handler: mux}
	ready := make(chan bool)
	go func() {
		fmt.Println("in the listen")
		ready <- true
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("listen and serve err: %v\n", err)
		}
	}()
	<-ready
	return server, nil
}

func TestRest(t *testing.T) {
	s := newDefaultTestServer(t)
	defer s.Stop()

	/*
		httpserver, err := runminigw(t, s.server.UdsAddress())
		assert.NoError(t, err)
		assert.NotNil(t, httpserver)
		defer func() {
			if err := httpserver.Shutdown(context.Background()); err != nil {
				fmt.Printf("httpserver.Shutdown failed: %v\n", err)
			} else {
				fmt.Println("== DONE ==")
			}
		}()
	*/

	resp, err := http.Get("http://localhost:9001/v1/identity:serverVersion")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotZero(t, resp.ContentLength)

	type helloVersion struct {
		// Version major number
		Major int32 `protobuf:"varint,1,opt,name=major,proto3" json:"major,omitempty"`
		// Version minor number
		Minor int32 `protobuf:"varint,2,opt,name=minor,proto3" json:"minor,omitempty"`
		// Version patch number
		Patch int32 `protobuf:"varint,3,opt,name=patch,proto3" json:"patch,omitempty"`
		// String representation of the version. Must be
		// in `major.minor.patch` format.
		Version string `protobuf:"bytes,4,opt,name=version,proto3" json:"version,omitempty"`
	}

	type serverVersionResponse struct {
		// Hello application version
		ServerVersion *helloVersion `json:"serverVersion,omitempty"`
	}

	var b serverVersionResponse
	err = lputils.GetJsonFromResponse(resp, &b)
	assert.NoError(t, err)
	assert.Equal(t, int32(appapi.HelloVersion_MAJOR), b.ServerVersion.Major)
	assert.Equal(t, int32(appapi.HelloVersion_MINOR), b.ServerVersion.Minor)
	assert.Equal(t, int32(appapi.HelloVersion_PATCH), b.ServerVersion.Patch)
}

func rateLimiterShowsDenial(t *testing.T, s *testServer) bool {
	denied := false
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			ctx := context.Background()
			conn, err := grpcclient.Connect(s.Address(), []grpc.DialOption{grpc.WithInsecure()})
			g := appapi.NewHelloGreeterClient(conn)

			for {
				_, err = g.SayHello(ctx, &appapi.SayHelloRequest{})
				if err != nil {
					serverError, ok := status.FromError(err)
					assert.True(t, ok)
					assert.Equal(t, serverError.Code(), codes.ResourceExhausted)

					if serverError.Code() == codes.ResourceExhausted {
						time.Sleep(time.Millisecond * time.Duration(i))
						denied = true
					}
					continue
				}
				break
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	return denied
}

func TestServerRateLimiter(t *testing.T) {
	c := newDefaultConfig(t)
	c.WithRateLimiter(rate.NewLimiter(2, 2))
	s := newTestServer(t, c)
	defer s.Stop()

	assert.True(t, rateLimiterShowsDenial(t, s))
}

func TestServerNoRateLimiter(t *testing.T) {
	s := newDefaultTestServer(t)
	defer s.Stop()

	assert.False(t, rateLimiterShowsDenial(t, s))
}
