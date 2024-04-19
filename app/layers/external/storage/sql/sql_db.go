package sql

import (
	"database/sql"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// This is the interface that we use to interact with an SQL Database.
// It's an adapter for Gorm. Concrete types gormAdapter and mocks.Gorm implement this.
type DB interface {
	GetSQL() *sql.DB
	AddError(err error) error
	AutoMigrate(dst ...interface{}) error
	Count(value *int64) DB
	Create(value interface{}) DB
	Debug() DB
	Delete(value interface{}, where ...interface{}) DB
	Error() error
	Find(out interface{}, where ...interface{}) DB
	First(out interface{}, where ...interface{}) DB
	FirstOrCreate(out interface{}, where ...interface{}) DB
	Group(query string) DB
	Joins(query string, args ...interface{}) DB
	Limit(value int) DB
	Model(value interface{}) DB
	Offset(value int) DB
	Order(value string) DB
	Or(query interface{}, args ...interface{}) DB
	Pluck(column string, value interface{}) DB
	Raw(sql string, values ...interface{}) DB
	Rows() (*sql.Rows, error)
	RowsAffected() int64
	Row() *sql.Row
	Save(value interface{}) DB
	Scan(dest interface{}) DB
	Where(query interface{}, args ...interface{}) DB
	Scopes(funcs ...func(DB) DB) DB
}

// Concrete type that implements the DB interface.
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

func (ga *gormAdapter) Count(value *int64) DB { return nga(ga.DB.Count(value)) }

func (ga *gormAdapter) Create(value any) DB { return nga(ga.DB.Create(value)) }

func (ga *gormAdapter) Debug() DB { return nga(ga.DB.Debug()) }

func (ga *gormAdapter) Delete(val any, where ...any) DB { return nga(ga.DB.Delete(val, where)) }

func (ga *gormAdapter) Error() error { return ga.DB.Error }

func (ga *gormAdapter) Find(out any, where ...any) DB { return nga(ga.DB.Find(out, where...)) }

func (ga *gormAdapter) First(out any, where ...any) DB { return nga(ga.DB.First(out, where...)) }

func (ga *gormAdapter) FirstOrCreate(out any, where ...any) DB {
	return nga(ga.DB.FirstOrCreate(out, where...))
}

func (ga *gormAdapter) Group(query string) DB { return nga(ga.DB.Group(query)) }

func (ga *gormAdapter) Joins(qry string, args ...any) DB { return nga(ga.DB.Joins(qry, args)) }

func (ga *gormAdapter) Limit(value int) DB { return nga(ga.DB.Limit(value)) }

func (ga *gormAdapter) Model(value any) DB { return nga(ga.DB.Model(value)) }

func (ga *gormAdapter) Offset(value int) DB { return nga(ga.DB.Offset(value)) }

func (ga *gormAdapter) Order(value string) DB { return nga(ga.DB.Order(value)) }

func (ga *gormAdapter) Or(query any, args ...any) DB { return nga(ga.DB.Or(query, args...)) }

func (ga *gormAdapter) Pluck(col string, val any) DB { return nga(ga.DB.Pluck(col, val)) }

func (ga *gormAdapter) Raw(sql string, vals ...any) DB { return nga(ga.DB.Raw(sql, vals...)) }

func (ga *gormAdapter) Rows() (*sql.Rows, error) { return ga.DB.Rows() }

func (ga *gormAdapter) RowsAffected() int64 { return ga.DB.RowsAffected }

func (ga *gormAdapter) Row() *sql.Row { return ga.DB.Row() }

func (ga *gormAdapter) Save(value any) DB { return nga(ga.DB.Save(value)) }

func (ga *gormAdapter) Scan(to any) DB { return nga(ga.DB.Scan(to)) }

func (ga *gormAdapter) Where(qry any, args ...any) DB { return nga(ga.DB.Where(qry, args...)) }

func (ga *gormAdapter) Scopes(f ...func(DB) DB) DB {
	adaptedFns := make([]func(*gorm.DB) *gorm.DB, len(f))
	for i, fn := range f {
		adaptedFns[i] = func(db *gorm.DB) *gorm.DB {
			return fn(&gormAdapter{db}).(*gormAdapter).DB // T0D0 horrible
		}
	}

	return nga(ga.DB.Scopes(adaptedFns...))
}
