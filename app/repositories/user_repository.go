package repositories

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - User Repository -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// GormUserRepository implements the UserRepository interface using GORM
type GormUserRepository struct {
	db core.DBOperations
}

// Verify that GormUserRepository implements the core.UserRepository interface
var _ core.UserRepository = (*GormUserRepository)(nil)

// NewGormUserRepository creates a new GormUserRepository
func NewGormUserRepository(db core.DBOperations) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// CreateUser creates a new user with the specified username and hashed password
func (r *GormUserRepository) CreateUser(ctx god.Ctx, username, hashedPwd string) (*models.User, error) {
	user := models.User{
		Username: username,
		Password: hashedPwd,
	}

	err := r.db.WithContext(ctx).CreateError(&user)
	if err != nil {
		return nil, &errs.DBErr{Err: err, Context: errs.FailedToCreateUser}
	}

	return &user, nil
}

// GetUserByID retrieves a user by their ID
func (r *GormUserRepository) GetUserByID(ctx god.Ctx, id int) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).FirstError(&user, id)
	if err != nil {
		return nil, &errs.DBErr{Err: err, Context: errs.UserNotFound}
	}

	return &user, nil
}

// GetUserByUsername retrieves a user by their username
func (r *GormUserRepository) GetUserByUsername(ctx god.Ctx, username string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).FirstError(&user, "username = ?", username)
	if err != nil {
		return nil, &errs.DBErr{Err: err, Context: errs.UserNotFound}
	}

	return &user, nil
}

// GetUsers retrieves a paginated list of users
func (r *GormUserRepository) GetUsers(ctx god.Ctx, page, pageSize int) ([]*models.User, int, error) {
	var users []*models.User
	var count int64

	// First, get the total count for pagination
	countErr := r.db.WithContext(ctx).(core.InnerDB).Model(&models.User{}).Count(&count).Error()

	if countErr != nil {
		return nil, 0, &errs.DBErr{Err: countErr, Context: errs.FailedToFetchUsers}
	}

	// Then, get the paginated users
	offset := (page - 1) * pageSize

	// We need to use a special type assertion because our simplified DB interface
	// doesn't support the Offset and Limit methods directly
	limitOffsetDB, ok := r.db.WithContext(ctx).(core.InnerDB)
	if !ok {
		return nil, 0, &errs.DBErr{Err: nil, Context: "DB implementation doesn't support pagination"}
	}

	err := limitOffsetDB.Offset(offset).Limit(pageSize).Find(&users)
	if err != nil {
		return nil, 0, &errs.DBErr{Err: err.Error(), Context: errs.FailedToFetchUsers}
	}

	return users, int(count), nil
}
