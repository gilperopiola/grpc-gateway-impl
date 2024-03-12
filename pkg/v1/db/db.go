package db

import (
	"fmt"
	"log"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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

	db.DB.AutoMigrate(allModels...)
}
