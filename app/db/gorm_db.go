package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/utils"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - GORM DB Implementation -     */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// GormDB is a wrapper around gorm.DB implementing the core.DBOperations interface
type GormDB struct {
	db *gorm.DB
}

// Verify that GormDB implements the core.DBOperations interface
var _ core.DBOperations = (*GormDB)(nil)

// NewGormDB creates a new GormDB instance
func NewGormDB(cfg *core.DBCfg, hashPwdFn func(string) string) (*GormDB, error) {
	// Configure the logger
	zapLogger := zap.L() // Use default logger as fallback
	if logs.GetZapLogger() != nil {
		zapLogger = logs.GetZapLogger()
	}

	dbLogger := newDBLogger(zapLogger, cfg.LogLevel)
	gormCfg := &gorm.Config{
		Logger:         dbLogger,
		TranslateError: true,
	}

	// Wrap connection function to match the signature of utils.RetryFunc
	connectToDB := func() (any, error) {
		gormDB, err := gorm.Open(mysql.Open(cfg.GetSQLConnString()), gormCfg)
		return gormDB, err
	}

	// Wrap database creation to match the signature of utils.RetryFuncNoErr
	createDB := func() {
		db, err := sql.Open("mysql", cfg.GetSQLConnStringNoDB())
		if err != nil {
			logs.LogResult("Error trying to connect to SQL instance ", err)
			return
		}
		defer db.Close()

		createStmt := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", cfg.Database)
		if _, err := db.Exec(createStmt); err != nil {
			logs.LogResult("Error trying to create SQL DB", err)
		}
	}

	// Try to connect to the database, with retries
	dbConn, err := utils.TryAndRetry(connectToDB, cfg.Retries, false, createDB)
	if err != nil {
		return nil, &errs.DBErr{Err: err, Context: errs.FailedDBConn}
	}
	logs.LogSimple("✅ DB Connection for " + cfg.Database + " established OK")

	gormDB := dbConn.(*gorm.DB)
	gormDBInstance := &GormDB{db: gormDB}

	// Perform post-connection setup
	if err := setupDBPostConnection(gormDBInstance, cfg, hashPwdFn); err != nil {
		return nil, err
	}

	return gormDBInstance, nil
}

// setupDBPostConnection handles post-connection setup like migrations and admin creation
func setupDBPostConnection(db *GormDB, cfg *core.DBCfg, hashPwdFn func(string) string) error {
	for _, model := range models.AllModels {
		if cfg.EraseAllData {
			tableName := ""
			if modelWithTable, ok := model.(models.Model); ok {
				tableName = modelWithTable.TableName()
			}
			logs.LogResult("Erasing DB table "+tableName, db.db.Unscoped().Delete(model).Error)
		}

		if cfg.MigrateModels {
			modelName := ""
			if modelWithTable, ok := model.(models.Model); ok {
				modelName = modelWithTable.TableName()
			}
			if err := db.db.AutoMigrate(model); err != nil {
				logs.LogResult("Error migrating model "+modelName, err)
				return err
			}
			logs.LogSimple("✅ DB Table " + modelName + " migrated OK")
		}
	}

	if cfg.InsertAdmin && cfg.InsertAdminPwd != "" {
		admin := models.User{
			Username: "admin",
			Password: hashPwdFn(cfg.InsertAdminPwd),
			Role:     models.AdminRole,
		}

		if err := db.db.FirstOrCreate(&admin).Error; err != nil {
			logs.LogResult("Error inserting admin user", err)
			return err
		}
		logs.LogResult("Inserting DB admin", nil)
	}

	sqlDB, err := db.db.DB()
	if err != nil {
		return &errs.DBErr{Err: err, Context: errs.FailedToGetSQLDB}
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}

// Implementation of core.DBOperations interface methods

func (g *GormDB) Find(out any, where ...any) error {
	return g.db.Find(out, where...).Error
}

func (g *GormDB) First(out any, where ...any) error {
	return g.db.First(out, where...).Error
}

func (g *GormDB) Create(value any) error {
	return g.db.Create(value).Error
}

func (g *GormDB) Save(value any) error {
	return g.db.Save(value).Error
}

func (g *GormDB) Delete(value any, where ...any) error {
	return g.db.Delete(value, where...).Error
}

func (g *GormDB) WithContext(ctx context.Context) core.DBOperations {
	return &GormDB{db: g.db.WithContext(ctx)}
}

func (g *GormDB) Transaction(fn func(tx core.DBOperations) error) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		return fn(&GormDB{db: tx})
	})
}

func (g *GormDB) Close() error {
	sqlDB, err := g.db.DB()
	if err != nil {
		return &errs.DBErr{Err: err, Context: errs.FailedToGetSQLDB}
	}
	return sqlDB.Close()
}

// Additional helpers for repositories to use

func (g *GormDB) Model(value any) *gorm.DB {
	return g.db.Model(value)
}

func (g *GormDB) Where(query any, args ...any) *gorm.DB {
	return g.db.Where(query, args...)
}

func (g *GormDB) Preload(query string, args ...any) *gorm.DB {
	return g.db.Preload(query, args...)
}

func (g *GormDB) Association(column string) *gorm.Association {
	return g.db.Association(column)
}

// GetDB returns the underlying *gorm.DB instance
// This should only be used in special cases where direct access is needed
func (g *GormDB) GetDB() *gorm.DB {
	return g.db
}
