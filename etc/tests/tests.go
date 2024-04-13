package tests

import (
	"errors"
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/modules"
	"github.com/gilperopiola/grpc-gateway-impl/app/storage"

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
	generateToken = func(g modules.TokenGenerator, id int, username, role string) string {
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
	hashPwd  = func(h modules.PwdHasher, pwd string) string {
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

	// Storage Errors
	errCreatingUser  = errors.New(storage.CreateUserErr)
	errGettingUser   = errors.New(storage.GetUserErr)
	errGettingUsers  = errors.New(storage.GetUsersErr)
	errCountingUsers = errors.New(storage.CountUsersErr)
	errNoOpts        = errors.New(storage.NoOptionsErr)

	// PB Requests & Responses
	signupRequest  = &pbs.SignupRequest{Username: username, Password: password}
	signupResponse = &pbs.SignupResponse{Id: int32(userID)}

	loginRequest = &pbs.LoginRequest{Username: username, Password: password}

	// gRPC
	grpcMethodName         = "modules.DefaultPublicMethods[0]"
	grpcSelfOnlyMethodName = "modules.DefaultSelfOnlyMethods[0]"
	grpcAdminMethodName    = "modules.DefaultAdminOnlyMethods[0]"
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
