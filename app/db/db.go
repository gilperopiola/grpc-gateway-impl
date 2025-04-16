package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/utils"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// LegacyDB is the old DB implementation, kept for backward compatibility
// It should be gradually phased out in favor of GormDB
type LegacyDB struct {
	core.InnerDB
}

// Verify that LegacyDB implements the core.DBOperations interface
var _ core.DBOperations = (*LegacyDB)(nil)

// These methods implement the DBOperations interface for backward compatibility
func (d *LegacyDB) Find(out any, where ...any) error {
	return d.InnerDB.Find(out, where...).Error()
}

func (d *LegacyDB) First(out any, where ...any) error {
	return d.InnerDB.First(out, where...).Error()
}

func (d *LegacyDB) Create(value any) error {
	return d.InnerDB.Create(value).Error()
}

func (d *LegacyDB) Save(value any) error {
	return d.InnerDB.Save(value).Error()
}

func (d *LegacyDB) Delete(value any, where ...any) error {
	return d.InnerDB.Delete(value, where...).Error()
}

func (d *LegacyDB) WithContext(ctx context.Context) core.DBOperations {
	return &LegacyDB{d.InnerDB.WithContext(ctx).(core.InnerDB)}
}

func (d *LegacyDB) Transaction(fn func(tx core.DBOperations) error) error {
	return fmt.Errorf("Transaction not implemented in legacy DB")
}

func (d *LegacyDB) Close() error {
	d.InnerDB.Close()
	return nil
}

// NewSQLDBConn creates a legacy DB connection
// This function is kept for backward compatibility but should be deprecated
// in favor of NewGormDB
func NewSQLDBConn(cfg *core.DBCfg, hashPwdFn func(string) string) core.DBOperations {

	// We use our Logger wrapped inside of a gorm adapter
	dbLogger := newDBLogger(zap.L(), cfg.LogLevel)
	gormCfg := &gorm.Config{Logger: dbLogger, TranslateError: true}

	// We wrap this to match the signature of [utils.RetryFunc]
	var connectToDB = func() (any, error) {
		gormDB, err := gorm.Open(mysql.Open(cfg.GetSQLConnString()), gormCfg)
		return &innerDB{gormDB}, err
	}

	// We wrap this to match the signature of [utils.RetryFuncNoErr]
	var createDB = func() {
		db, err := sql.Open("mysql", cfg.GetSQLConnStringNoDB())
		if err != nil {
			logs.LogResult("Error trying to connect to SQL instance ", err)
			return
		}
		defer db.Close()
		if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", cfg.Database)); err != nil {
			logs.LogResult("Error trying to create SQL DB", err)
		}
	}

	// We try to connect to the DB directly.
	// If it fails, we try to connect without specifying the DB and then creating it.
	// If it fails, we retry this process a number of times.
	dbConn, err := utils.TryAndRetry(connectToDB, cfg.Retries, false, createDB)
	logs.LogFatalIfErr(err, errs.FailedDBConn)
	logs.LogResult("DB Connection for "+cfg.Database, nil)

	innerDB := dbConn.(*innerDB)
	legacyPostDBConnActions(innerDB, cfg, hashPwdFn)
	return &LegacyDB{innerDB}
}

// legacyPostDBConnActions is a renamed version of the original function
// to avoid conflict with the new setupDBPostConnection function
func legacyPostDBConnActions(db *innerDB, cfg *core.DBCfg, hashPwdFn func(string) string) {
	for _, model := range models.AllModels {
		if cfg.EraseAllData {
			logs.LogResult("Erasing DB table   "+model.(models.Model).TableName(), db.Unscoped().Delete(model, nil).Error())
		}
		if cfg.MigrateModels {
			logs.LogResult("Migrating DB table "+model.(models.Model).TableName(), db.AutoMigrate(model))
		}
	}

	if cfg.InsertAdmin && cfg.InsertAdminPwd != "" {
		db.InsertAdmin(hashPwdFn(cfg.InsertAdminPwd))
		logs.LogResult("Inserting DB admin", nil)
	}

	sqlDB, err := db.DB.DB()
	logs.LogFatalIfErr(err, errs.FailedToGetSQLDB)

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
}

// Deprecated: Use methods directly on DBOperations interface instead
func (d *LegacyDB) GetDB() any { return d.InnerDB }

// Deprecated: Use Close() instead
func (d *LegacyDB) CloseDB() {
	d.Close()
}
