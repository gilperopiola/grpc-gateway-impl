package sqldb

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
)

var (
	CreateUserErr = errs.DBCreatingUser
	GetUserErr    = errs.DBGettingUser
	GetUsersErr   = errs.DBGettingUsers
	CountUsersErr = errs.DBCountingUsers
	NoOptionsErr  = errs.DBNoQueryOpts
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - SQL DB Tool: Users -       */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// CreateUser creates a new user in the database.
func (sdbt *sqlDBTool) CreateUser(ctx context.Context, username, hashedPwd string) (*core.User, error) {
	user := core.User{Username: username, Password: hashedPwd}

	if err := sdbt.DB.WithContext(ctx).Create(&user).Error(); err != nil {
		return nil, &errs.DBErr{err, CreateUserErr}
	}

	return &user, nil
}

// GetUser returns a user from the database.
// At least one option must be provided, otherwise an error will be returned.
func (sdbt *sqlDBTool) GetUser(ctx context.Context, opts ...any) (*core.User, error) {
	if len(opts) == 0 {
		return nil, &errs.DBErr{nil, NoOptionsErr}
	}

	query := sdbt.DB.Model(&core.User{}).WithContext(ctx)
	for _, opt := range opts {
		opt.(core.SQLDBOpt)(query)
	}

	var user core.User
	if err := query.First(&user).Error(); err != nil {
		return nil, &errs.DBErr{err, GetUserErr}
	}

	return &user, nil
}

// GetUsers returns a list of users from the database.
func (sdbt *sqlDBTool) GetUsers(ctx context.Context, page, pageSize int, opts ...any) (core.Users, int, error) {
	query := sdbt.DB.Model(&core.User{}).WithContext(ctx)
	for _, opt := range opts {
		opt.(core.SQLDBOpt)(query)
	}

	var matchingUsers int64
	if err := query.Count(&matchingUsers).Error(); err != nil {
		return nil, 0, &errs.DBErr{err, CountUsersErr}
	}

	if matchingUsers == 0 {
		return nil, 0, nil
	}

	var users core.Users
	if err := query.Offset(page * pageSize).Limit(pageSize).Find(&users).Error(); err != nil {
		return nil, 0, &errs.DBErr{err, GetUsersErr}
	}

	return users, int(matchingUsers), nil
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
