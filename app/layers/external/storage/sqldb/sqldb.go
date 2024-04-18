package sqldb

import (
	"database/sql"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/special_types"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Database special_types.SQLDB // Alias for special_types.SQLDB.

// gormAdapter is our concrete type that implements the special_types.SQLDB interface.
type gormAdapter struct {
	*gorm.DB
}

var _ Database = &gormAdapter{}

// newGormAdapter wraps *gorm.DB and returns a new concrete *gormAdapter that implements special_types.SQLDB interface.
func newGormAdapter(gormDB *gorm.DB) *gormAdapter {
	return &gormAdapter{gormDB}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - Adapter Methods -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Short version of newGormAdapter. Used for brevity.
var n = newGormAdapter

func (ga *gormAdapter) Count(value *int64) special_types.SQLDB { return n(ga.DB.Count(value)) }

func (ga *gormAdapter) Create(value any) special_types.SQLDB { return n(ga.DB.Create(value)) }

func (ga *gormAdapter) Debug() special_types.SQLDB { return n(ga.DB.Debug()) }

func (ga *gormAdapter) Delete(val any, where ...any) special_types.SQLDB {
	return n(ga.DB.Delete(val, where))
}

func (ga *gormAdapter) Error() error { return ga.DB.Error }

func (ga *gormAdapter) Find(out any, where ...any) special_types.SQLDB {
	return n(ga.DB.Find(out, where...))
}

func (ga *gormAdapter) First(out any, where ...any) special_types.SQLDB {
	return n(ga.DB.First(out, where...))
}

func (ga *gormAdapter) FirstOrCreate(out any, where ...any) special_types.SQLDB {
	return n(ga.DB.FirstOrCreate(out, where...))
}

func (ga *gormAdapter) GetSQL() *sql.DB {
	sqlDB, err := ga.DB.DB()
	if err != nil {
		zap.S().Errorf(errs.ErrMsgGettingSqlDB, err)
	}
	return sqlDB
}

func (ga *gormAdapter) Group(query string) special_types.SQLDB { return n(ga.DB.Group(query)) }

func (ga *gormAdapter) Joins(qry string, args ...any) special_types.SQLDB {
	return n(ga.DB.Joins(qry, args))
}

func (ga *gormAdapter) Limit(value int) special_types.SQLDB { return n(ga.DB.Limit(value)) }

func (ga *gormAdapter) Model(value any) special_types.SQLDB { return n(ga.DB.Model(value)) }

func (ga *gormAdapter) Offset(value int) special_types.SQLDB { return n(ga.DB.Offset(value)) }

func (ga *gormAdapter) Order(value string) special_types.SQLDB { return n(ga.DB.Order(value)) }

func (ga *gormAdapter) Or(query any, args ...any) special_types.SQLDB {
	return n(ga.DB.Or(query, args...))
}

func (ga *gormAdapter) Pluck(col string, val any) special_types.SQLDB {
	return n(ga.DB.Pluck(col, val))
}

func (ga *gormAdapter) Raw(sql string, vals ...any) special_types.SQLDB {
	return n(ga.DB.Raw(sql, vals...))
}

func (ga *gormAdapter) Rows() (*sql.Rows, error) { return ga.DB.Rows() }

func (ga *gormAdapter) RowsAffected() int64 { return ga.DB.RowsAffected }

func (ga *gormAdapter) Row() *sql.Row { return ga.DB.Row() }

func (ga *gormAdapter) Save(value any) special_types.SQLDB { return n(ga.DB.Save(value)) }

func (ga *gormAdapter) Scan(to any) special_types.SQLDB { return n(ga.DB.Scan(to)) }

func (ga *gormAdapter) Scopes(f ...func(special_types.SQLDB) special_types.SQLDB) special_types.SQLDB {
	adaptedFns := make([]func(*gorm.DB) *gorm.DB, len(f))
	for i, fn := range f {
		adaptedFns[i] = func(db *gorm.DB) *gorm.DB {
			return fn(&gormAdapter{db}).(*gormAdapter).DB // T0D0 horrible
		}
	}

	return n(ga.DB.Scopes(adaptedFns...))
}

func (ga *gormAdapter) Where(qry any, args ...any) special_types.SQLDB {
	return n(ga.DB.Where(qry, args...))
}
