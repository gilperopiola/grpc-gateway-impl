package tests

import (
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/mocks"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/db"

	"github.com/stretchr/testify/assert"
)

/* ----------------------------------- */
/*     - Testing Repository Setup -    */
/* ----------------------------------- */

// I'm taking some risks here with my testing approach. It's not what I usually do, but I'm trying to
// keep the tests simple and just explore freely on this project.

// setupNewTest returns a new testing repository with a gorm mock inside.
func setupNewTest() TestRepository {
	gormMock := &mocks.GormMock{}                     // ------------------------------------------------ a GormMock. Can I remove the DBWrapper? T0D0
	dbWrapper := db.NewDatabaseWrapper(nil, gormMock) // -------------------- a DBWrapper which contains
	repository := repository.NewRepository(dbWrapper) // Repository contains
	return &testRepository{repository, gormMock}
}

// TestRepository adds the GetGormMock method to the Repository interface
// to retrieve the underlying mock if needed.
type TestRepository interface {
	repository.Repository
	GetGormMock() *mocks.GormMock
}

// testRepository is a wrapper around the Repository interface that adds a GetGormMock method.
type testRepository struct {
	repository.Repository
	gormMock *mocks.GormMock
}

// GetGormMock returns the underlying gorm mock.
func (r *testRepository) GetGormMock() *mocks.GormMock {
	return r.gormMock
}

func TestNewRepository(t *testing.T) {
	repository := repository.NewRepository(db.NewDatabaseWrapper(nil, &mocks.GormMock{}))
	assert.NotNil(t, repository)
}
