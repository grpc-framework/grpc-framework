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
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	appapi "github.com/grpc-framework/grpc-framework/v2/example/apis/hello/apiv1"
	appserver "github.com/grpc-framework/grpc-framework/v2/example/pkg/server"
	"github.com/grpc-framework/grpc-framework/v2/pkg/auth"
	"github.com/grpc-framework/grpc-framework/v2/pkg/auth/role"
	grpcclient "github.com/grpc-framework/grpc-framework/v2/pkg/grpc/client"
)

const (
	testIssuer       = "trusted.example.com"
	testSharedSecret = "test-shared-secret"
)

func newTestToken(t *testing.T, issuer, secret string, ttl time.Duration) string {
	t.Helper()
	sig, err := auth.NewSignatureSharedSecret(secret)
	require.NoError(t, err)
	token, err := auth.Token(&auth.Claims{
		Issuer:  issuer,
		Subject: "jim.stevens",
		Name:    "Jim Stevens",
		Email:   "jim@example.com",
		Roles:   []string{role.SystemAdminRoleName},
	}, sig, &auth.Options{Expiration: time.Now().Add(ttl).Unix()})
	require.NoError(t, err)
	return token
}

func newTestSecurityConfig(t *testing.T) *SecurityConfig {
	t.Helper()
	authenticator, err := auth.NewJwtAuthenticator(&auth.JwtAuthConfig{
		SharedSecret: []byte(testSharedSecret),
	})
	require.NoError(t, err)
	return &SecurityConfig{
		Authenticators: map[string]auth.Authenticator{testIssuer: authenticator},
		Role:           role.NewDefaultGenericRoleManager(),
	}
}

func contextWithAuthorization(value string) context.Context {
	return metadata.NewIncomingContext(context.Background(),
		metadata.Pairs("authorization", value))
}

func TestAuth(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		wantCode    codes.Code
		wantMessage string
		wantUser    string
		wantGuest   bool
	}{
		{
			name:      "no authorization header is a guest",
			ctx:       context.Background(),
			wantCode:  codes.OK,
			wantGuest: true,
		},
		{
			name:     "valid token",
			ctx:      contextWithAuthorization("bearer " + newTestToken(t, testIssuer, testSharedSecret, time.Hour)),
			wantCode: codes.OK,
			wantUser: "jim.stevens",
		},
		{
			// Regression: the untrusted-issuer path passed a nil error to the
			// audit logger, which dereferenced it and crashed the server.
			name:        "untrusted issuer",
			ctx:         contextWithAuthorization("bearer " + newTestToken(t, "evil.example.com", testSharedSecret, time.Hour)),
			wantCode:    codes.Unauthenticated,
			wantMessage: "evil.example.com is not a trusted issuer",
		},
		{
			name:        "missing bearer scheme",
			ctx:         contextWithAuthorization(newTestToken(t, testIssuer, testSharedSecret, time.Hour)),
			wantCode:    codes.Unauthenticated,
			wantMessage: "Invalid or missing authentication token",
		},
		{
			name:        "not a JWT",
			ctx:         contextWithAuthorization("bearer not-a-jwt"),
			wantCode:    codes.Unauthenticated,
			wantMessage: "Unable to obtain token issuer from authorization token",
		},
		{
			name:        "wrong signing key",
			ctx:         contextWithAuthorization("bearer " + newTestToken(t, testIssuer, "some-other-secret", time.Hour)),
			wantCode:    codes.Unauthenticated,
			wantMessage: "Unable to authenticate token",
		},
		{
			name:        "expired token",
			ctx:         contextWithAuthorization("bearer " + newTestToken(t, testIssuer, testSharedSecret, -time.Hour)),
			wantCode:    codes.Unauthenticated,
			wantMessage: "Unable to authenticate token",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var auditLog bytes.Buffer
			s := &GrpcFrameworkServer{
				config:         ServerConfig{Security: newTestSecurityConfig(t)},
				auditLogOutput: &auditLog,
			}

			var (
				ctx context.Context
				err error
			)
			require.NotPanics(t, func() { ctx, err = s.auth(tc.ctx) })

			if tc.wantCode != codes.OK {
				assert.Nil(t, ctx)
				st, ok := status.FromError(err)
				require.True(t, ok, "expected a gRPC status error, got %v", err)
				assert.Equal(t, tc.wantCode, st.Code())
				assert.Equal(t, tc.wantMessage, st.Message())
				assert.Contains(t, auditLog.String(), tc.wantMessage)
				return
			}

			require.NoError(t, err)
			userinfo, ok := auth.NewUserInfoFromContext(ctx)
			require.True(t, ok, "expected user information in the context")
			assert.Equal(t, tc.wantGuest, userinfo.IsGuest())
			assert.Equal(t, tc.wantUser, userinfo.Username)
		})
	}
}

// TestServerSurvivesUntrustedIssuer sends a token from an untrusted issuer
// through the full interceptor chain and checks that the server rejects it
// and keeps serving.
func TestServerSurvivesUntrustedIssuer(t *testing.T) {
	config := &ServerConfig{
		Name:     "testServer",
		Net:      "tcp",
		Address:  "127.0.0.1:0",
		Security: newTestSecurityConfig(t),
	}
	config.RegisterGrpcServers(func(gs *grpc.Server) {
		appapi.RegisterHelloIdentityServer(gs, &appserver.HelloGreeter{})
	})

	s, err := New(config)
	require.NoError(t, err)
	require.NoError(t, s.Start())
	defer s.Stop()

	conn, err := grpcclient.Connect(s.Address(),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	require.NoError(t, err)
	defer conn.Close()
	identity := appapi.NewHelloIdentityClient(conn)

	call := func(token string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+token)
		_, err := identity.ServerVersion(ctx, &appapi.ServerVersionRequest{})
		return err
	}

	err = call(newTestToken(t, "evil.example.com", testSharedSecret, time.Hour))
	assert.Equal(t, codes.Unauthenticated, status.Code(err))

	// The server must still be up and serving authenticated callers.
	assert.NoError(t, call(newTestToken(t, testIssuer, testSharedSecret, time.Hour)))
}
