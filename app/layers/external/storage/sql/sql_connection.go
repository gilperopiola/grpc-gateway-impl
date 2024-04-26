package sql

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - SQL Database -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// NewGormDB returns a new configured instance of *gormAdapter, which implements sql.
func NewGormDB(cfg *core.DBCfg) *sqlAdapter {
	gormDB, err := gormConnect(cfg.GetSQLConnString(), newGormConfig(zap.L(), cfg.LogLevel))
	core.LogPanicIfErr(err)

	if cfg.MigrateModels {
		gormDB.AutoMigrate(core.AllModels...)
	}

	if cfg.InsertAdmin && cfg.InsertAdminPwd != "" {
		insertAdmin(gormDB, cfg.InsertAdminPwd)
	}

	return gormDB
}

// gormConnect calls gorm.Open and wraps the returned *gorm.DB with our concrete type that implements sql.
func gormConnect(dsn string, opts ...gorm.Option) (*sqlAdapter, error) {
	gormDB, err := gorm.Open(mysql.Open(dsn), opts...)
	return newGormAdapter(gormDB), err
}

// newGormConfig returns a new GormConfig with the given log level.
func newGormConfig(l *zap.Logger, logLevel int) *gorm.Config {
	return &gorm.Config{
		Logger:         newSQLLogger(l, logLevel),
		TranslateError: true, // "translate errors to gorm errors". no idea.
	}
}

// insertAdmin inserts the admin user into the database if it doesn't exist.
func insertAdmin(db core.SQLDatabaseAPI, pwd string) {
	admin := core.User{Username: "admin", Password: pwd, Role: core.AdminRole}
	if err := db.FirstOrCreate(&admin).Error(); err != nil {
		zap.S().Warnf(errs.FailedToInsertDBAdmin, err)
	}
}
