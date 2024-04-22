package core

import (
	"context"
	"crypto/x509"
	"database/sql"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Our App struct holds Servers, which trigger calls to our Service API on each gRPC/HTTP request, after applying interceptors.
// So: App -> Servers -> Interceptors -> Service.
//
// Our Service, assisted by our set of Tools (like TokenGenerator, PwdHasher, etc), interacts with our External Layer,
// which in turn just holds our Storage and Clients. Storage holds our SQL Database. Clients... is empty (for now).
// So: Service (with Tools) -> External Layer -> Storage -> SQL Database.
//
// While Storage is a high-level DB API, the SQL Database itself (DatabaseAPI interface) is much more low-level.
// Storage actually uses the SQL Database under the hood.
//
// To sum it all up:
// * App -> Servers -> Interceptors -> Service (with Tools) -> External Layer -> Storage -> SQL Database.

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Interfaces -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Embed all PB Services here! For now we only have 1.
// This interface is kind of our entire API. It has a method for each gRPC/HTTP endpoint we have.
type ServiceAPI interface {
	pbs.UsersServiceServer
}

// StorageAPI is our high-level API to interact with our DBs and similar.
// This would be named Repository instead of Storage only if it weren't so long. And Repo sucks.
type StorageAPI interface {
	CreateUser(username, hashedPwd string) (*models.User, error)
	GetUser(opts ...DBQueryOpt) (*models.User, error)
	GetUsers(page, pageSize int, opts ...DBQueryOpt) (models.Users, int, error)
}

// This is the interface we use to interact with our SQL Database.
// It's an adapter for Gorm. Concrete types gormAdapter and mocks.Gorm implement this.
type SQLDatabaseAPI interface {
	GetSQL() *sql.DB
	AddError(err error) error
	AutoMigrate(dst ...interface{}) error
	Count(value *int64) SQLDatabaseAPI
	Create(value interface{}) SQLDatabaseAPI
	Debug() SQLDatabaseAPI
	Delete(value interface{}, where ...interface{}) SQLDatabaseAPI
	Error() error
	Find(out interface{}, where ...interface{}) SQLDatabaseAPI
	First(out interface{}, where ...interface{}) SQLDatabaseAPI
	FirstOrCreate(out interface{}, where ...interface{}) SQLDatabaseAPI
	Group(query string) SQLDatabaseAPI
	Joins(query string, args ...interface{}) SQLDatabaseAPI
	Limit(value int) SQLDatabaseAPI
	Model(value interface{}) SQLDatabaseAPI
	Offset(value int) SQLDatabaseAPI
	Order(value string) SQLDatabaseAPI
	Or(query interface{}, args ...interface{}) SQLDatabaseAPI
	Pluck(column string, value interface{}) SQLDatabaseAPI
	Raw(sql string, values ...interface{}) SQLDatabaseAPI
	Rows() (*sql.Rows, error)
	RowsAffected() int64
	Row() *sql.Row
	Save(value interface{}) SQLDatabaseAPI
	Scan(dest interface{}) SQLDatabaseAPI
	Where(query interface{}, args ...interface{}) SQLDatabaseAPI
	Scopes(funcs ...func(SQLDatabaseAPI) SQLDatabaseAPI) SQLDatabaseAPI
}

// DBQueryOpt is any function which takes a DatabaseAPI instance and modifies it.
// We use it to apply different settings to our queries.
type DBQueryOpt func(SQLDatabaseAPI)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Lets you obtain all tools.Tools without having to import the tools package.
type Toolbox interface {
	GetRequestsValidator() RequestsValidator
	GetAuthenticator() TokenAuthenticator
	GetRateLimiter() *rate.Limiter
	GetPwdHasher() PwdHasher
	GetTLSServerCert() *x509.CertPool
	GetTLSServerCreds() credentials.TransportCredentials
	GetTLSClientCreds() credentials.TransportCredentials
}

type RequestsValidator interface {
	ValidateRequest(ctx context.Context, req any, i *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error)
}

type TokenGenerator interface {
	GenerateToken(id int, username string, role models.Role) (string, error)
}

type TokenValidator interface {
	ValidateToken(ctx context.Context, req any, i *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error)
}

type TokenAuthenticator interface {
	TokenGenerator
	TokenValidator
}

type PwdHasher interface {
	Hash(pwd string) string
	Compare(plainPwd, hashedPwd string) bool
}
