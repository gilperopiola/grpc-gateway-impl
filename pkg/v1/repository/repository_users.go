package repository

import (
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/options"
)

/* ----------------------------------- */
/*         - Users Repository -        */
/* ----------------------------------- */

// CreateUser creates a new user in the database.
func (r *repository) CreateUser(username, hashedPwd string) (*models.User, error) {
	user := models.User{Username: username, Password: hashedPwd}
	if err := r.DBWrapper.DB.Create(&user).Error(); err != nil {
		return nil, ErrCreatingUser(err)
	}
	return &user, nil
}

// GetUser returns a user from the database.
func (r *repository) GetUser(opts ...options.QueryOption) (*models.User, error) {
	var user models.User

	query := r.DBWrapper.DB.Model(&models.User{})
	for _, opt := range opts {
		opt(query)
	}

	if err := query.First(&user).Error(); err != nil {
		return nil, ErrGettingUser(err)
	}

	return &user, nil
}

// GetUsers returns a list of users from the database.
func (r *repository) GetUsers(page, pageSize int, opts ...options.QueryOption) (models.Users, int, error) {
	var users models.Users
	var totalMatchingUsers int64

	query := r.DBWrapper.DB.Model(&models.User{})
	for _, opt := range opts {
		opt(query)
	}

	if err := query.Count(&totalMatchingUsers).Error(); err != nil || totalMatchingUsers == 0 {
		return nil, 0, err
	}

	if err := query.Offset(page * pageSize).Limit(pageSize).Find(&users).Error(); err != nil {
		return nil, 0, ErrGettingUsers(err)
	}

	return users, int(totalMatchingUsers), nil
}

/* ----------------------------------- */
/*     - Users Repository Errors -     */
/* ----------------------------------- */

var (
	ErrCreatingUser = func(err error) error {
		return fmt.Errorf("repository error -> creating user -> %w", err)
	}
	ErrGettingUser = func(err error) error {
		return fmt.Errorf("repository error -> getting user -> %w", err)
	}
	ErrGettingUsers = func(err error) error {
		return fmt.Errorf("repository error -> getting users -> %w", err)
	}
)
