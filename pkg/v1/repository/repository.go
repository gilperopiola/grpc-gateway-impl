package repository

import (
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/db"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/options"
)

/* ----------------------------------- */
/*            - Repository -           */
/* ----------------------------------- */

// Repository is the interface that wraps the basic methods to interact with the database.
type Repository interface {
	CreateUser(username, hashedPwd string) (*models.User, error)
	GetUser(opts ...options.QueryOption) (*models.User, error)
	GetUsers(page, pageSize int, opts ...options.QueryOption) (models.Users, int, error)
}

// repository is our concrete implementation of the Repository interface.
type repository struct {
	*db.DatabaseWrapper
}

// NewRepository returns a new instance of the repository.
func NewRepository(database *db.DatabaseWrapper) *repository {
	return &repository{DatabaseWrapper: database}
}

func getTotalPages(totalRecords, pageSize int) int {
	totalPages := totalRecords / pageSize
	if totalRecords%pageSize > 0 {
		totalPages++
	}
	return totalPages
}
