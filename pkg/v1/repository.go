package v1

import (
	"errors"
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/db"

	"gorm.io/gorm"
)

// Repository is the interface that wraps the basic methods to interact with the database.
type Repository interface {
	CreateUser(username, hashedPwd string) (*db.User, error)
	GetUser(userID int, username string) (*db.User, error)
}

// repository is our concrete implementation of the Repository interface.
type repository struct {
	*db.Database
}

// NewRepository returns a new instance of the repository.
func NewRepository(database *db.Database) *repository {
	return &repository{Database: database}
}

// CreateUser creates a new user in the database.
func (r *repository) CreateUser(username, hashedPwd string) (*db.User, error) {
	user := db.User{Username: username, Password: hashedPwd}
	if err := r.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}
	return &user, nil
}

// GetUser returns a user from the database.
func (r *repository) GetUser(userID int, username string) (*db.User, error) {
	var user db.User
	err := r.DB.Where("id = ? OR username = ?", userID, username).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return &user, err // err can either be nil or gorm.ErrRecordNotFound
}
