package sql

import (
	"context"
	"database/sql"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"gorm.io/gorm"
)

var _ core.SQLDatabaseAPI = (*sqlAdapter)(nil)

// Concrete type for interacting with a SQL DB through GORM.
type sqlAdapter struct {
	*gorm.DB
}

func newGormAdapter(gormDB *gorm.DB) *sqlAdapter {
	return &sqlAdapter{gormDB}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/* - Concrete Implementation Methods - */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var nGA = newGormAdapter // Alias.

func (ga *sqlAdapter) Close() {
	sqlDB, err := ga.DB.DB()
	core.LogIfErr(err, errs.FailedToGetSQLDB)

	err = sqlDB.Close()
	core.LogIfErr(err, errs.FailedToCloseSQLDB)
}

func (ga *sqlAdapter) Count(value *int64) core.SQLDatabaseAPI { return nGA(ga.DB.Count(value)) }

func (ga *sqlAdapter) Create(value any) core.SQLDatabaseAPI { return nGA(ga.DB.Create(value)) }

func (ga *sqlAdapter) Debug() core.SQLDatabaseAPI { return nGA(ga.DB.Debug()) }

func (ga *sqlAdapter) Error() error { return ga.DB.Error }

func (ga *sqlAdapter) Group(query string) core.SQLDatabaseAPI { return nGA(ga.DB.Group(query)) }

func (ga *sqlAdapter) Limit(value int) core.SQLDatabaseAPI { return nGA(ga.DB.Limit(value)) }

func (ga *sqlAdapter) Model(value any) core.SQLDatabaseAPI { return nGA(ga.DB.Model(value)) }

func (ga *sqlAdapter) Offset(value int) core.SQLDatabaseAPI { return nGA(ga.DB.Offset(value)) }

func (ga *sqlAdapter) Order(value string) core.SQLDatabaseAPI { return nGA(ga.DB.Order(value)) }

func (ga *sqlAdapter) Rows() (*sql.Rows, error) { return ga.DB.Rows() }

func (ga *sqlAdapter) RowsAffected() int64 { return ga.DB.RowsAffected }

func (ga *sqlAdapter) Row() *sql.Row { return ga.DB.Row() }

func (ga *sqlAdapter) Save(value any) core.SQLDatabaseAPI { return nGA(ga.DB.Save(value)) }

func (ga *sqlAdapter) Scan(to any) core.SQLDatabaseAPI { return nGA(ga.DB.Scan(to)) }

func (ga *sqlAdapter) Delete(val any, where ...any) core.SQLDatabaseAPI {
	return nGA(ga.DB.Delete(val, where))
}

func (ga *sqlAdapter) Find(out any, where ...any) core.SQLDatabaseAPI {
	return nGA(ga.DB.Find(out, where...))
}

func (ga *sqlAdapter) First(out any, where ...any) core.SQLDatabaseAPI {
	return nGA(ga.DB.First(out, where...))
}

func (ga *sqlAdapter) FirstOrCreate(out any, where ...any) core.SQLDatabaseAPI {
	return nGA(ga.DB.FirstOrCreate(out, where...))
}

func (ga *sqlAdapter) Joins(qry string, args ...any) core.SQLDatabaseAPI {
	return nGA(ga.DB.Joins(qry, args))
}

func (ga *sqlAdapter) Or(query any, args ...any) core.SQLDatabaseAPI {
	return nGA(ga.DB.Or(query, args...))
}

func (ga *sqlAdapter) Pluck(col string, val any) core.SQLDatabaseAPI {
	return nGA(ga.DB.Pluck(col, val))
}

func (ga *sqlAdapter) Raw(sql string, vals ...any) core.SQLDatabaseAPI {
	return nGA(ga.DB.Raw(sql, vals...))
}

func (ga *sqlAdapter) Scopes(fns ...func(core.SQLDatabaseAPI) core.SQLDatabaseAPI) core.SQLDatabaseAPI {
	adaptedFns := make([]func(*gorm.DB) *gorm.DB, len(fns))
	for i, fn := range fns {
		adaptedFns[i] = func(db *gorm.DB) *gorm.DB {
			return fn(&sqlAdapter{db}).(*sqlAdapter).DB // Messy. T0D0.
		}
	}
	return nGA(ga.DB.Scopes(adaptedFns...))
}

func (ga *sqlAdapter) WithContext(ctx context.Context) core.SQLDatabaseAPI {
	// Calling the actual gorm WithContext func makes our SQLOptions fail to apply for some reason. T0D0.
	return nGA(ga.DB)
}

func (ga *sqlAdapter) Where(qry any, args ...any) core.SQLDatabaseAPI {
	return nGA(ga.DB.Where(qry, args...))
}
