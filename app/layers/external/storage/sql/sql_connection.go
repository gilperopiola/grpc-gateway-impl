package sql

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - SQL Database -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// NewGormDB returns a new configured instance of *gormAdapter, which implements sql.
func NewGormDB(dbCfg DBConfigAccessor) *gormAdapter {
	gormDB, err := gormConnect(dbCfg.GetSQLConnString(), gormConfig(zap.L(), dbCfg.GetLogLevel()))
	if err != nil {
		zap.S().Fatalf(errs.FatalErrMsgConnectingDB, err)
	}

	// Migrate the models to the database.
	if dbCfg.ShouldMigrate() {
		gormDB.AutoMigrate(models.AllModels...)
	}

	// Insert the admin user if it doesn't exist.
	if dbCfg.ShouldInsertAdmin() && dbCfg.GetAdminPwd() != "" {
		InsertAdmin(gormDB, dbCfg.GetAdminPwd())
	}

	return gormDB
}

// Used to avoid circular dependencies with the core package.
type DBConfigAccessor interface {
	GetSQLConnString() string
	GetLogLevel() int
	ShouldMigrate() bool
	ShouldInsertAdmin() bool
	GetAdminPwd() string
}

// gormConnect calls gorm.Open and wraps the returned *gorm.DB with our concrete type that implements sql.
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
func InsertAdmin(db core.SQLDatabaseAPI, adminPwd string) {
	admin := models.User{Username: "admin", Password: adminPwd, Role: models.AdminRole}
	if err := db.FirstOrCreate(&admin).Error(); err != nil {
		zap.S().Warnf(errs.InsertingDBAdmin, err)
	}
}
