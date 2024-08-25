package sqldb

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/types/models"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - SQL DB Tool: Groups -       */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (sdbt *sqlDBTool) CreateGroup(ctx god.Ctx, name string, ownerID int, invitedUserIDs []int) (*models.Group, error) {
	group := models.Group{Name: name, OwnerID: ownerID}

	if err := sdbt.DB.WithContext(ctx).Create(&group).Error(); err != nil {
		return nil, &errs.DBErr{err, "T0D0"}
	}

	for _, invitedUserID := range invitedUserIDs {
		if err := sdbt.DB.WithContext(ctx).Model(&group).Association("Invited").Append(&models.User{ID: invitedUserID}); err != nil {
			return nil, &errs.DBErr{err, "T0D0"}
		}
	}

	return &group, nil
}

func (sdbt *sqlDBTool) GetGroup(ctx god.Ctx, opts ...any) (*models.Group, error) {
	if len(opts) == 0 {
		return nil, &errs.DBErr{nil, NoOptionsErr}
	}

	query := sdbt.DB.Model(&models.Group{}).WithContext(ctx)
	for _, opt := range opts {
		opt.(core.SqlDBOpt)(query)
	}

	var group models.Group
	if err := query.First(&group).Error(); err != nil {
		return nil, &errs.DBErr{err, "T0D0"}
	}

	return &group, nil
}
