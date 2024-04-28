package sqldb

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

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
func NewSQLDB(cfg *core.DBCfg) core.SQLDB {
	gormCfg := &gorm.Config{
		Logger:         newSQLDBLogger(zap.L(), cfg.LogLevel),
		TranslateError: true,
	}

	gormDB, err := connectToSQLDBWithGorm(cfg.GetSQLConnString(), gormCfg)
	core.LogPanicIfErr(err)

	if cfg.MigrateModels {
		gormDB.AutoMigrate(core.AllModels...)
	}

	if cfg.InsertAdmin && cfg.InsertAdminPwd != "" {
		gormDB.InsertAdmin(cfg.InsertAdminPwd)
	}

	return gormDB
}

// Calls gorm.Open and wraps the returned *gorm.DB with our concrete type.
func connectToSQLDBWithGorm(dsn string, opts ...gorm.Option) (*sqlDB, error) {
	gormDB, err := gorm.Open(mysql.Open(dsn), opts...)
	return &sqlDB{gormDB}, err
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - SQL DB Methods -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (sdb *sqlDB) GetInnerDB() any { return sdb.DB }

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
	admin := core.User{Username: "admin", Password: hashedPwd, Role: core.AdminRole}
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

func (sdb *sqlDB) WithContext(ctx context.Context) core.SQLDB {
	// Calling the actual gorm WithContext func makes our SQLOptions fail to apply for some reason. T0D0.
	return &sqlDB{sdb.DB}
}

func (sdb *sqlDB) Where(qry any, args ...any) core.SQLDB {
	return &sqlDB{sdb.DB.Where(qry, args...)}
}
