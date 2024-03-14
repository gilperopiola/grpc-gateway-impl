package db

import (
	"fmt"
	"log"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/* ----------------------------------- */
/*             - Database -            */
/* ----------------------------------- */

// Database holds the connection to the database.
type Database struct {
	*gorm.DB
}

// NewDatabase returns a new instance of the Database.
func NewDatabase(c *cfg.DBConfig) *Database {
	database := &Database{}
	database.Connect(c)
	return database
}

// Connect connects to the database using the given configuration.
func (db *Database) Connect(c *cfg.DBConfig) {
	connectionStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", c.Username, c.Password, c.Hostname, c.Port, c.Schema, c.Params)

	var err error
	if db.DB, err = gorm.Open(mysql.Open(connectionStr)); err != nil {
		log.Fatalf(errs.FatalErrMsgConnectingDB, err)
	}

	// Migrate the models to the database.
	db.DB.AutoMigrate(allModels...)

	// Insert the admin user if it doesn't exist.
	db.InsertAdmin(c.AdminPassword)
}

// InsertAdmin inserts the admin user if it doesn't exist.
func (db *Database) InsertAdmin(adminPwd string) {
	if err := db.DB.Create(&User{Username: "admin", Password: adminPwd, Role: AdminRole}).Error; err != nil {
		log.Printf("error inserting admin: %v\n", err)
	}
}

// Close closes the database connection.
func (db *Database) Close() {
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Printf("error getting db connection: %v\n", err)
		return
	}
	sqlDB.Close()
}
