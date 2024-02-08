package service

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/entities"
)

type ServiceLayer interface {
	Signup(ctx context.Context, in entities.SignupRequest) (entities.SignupResponse, error)
	Login(ctx context.Context, in entities.LoginRequest) (entities.LoginResponse, error)
}

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) Signup(ctx context.Context, in entities.SignupRequest) (entities.SignupResponse, error) {
	return entities.SignupResponse{}, nil
}

func (s *service) Login(ctx context.Context, in entities.LoginRequest) (entities.LoginResponse, error) {
	return entities.LoginResponse{}, nil
}
