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

	api "github.com/grpc-framework/grpc-framework/v2/example/apis/hello/apiv1"
	"github.com/sirupsen/logrus"
)

type HelloGreeter struct {
	api.HelloGreeterServer
	api.HelloIdentityServer
}

func (h *HelloGreeter) SayHello(
	ctx context.Context,
	req *api.SayHelloRequest,
) (*api.SayHelloResponse, error) {
	logrus.Info("Received a request in SayHello()")

	return &api.SayHelloResponse{
		Message: fmt.Sprintf("Hello, %s", req.GetName()),
	}, nil
}

func (h *HelloGreeter) ServerVersion(
	ctx context.Context,
	req *api.ServerVersionRequest,
) (*api.ServerVersionResponse, error) {
	logrus.Info("Received request for version")
	return &api.ServerVersionResponse{
		ServerVersion: &api.HelloVersion{
			Major: int32(api.HelloVersion_MAJOR),
			Minor: int32(api.HelloVersion_MINOR),
			Patch: int32(api.HelloVersion_PATCH),
			Version: fmt.Sprintf("%d.%d.%d",
				api.HelloVersion_MAJOR,
				api.HelloVersion_MINOR,
				api.HelloVersion_PATCH),
		},
	}, nil
}
