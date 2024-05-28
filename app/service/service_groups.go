package service

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/toolbox/db_tool/sqldb"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Groups Service -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (s *Service) CreateGroup(ctx core.Ctx, req *pbs.CreateGroupRequest) (*pbs.CreateGroupResponse, error) {
	groupOwnerID, err := core.ToIntAndErr(s.Toolbox.ExtractMetadata(ctx, "user_id"))
	if err != nil {
		return nil, errs.GRPCInternal(err.Error())
	}

	invitedUserIDs := core.Int32Slice(req.InvitedUserIds).ToIntSlice()

	group, err := s.Toolbox.CreateGroup(ctx, req.Name, groupOwnerID, invitedUserIDs)
	if err != nil {
		return nil, errCallingGroupsDB(ctx, err)
	}

	return &pbs.CreateGroupResponse{Group: group.ToGroupInfoPB()}, nil
}

func (s *Service) GetGroup(ctx core.Ctx, req *pbs.GetGroupRequest) (*pbs.GetGroupResponse, error) {
	group, err := s.Toolbox.GetGroup(ctx, sqldb.WithID(req.GroupId))
	if err != nil {
		if s.Toolbox.IsNotFound(err) {
			return nil, errGroupNotFound()
		}
		return nil, errCallingGroupsDB(ctx, err)
	}

	return &pbs.GetGroupResponse{Group: group.ToGroupInfoPB()}, nil
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var (
	errGroupNotFound   = func() error { return errs.GRPCNotFound("group") }
	errCallingGroupsDB = func(ctx core.Ctx, err error) error {
		return errs.GRPCGroupsDBCall(err, core.RouteNameFromCtx(ctx), core.LogUnexpectedErr)
	}
)
