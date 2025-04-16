package db

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
// Deprecated: Use repositories.UserRepository instead
func (this *LegacyDB) DBCreateUser(ctx god.Ctx, username, hashedPwd string) (*models.User, error) {
	user := models.User{Username: username, Password: hashedPwd}

	if err := this.InnerDB.WithContext(ctx).Create(&user).Error(); err != nil {
		return nil, &errs.DBErr{Err: err, Context: CreateUserErr}
	}

	return &user, nil
}

// DBGetUser gets a user from the database.
// Deprecated: Use repositories.UserRepository instead
func (this *LegacyDB) DBGetUser(ctx god.Ctx, opts ...any) (*models.User, error) {
	if len(opts) == 0 {
		return nil, &errs.DBErr{Err: nil, Context: NoOptionsErr}
	}

	query := this.InnerDB.Model(&models.User{}).WithContext(ctx)
	for _, opt := range opts {
		opt.(core.SqlDBOpt)(query)
	}

	var user models.User
	if err := query.First(&user).Error(); err != nil {
		return nil, &errs.DBErr{Err: err, Context: GetUserErr}
	}

	return &user, nil
}

// DBGetUsers gets users from the database.
// Deprecated: Use repositories.UserRepository instead
func (this *LegacyDB) DBGetUsers(ctx god.Ctx, page, pageSize int, opts ...any) ([]*models.User, int, error) {
	var users []*models.User
	var count int64

	query := this.InnerDB.Model(&models.User{}).WithContext(ctx)
	for _, opt := range opts {
		opt.(core.SqlDBOpt)(query)
	}

	if err := query.Count(&count).Error(); err != nil {
		return nil, 0, &errs.DBErr{Err: err, Context: CountUsersErr}
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error(); err != nil {
		return nil, 0, &errs.DBErr{Err: err, Context: GetUsersErr}
	}

	return users, int(count), nil
}
