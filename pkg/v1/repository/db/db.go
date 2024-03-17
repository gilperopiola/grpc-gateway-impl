package db

import (
	"fmt"
	"log"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/* ----------------------------------- */
/*             - Database -            */
/* ----------------------------------- */

// DatabaseWrapper holds the connection to the database.
type DatabaseWrapper struct {
	*gorm.DB
}

// NewDatabaseWrapper returns a new instance of the Database.
func NewDatabaseWrapper(c *cfg.DBConfig) *DatabaseWrapper {
	database := &DatabaseWrapper{}
	database.Connect(c)
	return database
}

// Connect connects to the database using the given configuration.
func (dbw *DatabaseWrapper) Connect(c *cfg.DBConfig) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", c.Username, c.Password, c.Hostname, c.Port, c.Schema, c.Params)

	var err error
	if dbw.DB, err = gorm.Open(mysql.Open(connStr)); err != nil {
		log.Fatalf(errs.FatalErrMsgConnectingDB, err)
	}

	// Migrate the models to the database.
	dbw.DB.AutoMigrate(models.AllModels...)

	// Insert the admin user if it doesn't exist.
	dbw.InsertAdmin(c.AdminPassword)
}

// InsertAdmin inserts the admin user if it doesn't exist.
func (dbw *DatabaseWrapper) InsertAdmin(adminPwd string) {
	admin := models.User{
		Username: "admin",
		Password: adminPwd,
		Role:     models.AdminRole,
	}

	if err := dbw.DB.Create(&admin).Error; err != nil {
		log.Printf(errs.ErrMsgInsertingAdmin, err)
	}
}

// Close closes the database connection.
func (dbw *DatabaseWrapper) Close() {
	sqlDB, err := dbw.DB.DB()
	if err != nil {
		log.Printf(errs.ErrMsgGettingDBConn, err)
		return
	}

	sqlDB.Close()
}
