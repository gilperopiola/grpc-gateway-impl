package db

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

/* ----------------------------------- */
/*            - Repository -           */
/* ----------------------------------- */

// Repository is the interface that wraps the basic methods to interact with the database.
type Repository interface {
	CreateUser(username, hashedPwd string) (*User, error)
	GetUser(opts ...QueryOption) (*User, error)
	GetUsers(page, pageSize int, opts ...QueryOption) ([]*User, int, error)
}

// repository is our concrete implementation of the Repository interface.
type repository struct {
	*Database
}

// NewRepository returns a new instance of the repository.
func NewRepository(database *Database) *repository {
	return &repository{Database: database}
}

/* ----------------------------------- */
/*         - Users Repository -        */
/* ----------------------------------- */

// CreateUser creates a new user in the database.
func (r *repository) CreateUser(username, hashedPwd string) (*User, error) {
	user := User{Username: username, Password: hashedPwd}
	if err := r.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	return &user, nil
}

// GetUser returns a user from the database.
func (r *repository) GetUser(opts ...QueryOption) (*User, error) {
	var user User

	query := r.DB.Model(&User{})
	for _, opt := range opts {
		query = opt(query)
	}

	err := query.First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &user, err // err can either be nil or gorm.ErrRecordNotFound
}

// GetUsers returns a list of users from the database.
func (r *repository) GetUsers(page, pageSize int, opts ...QueryOption) ([]*User, int, error) {
	var users []*User
	var totalRecords int64

	query := r.DB.Model(&User{})
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

func getTotalPages(totalRecords, pageSize int) int {
	totalPages := totalRecords / pageSize
	if totalRecords%pageSize > 0 {
		totalPages++
	}
	return totalPages
}
