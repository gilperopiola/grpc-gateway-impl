package repositories

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - Group Repository -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// GormGroupRepository implements the GroupRepository interface using GORM
type GormGroupRepository struct {
	db core.DBOperations
}

// Verify that GormGroupRepository implements the core.GroupRepository interface
var _ core.GroupRepository = (*GormGroupRepository)(nil)

// NewGormGroupRepository creates a new GormGroupRepository
func NewGormGroupRepository(db core.DBOperations) *GormGroupRepository {
	return &GormGroupRepository{db: db}
}

// CreateGroup creates a new group with the specified owner and invited users
func (r *GormGroupRepository) CreateGroup(ctx god.Ctx, name string, ownerID int, invitedUserIDs []int) (*models.Group, error) {
	group := models.Group{
		Name:    name,
		OwnerID: ownerID,
	}

	err := r.db.WithContext(ctx).CreateError(&group)
	if err != nil {
		return nil, &errs.DBErr{Err: err, Context: errs.FailedToCreateGroup}
	}

	// Add invited users to the group
	if len(invitedUserIDs) > 0 {
		err = r.addInvitedUsersToGroup(ctx, &group, invitedUserIDs)
		if err != nil {
			return nil, err
		}
	}

	return &group, nil
}

// GetGroupByID retrieves a group by its ID
func (r *GormGroupRepository) GetGroupByID(ctx god.Ctx, id int) (*models.Group, error) {
	var group models.Group

	err := r.db.WithContext(ctx).FirstError(&group, id)
	if err != nil {
		return nil, &errs.DBErr{Err: err, Context: errs.GroupNotFound}
	}

	return &group, nil
}

// GetGroupsByUserID retrieves all groups where the specified user is a member or owner
func (r *GormGroupRepository) GetGroupsByUserID(ctx god.Ctx, userID int) ([]*models.Group, error) {
	var groups []*models.Group

	// Find groups where the user is the owner
	err := r.db.WithContext(ctx).FindError(&groups, "owner_id = ?", userID)
	if err != nil {
		return nil, &errs.DBErr{Err: err, Context: errs.FailedToFetchGroups}
	}

	// TODO: Find groups where the user is a member (using the many-to-many relationship)
	// This would require a direct GORM query with joins, which is outside our simplified DB interface
	// We'll need to enhance the DB interface or use the helper methods

	return groups, nil
}

// addInvitedUsersToGroup is a helper method to add invited users to a group
func (r *GormGroupRepository) addInvitedUsersToGroup(ctx god.Ctx, group *models.Group, invitedUserIDs []int) error {
	// This would typically require direct access to GORM's Association method
	// In a real implementation, we might need to enhance our DB interface
	// or have a specialized method for many-to-many operations

	// For now, this is a simplified demonstration using a helper method
	// This assumes the GormDB implementation has the Association method exposed
	gormDB, ok := r.db.(interface {
		Association(column string) interface {
			Append(...interface{}) error
		}
	})

	if !ok {
		return &errs.DBErr{Err: nil, Context: "DB implementation doesn't support associations"}
	}

	for _, invitedUserID := range invitedUserIDs {
		err := gormDB.Association("Invited").Append(&models.User{ID: invitedUserID})
		if err != nil {
			return &errs.DBErr{Err: err, Context: errs.FailedToAddUserToGroup}
		}
	}

	return nil
}
