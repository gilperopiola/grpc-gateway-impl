package tests

import (
	"errors"
	"testing"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/common"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* ----------------------------------- */
/*         - Tests Variables -         */
/* ----------------------------------- */

var (

	// JWT Authenticator
	jwtSecret     = "jwt_secret"
	generateToken = func(g common.TokenGenerator, id int, username, role string) string {
		token, _ := g.Generate(id, username, models.Role(role))
		return token
	}

	validBearer        = "Bearer valid"
	validBearerTrimmed = "valid"
	malformedBearer    = "Not-Bearer invalid"

	validAuthMetadata = metadata.Pairs("authorization", validBearer)
	emptyMetadata     = metadata.New(nil)

	// Pwd Hasher
	hashSalt = "hash_salt"
	hashPwd  = func(h common.PwdHasher, pwd string) string {
		return h.Hash(pwd)
	}

	// User Model
	userID   = 1
	username = "username"
	password = "password"

	user       = &models.User{Username: username, Password: password}
	userWithID = &models.User{Username: username, Password: password, ID: 1}
	userEmpty  = &models.User{}

	users    = models.Users{user}
	usersNil = models.Users(nil)

	// Repository Errors
	errCreatingUser  = errors.New(repository.CreateUserErr)
	errGettingUser   = errors.New(repository.GetUserErr)
	errGettingUsers  = errors.New(repository.GetUsersErr)
	errCountingUsers = errors.New(repository.CountUsersErr)
	errNoOpts        = errors.New(repository.NoOptionsErr)

	// PB Requests & Responses
	signupRequest  = &usersPB.SignupRequest{Username: username, Password: password}
	signupResponse = &usersPB.SignupResponse{Id: int32(userID)}

	loginRequest = &usersPB.LoginRequest{Username: username, Password: password}

	// gRPC
	grpcMethodName         = common.DefaultPublicMethods[0]
	grpcSelfOnlyMethodName = common.DefaultSelfOnlyMethods[0]
	grpcAdminMethodName    = common.DefaultAdminOnlyMethods[0]
)

/* ----------------------------------- */
/*          - Tests Helpers -          */
/* ----------------------------------- */

func assertDBError(t *testing.T, expectedErr, err error) {
	if expectedErr == nil {
		assert.NoError(t, err)
		return
	}
	var dbErr *errs.DBError
	if assert.ErrorAs(t, err, &dbErr) {
		assert.Equal(t, expectedErr, dbErr.Unwrap())
	}
}

func assertSvcError(t *testing.T, expectedErr, err error) {
	if expectedErr == nil {
		assert.NoError(t, err)
		return
	}
	assert.ErrorIs(t, err, expectedErr)
}

func copyPB(response protoreflect.ProtoMessage) protoreflect.ProtoMessage {
	if response == nil {
		return nil
	}
	return proto.Clone(response)
}

func copyUser(user *models.User) *models.User {
	if user == nil {
		return nil
	}
	copy := *user
	return &copy
}

func copyUsers(users models.Users) models.Users {
	if users == nil {
		return nil
	}
	copy := make(models.Users, len(users))
	for i, user := range users {
		copy[i] = copyUser(user)
	}
	return copy
}
