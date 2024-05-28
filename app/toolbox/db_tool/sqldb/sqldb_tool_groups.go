package sqldb

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - SQL DB Tool: Groups -       */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (sdbt *sqlDBTool) CreateGroup(ctx core.Ctx, name string, ownerID int, invitedUserIDs []int) (*core.Group, error) {
	group := core.Group{Name: name, OwnerID: ownerID}

	if err := sdbt.DB.WithContext(ctx).Create(&group).Error(); err != nil {
		return nil, &errs.DBErr{err, "T0D0"}
	}

	for _, invitedUserID := range invitedUserIDs {
		if err := sdbt.DB.WithContext(ctx).Model(&group).Association("Invited").Append(&core.User{ID: invitedUserID}); err != nil {
			return nil, &errs.DBErr{err, "T0D0"}
		}
	}

	return &group, nil
}

func (sdbt *sqlDBTool) GetGroup(ctx core.Ctx, opts ...any) (*core.Group, error) {
	if len(opts) == 0 {
		return nil, &errs.DBErr{nil, NoOptionsErr}
	}

	query := sdbt.DB.Model(&core.Group{}).WithContext(ctx)
	for _, opt := range opts {
		opt.(core.SQLDBOpt)(query)
	}

	var group core.Group
	if err := query.First(&group).Error(); err != nil {
		return nil, &errs.DBErr{err, "T0D0"}
	}

	return &group, nil
}
