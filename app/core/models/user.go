package models

import (
	"time"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - User Model -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

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

type Users []*User

type UsersInGroup struct {
	UserID  int
	GroupID int
}
