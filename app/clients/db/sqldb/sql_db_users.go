package sqldb

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
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

// DBCreateUser creates a new user in the database.
func (sdbt *DB) DBCreateUser(ctx god.Ctx, username, hashedPwd string) (*models.User, error) {
	user := models.User{Username: username, Password: hashedPwd}

	if err := sdbt.InnerDB.WithContext(ctx).Create(&user).Error(); err != nil {
		return nil, &errs.DBErr{err, CreateUserErr}
	}

	return &user, nil
}

// DBGetUser returns a user from the database.
// At least one option must be provided, otherwise an error will be returned.
func (sdbt *DB) DBGetUser(ctx god.Ctx, opts ...any) (*models.User, error) {
	if len(opts) == 0 {
		return nil, &errs.DBErr{nil, NoOptionsErr}
	}

	query := sdbt.InnerDB.Model(&models.User{}).WithContext(ctx)
	for _, opt := range opts {
		opt.(core.SqlDBOpt)(query)
	}

	var user models.User
	if err := query.First(&user).Error(); err != nil {
		return nil, &errs.DBErr{err, GetUserErr}
	}

	return &user, nil
}

// DBGetUsers returns a list of users from the database.
func (sdbt *DB) DBGetUsers(ctx god.Ctx, page, pageSize int, opts ...any) ([]*models.User, int, error) {
	query := sdbt.InnerDB.Model(&models.User{}).WithContext(ctx)
	for _, opt := range opts {
		opt.(core.SqlDBOpt)(query)
	}

	var matchingUsers int64
	if err := query.Count(&matchingUsers).Error(); err != nil {
		return nil, 0, &errs.DBErr{err, CountUsersErr}
	}

	if matchingUsers == 0 {
		return nil, 0, nil
	}

	var users []*models.User
	if err := query.Offset(page * pageSize).Limit(pageSize).Find(&users).Error(); err != nil {
		return nil, 0, &errs.DBErr{err, GetUsersErr}
	}

	return users, int(matchingUsers), nil
}
