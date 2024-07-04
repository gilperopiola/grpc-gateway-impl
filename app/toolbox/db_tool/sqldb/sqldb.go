package sqldb

import (
	"database/sql"
	"fmt"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _ core.SQLDB = (*sqlDB)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - SQL Database -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// The SQL DB Tool holds a SQL Database object/connection.
//
// -> DB Tool = High Level Operations (e.g. CreateUser, GetUser, GetUsers)
// -> DB = Low Level Operations (e.g. Insert, Find, Count)

type sqlDB struct {
	*gorm.DB
}

// Returns a new connection to a SQL Database. It uses Gorm.
func NewSQLDB(cfg *core.DBCfg, retrier core.Retrier) core.SQLDB {
	gormCfg := newGormCfg(cfg.LogLevel)

	connectToDB := func() (any, error) {
		connectionString := cfg.GetSQLConnString()
		gormDB, err := gorm.Open(mysql.Open(connectionString), gormCfg)
		return &sqlDB{gormDB}, err
	}

	onConnectionFailureDo := func() {
		if db, err := sql.Open("mysql", cfg.GetSQLConnStringNoSchema()); err == nil {
			defer db.Close()

			_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", cfg.Schema))
			core.LogResult(cfg.Schema+" DB creation", err)
		}
	}

	sqlDBAny, err := retrier.TryToConnectToDB(connectToDB, onConnectionFailureDo)
	core.LogFatalIfErr(err)

	gormDB := sqlDBAny.(*sqlDB)

	if cfg.EraseAllData {
		gormDB.Unscoped().Delete(models.AllDBModels, nil)
	}

	if cfg.MigrateModels {
		gormDB.AutoMigrate(models.AllDBModels...)
	}

	if cfg.InsertAdmin && cfg.InsertAdminPwd != "" {
		gormDB.InsertAdmin(cfg.InsertAdminPwd)
	}

	return gormDB
}

func newGormCfg(logLevel int) *gorm.Config {
	return &gorm.Config{
		Logger:         newSQLDBLogger(zap.L(), logLevel),
		TranslateError: true,
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - SQL DB Methods -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (sdb *sqlDB) GetInnerDB() any { return sdb.DB }

func (sdb *sqlDB) Association(column string) *gorm.Association { return sdb.DB.Association(column) }

func (sdb *sqlDB) Count(value *int64) core.SQLDB { return &sqlDB{sdb.DB.Count(value)} }

func (sdb *sqlDB) Create(value any) core.SQLDB { return &sqlDB{sdb.DB.Create(value)} }

func (sdb *sqlDB) Debug() core.SQLDB { return &sqlDB{sdb.DB.Debug()} }

func (sdb *sqlDB) Error() error { return sdb.DB.Error }

func (sdb *sqlDB) Group(query string) core.SQLDB { return &sqlDB{sdb.DB.Group(query)} }

func (sdb *sqlDB) Limit(value int) core.SQLDB { return &sqlDB{sdb.DB.Limit(value)} }

func (sdb *sqlDB) Model(value any) core.SQLDB { return &sqlDB{sdb.DB.Model(value)} }

func (sdb *sqlDB) Offset(value int) core.SQLDB { return &sqlDB{sdb.DB.Offset(value)} }

func (sdb *sqlDB) Order(value string) core.SQLDB { return &sqlDB{sdb.DB.Order(value)} }

func (sdb *sqlDB) RowsAffected() int64 { return sdb.DB.RowsAffected }

func (sdb *sqlDB) Save(value any) core.SQLDB { return &sqlDB{sdb.DB.Save(value)} }

func (sdb *sqlDB) Scan(to any) core.SQLDB { return &sqlDB{sdb.DB.Scan(to)} }

func (sdb *sqlDB) Close() {
	innerSQLDB, err := sdb.DB.DB()
	core.LogIfErr(err, errs.FailedToGetSQLDB)

	err = innerSQLDB.Close()
	core.LogIfErr(err, errs.FailedToCloseSQLDB)
}

func (sdb *sqlDB) Delete(val any, where ...any) core.SQLDB {
	return &sqlDB{sdb.DB.Delete(val, where)}
}

func (sdb *sqlDB) Find(out any, where ...any) core.SQLDB {
	return &sqlDB{sdb.DB.Find(out, where...)}
}

func (sdb *sqlDB) First(out any, where ...any) core.SQLDB {
	return &sqlDB{sdb.DB.First(out, where...)}
}

func (sdb *sqlDB) FirstOrCreate(out any, where ...any) core.SQLDB {
	return &sqlDB{sdb.DB.FirstOrCreate(out, where...)}
}

func (sdb *sqlDB) InsertAdmin(hashedPwd string) {
	admin := models.User{Username: "admin", Password: hashedPwd, Role: models.AdminRole}
	err := sdb.DB.FirstOrCreate(&admin).Error
	core.WarnIfErr(err, errs.FailedToInsertDBAdmin)
}

func (sdb *sqlDB) Joins(qry string, args ...any) core.SQLDB {
	return &sqlDB{sdb.DB.Joins(qry, args)}
}

func (sdb *sqlDB) Or(query any, args ...any) core.SQLDB {
	return &sqlDB{sdb.DB.Or(query, args...)}
}

func (sdb *sqlDB) Pluck(col string, val any) core.SQLDB {
	return &sqlDB{sdb.DB.Pluck(col, val)}
}

func (sdb *sqlDB) Raw(sql string, vals ...any) core.SQLDB {
	return &sqlDB{sdb.DB.Raw(sql, vals...)}
}

func (sdb *sqlDB) Scopes(fns ...func(core.SQLDB) core.SQLDB) core.SQLDB {
	adaptedFns := make([]func(*gorm.DB) *gorm.DB, len(fns))
	for i, fn := range fns {
		adaptedFns[i] = func(db *gorm.DB) *gorm.DB {
			return fn(&sqlDB{db}).(*sqlDB).DB // Messy. T0D0.
		}
	}
	return &sqlDB{sdb.DB.Scopes(adaptedFns...)}
}

func (sdb *sqlDB) WithContext(ctx god.Ctx) core.SQLDB {
	// Calling the actual gorm WithContext func makes our SQLOptions fail to apply for some reason. T0D0.
	return &sqlDB{sdb.DB}
}

func (sdb *sqlDB) Where(qry any, args ...any) core.SQLDB {
	return &sqlDB{sdb.DB.Where(qry, args...)}
}
