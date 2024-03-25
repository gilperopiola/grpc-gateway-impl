package mocks

import (
	"database/sql"
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/db"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Gorm is a mock for the db.GormAdapter interface.
type Gorm struct {
	mock.Mock
}

func (m *Gorm) GetSQL() *sql.DB {
	args := m.Called()
	return args.Get(0).(*sql.DB)
}

func (m *Gorm) ErrorWillBe(err error) {
	m.On("Error").Return(err).Once()
}

func (m *Gorm) Create(value interface{}) db.GormAdapter {
	m.Called(value)
	return m
}

func (m *Gorm) OnCreateUser(userBefore, userAfter *models.User) *Gorm {
	updateUser := func(args mock.Arguments) {
		user := args.Get(0).(*models.User)
		if userAfter == nil {
			user = nil
			return
		}
		*user = *userAfter
	}
	m.On("Create", userBefore).Run(updateUser).Return(m).Once()
	return m
}

func (m *Gorm) OnModel(model interface{}) *Gorm {
	m.On("Model", model).Return(m).Once()
	return m
}

func (m *Gorm) OnWhereIDOrUsername(id int, username string) *Gorm {
	if id != 0 {
		m.On("Where", "id = ?", []interface{}{fmt.Sprint(id)}).Return(m).Once()
	}
	if username != "" {
		m.On("Where", "username = ?", []interface{}{username}).Return(m).Once()
	}
	return m
}

func (m *Gorm) OnFirstUser(userBefore, userAfter *models.User) *Gorm {
	updateUser := func(args mock.Arguments) {
		user := args.Get(0).(*models.User)
		if userAfter == nil {
			user = nil
			return
		}
		*user = *userAfter
	}
	m.On("First", userBefore, []interface{}(nil)).Run(updateUser).Return(m).Once()
	return m
}

func (m *Gorm) First(out interface{}, where ...interface{}) db.GormAdapter {
	m.Called(out, where)
	return m
}

func (m *Gorm) Count(value *int64) db.GormAdapter {
	m.Called(value)
	return m
}

func (m *Gorm) Where(query interface{}, args ...interface{}) db.GormAdapter {
	m.Called(query, args)
	return m
}

func (m *Gorm) AutoMigrate(dst ...interface{}) error {
	args := m.Called(dst)
	return args.Error(0)
}

func (m *Gorm) Find(out interface{}, where ...interface{}) db.GormAdapter {
	m.Called(out, where)
	return m
}

func (m *Gorm) Or(query interface{}, args ...interface{}) db.GormAdapter {
	m.Called(query, args)
	return m
}

func (m *Gorm) Not(query interface{}, args ...interface{}) db.GormAdapter {
	m.Called(query, args)
	return m
}

func (m *Gorm) Limit(value int) db.GormAdapter {
	m.Called(value)
	return m
}

func (m *Gorm) Offset(value int) db.GormAdapter {
	m.Called(value)
	return m
}

func (m *Gorm) Order(value string) db.GormAdapter {
	m.Called(value)
	return m
}

func (m *Gorm) Select(query interface{}, args ...interface{}) db.GormAdapter {
	m.Called(query, args)
	return m
}

func (m *Gorm) Omit(columns ...string) db.GormAdapter {
	m.Called(columns)
	return m
}

func (m *Gorm) Group(query string) db.GormAdapter {
	m.Called(query)
	return m
}

func (m *Gorm) Having(query string, values ...interface{}) db.GormAdapter {
	m.Called(query, values)
	return m
}

func (m *Gorm) Joins(query string, args ...interface{}) db.GormAdapter {
	m.Called(query, args)
	return m
}

func (m *Gorm) Scopes(funcs ...func(*gorm.DB) *gorm.DB) db.GormAdapter {
	m.Called(funcs)
	return m
}

func (m *Gorm) Unscoped() db.GormAdapter {
	m.Called()
	return m
}

func (m *Gorm) Attrs(attrs ...interface{}) db.GormAdapter {
	m.Called(attrs)
	return m
}

func (m *Gorm) Assign(attrs ...interface{}) db.GormAdapter {
	m.Called(attrs)
	return m
}

func (m *Gorm) Last(out interface{}, where ...interface{}) db.GormAdapter {
	m.Called(out, where)
	return m
}

func (m *Gorm) Scan(dest interface{}) db.GormAdapter {
	m.Called(dest)
	return m
}

func (m *Gorm) Row() *sql.Row {
	args := m.Called()
	return args.Get(0).(*sql.Row)
}

func (m *Gorm) Rows() (*sql.Rows, error) {
	args := m.Called()
	return args.Get(0).(*sql.Rows), args.Error(1)
}

func (m *Gorm) ScanRows(rows *sql.Rows, result interface{}) error {
	args := m.Called(rows, result)
	return args.Error(0)
}

func (m *Gorm) Pluck(column string, value interface{}) db.GormAdapter {
	m.Called(column, value)
	return m
}
func (m *Gorm) FirstOrInit(out interface{}, where ...interface{}) db.GormAdapter {
	m.Called(out, where)
	return m
}

func (m *Gorm) FirstOrCreate(out interface{}, where ...interface{}) db.GormAdapter {
	m.Called(out, where)
	return m
}

func (m *Gorm) UpdateColumns(values interface{}) db.GormAdapter {
	m.Called(values)
	return m
}

func (m *Gorm) Save(value interface{}) db.GormAdapter {
	m.Called(value)
	return m
}

func (m *Gorm) Delete(value interface{}, where ...interface{}) db.GormAdapter {
	m.Called(value, where)
	return m
}

func (m *Gorm) Raw(sql string, values ...interface{}) db.GormAdapter {
	m.Called(sql, values)
	return m
}

func (m *Gorm) Exec(sql string, values ...interface{}) db.GormAdapter {
	m.Called(sql, values)
	return m
}

func (m *Gorm) Model(value interface{}) db.GormAdapter {
	m.Called(value)
	return m
}

func (m *Gorm) Table(name string) db.GormAdapter {
	m.Called(name)
	return m
}

func (m *Gorm) Debug() db.GormAdapter {
	m.Called()
	return m
}

func (m *Gorm) Begin() db.GormAdapter {
	m.Called()
	return m
}

func (m *Gorm) Commit() db.GormAdapter {
	m.Called()
	return m
}

func (m *Gorm) Rollback() db.GormAdapter {
	m.Called()
	return m
}

func (m *Gorm) Get(name string) (interface{}, bool) {
	args := m.Called(name)
	return args.Get(0), args.Bool(1)
}

func (m *Gorm) AddError(err error) error {
	args := m.Called(err)
	return args.Error(0)
}

func (m *Gorm) RowsAffected() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *Gorm) Error() error {
	args := m.Called()
	return args.Error(0)
}
