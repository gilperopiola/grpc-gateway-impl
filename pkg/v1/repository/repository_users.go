package repository

import (
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/options"
)

/* ----------------------------------- */
/*         - Users Repository -        */
/* ----------------------------------- */

// CreateUser creates a new user in the database.
func (r *repository) CreateUser(username, hashedPwd string) (*models.User, error) {
	user := models.User{Username: username, Password: hashedPwd}
	if err := r.DB.Create(&user).Error(); err != nil {
		return nil, &errs.DBError{errCreate, err}
	}
	return &user, nil
}

// GetUser returns a user from the database.
// At least one option must be provided, otherwise an error will be returned.
func (r *repository) GetUser(opts ...options.QueryOpt) (*models.User, error) {
	if len(opts) == 0 {
		return nil, &errs.DBError{ErrNoOpts, nil}
	}

	query := r.DB.Model(&models.User{})
	for _, opt := range opts {
		opt(query)
	}

	var user models.User
	if err := query.First(&user).Error(); err != nil {
		return nil, &errs.DBError{errGet, err}
	}

	return &user, nil
}

// GetUsers returns a list of users from the database.
func (r *repository) GetUsers(page, pageSize int, opts ...options.QueryOpt) (models.Users, int, error) {
	query := r.DB.Model(&models.User{})
	for _, opt := range opts {
		opt(query)
	}

	var totalMatchingUsers int64
	if err := query.Count(&totalMatchingUsers).Error(); err != nil {
		return nil, 0, &errs.DBError{errCount, err}
	}

	if totalMatchingUsers == 0 {
		return nil, 0, nil
	}

	var users models.Users
	if err := query.Offset(page * pageSize).Limit(pageSize).Find(&users).Error(); err != nil {
		return nil, 0, &errs.DBError{errGetMany, err}
	}

	return users, int(totalMatchingUsers), nil
}

/* ----------------------------------- */
/*     - Users Repository Errors -     */
/* ----------------------------------- */

var (
	errCreate  = errs.ErrMsgRepoCreatingUser
	errGet     = errs.ErrMsgRepoGettingUser
	errGetMany = errs.ErrMsgRepoGettingUsers
	errCount   = errs.ErrMsgRepoCountingUsers
	ErrNoOpts  = errs.ErrMsgRepoNoQueryOpts
)
