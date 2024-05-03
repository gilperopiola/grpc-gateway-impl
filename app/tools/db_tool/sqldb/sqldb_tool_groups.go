package sqldb

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - SQL DB Tool: Groups -       */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (sdbt *sqlDBTool) CreateGroup(ctx context.Context, name string, ownerID int) (*core.Group, error) {
	group := core.Group{Name: name, OwnerID: ownerID}

	if err := sdbt.DB.WithContext(ctx).Create(&group).Error(); err != nil {
		return nil, &errs.DBErr{err, "T0D0"}
	}

	return &group, nil
}

func (sdbt *sqlDBTool) GetGroup(ctx context.Context, opts ...any) (*core.Group, error) {
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
