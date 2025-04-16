package db

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - SQL DB Tool: Groups -       */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Deprecated: Use repositories.GroupRepository instead
func (this *LegacyDB) DBCreateGroup(ctx god.Ctx, name string, ownerID int, invitedUserIDs []int) (*models.Group, error) {
	group := models.Group{Name: name, OwnerID: ownerID}

	if err := this.InnerDB.WithContext(ctx).Create(&group).Error(); err != nil {
		return nil, &errs.DBErr{Err: err, Context: "Failed to create group"}
	}

	for _, invitedUserID := range invitedUserIDs {
		if err := this.InnerDB.WithContext(ctx).Model(&group).Association("Invited").Append(&models.User{ID: invitedUserID}); err != nil {
			return nil, &errs.DBErr{Err: err, Context: "Failed to add invited user to group"}
		}
	}

	return &group, nil
}

// Deprecated: Use repositories.GroupRepository instead
func (this *LegacyDB) DBGetGroup(ctx god.Ctx, opts ...any) (*models.Group, error) {
	if len(opts) == 0 {
		return nil, &errs.DBErr{Err: nil, Context: NoOptionsErr}
	}

	query := this.InnerDB.Model(&models.Group{}).WithContext(ctx)
	for _, opt := range opts {
		opt.(core.SqlDBOpt)(query)
	}

	var group models.Group
	if err := query.First(&group).Error(); err != nil {
		return nil, &errs.DBErr{Err: err, Context: "Failed to get group"}
	}

	return &group, nil
}
