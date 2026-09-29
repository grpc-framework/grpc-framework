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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io"
	"math/big"
	mathrand "math/rand/v2"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	appapi "github.com/grpc-framework/grpc-framework/v2/example/apis/hello/apiv1"
	appserver "github.com/grpc-framework/grpc-framework/v2/example/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const sayHelloPath = "/v1/greeter:sayHello"

// freePort returns a TCP port that is not bound when it returns. The port
// is below the ephemeral range, so that the OS does not give it to a
// listener on port 0, such as the gRPC server, or to a client connection
// before the test binds it.
func freePort(t *testing.T) string {
	for i := 0; i < 100; i++ {
		port := strconv.Itoa(20000 + mathrand.IntN(10000))
		ln, err := net.Listen("tcp", ":"+port)
		if err == nil {
			ln.Close()
			return port
		}
	}
	t.Fatal("unable to find a free port")
	return ""
}

// newRestConfig returns a configuration with the REST gateway on port and
// counts the calls to the gRPC handlers in calls.
func newRestConfig(port string, calls *atomic.Int32) *ServerConfig {
	config := &ServerConfig{
		Name:         "testServer",
		Net:          "tcp",
		Address:      "127.0.0.1:0",
		Socket:       grpcSocket,
		AuditOutput:  io.Discard,
		AccessOutput: io.Discard,
	}
	config.WithDefaultRestServer(port).
		RegisterGrpcServers(func(gs *grpc.Server) {
			appapi.RegisterHelloGreeterServer(gs, &appserver.HelloGreeter{})
			appapi.RegisterHelloIdentityServer(gs, &appserver.HelloGreeter{})
		}).
		RegisterRestHandlers(
			appapi.RegisterHelloGreeterHandler,
			appapi.RegisterHelloIdentityHandler,
		).
		WithServerUnaryInterceptors(func(
			ctx context.Context,
			req interface{},
			info *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler,
		) (interface{}, error) {
			if calls != nil {
				calls.Add(1)
			}
			return handler(ctx, req)
		})
	return config
}

// startServer starts a server with config and stops it when the test ends
func startServer(t *testing.T, config *ServerConfig) *Server {
	os.Remove(config.Socket)
	s, err := New(config)
	require.NoError(t, err)
	require.NoError(t, s.Start())
	t.Cleanup(func() {
		assert.NoError(t, s.Stop())
	})
	return s
}

// sayHelloBody returns a SayHello JSON request of exactly size bytes
func sayHelloBody(t *testing.T, size int) string {
	const empty = `{"name":""}`
	require.GreaterOrEqual(t, size, len(empty))
	body := `{"name":"` + strings.Repeat("a", size-len(empty)) + `"}`
	require.Len(t, body, size)
	return body
}

func postSayHello(t *testing.T, port string, body io.Reader) *http.Response {
	resp, err := http.Post("http://127.0.0.1:"+port+sayHelloPath, "application/json", body)
	require.NoError(t, err)
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	return resp
}

func TestRestMaxRequestBodyBytes(t *testing.T) {
	const limit = 1024
	var calls atomic.Int32
	port := freePort(t)
	config := newRestConfig(port, &calls)
	config.RestConfig.MaxRequestBodyBytes = limit
	startServer(t, config)

	// At the limit
	resp := postSayHello(t, port, strings.NewReader(sayHelloBody(t, limit)))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int32(1), calls.Load())

	// One byte over the limit
	resp = postSayHello(t, port, strings.NewReader(sayHelloBody(t, limit+1)))
	assert.Equal(t, http.StatusRequestEntityTooLarge, resp.StatusCode)
	assert.Equal(t, int32(1), calls.Load())

	// One byte over the limit, sent chunked with no Content-Length.
	// The gateway fails to read the body.
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:"+port+sayHelloPath,
		strings.NewReader(sayHelloBody(t, limit+1)))
	require.NoError(t, err)
	req.ContentLength = -1
	req.TransferEncoding = []string{"chunked"}
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, int32(1), calls.Load())
}

func TestRestNoMaxRequestBodyBytes(t *testing.T) {
	var calls atomic.Int32
	port := freePort(t)
	startServer(t, newRestConfig(port, &calls))

	resp := postSayHello(t, port, strings.NewReader(sayHelloBody(t, 1<<20)))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int32(1), calls.Load())
}

func TestRestReadHeaderTimeout(t *testing.T) {
	port := freePort(t)
	config := newRestConfig(port, nil)
	config.RestConfig.ReadHeaderTimeout = 200 * time.Millisecond
	startServer(t, config)

	conn, err := net.Dial("tcp", "127.0.0.1:"+port)
	require.NoError(t, err)
	defer conn.Close()

	// Send part of the headers and never finish them
	_, err = conn.Write([]byte("GET /v1/identity:serverVersion HTTP/1.1\r\nHost: localhost\r\n"))
	require.NoError(t, err)

	start := time.Now()
	require.NoError(t, conn.SetReadDeadline(start.Add(5*time.Second)))
	_, err = conn.Read(make([]byte, 1))
	require.Error(t, err)
	assert.False(t, errors.Is(err, os.ErrDeadlineExceeded), "server did not close the connection")
	assert.Less(t, time.Since(start), 5*time.Second)
}

func TestRestTimeoutsConfigured(t *testing.T) {
	port := freePort(t)
	config := newRestConfig(port, nil)
	config.RestConfig.ReadHeaderTimeout = 1 * time.Second
	config.RestConfig.ReadTimeout = 2 * time.Second
	config.RestConfig.WriteTimeout = 3 * time.Second
	config.RestConfig.IdleTimeout = 4 * time.Second
	s := startServer(t, config)

	server := s.restGateway.server
	assert.Equal(t, 1*time.Second, server.ReadHeaderTimeout)
	assert.Equal(t, 2*time.Second, server.ReadTimeout)
	assert.Equal(t, 3*time.Second, server.WriteTimeout)
	assert.Equal(t, 4*time.Second, server.IdleTimeout)
}

func TestRestMiddlewareRunsBeforeBodyIsRead(t *testing.T) {
	const limit = 1024
	var calls, middlewareCalls atomic.Int32
	port := freePort(t)
	config := newRestConfig(port, &calls)
	config.RestConfig.MaxRequestBodyBytes = limit
	config.RestConfig.Middleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalls.Add(1)
			if r.Method == http.MethodPost {
				http.Error(w, "rate limited", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
	startServer(t, config)

	// Refused by the middleware, not by the size limit, and never
	// seen by the handler
	resp := postSayHello(t, port, strings.NewReader(sayHelloBody(t, 2*limit)))
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	resp = postSayHello(t, port, strings.NewReader(sayHelloBody(t, limit)))
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	assert.Equal(t, int32(0), calls.Load())

	// Passed on by the middleware
	resp, err := http.Get("http://127.0.0.1:" + port + "/v1/identity:serverVersion")
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int32(1), calls.Load())
	assert.Equal(t, int32(3), middlewareCalls.Load())
}

func TestRestStartPortInUse(t *testing.T) {
	port := freePort(t)
	ln, err := net.Listen("tcp", ":"+port)
	require.NoError(t, err)

	config := newRestConfig(port, nil)
	os.Remove(config.Socket)
	s, err := New(config)
	require.NoError(t, err)
	err = s.Start()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "REST gateway unable to listen on :"+port)

	// Once the port is free, a new server starts on it
	require.NoError(t, ln.Close())
	s = startServer(t, newRestConfig(port, nil))
	resp, err := http.Get("http://127.0.0.1:" + port + "/v1/identity:serverVersion")
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRestStartStopRepeatedly(t *testing.T) {
	port := freePort(t)
	for i := 0; i < 50; i++ {
		config := newRestConfig(port, nil)
		os.Remove(config.Socket)
		s, err := New(config)
		require.NoError(t, err)
		require.NoError(t, s.Start(), "start %d", i)
		require.NoError(t, s.Stop(), "stop %d", i)
	}

	// The port is free after Stop returns
	ln, err := net.Listen("tcp", ":"+port)
	require.NoError(t, err)
	ln.Close()
}

func TestRestGatewayStartStopRepeatedly(t *testing.T) {
	s := startServer(t, newRestConfig(freePort(t), nil))
	config := s.config
	config.RestConfig.Port = freePort(t)
	gw, err := NewRestGateway(&config, s.udsServer)
	require.NoError(t, err)

	// Stop right after Start, so that Stop can run before the
	// gateway serves the port
	for i := 0; i < 50; i++ {
		require.NoError(t, gw.Start(), "start %d", i)
		require.NoError(t, gw.Stop(), "stop %d", i)
	}

	// The port is free after Stop returns
	ln, err := net.Listen("tcp", ":"+config.RestConfig.Port)
	require.NoError(t, err)
	ln.Close()
}

func TestRestGatewayStopWhenNotStarted(t *testing.T) {
	gw, err := NewRestGateway(newRestConfig(freePort(t), nil), nil)
	require.NoError(t, err)
	assert.NoError(t, gw.Stop())
}

// writeSelfSignedCert writes a certificate for 127.0.0.1 and its key
// to dir and returns their paths
func writeSelfSignedCert(t *testing.T, dir string) (string, string) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "127.0.0.1"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)
	keyDer, err := x509.MarshalECPrivateKey(key)
	require.NoError(t, err)

	certFile := filepath.Join(dir, "cert.pem")
	keyFile := filepath.Join(dir, "key.pem")
	require.NoError(t, os.WriteFile(certFile,
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0600))
	require.NoError(t, os.WriteFile(keyFile,
		pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDer}), 0600))
	return certFile, keyFile
}

// newTLSGateway returns a gateway using TLS in front of the unix domain
// socket gRPC server of s
func newTLSGateway(t *testing.T, s *Server, port, certFile, keyFile string) *RestGateway {
	config := s.config
	config.RestConfig.Port = port
	config.Security = &SecurityConfig{Tls: &TLSConfig{CertFile: certFile, KeyFile: keyFile}}
	gw, err := NewRestGateway(&config, s.udsServer)
	require.NoError(t, err)
	return gw
}

func TestRestTLS(t *testing.T) {
	s := startServer(t, newRestConfig(freePort(t), nil))
	certFile, keyFile := writeSelfSignedCert(t, t.TempDir())
	port := freePort(t)
	gw := newTLSGateway(t, s, port, certFile, keyFile)
	require.NoError(t, gw.Start())
	defer func() {
		assert.NoError(t, gw.Stop())
	}()

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	resp, err := client.Get("https://127.0.0.1:" + port + "/v1/identity:serverVersion")
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRestTLSBadCertificate(t *testing.T) {
	s := startServer(t, newRestConfig(freePort(t), nil))
	missing := filepath.Join(t.TempDir(), "missing.pem")
	port := freePort(t)
	gw := newTLSGateway(t, s, port, missing, missing)

	err := gw.Start()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "REST gateway unable to load TLS certificate")
	assert.NoError(t, gw.Stop())

	// The port was not bound
	ln, err := net.Listen("tcp", ":"+port)
	require.NoError(t, err)
	ln.Close()
}
