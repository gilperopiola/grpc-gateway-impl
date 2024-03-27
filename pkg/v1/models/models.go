package models

import (
	"time"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type PBResponse interface {
	Descriptor() ([]byte, []int)
	ProtoMessage()
	ProtoReflect() protoreflect.Message
	Reset()
	String() string
}

/* ----------------------------------- */
/*              - Models -             */
/* ----------------------------------- */

// AllModels is used to migrate all models at once.
var AllModels = []interface{}{
	User{},
}

type Users []*User

type User struct {
	ID        int    `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Role      Role   `gorm:"default:'default'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

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
	for _, user := range us {
		usersInfo = append(usersInfo, user.ToUserInfo())
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
