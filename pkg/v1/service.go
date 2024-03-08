package v1

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

/* ----------------------------------- */
/*           - v1 Service -            */
/* ----------------------------------- */

// Service is the interface that defines the methods of the API.
type Service interface {
	usersPB.UsersServiceServer
}

// Service is our concrete implementation of the gRPC API defined in the .proto files.
type service struct {
	DB Repository
	*usersPB.UnimplementedUsersServiceServer
}

// NewService returns a new instance of the Service.
func NewService(db Repository) *service {
	return &service{DB: db}
}

/* ----------------------------------- */
/*        - Service Handlers -         */
/* ----------------------------------- */

// Signup is the handler / entrypoint for the Signup API method. Both gRPC and HTTP.
func (s *service) Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {
	user, err := s.DB.GetUser(0, in.Username)
	if err != nil || user == nil {
		return nil, status.Errorf(codes.Unknown, "Service.Signup: %v", err)
	}

	if user.Username != "" {
		return nil, status.Errorf(codes.AlreadyExists, "Service.Signup: user already exists")
	}

	user.Username = in.Username
	user.Password = hashPassword(in.Password, "some-salt")

	if user, err = s.DB.CreateUser(*user); err != nil {
		return nil, status.Errorf(codes.Unknown, "Service.Signup: %v", err)
	}

	return signupOKResponse(user.ID)
}

func signupOKResponse(id int) (*usersPB.SignupResponse, error) {
	return &usersPB.SignupResponse{Id: int32(id)}, nil
}

// Login is the handler / entrypoint for the Login API method. Both gRPC and HTTP.
func (s *service) Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error) {

	user, err := s.DB.GetUser(0, in.Username)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Service.Login: %v", err)
	}
	if user != nil && user.Username == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Service.Login: user does not exist")
	}

	if user.Password != hashPassword(in.Password, "some-salt") {
		return nil, status.Errorf(codes.Unauthenticated, "Service.Login: password does not match")
	}

	tokenString, _ := GenerateToken(user.ID, user.Username, "", "user", "some", 7)

	return loginOKResponse(tokenString)
}

func loginOKResponse(token string) (*usersPB.LoginResponse, error) {
	return &usersPB.LoginResponse{Token: token}, nil
}

// GetUser is the handler / entrypoint for the GetUser API method. Both gRPC and HTTP.
func (s *service) GetUser(ctx context.Context, in *usersPB.GetUserRequest) (*usersPB.GetUserResponse, error) {
	user, err := s.DB.GetUser(int(in.UserId), "")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "Service.GetUser: user not found")
		}
		return nil, status.Errorf(codes.Unknown, "Service.GetUser: error retrieving user: %v", err)
	}

	return &usersPB.GetUserResponse{
		Username: user.Username,
	}, nil
}

// simulateRandomErr returns an error 1 out of 5 times.
func simulateRandomErr() error {
	if rand.Intn(5) == 1 {
		return fmt.Errorf("random error")
	}
	return nil
}

// hashPassword returns a base64 encoded sha256 hash of the pwd + salt
func hashPassword(pwd string, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(pwd + salt))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
