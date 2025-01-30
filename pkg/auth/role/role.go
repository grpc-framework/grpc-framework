/*
Package role manages roles in Kvdb and provides validation
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
package role

import (
	"context"

	authv1 "github.com/grpc-framework/grpc-framework/v2/apis/auth/apiv1"
)

// RoleManager provides an implementation of the SDK Role handler
// and the necessary verification methods
type RoleManager interface {
	// Verify returns no error if the role exists and is allowed
	// to run the requested method
	Verify(ctx context.Context, roles []string, method string) error
}

// RoleStore provides an interface to storing roles
type RoleStore interface {
	// Get gets a role and its rules
	Get(roleName string) (*authv1.Role, bool)
	// Set saves a new role
	Set(role *authv1.Role) error
	// Delete deletes a saved role
	Delete(roleName string)
	// List returns a list of all the roles saved
	List() ([]*authv1.Role, error)
}
