package repository

import (
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/db"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/options"
)

/* ----------------------------------- */
/*          - v1 Repository -          */
/* ----------------------------------- */

// Repository is the interface that wraps the basic methods to interact with the Database.
// App -> Service -> Repository -> DB.
type Repository interface {
	CreateUser(username, hashedPwd string) (*models.User, error)
	GetUser(opts ...options.QueryOpt) (*models.User, error)
	GetUsers(page, pageSize int, opts ...options.QueryOpt) (models.Users, int, error)
}

// repository is our concrete implementation of the Repository interface.
type repository struct {
	DB db.GormAdapter
}

// NewRepository returns a new instance of the repository.
func NewRepository(db db.GormAdapter) *repository {
	return &repository{db}
}
