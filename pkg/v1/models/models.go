package models

import (
	"time"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
)

/* ----------------------------------- */
/*              - Models -             */
/* ----------------------------------- */

// AllModels is used to migrate all models at once.
var AllModels = []interface{}{
	User{},
}

type User struct {
	ID        int    `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Role      Role   `gorm:"default:'default'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Users []*User

type Role string

const (
	DefaultRole Role = "default"
	AdminRole   Role = "admin"
)

/* ----------------------------------- */
/*          - Model Methods -          */
/* ----------------------------------- */

func (us Users) ToUserInfo() []*usersPB.UserInfo {
	usersInfo := make([]*usersPB.UserInfo, 0, len(us))
	for _, u := range us {
		usersInfo = append(usersInfo, u.ToUserInfo())
	}
	return usersInfo
}

func (u User) ToUserInfo() *usersPB.UserInfo {
	return &usersPB.UserInfo{
		Id:        int32(u.ID),
		Username:  u.Username,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
