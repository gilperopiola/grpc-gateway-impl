package storage

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/storage/options"
)

/* ----------------------------------- */
/*         - Users Storage -        */
/* ----------------------------------- */

// CreateUser creates a new user in the database.
func (r *storage) CreateUser(username, hashedPwd string) (*models.User, error) {
	user := models.User{Username: username, Password: hashedPwd}
	if err := r.DB.Create(&user).Error(); err != nil {
		return nil, &errs.DBError{err, CreateUserErr}
	}
	return &user, nil
}

// GetUser returns a user from the database.
// At least one option must be provided, otherwise an error will be returned.
func (r *storage) GetUser(opts ...options.QueryOpt) (*models.User, error) {
	if len(opts) == 0 {
		return nil, &errs.DBError{nil, NoOptionsErr}
	}

	query := r.DB.Model(&models.User{})
	for _, opt := range opts {
		opt(query)
	}

	var user models.User
	if err := query.First(&user).Error(); err != nil {
		return nil, &errs.DBError{err, GetUserErr}
	}

	return &user, nil
}

// GetUsers returns a list of users from the database.
func (r *storage) GetUsers(page, pageSize int, opts ...options.QueryOpt) (models.Users, int, error) {
	query := r.DB.Model(&models.User{})
	for _, opt := range opts {
		opt(query)
	}

	var totalMatchingUsers int64
	if err := query.Count(&totalMatchingUsers).Error(); err != nil {
		return nil, 0, &errs.DBError{err, CountUsersErr}
	}

	if totalMatchingUsers == 0 {
		return nil, 0, nil
	}

	var users models.Users
	if err := query.Offset(page * pageSize).Limit(pageSize).Find(&users).Error(); err != nil {
		return nil, 0, &errs.DBError{err, GetUsersErr}
	}

	return users, int(totalMatchingUsers), nil
}

/* ----------------------------------- */
/*     - Users Storage Errors -     */
/* ----------------------------------- */

var (
	CreateUserErr = errs.ErrMsgRepoCreatingUser
	GetUserErr    = errs.ErrMsgRepoGettingUser
	GetUsersErr   = errs.ErrMsgRepoGettingUsers
	CountUsersErr = errs.ErrMsgRepoCountingUsers
	NoOptionsErr  = errs.ErrMsgRepoNoQueryOpts
)
