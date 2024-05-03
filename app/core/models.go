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
	Group{},
	UsersInGroup{},
}

/* -~-~-~-~-~- Users ~-~-~-~-~-~ */

type User struct {
	ID        int       `gorm:"primaryKey" bson:"_id"`
	Username  string    `gorm:"unique;not null" bson:"username"`
	Password  string    `gorm:"not null" bson:"password"`
	Role      Role      `gorm:"default:'default'" bson:"role"`
	Groups    []Group   `gorm:"many2many:users_in_groups" bson:"groups"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	Deleted   bool      `bson:"deleted"`
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

/* -~-~-~-~-~- Groups ~-~-~-~-~-~ */

type Group struct {
	ID        int       `gorm:"primaryKey" bson:"id"`
	OwnerID   int       `gorm:"not null" bson:"owner_id"`
	Name      string    `gorm:"not null" bson:"name"`
	Users     []User    `gorm:"many2many:users_in_groups" bson:"users"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	Deleted   bool      `bson:"deleted"`
}

type UsersInGroup struct {
	UserID  int
	GroupID int
}
