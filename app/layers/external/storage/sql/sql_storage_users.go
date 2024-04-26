package sql

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - Users Storage -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// CreateUser creates a new user in the database.
func (r *SQLStorage) CreateUser(ctx context.Context, username, hashedPwd string) (*core.User, error) {
	user := core.User{Username: username, Password: hashedPwd}

	if err := r.DB.WithContext(ctx).Create(&user).Error(); err != nil {
		return nil, &errs.DBError{err, CreateUserErr}
	}

	return &user, nil
}

// GetUser returns a user from the database.
// At least one option must be provided, otherwise an error will be returned.
func (r *SQLStorage) GetUser(ctx context.Context, opts ...any) (*core.User, error) {
	if len(opts) == 0 {
		return nil, &errs.DBError{nil, NoOptionsErr}
	}

	query := r.DB.Model(&core.User{}).WithContext(ctx)
	for _, opt := range opts {
		opt.(core.SQLQueryOpt)(query)
	}

	var user core.User
	if err := query.First(&user).Error(); err != nil {
		return nil, &errs.DBError{err, GetUserErr}
	}

	return &user, nil
}

// GetUsers returns a list of users from the database.
func (r *SQLStorage) GetUsers(ctx context.Context, page, pageSize int, opts ...any) (core.Users, int, error) {
	query := r.DB.Model(&core.User{}).WithContext(ctx)
	for _, opt := range opts {
		opt.(core.SQLQueryOpt)(query)
	}

	var totalMatchingUsers int64
	if err := query.Count(&totalMatchingUsers).Error(); err != nil {
		return nil, 0, &errs.DBError{err, CountUsersErr}
	}

	if totalMatchingUsers == 0 {
		return nil, 0, nil
	}

	var users core.Users
	if err := query.Offset(page * pageSize).Limit(pageSize).Find(&users).Error(); err != nil {
		return nil, 0, &errs.DBError{err, GetUsersErr}
	}

	return users, int(totalMatchingUsers), nil
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - Users Storage Errors -      */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var (
	CreateUserErr = errs.DBCreatingUser
	GetUserErr    = errs.DBGettingUser
	GetUsersErr   = errs.DBGettingUsers
	CountUsersErr = errs.DBCountingUsers
	NoOptionsErr  = errs.DBNoQueryOpts
)
