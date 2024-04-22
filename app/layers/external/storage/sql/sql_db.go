package sql

import (
	"database/sql"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Concrete type that implements the core.DBMethods interface.
type gormAdapter struct {
	*gorm.DB
}

// newGormAdapter wraps *gorm.DB and returns a new concrete *gormAdapter that implements DB.
func newGormAdapter(gormDB *gorm.DB) *gormAdapter {
	return &gormAdapter{gormDB}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/* - Concrete Implementation Methods - */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var nga = newGormAdapter // Short version of newGormAdapter. Used for brevity.

func (ga *gormAdapter) GetSQL() *sql.DB {
	sqlDB, err := ga.DB.DB()
	if err != nil {
		zap.S().Errorf(errs.ErrMsgGettingSQL, err)
	}
	return sqlDB
}

func (ga *gormAdapter) Count(value *int64) core.SQLDatabaseAPI { return nga(ga.DB.Count(value)) }

func (ga *gormAdapter) Create(value any) core.SQLDatabaseAPI { return nga(ga.DB.Create(value)) }

func (ga *gormAdapter) Debug() core.SQLDatabaseAPI { return nga(ga.DB.Debug()) }

func (ga *gormAdapter) Delete(val any, where ...any) core.SQLDatabaseAPI {
	return nga(ga.DB.Delete(val, where))
}

func (ga *gormAdapter) Error() error { return ga.DB.Error }

func (ga *gormAdapter) Find(out any, where ...any) core.SQLDatabaseAPI {
	return nga(ga.DB.Find(out, where...))
}

func (ga *gormAdapter) First(out any, where ...any) core.SQLDatabaseAPI {
	return nga(ga.DB.First(out, where...))
}

func (ga *gormAdapter) FirstOrCreate(out any, where ...any) core.SQLDatabaseAPI {
	return nga(ga.DB.FirstOrCreate(out, where...))
}

func (ga *gormAdapter) Group(query string) core.SQLDatabaseAPI { return nga(ga.DB.Group(query)) }

func (ga *gormAdapter) Joins(qry string, args ...any) core.SQLDatabaseAPI {
	return nga(ga.DB.Joins(qry, args))
}

func (ga *gormAdapter) Limit(value int) core.SQLDatabaseAPI { return nga(ga.DB.Limit(value)) }

func (ga *gormAdapter) Model(value any) core.SQLDatabaseAPI { return nga(ga.DB.Model(value)) }

func (ga *gormAdapter) Offset(value int) core.SQLDatabaseAPI { return nga(ga.DB.Offset(value)) }

func (ga *gormAdapter) Order(value string) core.SQLDatabaseAPI { return nga(ga.DB.Order(value)) }

func (ga *gormAdapter) Or(query any, args ...any) core.SQLDatabaseAPI {
	return nga(ga.DB.Or(query, args...))
}

func (ga *gormAdapter) Pluck(col string, val any) core.SQLDatabaseAPI {
	return nga(ga.DB.Pluck(col, val))
}

func (ga *gormAdapter) Raw(sql string, vals ...any) core.SQLDatabaseAPI {
	return nga(ga.DB.Raw(sql, vals...))
}

func (ga *gormAdapter) Rows() (*sql.Rows, error) { return ga.DB.Rows() }

func (ga *gormAdapter) RowsAffected() int64 { return ga.DB.RowsAffected }

func (ga *gormAdapter) Row() *sql.Row { return ga.DB.Row() }

func (ga *gormAdapter) Save(value any) core.SQLDatabaseAPI { return nga(ga.DB.Save(value)) }

func (ga *gormAdapter) Scan(to any) core.SQLDatabaseAPI { return nga(ga.DB.Scan(to)) }

func (ga *gormAdapter) Where(qry any, args ...any) core.SQLDatabaseAPI {
	return nga(ga.DB.Where(qry, args...))
}

func (ga *gormAdapter) Scopes(f ...func(core.SQLDatabaseAPI) core.SQLDatabaseAPI) core.SQLDatabaseAPI {
	adaptedFns := make([]func(*gorm.DB) *gorm.DB, len(f))
	for i, fn := range f {
		adaptedFns[i] = func(db *gorm.DB) *gorm.DB {
			return fn(&gormAdapter{db}).(*gormAdapter).DB // T0D0 this is horrible
		}
	}

	return nga(ga.DB.Scopes(adaptedFns...))
}
