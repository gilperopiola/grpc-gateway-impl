package models

import (
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - User Model -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type User struct {
	ID        int         `gorm:"primaryKey" bson:"_id"`
	Username  string      `gorm:"unique;not null" bson:"username"`
	Password  string      `gorm:"not null" bson:"password"`
	Role      shared.Role `gorm:"default:'default'" bson:"role"`
	Groups    []Group     `gorm:"many2many:users_in_groups" bson:"groups"`
	CreatedAt time.Time   `bson:"created_at"`
	UpdatedAt time.Time   `bson:"updated_at"`
	Deleted   bool        `bson:"deleted"`
}

func (User) TableName() string {
	return "users"
}

type Users []*User

type UsersInGroup struct {
	UserID    int       `gorm:"primaryKey;column:user_id;index;" bson:"user_id"`
	GroupID   int       `gorm:"primaryKey;column:group_id;index;" bson:"group_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" bson:"created_at"`
}

func (UsersInGroup) TableName() string {
	return "users_in_groups"
}
