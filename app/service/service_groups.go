package service

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/utils"
)

type GroupSvc struct {
	pbs.UnimplementedGroupsServiceServer
	Clients core.Clients
	Tools   core.Tools
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Groups Service -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (s *GroupSvc) CreateGroup(ctx god.Ctx, req *pbs.CreateGroupRequest) (*pbs.CreateGroupResponse, error) {
	groupOwnerID, err := god.ToIntAndErr(s.Tools.GetFromCtx(ctx, "user_id"))
	if err != nil {
		return nil, errs.GRPCFromDB(err, core.GetRouteFromCtx(ctx).Name)
	}

	invitedUserIDs := utils.Int32Slice(req.InvitedUserIds).ToIntSlice()
	group, err := s.Clients.DBCreateGroup(ctx, req.Name, groupOwnerID, invitedUserIDs)
	if err != nil {
		return nil, errs.GRPCFromDB(err, core.GetRouteFromCtx(ctx).Name)
	}

	return &pbs.CreateGroupResponse{Group: s.Tools.GroupToGroupInfoPB(group)}, nil
}

func (s *GroupSvc) GetGroup(ctx god.Ctx, req *pbs.GetGroupRequest) (*pbs.GetGroupResponse, error) {
	group, err := s.Clients.DBGetGroup(ctx, core.WithID(req.GroupId))
	if err != nil {
		if errs.IsDBNotFound(err) {
			return nil, errs.GRPCNotFound("group", int(req.GroupId))
		}
		return nil, errs.GRPCFromDB(err, core.GetRouteFromCtx(ctx).Name)
	}

	return &pbs.GetGroupResponse{Group: s.Tools.GroupToGroupInfoPB(group)}, nil
}
