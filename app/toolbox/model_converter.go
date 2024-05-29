package toolbox

import (
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
)

var _ core.ModelConverter = (*modelConverter)(nil)

type modelConverter struct {
}

func NewModelConverter() core.ModelConverter {
	return &modelConverter{}
}

/* Users */

func (mc modelConverter) UserToUserInfoPB(user *core.User) *pbs.UserInfo {
	return &pbs.UserInfo{
		Id:        int32(user.ID),
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

func (mc modelConverter) UsersToUsersInfoPB(users core.Users) []*pbs.UserInfo {
	usersInfo := make([]*pbs.UserInfo, 0, len(users))
	for _, u := range users {
		usersInfo = append(usersInfo, mc.UserToUserInfoPB(u))
	}
	return usersInfo
}

/* Groups */

func (mc modelConverter) GroupToGroupInfoPB(group *core.Group) *pbs.GroupInfo {
	return &pbs.GroupInfo{
		Id:        int32(group.ID),
		Name:      group.Name,
		Owner:     &pbs.UserInfo{Id: int32(group.OwnerID)},
		CreatedAt: group.CreatedAt.Format(time.RFC3339),
		UpdatedAt: group.UpdatedAt.Format(time.RFC3339),
	}
}

func (mc modelConverter) GroupsToGroupsInfoPB(groups core.Groups) []*pbs.GroupInfo {
	groupsInfo := make([]*pbs.GroupInfo, 0, len(groups))
	for _, g := range groups {
		groupsInfo = append(groupsInfo, mc.GroupToGroupInfoPB(g))
	}
	return groupsInfo
}
