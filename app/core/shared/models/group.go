package models

import (
	"time"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - Group Model -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

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

func (Group) TableName() string {
	return "groups"
}
