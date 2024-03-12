package db

import (
	"time"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
)

var allModels = []interface{}{
	User{},
}

type User struct {
	ID        int    `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) ToUserInfo() *usersPB.UserInfo {
	return &usersPB.UserInfo{
		Id:        int32(u.ID),
		Username:  u.Username,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
