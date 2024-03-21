package mocks

import (
	"database/sql"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/db"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// GormMock is a mock for the db.GormAdapter interface.
type GormMock struct {
	mock.Mock
}

func (m *GormMock) AutoMigrate(dst ...interface{}) error {
	args := m.Called(dst)
	return args.Error(0)
}

func (m *GormMock) GetSQLDB() *sql.DB {
	args := m.Called()
	return args.Get(0).(*sql.DB)
}

func (m *GormMock) Where(query interface{}, args ...interface{}) db.GormAdapter {
	m.Called(query, args)
	return m
}

func (m *GormMock) Or(query interface{}, args ...interface{}) db.GormAdapter {
	m.Called(query, args)
	return m
}

func (m *GormMock) Not(query interface{}, args ...interface{}) db.GormAdapter {
	m.Called(query, args)
	return m
}

func (m *GormMock) Limit(value int) db.GormAdapter {
	m.Called(value)
	return m
}

func (m *GormMock) Offset(value int) db.GormAdapter {
	m.Called(value)
	return m
}

func (m *GormMock) Order(value string) db.GormAdapter {
	m.Called(value)
	return m
}

func (m *GormMock) Select(query interface{}, args ...interface{}) db.GormAdapter {
	m.Called(query, args)
	return m
}

func (m *GormMock) Omit(columns ...string) db.GormAdapter {
	m.Called(columns)
	return m
}

func (m *GormMock) Group(query string) db.GormAdapter {
	m.Called(query)
	return m
}

func (m *GormMock) Having(query string, values ...interface{}) db.GormAdapter {
	m.Called(query, values)
	return m
}

func (m *GormMock) Joins(query string, args ...interface{}) db.GormAdapter {
	m.Called(query, args)
	return m
}

func (m *GormMock) Scopes(funcs ...func(*gorm.DB) *gorm.DB) db.GormAdapter {
	m.Called(funcs)
	return m
}

func (m *GormMock) Unscoped() db.GormAdapter {
	m.Called()
	return m
}

func (m *GormMock) Attrs(attrs ...interface{}) db.GormAdapter {
	m.Called(attrs)
	return m
}

func (m *GormMock) Assign(attrs ...interface{}) db.GormAdapter {
	m.Called(attrs)
	return m
}

func (m *GormMock) First(out interface{}, where ...interface{}) db.GormAdapter {
	m.Called(out, where)
	return m
}

func (m *GormMock) Last(out interface{}, where ...interface{}) db.GormAdapter {
	m.Called(out, where)
	return m
}

func (m *GormMock) Find(out interface{}, where ...interface{}) db.GormAdapter {
	m.Called(out, where)
	return m
}

func (m *GormMock) Scan(dest interface{}) db.GormAdapter {
	m.Called(dest)
	return m
}

func (m *GormMock) Row() *sql.Row {
	args := m.Called()
	return args.Get(0).(*sql.Row)
}

func (m *GormMock) Rows() (*sql.Rows, error) {
	args := m.Called()
	return args.Get(0).(*sql.Rows), args.Error(1)
}

func (m *GormMock) ScanRows(rows *sql.Rows, result interface{}) error {
	args := m.Called(rows, result)
	return args.Error(0)
}

func (m *GormMock) Pluck(column string, value interface{}) db.GormAdapter {
	m.Called(column, value)
	return m
}

func (m *GormMock) Count(value *int64) db.GormAdapter {
	m.Called(value)
	return m
}

func (m *GormMock) FirstOrInit(out interface{}, where ...interface{}) db.GormAdapter {
	m.Called(out, where)
	return m
}

func (m *GormMock) FirstOrCreate(out interface{}, where ...interface{}) db.GormAdapter {
	m.Called(out, where)
	return m
}

func (m *GormMock) UpdateColumns(values interface{}) db.GormAdapter {
	m.Called(values)
	return m
}

func (m *GormMock) Save(value interface{}) db.GormAdapter {
	m.Called(value)
	return m
}

func (m *GormMock) Create(value interface{}) db.GormAdapter {
	m.Called(value)
	return m
}

func (m *GormMock) Delete(value interface{}, where ...interface{}) db.GormAdapter {
	m.Called(value, where)
	return m
}

func (m *GormMock) Raw(sql string, values ...interface{}) db.GormAdapter {
	m.Called(sql, values)
	return m
}

func (m *GormMock) Exec(sql string, values ...interface{}) db.GormAdapter {
	m.Called(sql, values)
	return m
}

func (m *GormMock) Model(value interface{}) db.GormAdapter {
	m.Called(value)
	return m
}

func (m *GormMock) Table(name string) db.GormAdapter {
	m.Called(name)
	return m
}

func (m *GormMock) Debug() db.GormAdapter {
	m.Called()
	return m
}

func (m *GormMock) Begin() db.GormAdapter {
	m.Called()
	return m
}

func (m *GormMock) Commit() db.GormAdapter {
	m.Called()
	return m
}

func (m *GormMock) Rollback() db.GormAdapter {
	m.Called()
	return m
}

func (m *GormMock) Get(name string) (interface{}, bool) {
	args := m.Called(name)
	return args.Get(0), args.Bool(1)
}

func (m *GormMock) AddError(err error) error {
	args := m.Called(err)
	return args.Error(0)
}

func (m *GormMock) RowsAffected() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *GormMock) Error() error {
	args := m.Called()
	return args.Error(0)
}
