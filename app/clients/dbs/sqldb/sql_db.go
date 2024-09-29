package sqldb

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

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*      - High Level SQL DB Conn -     */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type DB struct {
	DB core.BaseSQLDB
}

func NewSQLDBConnection(cfg *core.DBCfg) core.DB {

	// We use our Logger wrapped inside of a gorm adapter
	gormCfg := &gorm.Config{
		Logger:         newDBLogger(zap.L(), cfg.LogLevel),
		TranslateError: true,
	}

	var connectToDB = func() (any, error) {
		gormDB, err := gorm.Open(mysql.Open(cfg.GetSQLConnectionString()), gormCfg)
		return &baseSQLDB{gormDB}, err
	}

	var createDB = func() {
		if db, err := sql.Open("mysql", cfg.GetSQLConnectionStringNoDB()); err == nil {
			defer db.Close()
			db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", cfg.Database))
		}
	}

	// We try to connect to the DB directly.
	// If it fails, we try to connect without DB, then creating it.
	// Then we retry.
	dbConn, err := utils.Retry(connectToDB, utils.BasicRetryCfg(cfg.Retries, createDB))
	logs.LogFatalIfErr(err, errs.FailedDBConn)

	db := &DB{dbConn.(*baseSQLDB)}

	if cfg.EraseAllData {
		db.DB.Unscoped().Delete(models.AllModels, nil)
	}

	if cfg.MigrateModels {
		db.DB.AutoMigrate(models.AllModels...)
	}

	if cfg.InsertAdmin && cfg.InsertAdminPwd != "" {
		db.DB.InsertAdmin(cfg.InsertAdminPwd)
	}

	return db
}

func (this DB) GetDB() any { return this.DB }
func (this DB) CloseDB()   { this.DB.Close() }
