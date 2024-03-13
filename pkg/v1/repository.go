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
	GetUsers(page, pageSize int, filter string) ([]*db.User, int, error)
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
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	return &user, nil
}

// GetUser returns a user from the database.
func (r *repository) GetUser(userID int, username string) (*db.User, error) {
	var user db.User
	err := r.DB.Where("id = ? OR username = ?", userID, username).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &user, err // err can either be nil or gorm.ErrRecordNotFound
}

// GetUsers returns a list of users from the database.
func (r *repository) GetUsers(page, pageSize int, filter string) ([]*db.User, int, error) {
	var users []*db.User
	var totalRecords int64

	query := r.DB.Model(&db.User{})
	if filter != "" {
		query = query.Where("username LIKE ?", "%"+filter+"%")
	}

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting users: %w", err)
	}

	if totalRecords == 0 {
		return users, 0, nil
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
