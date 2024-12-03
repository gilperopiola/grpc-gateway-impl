package tools

import (
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/pbs"
)

// Would this be better off inside of utils/converter or something like that?

var _ core.ModelConverter = &modelConverter{}

type modelConverter struct{}

func NewModelConverter() core.ModelConverter {
	return &modelConverter{}
}

// 🔻 Users 🔻

func (this modelConverter) UserToUserInfoPB(user *models.User) *pbs.UserInfo {
	return &pbs.UserInfo{
		Id:        int32(user.ID),
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

func (this modelConverter) UsersToUsersInfoPB(users []*models.User) []*pbs.UserInfo {
	usersInfo := make([]*pbs.UserInfo, 0, len(users))
	for _, u := range users {
		usersInfo = append(usersInfo, this.UserToUserInfoPB(u))
	}
	return usersInfo
}

// 🔻 Groups 🔻

func (this modelConverter) GroupToGroupInfoPB(group *models.Group) *pbs.GroupInfo {
	return &pbs.GroupInfo{
		Id:        int32(group.ID),
		Name:      group.Name,
		Owner:     &pbs.UserInfo{Id: int32(group.OwnerID)},
		CreatedAt: group.CreatedAt.Format(time.RFC3339),
		UpdatedAt: group.UpdatedAt.Format(time.RFC3339),
	}
}

func (this modelConverter) GroupsToGroupsInfoPB(groups []*models.Group) []*pbs.GroupInfo {
	groupsInfo := make([]*pbs.GroupInfo, 0, len(groups))
	for _, group := range groups {
		groupsInfo = append(groupsInfo, this.GroupToGroupInfoPB(group))
	}
	return groupsInfo
}
