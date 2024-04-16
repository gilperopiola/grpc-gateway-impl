package special_types

import (
	"database/sql"

	"gorm.io/gorm"
)

type ServerLayer struct {
	GRPCServer Server
	HTTPServer Server
}

// Server abstracts the gRPC Server & HTTP Gateway.
// We used to need it to avoid import cycles, idk now.
// I think it's still useful somehow, but I'm not sure how.
type Server interface {
	Run()
	Shutdown()
}

// SQLDB is our adapter interface for Gorm.
// Concrete types gormAdapter and mocks.Gorm implement this.

type SQLDB interface {
	AddError(err error) error
	AutoMigrate(dst ...interface{}) error
	Count(value *int64) SQLDB
	Create(value interface{}) SQLDB
	Debug() SQLDB
	Delete(value interface{}, where ...interface{}) SQLDB
	Error() error
	Find(out interface{}, where ...interface{}) SQLDB
	First(out interface{}, where ...interface{}) SQLDB
	FirstOrCreate(out interface{}, where ...interface{}) SQLDB
	GetSQL() *sql.DB
	Group(query string) SQLDB
	Joins(query string, args ...interface{}) SQLDB
	Limit(value int) SQLDB
	Model(value interface{}) SQLDB
	Offset(value int) SQLDB
	Order(value string) SQLDB
	Or(query interface{}, args ...interface{}) SQLDB
	Pluck(column string, value interface{}) SQLDB
	Raw(sql string, values ...interface{}) SQLDB
	Rows() (*sql.Rows, error)
	RowsAffected() int64
	Row() *sql.Row
	Save(value interface{}) SQLDB
	Scan(dest interface{}) SQLDB
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) SQLDB
	Where(query interface{}, args ...interface{}) SQLDB
}
