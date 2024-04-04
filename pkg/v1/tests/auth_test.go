package tests

import (
	"context"
	"fmt"
	"testing"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/common"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestAuthGetBearer(t *testing.T) {
	tests := []struct {
		name           string
		ctx            context.Context
		expectedBearer string
		expectedErr    bool
	}{
		{
			name:           "tc_auth_get_bearer_ok",
			ctx:            metadata.NewIncomingContext(context.Background(), validAuthMetadata),
			expectedBearer: validBearerTrimmed,
			expectedErr:    false,
		},
		{
			name:        "tc_auth_get_bearer_no_metadata",
			ctx:         context.Background(),
			expectedErr: true,
		},
		{
			name:        "tc_auth_get_bearer_empty_metadata",
			ctx:         metadata.NewIncomingContext(context.Background(), emptyMetadata),
			expectedErr: true,
		},
		{
			name:        "tc_auth_get_bearer_malformed_token",
			ctx:         metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", malformedBearer)),
			expectedErr: true,
		},
	}

	jwtAuthenticator := common.NewJWTAuthenticator(jwtSecret, 10)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bearer, err := jwtAuthenticator.GetBearer(test.ctx)
			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedBearer, bearer)
			}
		})
	}
}

func TestAuthCtxWithUserInfo(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		username string
	}{
		{
			name:     "tc_auth_ctx_with_user_info_ok",
			userID:   fmt.Sprint(userID),
			username: username,
		},
		{
			name:     "tc_auth_ctx_with_user_info_empty",
			userID:   "",
			username: "",
		},
	}

	jwtAuthenticator := common.NewJWTAuthenticator(jwtSecret, 10)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = jwtAuthenticator.CtxWithUserInfo(ctx, test.userID, test.username)

			userID, _ := ctx.Value(&common.CtxKeyUserID{}).(string)
			username, _ := ctx.Value(&common.CtxKeyUsername{}).(string)

			assert.Equal(t, test.userID, userID)
			assert.Equal(t, test.username, username)
		})
	}
}

func TestAuthIsMethodAllowed(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		userID        string
		userRole      models.Role
		req           interface{}
		expectAllowed bool
	}{
		{
			name:          "tc_auth_is_method_allowed_public_ok",
			method:        grpcMethodName,
			expectAllowed: true,
		},
		{
			name:          "tc_auth_is_method_allowed_self_only_ok",
			method:        grpcSelfOnlyMethodName,
			userID:        fmt.Sprint(userID),
			req:           &usersPB.GetUserRequest{UserId: int32(userID)},
			expectAllowed: true,
		},
		{
			name:          "tc_auth_is_method_allowed_admin_ok",
			method:        grpcAdminMethodName,
			userRole:      models.AdminRole,
			expectAllowed: true,
		},
		{
			name:          "tc_auth_is_method_allowed_self_only_failure",
			method:        grpcSelfOnlyMethodName,
			userID:        fmt.Sprint(userID),
			req:           &usersPB.GetUserRequest{UserId: int32(userID + 1)},
			expectAllowed: false,
		},
		{
			name:          "tc_auth_is_method_allowed_admin_failure",
			method:        grpcAdminMethodName,
			userRole:      models.DefaultRole,
			expectAllowed: false,
		},
	}

	jwtAuthenticator := common.NewJWTAuthenticator(jwtSecret, 10)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := jwtAuthenticator.IsMethodAllowed(test.method, test.userID, test.userRole, test.req)

			if test.expectAllowed {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
