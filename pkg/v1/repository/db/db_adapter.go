package db

import (
	"database/sql"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GormAdapter is our adapter interface for Gorm. We have a concrete gormAdapter that implements it.
// With this we can mock *gorm.DB in our tests.
type GormAdapter interface {
	AddError(err error) error
	AutoMigrate(dst ...interface{}) error
	Count(value *int64) GormAdapter
	Create(value interface{}) GormAdapter
	Debug() GormAdapter
	Delete(value interface{}, where ...interface{}) GormAdapter
	Error() error
	Find(out interface{}, where ...interface{}) GormAdapter
	First(out interface{}, where ...interface{}) GormAdapter
	FirstOrCreate(out interface{}, where ...interface{}) GormAdapter
	GetSQL() *sql.DB
	Group(query string) GormAdapter
	Joins(query string, args ...interface{}) GormAdapter
	Limit(value int) GormAdapter
	Model(value interface{}) GormAdapter
	Offset(value int) GormAdapter
	Order(value string) GormAdapter
	Or(query interface{}, args ...interface{}) GormAdapter
	Pluck(column string, value interface{}) GormAdapter
	Raw(sql string, values ...interface{}) GormAdapter
	Rows() (*sql.Rows, error)
	RowsAffected() int64
	Row() *sql.Row
	Save(value interface{}) GormAdapter
	Scan(dest interface{}) GormAdapter
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) GormAdapter
	Where(query interface{}, args ...interface{}) GormAdapter
}

// gormAdapter is our concrete type that implements the GormAdapter interface.
type gormAdapter struct {
	*gorm.DB
}

// newGormAdapter wraps *gorm.DB and returns a new concrete *gormAdapter as a GormAdapter interface.
func newGormAdapter(gormDB *gorm.DB) GormAdapter {
	return &gormAdapter{gormDB}
}

// n is a short version of newGormAdapter. Used for brevity.
var n = newGormAdapter

/* ----------------------------------- */
/*         - Adapter Methods -         */
/* ----------------------------------- */

func (ga *gormAdapter) AddError(err error) error { return ga.DB.AddError(err) }

func (ga *gormAdapter) AutoMigrate(dst ...any) error { return ga.DB.AutoMigrate(dst...) }

func (ga *gormAdapter) Count(value *int64) GormAdapter { return n(ga.DB.Count(value)) }

func (ga *gormAdapter) Create(value any) GormAdapter { return n(ga.DB.Create(value)) }

func (ga *gormAdapter) Debug() GormAdapter { return n(ga.DB.Debug()) }

func (ga *gormAdapter) Delete(val any, where ...any) GormAdapter { return n(ga.DB.Delete(val, where)) }

func (ga *gormAdapter) Error() error { return ga.DB.Error }

func (ga *gormAdapter) Find(out any, where ...any) GormAdapter { return n(ga.DB.Find(out, where...)) }

func (ga *gormAdapter) First(out any, where ...any) GormAdapter { return n(ga.DB.First(out, where...)) }

func (ga *gormAdapter) FirstOrCreate(out any, where ...any) GormAdapter {
	return n(ga.DB.FirstOrCreate(out, where...))
}

func (ga *gormAdapter) GetSQL() *sql.DB {
	sqlDB, err := ga.DB.DB()
	if err != nil {
		zap.S().Errorf(errs.ErrMsgGettingSqlDB, err)
	}
	return sqlDB
}

func (ga *gormAdapter) Group(query string) GormAdapter { return n(ga.DB.Group(query)) }

func (ga *gormAdapter) Joins(qry string, args ...any) GormAdapter { return n(ga.DB.Joins(qry, args)) }

func (ga *gormAdapter) Limit(value int) GormAdapter { return n(ga.DB.Limit(value)) }

func (ga *gormAdapter) Model(value any) GormAdapter { return n(ga.DB.Model(value)) }

func (ga *gormAdapter) Offset(value int) GormAdapter { return n(ga.DB.Offset(value)) }

func (ga *gormAdapter) Order(value string) GormAdapter { return n(ga.DB.Order(value)) }

func (ga *gormAdapter) Or(query any, args ...any) GormAdapter { return n(ga.DB.Or(query, args...)) }

func (ga *gormAdapter) Pluck(col string, val any) GormAdapter { return n(ga.DB.Pluck(col, val)) }

func (ga *gormAdapter) Raw(sql string, vals ...any) GormAdapter { return n(ga.DB.Raw(sql, vals...)) }

func (ga *gormAdapter) Rows() (*sql.Rows, error) { return ga.DB.Rows() }

func (ga *gormAdapter) RowsAffected() int64 { return ga.DB.RowsAffected }

func (ga *gormAdapter) Row() *sql.Row { return ga.DB.Row() }

func (ga *gormAdapter) Save(value any) GormAdapter { return n(ga.DB.Save(value)) }

func (ga *gormAdapter) Scan(to any) GormAdapter { return n(ga.DB.Scan(to)) }

func (ga *gormAdapter) Scopes(f ...func(*gorm.DB) *gorm.DB) GormAdapter { return n(ga.DB.Scopes(f...)) }

func (ga *gormAdapter) Where(qry any, args ...any) GormAdapter { return n(ga.DB.Where(qry, args...)) }
