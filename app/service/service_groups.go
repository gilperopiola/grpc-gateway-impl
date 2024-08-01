package service

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/db_tool/sqldb"
)

type GroupsSubService struct {
	pbs.UnimplementedGroupsServiceServer
	Tools core.Tools
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Groups Service -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (s *GroupsSubService) CreateGroup(ctx god.Ctx, req *pbs.CreateGroupRequest) (*pbs.CreateGroupResponse, error) {
	groupOwnerID, err := god.ToIntAndErr(s.Tools.GetFromCtx(ctx, "user_id"))
	if err != nil {
		return nil, errs.GRPCInternal(err.Error())
	}

	invitedUserIDs := god.Int32Slice(req.InvitedUserIds).ToIntSlice()

	group, err := s.Tools.CreateGroup(ctx, req.Name, groupOwnerID, invitedUserIDs)
	if err != nil {
		return nil, errCallingGroupsDB(ctx, err)
	}

	return &pbs.CreateGroupResponse{Group: s.Tools.GroupToGroupInfoPB(group)}, nil
}

func (s *GroupsSubService) GetGroup(ctx god.Ctx, req *pbs.GetGroupRequest) (*pbs.GetGroupResponse, error) {
	group, err := s.Tools.GetGroup(ctx, sqldb.WithID(req.GroupId))
	if err != nil {
		if s.Tools.IsNotFound(err) {
			return nil, errGroupNotFound()
		}
		return nil, errCallingGroupsDB(ctx, err)
	}

	return &pbs.GetGroupResponse{Group: s.Tools.GroupToGroupInfoPB(group)}, nil
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var (
	errGroupNotFound   = func() error { return errs.GRPCNotFound("group") }
	errCallingGroupsDB = func(ctx god.Ctx, err error) error {
		core.LogUnexpected(err)
		return errs.GRPCGroupsDBCall(err, core.RouteNameFromCtx(ctx))
	}
)
