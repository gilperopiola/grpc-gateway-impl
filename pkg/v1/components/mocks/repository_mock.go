package mocks

import (
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/options"

	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

type RepoGetUserReturns struct {
	User *models.User
	Err  error
}

type RepoCreateUserReturns struct {
	User *models.User
	Err  error
}

type RepoGetUsersReturns struct {
	Users models.Users
	Total int
	Err   error
}

func (m *Repository) CreateUser(username, hashedPwd string) (*models.User, error) {
	args := m.Called(username, hashedPwd)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *Repository) GetUser(opts ...options.QueryOpt) (*models.User, error) {
	args := m.Called(opts)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *Repository) GetUsers(page, pageSize int, opts ...options.QueryOpt) (models.Users, int, error) {
	args := m.Called(page, pageSize, opts)
	return args.Get(0).(models.Users), args.Int(1), args.Error(2)
}
