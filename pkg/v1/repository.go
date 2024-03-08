package v1

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID        int       `gorm:"primaryKey" bson:"id"`
	Username  string    `gorm:"unique;not null" bson:"username"`
	Password  string    `gorm:"not null" bson:"password"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type Database struct {
	*gorm.DB
}

func NewDatabase() *Database {
	database := &Database{}
	database.Connect()
	return database
}

type Repository interface {
	CreateUser(user User) (*User, error)
	GetUser(userID int, username string) (*User, error)
}

type repository struct {
	*Database
}

func NewRepository(database *Database) *repository {
	return &repository{Database: database}
}

func (database *Database) Connect() {
	var (
		username = "root"
		password = ""
		hostname = "localhost"
		port     = "3306"
		schema   = "grpc-gateway-impl"
		params   = "?charset=utf8&parseTime=True&loc=Local"

		err error
	)

	if database.DB, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", username, password, hostname, port, schema, params))); err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	database.DB.AutoMigrate(User{})
}

func (r *repository) CreateUser(user User) (*User, error) {
	if err := r.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}
	return &user, nil
}

func (r *repository) GetUser(userID int, username string) (*User, error) {
	var user User
	err := r.DB.Where("id = ? OR username = ?", userID, username).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error retrieving user: %v", err)
	}
	return &user, err // err can either be nil or gorm.ErrRecordNotFound
}
