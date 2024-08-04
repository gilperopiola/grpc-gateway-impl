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

var _ core.SqlDB = (*sqlDB)(nil)

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
func NewSqlDB(cfg *core.DBCfg) core.SqlDB {

	// We use our Logger but wrapped inside of an Adapter for gorm.
	// Translate mySQL errors to gorm.
	gormCfg := &gorm.Config{
		Logger:         newSqlDBLogger(zap.L(), cfg.LogLevel),
		TranslateError: true,
	}

	// We try to connect to the DB directly.
	// If it fails, we try to connect without a schema and then create the DB.
	// We wait a bit, then we try to connect directly again.
	//
	// This process is retried a few times (with exponential backoff) until
	// the connection succeeds or we run out of retries.
	connectToDB := connectToDBAdapter(cfg.GetSQLConnString(), gormCfg)
	createDB := createDBAdapter(cfg)

	result, err := core.FallbackAndRetry(connectToDB, createDB, cfg.Retries)
	core.LogFatalIfErr(err, errs.FailedToConnectToDB)

	sqlDB := result.(*sqlDB)

	if cfg.EraseAllData {
		sqlDB.Unscoped().Delete(models.AllDBModels, nil)
	}

	if cfg.MigrateModels {
		sqlDB.AutoMigrate(models.AllDBModels...)
	}

	if cfg.InsertAdmin && cfg.InsertAdminPwd != "" {
		sqlDB.InsertAdmin(cfg.InsertAdminPwd)
	}

	return sqlDB
}

// Adapts our connectToDB func so it can be used with the Retrier.
var connectToDBAdapter = func(connString string, gormCfg *gorm.Config) func() (any, error) {
	return func() (any, error) {
		gormDB, err := gorm.Open(mysql.Open(connString), gormCfg)
		return &sqlDB{gormDB}, err
	}
}

// Adapts our createDB func so it can be used with the Retrier.
var createDBAdapter = func(cfg *core.DBCfg) func() {
	return func() {
		if db, err := sql.Open("mysql", cfg.GetSQLConnStringNoSchema()); err == nil {
			defer db.Close()
			db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", cfg.Schema))
		}
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - SQL DB Methods -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (sdb *sqlDB) GetInnerDB() any { return sdb.DB }

func (sdb *sqlDB) Association(column string) core.SqlDBAssoc { return sdb.DB.Association(column) }

func (sdb *sqlDB) Count(value *int64) core.SqlDB { return &sqlDB{sdb.DB.Count(value)} }

func (sdb *sqlDB) Create(value any) core.SqlDB { return &sqlDB{sdb.DB.Create(value)} }

func (sdb *sqlDB) Debug() core.SqlDB { return &sqlDB{sdb.DB.Debug()} }

func (sdb *sqlDB) Error() error { return sdb.DB.Error }

func (sdb *sqlDB) Group(query string) core.SqlDB { return &sqlDB{sdb.DB.Group(query)} }

func (sdb *sqlDB) Limit(value int) core.SqlDB { return &sqlDB{sdb.DB.Limit(value)} }

func (sdb *sqlDB) Model(value any) core.SqlDB { return &sqlDB{sdb.DB.Model(value)} }

func (sdb *sqlDB) Offset(value int) core.SqlDB { return &sqlDB{sdb.DB.Offset(value)} }

func (sdb *sqlDB) Order(value string) core.SqlDB { return &sqlDB{sdb.DB.Order(value)} }

func (sdb *sqlDB) RowsAffected() int64 { return sdb.DB.RowsAffected }

func (sdb *sqlDB) Save(value any) core.SqlDB { return &sqlDB{sdb.DB.Save(value)} }

func (sdb *sqlDB) Scan(to any) core.SqlDB { return &sqlDB{sdb.DB.Scan(to)} }

func (sdb *sqlDB) Close() {
	innerSQLDB, err := sdb.DB.DB()
	core.LogIfErr(err, errs.FailedToGetSQLDB)

	err = innerSQLDB.Close()
	core.LogIfErr(err, errs.FailedToCloseSQLDB)
}

func (sdb *sqlDB) Delete(val any, where ...any) core.SqlDB {
	return &sqlDB{sdb.DB.Delete(val, where)}
}

func (sdb *sqlDB) Find(out any, where ...any) core.SqlDB {
	return &sqlDB{sdb.DB.Find(out, where...)}
}

func (sdb *sqlDB) First(out any, where ...any) core.SqlDB {
	return &sqlDB{sdb.DB.First(out, where...)}
}

func (sdb *sqlDB) FirstOrCreate(out any, where ...any) core.SqlDB {
	return &sqlDB{sdb.DB.FirstOrCreate(out, where...)}
}

func (sdb *sqlDB) InsertAdmin(hashedPwd string) {
	admin := models.User{Username: "admin", Password: hashedPwd, Role: models.AdminRole}
	err := sdb.DB.FirstOrCreate(&admin).Error
	core.WarnIfErr(err, errs.FailedToInsertDBAdmin)
}

func (sdb *sqlDB) Joins(qry string, args ...any) core.SqlDB {
	return &sqlDB{sdb.DB.Joins(qry, args)}
}

func (sdb *sqlDB) Or(query any, args ...any) core.SqlDB {
	return &sqlDB{sdb.DB.Or(query, args...)}
}

func (sdb *sqlDB) Pluck(col string, val any) core.SqlDB {
	return &sqlDB{sdb.DB.Pluck(col, val)}
}

func (sdb *sqlDB) Raw(sql string, vals ...any) core.SqlDB {
	return &sqlDB{sdb.DB.Raw(sql, vals...)}
}

func (sdb *sqlDB) Row() core.SqlRow { return sdb.DB.Row() }

func (sdb *sqlDB) Rows() (core.SqlRows, error) { return sdb.DB.Rows() }

func (sdb *sqlDB) Scopes(fns ...func(core.SqlDB) core.SqlDB) core.SqlDB {
	adaptedFns := make([]func(*gorm.DB) *gorm.DB, len(fns))
	for i, fn := range fns {
		adaptedFns[i] = func(db *gorm.DB) *gorm.DB {
			return fn(&sqlDB{db}).(*sqlDB).DB // Messy. T0D0.
		}
	}
	return &sqlDB{sdb.DB.Scopes(adaptedFns...)}
}

func (sdb *sqlDB) WithContext(ctx god.Ctx) core.SqlDB {
	// Calling the actual gorm WithContext func makes our SQLOptions fail to apply for some reason. T0D0.
	return &sqlDB{sdb.DB}
}

func (sdb *sqlDB) Where(qry any, args ...any) core.SqlDB {
	return &sqlDB{sdb.DB.Where(qry, args...)}
}
