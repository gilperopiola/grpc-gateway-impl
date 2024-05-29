package core

import (
	"time"
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

type Users []*User

/* -~-~-~-~-~- Groups ~-~-~-~-~-~ */

type Group struct {
	ID        int       `gorm:"primaryKey" bson:"id"`
	OwnerID   int       `gorm:"not null" bson:"owner_id"`
	Name      string    `gorm:"not null" bson:"name"`
	Members   []User    `gorm:"many2many:users_in_groups" bson:"members"`
	Invited   []User    `gorm:"many2many:users_in_groups" bson:"invited"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	Deleted   bool      `bson:"deleted"`
}

type Groups []*Group

/* -~-~-~-~-~- UsersInGroup ~-~-~-~-~-~ */

type UsersInGroup struct {
	UserID  int
	GroupID int
}
