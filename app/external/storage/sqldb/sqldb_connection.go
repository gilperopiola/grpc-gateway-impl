package sqldb

import (
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/special_types"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/* ----------------------------------- */
/*           - SQL Database -          */
/* ----------------------------------- */

// NewGormDB returns a new configured instance of *gormAdapter, which implements SQLDB.
func NewGormDB(cfg *core.DatabaseCfg) *gormAdapter {
	gormDB, err := gormConnect(sqlConnString(cfg), gormConfig(zap.L(), cfg.LogLevel))
	if err != nil {
		zap.S().Fatalf(errs.FatalErrMsgConnectingDB, err)
	}

	// Migrate the models to the database.
	gormDB.AutoMigrate(models.AllModels...)

	// Insert the admin user if it doesn't exist.
	if cfg.InsertAdmin && cfg.AdminPwd != "" {
		InsertAdmin(gormDB, cfg.AdminPwd)
	}

	return gormDB
}

// gormConnect calls gorm.Open and wraps the returned *gorm.DB with our concrete type that implements SQLDB.
func gormConnect(dsn string, opts ...gorm.Option) (*gormAdapter, error) {
	gormDB, err := gorm.Open(mysql.Open(dsn), opts...)
	return newGormAdapter(gormDB), err
}

// gormConfig returns a new GormConfig with the given log level.
func gormConfig(l *zap.Logger, logLevel int) *gorm.Config {
	return &gorm.Config{
		Logger:         newGormLoggerAdapter(l, logLevel),
		TranslateError: true, // translate errors to gorm errors. idrk.
	}
}

// InsertAdmin inserts the admin user into the database if it doesn't exist.
func InsertAdmin(db special_types.SQLDB, adminPwd string) {
	admin := models.User{Username: "admin", Password: adminPwd, Role: models.AdminRole}
	if err := db.FirstOrCreate(&admin).Error(); err != nil {
		zap.S().Warnf(errs.InsertingDBAdmin, err)
	}
}

// sqlConnString returns the string needed to connect to a SQL DB.
func sqlConnString(c *core.DatabaseCfg) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", c.Username, c.Password, c.Hostname, c.Port, c.Schema, c.Params)
}
