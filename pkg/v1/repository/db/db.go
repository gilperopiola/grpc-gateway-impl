package db

import (
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/* ----------------------------------- */
/*             - Database -            */
/* ----------------------------------- */

// Because gorm is designed so badly, we do not only need to make an adapter for its logger, but also for the database itself.

// DBWrapper holds the connection to the database.
type DBWrapper struct {
	DB GormAdapter
}

// NewDatabaseWrapper returns a new instance of the Database.
func NewDatabaseWrapper(c *cfg.DBConfig) *DBWrapper {
	database := &DBWrapper{}
	database.Connect(c)
	return database
}

// Connect connects to the database using the given configuration.
func (dbw *DBWrapper) Connect(c *cfg.DBConfig) {
	var (
		connErr    error
		connStr    = getSQLConnectionString(c)
		connConfig = &gorm.Config{
			Logger: newGormLoggerAdapter(zap.L()),
		}
	)

	if dbw.DB, connErr = openGormAdapter(mysql.Open(connStr), connConfig); connErr != nil {
		zap.S().Fatalf(errs.FatalErrMsgConnectingDB, connErr)
	}

	// Migrate the models to the database.
	dbw.DB.AutoMigrate(models.AllModels...)

	// Insert the admin user if it doesn't exist.
	dbw.InsertAdmin(c.AdminPassword)
}

// InsertAdmin inserts the admin user if it doesn't exist.
func (dbw *DBWrapper) InsertAdmin(adminPwd string) {
	admin := models.User{
		Username: "admin",
		Password: adminPwd,
		Role:     models.AdminRole,
	}

	if err := dbw.DB.FirstOrCreate(&admin).Error(); err != nil {
		zap.S().Warnf(errs.ErrMsgInsertingAdmin, err)
	}
}

// Close closes the database connection.
func (dbw *DBWrapper) Close() {
	sqlDB := dbw.DB.GetSQLDB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func getSQLConnectionString(c *cfg.DBConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", c.Username, c.Password, c.Hostname, c.Port, c.Schema, c.Params)
}
