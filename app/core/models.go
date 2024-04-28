package core

import (
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - Database Models -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Used to migrate all models at once.
var AllModels = []interface{}{
	User{},
}

/* -~-~-~-~-~- User ~-~-~-~-~-~ */

type User struct {
	ID        int    `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Role      Role   `gorm:"default:'default'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) ToUserInfoPB() *pbs.UserInfo {
	return &pbs.UserInfo{
		Id:        int32(u.ID),
		Username:  u.Username,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}

type Users []*User

func (us Users) ToUsersInfoPB() []*pbs.UserInfo {
	usersInfo := make([]*pbs.UserInfo, 0, len(us))
	for _, u := range us {
		usersInfo = append(usersInfo, u.ToUserInfoPB())
	}
	return usersInfo
}
