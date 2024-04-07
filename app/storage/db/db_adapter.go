package db

import (
	"database/sql"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Database DBAdapter // Alias for DBAdapter.

// DBAdapter is our adapter interface for Gorm. We have a concrete gormAdapter that implements it.
// With this we can mock *gorm.DB in our tests.
type DBAdapter interface {
	AddError(err error) error
	AutoMigrate(dst ...interface{}) error
	Count(value *int64) DBAdapter
	Create(value interface{}) DBAdapter
	Debug() DBAdapter
	Delete(value interface{}, where ...interface{}) DBAdapter
	Error() error
	Find(out interface{}, where ...interface{}) DBAdapter
	First(out interface{}, where ...interface{}) DBAdapter
	FirstOrCreate(out interface{}, where ...interface{}) DBAdapter
	GetSQL() *sql.DB
	Group(query string) DBAdapter
	Joins(query string, args ...interface{}) DBAdapter
	Limit(value int) DBAdapter
	Model(value interface{}) DBAdapter
	Offset(value int) DBAdapter
	Order(value string) DBAdapter
	Or(query interface{}, args ...interface{}) DBAdapter
	Pluck(column string, value interface{}) DBAdapter
	Raw(sql string, values ...interface{}) DBAdapter
	Rows() (*sql.Rows, error)
	RowsAffected() int64
	Row() *sql.Row
	Save(value interface{}) DBAdapter
	Scan(dest interface{}) DBAdapter
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) DBAdapter
	Where(query interface{}, args ...interface{}) DBAdapter
}

// gormAdapter is our concrete type that implements the DBAdapter interface.
type gormAdapter struct {
	*gorm.DB
}

// newGormAdapter wraps *gorm.DB and returns a new concrete *gormAdapter as a DBAdapter interface.
func newGormAdapter(gormDB *gorm.DB) *gormAdapter {
	return &gormAdapter{gormDB}
}

// n is a short version of newGormAdapter. Used for brevity.
var n = newGormAdapter

/* ----------------------------------- */
/*         - Adapter Methods -         */
/* ----------------------------------- */

func (ga *gormAdapter) AddError(err error) error { return ga.DB.AddError(err) }

func (ga *gormAdapter) AutoMigrate(dst ...any) error { return ga.DB.AutoMigrate(dst...) }

func (ga *gormAdapter) Count(value *int64) DBAdapter { return n(ga.DB.Count(value)) }

func (ga *gormAdapter) Create(value any) DBAdapter { return n(ga.DB.Create(value)) }

func (ga *gormAdapter) Debug() DBAdapter { return n(ga.DB.Debug()) }

func (ga *gormAdapter) Delete(val any, where ...any) DBAdapter { return n(ga.DB.Delete(val, where)) }

func (ga *gormAdapter) Error() error { return ga.DB.Error }

func (ga *gormAdapter) Find(out any, where ...any) DBAdapter { return n(ga.DB.Find(out, where...)) }

func (ga *gormAdapter) First(out any, where ...any) DBAdapter { return n(ga.DB.First(out, where...)) }

func (ga *gormAdapter) FirstOrCreate(out any, where ...any) DBAdapter {
	return n(ga.DB.FirstOrCreate(out, where...))
}

func (ga *gormAdapter) GetSQL() *sql.DB {
	sqlDB, err := ga.DB.DB()
	if err != nil {
		zap.S().Errorf(errs.ErrMsgGettingSqlDB, err)
	}
	return sqlDB
}

func (ga *gormAdapter) Group(query string) DBAdapter { return n(ga.DB.Group(query)) }

func (ga *gormAdapter) Joins(qry string, args ...any) DBAdapter { return n(ga.DB.Joins(qry, args)) }

func (ga *gormAdapter) Limit(value int) DBAdapter { return n(ga.DB.Limit(value)) }

func (ga *gormAdapter) Model(value any) DBAdapter { return n(ga.DB.Model(value)) }

func (ga *gormAdapter) Offset(value int) DBAdapter { return n(ga.DB.Offset(value)) }

func (ga *gormAdapter) Order(value string) DBAdapter { return n(ga.DB.Order(value)) }

func (ga *gormAdapter) Or(query any, args ...any) DBAdapter { return n(ga.DB.Or(query, args...)) }

func (ga *gormAdapter) Pluck(col string, val any) DBAdapter { return n(ga.DB.Pluck(col, val)) }

func (ga *gormAdapter) Raw(sql string, vals ...any) DBAdapter { return n(ga.DB.Raw(sql, vals...)) }

func (ga *gormAdapter) Rows() (*sql.Rows, error) { return ga.DB.Rows() }

func (ga *gormAdapter) RowsAffected() int64 { return ga.DB.RowsAffected }

func (ga *gormAdapter) Row() *sql.Row { return ga.DB.Row() }

func (ga *gormAdapter) Save(value any) DBAdapter { return n(ga.DB.Save(value)) }

func (ga *gormAdapter) Scan(to any) DBAdapter { return n(ga.DB.Scan(to)) }

func (ga *gormAdapter) Scopes(f ...func(*gorm.DB) *gorm.DB) DBAdapter { return n(ga.DB.Scopes(f...)) }

func (ga *gormAdapter) Where(qry any, args ...any) DBAdapter { return n(ga.DB.Where(qry, args...)) }
