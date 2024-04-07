package storage

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/storage/db"
	"github.com/gilperopiola/grpc-gateway-impl/app/storage/options"
)

/* ----------------------------------- */
/*          - v1 Storage -          */
/* ----------------------------------- */

// Storage is the interface that wraps the basic methods to interact with the Database.
// App -> Service -> Storage -> DB.
type Storage interface {
	CreateUser(username, hashedPwd string) (*models.User, error)
	GetUser(opts ...options.QueryOpt) (*models.User, error)
	GetUsers(page, pageSize int, opts ...options.QueryOpt) (models.Users, int, error)
}

// storage is our concrete implementation of the Storage interface.
type storage struct {
	DB db.DBAdapter
}

// NewStorage returns a new instance of the storage.
func NewStorage(db db.DBAdapter) *storage {
	return &storage{db}
}
