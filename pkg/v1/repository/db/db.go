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

// NewDB returns a new instance of the Database, either by connecting with a given configuration
// or by using an already configured database. This is useful for testing.
func NewDB(c *cfg.DBCfg) GormAdapter {
	gormDB, err := gormConnect(connectionString(c), gormConfigOption(zap.L(), c.GormLogLevel))
	if err != nil {
		zap.S().Fatalf(errs.FatalErrMsgConnectingDB, err)
	}

	// Migrate the models to the database.
	gormDB.AutoMigrate(models.AllModels...)

	// Insert the admin user if it doesn't exist.
	if c.InsertAdmin && c.AdminPwd != "" {
		gormInsertAdmin(gormDB, c.AdminPwd)
	}

	return gormDB
}

// gormConnect calls gorm.Open and wraps the returned *gorm.DB with our concrete type that implements GormAdapter.
func gormConnect(dsn string, opts ...gorm.Option) (GormAdapter, error) {
	gormDB, err := gorm.Open(mysql.Open(dsn), opts...)
	return newGormAdapter(gormDB), err
}

// gormConfigOption returns a new GormConfig with the given log level.
func gormConfigOption(l *zap.Logger, logLevel int) *gorm.Config {
	return &gorm.Config{
		Logger:         newGormLoggerAdapter(l, logLevel),
		TranslateError: true,
	}
}

// gormInsertAdmin inserts the admin user into the database if it doesn't exist.
func gormInsertAdmin(gormDB GormAdapter, adminPwd string) {
	admin := models.User{Username: "admin", Password: adminPwd, Role: models.AdminRole}
	if err := gormDB.FirstOrCreate(&admin).Error(); err != nil {
		zap.S().Warnf(errs.ErrMsgInsertingAdmin, err)
	}
}

// connectionString returns the connection string for the database.
// Only MySQL is supported for now.
func connectionString(c *cfg.DBCfg) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", c.Username, c.Password, c.Hostname, c.Port, c.Schema, c.Params)
}
