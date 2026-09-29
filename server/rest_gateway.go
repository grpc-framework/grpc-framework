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
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-framework/grpc-framework/v2/pkg/correlation"
	grpcclient "github.com/grpc-framework/grpc-framework/v2/pkg/grpc/client"
)

type RestGateway struct {
	config     ServerConfig
	grpcServer *GrpcFrameworkServer

	// lock protects the fields below, set by Start and cleared by Stop
	lock   sync.Mutex
	server *http.Server
	conn   *grpc.ClientConn
	// done is closed when the goroutine serving the port returns
	done chan struct{}
}

func NewRestGateway(config *ServerConfig, grpcServer *GrpcFrameworkServer) (*RestGateway, error) {
	return &RestGateway{
		config:     *config,
		grpcServer: grpcServer,
	}, nil
}

// Start binds the REST port and serves it in the background. The port is
// bound when Start returns; if it cannot be bound, Start returns the error.
// It returns an error if the gateway is already running.
func (s *RestGateway) Start() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.server != nil {
		return fmt.Errorf("REST gateway already running")
	}

	handler, conn, err := s.restServerSetupHandlers()
	if err != nil {
		return err
	}

	// From the outside in:
	//  - CORS, so that every response, refusals included, has the CORS
	//    headers, and preflight requests are answered before Middleware.
	//  - Middleware, so that it can refuse a request, for example when
	//    rate limited, before the body is read.
	//  - The request body size limit.
	restConfig := s.config.RestConfig
	if restConfig.MaxRequestBodyBytes > 0 {
		handler = maxBytesHandler(handler, restConfig.MaxRequestBodyBytes)
	}
	if restConfig.Middleware != nil {
		handler = restConfig.Middleware(handler)
	}
	handler = s.corsHandler(handler)

	server := &http.Server{
		Addr:              ":" + restConfig.Port,
		Handler:           handler,
		ReadHeaderTimeout: restConfig.ReadHeaderTimeout,
		ReadTimeout:       restConfig.ReadTimeout,
		WriteTimeout:      restConfig.WriteTimeout,
		IdleTimeout:       restConfig.IdleTimeout,
	}

	// Load the certificate here so that an error is returned to the caller
	if s.config.Security != nil && s.config.Security.Tls != nil {
		cert, err := tls.LoadX509KeyPair(s.config.Security.Tls.CertFile, s.config.Security.Tls.KeyFile)
		if err != nil {
			conn.Close()
			return fmt.Errorf("REST gateway unable to load TLS certificate: %w", err)
		}
		server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	}

	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		conn.Close()
		return fmt.Errorf("REST gateway unable to listen on %s: %w", server.Addr, err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		// Serve closes the listener when it returns, but ServeTLS can fail
		// before it calls Serve.
		defer ln.Close()

		var err error
		if server.TLSConfig != nil {
			err = server.ServeTLS(ln, "", "")
		} else {
			err = server.Serve(ln)
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Errorf("REST gRPC Gateway stopped: %v", err)
		}
	}()

	s.server = server
	s.conn = conn
	s.done = done
	logrus.Infof("gRPC REST Gateway started on port :%s", restConfig.Port)

	return nil
}

// Stop closes the REST gateway. The port is free when Stop returns.
// It does nothing if the gateway is not running.
func (s *RestGateway) Stop() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.server == nil {
		return nil
	}

	err := s.server.Close()

	// Close does not close a listener that Serve has not started using yet;
	// Serve closes it on its way out, so wait for it.
	<-s.done

	err = errors.Join(err, s.conn.Close())
	s.server = nil
	s.conn = nil
	s.done = nil

	return err
}

// maxBytesHandler refuses a request whose body is larger than n bytes.
// A request that declares a larger Content-Length gets 413 without being
// passed on; reading past n bytes of a body of unknown length fails.
func maxBytesHandler(h http.Handler, n int64) http.Handler {
	// MaxBytesHandler passes on a copy of the request. The request itself
	// must keep its body: before Go 1.27, the server checks the type of
	// the body to send "100 Continue" and to reuse the connection.
	limited := http.MaxBytesHandler(h, n)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > n {
			http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
			return
		}
		limited.ServeHTTP(w, r)
	})
}

// corsHandler wraps h with the configured CORS handler
func (s *RestGateway) corsHandler(h http.Handler) http.Handler {
	corsOptions := s.config.RestConfig.CorsOptions
	if !corsOptions.Enabled {
		return h
	}
	if corsOptions.CustomOptions == nil {
		logrus.Warn("REST Cors configuration missing; skipping")
		return h
	}
	return cors.New(*corsOptions.CustomOptions).Handler(h)
}

// restServerSetupHandlers sets up the handlers to the swagger ui and
// to the gRPC REST Gateway. It returns the connection to the gRPC server,
// which the caller must close.
func (s *RestGateway) restServerSetupHandlers() (http.Handler, *grpc.ClientConn, error) {

	// Create an HTTP server router
	mux := http.NewServeMux()

	// Swagger files using packr
	/* Packr on swagger

	swaggerUIBox := packr.NewBox("./swagger-ui")
	swaggerJSONBox := packr.NewBox("./api")
	mime.AddExtensionType(".svg", "image/svg+xml")

	// Handler to return swagger.json
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Write(swaggerJSONBox.Bytes("api.swagger.json"))
	})

	// Handler to access the swagger ui. The UI pulls the swagger
	// json file from /swagger.json
	// The link below MUST have th last '/'. It is really important.
	// This link is deprecated
	prefix := "/swagger-ui/"
	mux.Handle(prefix,
		http.StripPrefix(prefix, http.FileServer(swaggerUIBox)))
	// This is the new location
	prefix = "/sdk/"
	mux.Handle(prefix,
		http.StripPrefix(prefix, http.FileServer(swaggerUIBox)))
	*/

	if s.config.RestConfig.PrometheusConfig.Enabled {
		if s.config.RestConfig.PrometheusConfig.Path == "" {
			logrus.Warn("REST Prometheus path missing; skipping")
		} else {
			mux.Handle(s.config.RestConfig.PrometheusConfig.Path, promhttp.Handler())
		}
	}

	// Create a router just for HTTP REST gRPC Server Gateway
	gmux := runtime.NewServeMux()
	fmt.Printf("address is %s\n", s.grpcServer.Address())

	// Connect to gRPC unix domain socket
	conn, err := grpcclient.Connect(
		s.grpcServer.Address(),
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(correlation.ContextUnaryClientInterceptor),
		})
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to connect to gRPC handler: %v", err)
	}

	// Register the REST Gateway handlers
	for _, handler := range s.config.RestServerExtensions {
		err := handler(context.Background(), gmux, conn)
		if err != nil {
			conn.Close()
			return nil, nil, err
		}
	}

	// Pass all other unhandled paths to the gRPC gateway
	mux.Handle("/", gmux)

	return mux, conn, nil
}
