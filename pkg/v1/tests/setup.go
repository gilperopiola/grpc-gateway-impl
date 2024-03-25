package tests

import (
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/mocks"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/db"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/options"
)

type setupGormMockFn func(*mocks.Gorm)

var emptyGormMockFn = func(*mocks.Gorm) {
	// Use this when a test case doesn't ever call any method on the mock.
	// Please.
}

// NewTestRepository returns a new testing repository with a gorm mock inside.
func NewTestRepository(setupMock setupGormMockFn) (repository.Repository, *mocks.Gorm) {
	gormMock := &mocks.Gorm{}                         // ------------------------------------------------ a GormMock. Can I remove the DBWrapper? T0D0
	dbWrapper := db.NewDB(nil, gormMock)              // -------------------- a DBWrapper which contains
	repository := repository.NewRepository(dbWrapper) // Repository contains
	setupMock(gormMock)
	return repository, gormMock
}

func copyUserPtr(user *models.User) *models.User {
	if user == nil {
		return nil
	}
	copy := *user
	return &copy
}

func queryOptions(option options.QueryOpt) []options.QueryOpt {
	if option == nil {
		return nil
	}
	return []options.QueryOpt{option}
}
