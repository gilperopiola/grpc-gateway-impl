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

// We first have all the actual Mock Functions from the db.GormAdapter interface.
// Then we have the Mock Helpers that help us control the behavior of the mock.

/* ----------------------------------- */
/*          - Mock Functions -         */
/* ----------------------------------- */

func (m *Gorm) AddError(err error) error { args := m.Called(err); return args.Error(0) }

func (m *Gorm) AutoMigrate(dst ...interface{}) error { args := m.Called(dst); return args.Error(0) }

func (m *Gorm) Create(value interface{}) db.GormAdapter { m.Called(value); return m }

func (m *Gorm) Count(value *int64) db.GormAdapter { m.Called(value); return m }

func (m *Gorm) Debug() db.GormAdapter { m.Called(); return m }

func (m *Gorm) Delete(value interface{}, where ...interface{}) db.GormAdapter {
	m.Called(value, where)
	return m
}

func (m *Gorm) Error() error { args := m.Called(); return args.Error(0) }

func (m *Gorm) Find(out interface{}, where ...interface{}) db.GormAdapter {
	m.Called(out, where)
	return m
}

func (m *Gorm) First(out interface{}, where ...interface{}) db.GormAdapter {
	m.Called(out, where)
	return m
}

func (m *Gorm) FirstOrCreate(out interface{}, where ...interface{}) db.GormAdapter {
	m.Called(out, where)
	return m
}

func (m *Gorm) GetSQL() *sql.DB { args := m.Called(); return args.Get(0).(*sql.DB) }

func (m *Gorm) Group(query string) db.GormAdapter { m.Called(query); return m }

func (m *Gorm) Joins(query string, args ...interface{}) db.GormAdapter {
	m.Called(query, args)
	return m
}

func (m *Gorm) Limit(value int) db.GormAdapter { m.Called(value); return m }

func (m *Gorm) Model(model interface{}) db.GormAdapter { m.Called(model); return m }

func (m *Gorm) Offset(value int) db.GormAdapter { m.Called(value); return m }

func (m *Gorm) Or(query interface{}, args ...interface{}) db.GormAdapter {
	m.Called(query, args)
	return m
}

func (m *Gorm) Order(value string) db.GormAdapter { m.Called(value); return m }

func (m *Gorm) Pluck(column string, value interface{}) db.GormAdapter {
	m.Called(column, value)
	return m
}

func (m *Gorm) Raw(sql string, values ...interface{}) db.GormAdapter { m.Called(sql, values); return m }

func (m *Gorm) Row() *sql.Row { args := m.Called(); return args.Get(0).(*sql.Row) }

func (m *Gorm) Rows() (*sql.Rows, error) {
	args := m.Called()
	return args.Get(0).(*sql.Rows), args.Error(1)
}

func (m *Gorm) RowsAffected() int64 { args := m.Called(); return args.Get(0).(int64) }

func (m *Gorm) Save(value interface{}) db.GormAdapter { m.Called(value); return m }

func (m *Gorm) Scan(dest interface{}) db.GormAdapter { m.Called(dest); return m }

func (m *Gorm) Scopes(funcs ...func(*gorm.DB) *gorm.DB) db.GormAdapter { m.Called(funcs); return m }

func (m *Gorm) Table(name string) db.GormAdapter { m.Called(name); return m }

func (m *Gorm) Where(query interface{}, args ...interface{}) db.GormAdapter {
	m.Called(query, args)
	return m
}

/* ----------------------------------- */
/*           - Mock Helpers  -         */
/* ----------------------------------- */

func (m *Gorm) OnModel(model interface{}) *Gorm {
	m.On("Model", model).Return(m).Once()
	return m
}

func (m *Gorm) OnCount(count int) *Gorm {
	updateCount := func(args mock.Arguments) {
		*args.Get(0).(*int64) = int64(count)
	}
	m.On("Count", mock.AnythingOfType("*int64")).Run(updateCount).Return(m).Once()
	return m
}

func (m *Gorm) OnOffset() *Gorm {
	m.On("Offset", mock.AnythingOfType("int")).Return(m).Once()
	return m
}

func (m *Gorm) OnLimit() *Gorm {
	m.On("Limit", mock.AnythingOfType("int")).Return(m).Once()
	return m
}

func (m *Gorm) ErrorWillBe(err error) {
	m.On("Error").Return(err).Once()
}

/* ----------------------------------- */
/*       - Users Mock Helpers  -       */
/* ----------------------------------- */

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

func (m *Gorm) OnFindUsers(usersBefore, usersAfter *models.Users) *Gorm {
	updateUsers := func(args mock.Arguments) {
		users := args.Get(0).(*models.Users)
		if usersAfter == nil {
			users = nil
			return
		}
		*users = *usersAfter
	}
	m.On("Find", usersBefore, []interface{}(nil)).Run(updateUsers).Return(m).Once()
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

func (m *Gorm) OnWhereUser(id int, username string) *Gorm {
	if id != 0 {
		m.On("Where", "id = ?", []interface{}{fmt.Sprint(id)}).Return(m).Once()
	}
	if username != "" {
		m.On("Where", "username = ?", []interface{}{username}).Return(m).Once()
	}
	return m
}
