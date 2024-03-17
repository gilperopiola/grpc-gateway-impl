package repository

import (
	"errors"
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/options"

	"gorm.io/gorm"
)

/* ----------------------------------- */
/*         - Users Repository -        */
/* ----------------------------------- */

// CreateUser creates a new user in the database.
func (r *repository) CreateUser(username, hashedPwd string) (*models.User, error) {
	user := models.User{Username: username, Password: hashedPwd}
	if err := r.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	return &user, nil
}

// GetUser returns a user from the database.
func (r *repository) GetUser(opts ...options.QueryOption) (*models.User, error) {
	var user models.User

	query := r.DB.Model(&models.User{})
	for _, opt := range opts {
		query = opt(query)
	}

	err := query.First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &user, err // err can be nil or gorm.ErrRecordNotFound
}

// GetUsers returns a list of users from the database.
func (r *repository) GetUsers(page, pageSize int, opts ...options.QueryOption) (models.Users, int, error) {
	var users models.Users
	var totalRecords int64

	query := r.DB.Model(&models.User{})
	for _, opt := range opts {
		query = opt(query)
	}

	if err := query.Count(&totalRecords).Error; err != nil || totalRecords == 0 {
		return nil, 0, err
	}

	if err := query.Offset(page * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("error getting users: %w", err)
	}

	return users, getTotalPages(int(totalRecords), pageSize), nil
}
