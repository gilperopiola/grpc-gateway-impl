package db

import (
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

var _ core.DB = &DB{}

type DB struct {
	core.InnerDB
}

func NewSQLDBConn(cfg *core.DBCfg, hashPwdFn func(string) string) core.DB {

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
	retryCfg := utils.RetryCfg{Times: cfg.Retries, OnFailure: createDB}
	dbConn, err := utils.RetryFunc(connectToDB, retryCfg)
	logs.LogFatalIfErr(err, errs.FailedDBConn)
	logs.LogResult("DB Connection for "+cfg.Database, nil)

	innerDB := dbConn.(*innerDB)
	postDBConnActions(innerDB, cfg, hashPwdFn)
	return &DB{innerDB}
}

func postDBConnActions(db *innerDB, cfg *core.DBCfg, hashPwdFn func(string) string) {
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

func (this *DB) GetDB() any { return this.InnerDB }
func (this *DB) CloseDB()   { this.InnerDB.Close() }
