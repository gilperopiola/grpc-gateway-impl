package repository

import (
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/db"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/options"
)

/* ----------------------------------- */
/*          - v1 Repository -          */
/* ----------------------------------- */

// Repository is the interface that wraps the basic methods to interact with the database.
// App -> Service -> Repository -> Database.
type Repository interface {
	CreateUser(username, hashedPwd string) (*models.User, error)
	GetUser(opts ...options.QueryOption) (*models.User, error)
	GetUsers(page, pageSize int, opts ...options.QueryOption) (models.Users, int, error)
}

// repository is our concrete implementation of the Repository interface.
type repository struct {
	*db.DBWrapper
}

// NewRepository returns a new instance of the repository.
func NewRepository(database *db.DBWrapper) *repository {
	return &repository{DBWrapper: database}
}
