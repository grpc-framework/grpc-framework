// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: api/hello.proto

/*
Package api is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package api

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = metadata.Join

func request_HelloGreeter_SayHello_0(ctx context.Context, marshaler runtime.Marshaler, client HelloGreeterClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq HelloGreeterSayHelloRequest
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.SayHello(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_HelloGreeter_SayHello_0(ctx context.Context, marshaler runtime.Marshaler, server HelloGreeterServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq HelloGreeterSayHelloRequest
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.SayHello(ctx, &protoReq)
	return msg, metadata, err

}

func request_HelloIdentity_ServerVersion_0(ctx context.Context, marshaler runtime.Marshaler, client HelloIdentityClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq HelloIdentityVersionRequest
	var metadata runtime.ServerMetadata

	msg, err := client.ServerVersion(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_HelloIdentity_ServerVersion_0(ctx context.Context, marshaler runtime.Marshaler, server HelloIdentityServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq HelloIdentityVersionRequest
	var metadata runtime.ServerMetadata

	msg, err := server.ServerVersion(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterHelloGreeterHandlerServer registers the http handlers for service HelloGreeter to "mux".
// UnaryRPC     :call HelloGreeterServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterHelloGreeterHandlerFromEndpoint instead.
func RegisterHelloGreeterHandlerServer(ctx context.Context, mux *runtime.ServeMux, server HelloGreeterServer) error {

	mux.Handle("POST", pattern_HelloGreeter_SayHello_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/hello.v1.HelloGreeter/SayHello", runtime.WithHTTPPathPattern("/v1/greeter:sayHello"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_HelloGreeter_SayHello_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_HelloGreeter_SayHello_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterHelloIdentityHandlerServer registers the http handlers for service HelloIdentity to "mux".
// UnaryRPC     :call HelloIdentityServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterHelloIdentityHandlerFromEndpoint instead.
func RegisterHelloIdentityHandlerServer(ctx context.Context, mux *runtime.ServeMux, server HelloIdentityServer) error {

	mux.Handle("GET", pattern_HelloIdentity_ServerVersion_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/hello.v1.HelloIdentity/ServerVersion", runtime.WithHTTPPathPattern("/v1/identity:serverVersion"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_HelloIdentity_ServerVersion_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_HelloIdentity_ServerVersion_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterHelloGreeterHandlerFromEndpoint is same as RegisterHelloGreeterHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterHelloGreeterHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.DialContext(ctx, endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterHelloGreeterHandler(ctx, mux, conn)
}

// RegisterHelloGreeterHandler registers the http handlers for service HelloGreeter to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterHelloGreeterHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterHelloGreeterHandlerClient(ctx, mux, NewHelloGreeterClient(conn))
}

// RegisterHelloGreeterHandlerClient registers the http handlers for service HelloGreeter
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "HelloGreeterClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "HelloGreeterClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "HelloGreeterClient" to call the correct interceptors.
func RegisterHelloGreeterHandlerClient(ctx context.Context, mux *runtime.ServeMux, client HelloGreeterClient) error {

	mux.Handle("POST", pattern_HelloGreeter_SayHello_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/hello.v1.HelloGreeter/SayHello", runtime.WithHTTPPathPattern("/v1/greeter:sayHello"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_HelloGreeter_SayHello_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_HelloGreeter_SayHello_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_HelloGreeter_SayHello_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"v1", "greeter"}, "sayHello"))
)

var (
	forward_HelloGreeter_SayHello_0 = runtime.ForwardResponseMessage
)

// RegisterHelloIdentityHandlerFromEndpoint is same as RegisterHelloIdentityHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterHelloIdentityHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.DialContext(ctx, endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterHelloIdentityHandler(ctx, mux, conn)
}

// RegisterHelloIdentityHandler registers the http handlers for service HelloIdentity to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterHelloIdentityHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterHelloIdentityHandlerClient(ctx, mux, NewHelloIdentityClient(conn))
}

// RegisterHelloIdentityHandlerClient registers the http handlers for service HelloIdentity
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "HelloIdentityClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "HelloIdentityClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "HelloIdentityClient" to call the correct interceptors.
func RegisterHelloIdentityHandlerClient(ctx context.Context, mux *runtime.ServeMux, client HelloIdentityClient) error {

	mux.Handle("GET", pattern_HelloIdentity_ServerVersion_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/hello.v1.HelloIdentity/ServerVersion", runtime.WithHTTPPathPattern("/v1/identity:serverVersion"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_HelloIdentity_ServerVersion_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_HelloIdentity_ServerVersion_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_HelloIdentity_ServerVersion_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"v1", "identity"}, "serverVersion"))
)

var (
	forward_HelloIdentity_ServerVersion_0 = runtime.ForwardResponseMessage
)
