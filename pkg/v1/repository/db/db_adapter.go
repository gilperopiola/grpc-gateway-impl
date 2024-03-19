package db

import (
	"database/sql"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GormAdapter is our adapter interface for GORM. We have a concrete gormAdapter that implements it.
// With this we can mock the *gorm.DB in our tests.
type GormAdapter interface {
	AutoMigrate(dst ...interface{}) error
	GetSQLDB() *sql.DB
	Where(query interface{}, args ...interface{}) GormAdapter
	Or(query interface{}, args ...interface{}) GormAdapter
	Not(query interface{}, args ...interface{}) GormAdapter
	Limit(value int) GormAdapter
	Offset(value int) GormAdapter
	Order(value string) GormAdapter
	Select(query interface{}, args ...interface{}) GormAdapter
	Omit(columns ...string) GormAdapter
	Group(query string) GormAdapter
	Having(query string, values ...interface{}) GormAdapter
	Joins(query string, args ...interface{}) GormAdapter
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) GormAdapter
	Unscoped() GormAdapter
	Attrs(attrs ...interface{}) GormAdapter
	Assign(attrs ...interface{}) GormAdapter
	First(out interface{}, where ...interface{}) GormAdapter
	Last(out interface{}, where ...interface{}) GormAdapter
	Find(out interface{}, where ...interface{}) GormAdapter
	Scan(dest interface{}) GormAdapter
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	ScanRows(rows *sql.Rows, result interface{}) error
	Pluck(column string, value interface{}) GormAdapter
	Count(value *int64) GormAdapter
	FirstOrInit(out interface{}, where ...interface{}) GormAdapter
	FirstOrCreate(out interface{}, where ...interface{}) GormAdapter
	UpdateColumns(values interface{}) GormAdapter
	Save(value interface{}) GormAdapter
	Create(value interface{}) GormAdapter
	Delete(value interface{}, where ...interface{}) GormAdapter
	Raw(sql string, values ...interface{}) GormAdapter
	Exec(sql string, values ...interface{}) GormAdapter
	Model(value interface{}) GormAdapter
	Table(name string) GormAdapter
	Debug() GormAdapter
	Begin() GormAdapter
	Commit() GormAdapter
	Rollback() GormAdapter
	Get(name string) (interface{}, bool)
	AddError(err error) error
	RowsAffected() int64
	Error() error
}

type gormAdapter struct {
	*gorm.DB
}

// openGormAdapter calls gorm.Open and wraps the returned gorm.DB with our concrete type that implements the GormAdapter interface.
// This way we can mock the GormAdapter in our tests.
func openGormAdapter(dialector gorm.Dialector, opts ...gorm.Option) (GormAdapter, error) {
	gormDB, err := gorm.Open(dialector, opts...)
	return newGormAdapter(gormDB), err
}

// newGormAdapter wraps *gorm.DB and returns a new concrete *gormAdapter as a GormAdapter interface.
func newGormAdapter(gormDB *gorm.DB) GormAdapter {
	return &gormAdapter{gormDB}
}

func (ga *gormAdapter) GetSQLDB() *sql.DB {
	sqlDB, err := ga.DB.DB()
	if err != nil {
		zap.S().Errorf(errs.ErrMsgGettingSQLDB, err)
		return nil
	}
	return sqlDB
}

func (ga *gormAdapter) AutoMigrate(dst ...interface{}) error {
	return ga.DB.AutoMigrate(dst...)
}

func (ga *gormAdapter) Where(query interface{}, args ...interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Where(query, args...))
}

func (ga *gormAdapter) Or(query interface{}, args ...interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Or(query, args...))
}

func (ga *gormAdapter) Not(query interface{}, args ...interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Not(query, args...))
}

func (ga *gormAdapter) Limit(value int) GormAdapter {
	return newGormAdapter(ga.DB.Limit(value))
}

func (ga *gormAdapter) Offset(value int) GormAdapter {
	return newGormAdapter(ga.DB.Offset(value))
}

func (ga *gormAdapter) Order(value string) GormAdapter {
	return newGormAdapter(ga.DB.Order(value))
}

func (ga *gormAdapter) Select(query interface{}, args ...interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Select(query, args...))
}

func (ga *gormAdapter) Omit(columns ...string) GormAdapter {
	return newGormAdapter(ga.DB.Omit(columns...))
}

func (ga *gormAdapter) Group(query string) GormAdapter {
	return newGormAdapter(ga.DB.Group(query))
}

func (ga *gormAdapter) Having(query string, values ...interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Having(query, values...))
}

func (ga *gormAdapter) Joins(query string, args ...interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Joins(query, args...))
}

func (ga *gormAdapter) Scopes(funcs ...func(*gorm.DB) *gorm.DB) GormAdapter {
	return newGormAdapter(ga.DB.Scopes(funcs...))
}

func (ga *gormAdapter) Unscoped() GormAdapter {
	return newGormAdapter(ga.DB.Unscoped())
}

func (ga *gormAdapter) Attrs(attrs ...interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Attrs(attrs...))
}

func (ga *gormAdapter) Assign(attrs ...interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Assign(attrs...))
}

func (ga *gormAdapter) First(out interface{}, where ...interface{}) GormAdapter {
	return newGormAdapter(ga.DB.First(out, where...))
}

func (ga *gormAdapter) Last(out interface{}, where ...interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Last(out, where...))
}

func (ga *gormAdapter) Find(out interface{}, where ...interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Find(out, where...))
}

func (ga *gormAdapter) Scan(dest interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Scan(dest))
}

func (ga *gormAdapter) Row() *sql.Row {
	return ga.DB.Row()
}

func (ga *gormAdapter) Rows() (*sql.Rows, error) {
	return ga.DB.Rows()
}

func (ga *gormAdapter) ScanRows(rows *sql.Rows, result interface{}) error {
	return ga.DB.ScanRows(rows, result)
}

func (ga *gormAdapter) Pluck(column string, value interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Pluck(column, value))
}

func (ga *gormAdapter) Count(value *int64) GormAdapter {
	return newGormAdapter(ga.DB.Count(value))
}

func (ga *gormAdapter) FirstOrInit(out interface{}, where ...interface{}) GormAdapter {
	return newGormAdapter(ga.DB.FirstOrInit(out, where...))
}

func (ga *gormAdapter) FirstOrCreate(out interface{}, where ...interface{}) GormAdapter {
	return newGormAdapter(ga.DB.FirstOrCreate(out, where...))
}

func (ga *gormAdapter) UpdateColumns(values interface{}) GormAdapter {
	return newGormAdapter(ga.DB.UpdateColumns(values))
}

func (ga *gormAdapter) Save(value interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Save(value))
}

func (ga *gormAdapter) Create(value interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Create(value))
}

func (ga *gormAdapter) Delete(value interface{}, where ...interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Delete(value, where...))
}

func (ga *gormAdapter) Raw(sql string, values ...interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Raw(sql, values...))
}

func (ga *gormAdapter) Exec(sql string, values ...interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Exec(sql, values...))
}

func (ga *gormAdapter) Model(value interface{}) GormAdapter {
	return newGormAdapter(ga.DB.Model(value))
}

func (ga *gormAdapter) Table(name string) GormAdapter {
	return newGormAdapter(ga.DB.Table(name))
}

func (ga *gormAdapter) Debug() GormAdapter {
	return newGormAdapter(ga.DB.Debug())
}

func (ga *gormAdapter) Begin() GormAdapter {
	return newGormAdapter(ga.DB.Begin())
}

func (ga *gormAdapter) Commit() GormAdapter {
	return newGormAdapter(ga.DB.Commit())
}

func (ga *gormAdapter) Rollback() GormAdapter {
	return newGormAdapter(ga.DB.Rollback())
}

func (ga *gormAdapter) Get(name string) (interface{}, bool) {
	return ga.DB.Get(name)
}

func (ga *gormAdapter) AddError(err error) error {
	return ga.DB.AddError(err)
}

func (ga *gormAdapter) RowsAffected() int64 {
	return ga.DB.RowsAffected
}

func (ga *gormAdapter) Error() error {
	return ga.DB.Error
}
