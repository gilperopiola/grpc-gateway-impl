package core

import (
	"context"
	"crypto/x509"
	"database/sql"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// When a GRPC/HTTP Request arrives, our Servers pass it through Interceptors, and then through the Service.
// So: App -> Servers -> Interceptors -> Service.
//
// Our Service, assisted by our set of Tools (like TokenGenerator, PwdHasher, etc), interacts with our External Layer,
// which in turn just holds our Storage and Clients. Storage holds our SQL Database. Clients is empty (for now).
// So: Service (with Tools) -> External Layer -> Storage -> SQL Database.
//
// To sum it all up:
// * App -> Servers -> Interceptors -> Service (with Tools) -> External Layer -> Storage -> SQL Database.

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Interfaces -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

/* -~-~-~-~- Service Layer -~-~-~-~- */

// Embed all PB Services here! For now we only have 1.
// This interface is kind of our entire API. It has a method for each GRPC/HTTP endpoint we have.
type ServiceLayer interface {
	pbs.UsersServiceServer
}

/* -~-~-~- External Layer -~-~-~- */

type ExternalLayer interface {
	GetStorage() StorageAPI
	GetClients() ClientsAPI
}

// High-level API to interact with our DBs (or any Storage we implement).
// This would be named Repository instead of Storage if it was shorter. And Repo sucks.
type StorageAPI interface {
	CreateUser(ctx context.Context, username, hashedPwd string) (*User, error)
	GetUser(ctx context.Context, opts ...any) (*User, error)
	GetUsers(ctx context.Context, page, pageSize int, opts ...any) (Users, int, error)
	CloseDB() // Either SQL or Mongo.
}

type ClientsAPI interface {
	// We still don't have any Clients.
}

/* -~-~-~ SQL Database ~-~-~- */

// Low-level API for our SQL Database.
// It's an adapter for Gorm. Concrete types sqlAdapter and mocks.Gorm implement this.
type SQLDatabaseAPI interface {
	AddError(err error) error
	AutoMigrate(dst ...any) error
	Close()
	Count(value *int64) SQLDatabaseAPI
	Create(value any) SQLDatabaseAPI
	Debug() SQLDatabaseAPI
	Delete(value any, where ...any) SQLDatabaseAPI
	Error() error
	Find(out any, where ...any) SQLDatabaseAPI
	First(out any, where ...any) SQLDatabaseAPI
	FirstOrCreate(out any, where ...any) SQLDatabaseAPI
	Group(query string) SQLDatabaseAPI
	Joins(query string, args ...any) SQLDatabaseAPI
	Limit(value int) SQLDatabaseAPI
	Model(value any) SQLDatabaseAPI
	Offset(value int) SQLDatabaseAPI
	Or(query any, args ...any) SQLDatabaseAPI
	Order(value string) SQLDatabaseAPI
	Pluck(column string, value any) SQLDatabaseAPI
	Raw(sql string, values ...any) SQLDatabaseAPI
	Rows() (*sql.Rows, error)
	RowsAffected() int64
	Row() *sql.Row
	Save(value any) SQLDatabaseAPI
	Scan(dest any) SQLDatabaseAPI
	Scopes(funcs ...func(SQLDatabaseAPI) SQLDatabaseAPI) SQLDatabaseAPI
	WithContext(ctx context.Context) SQLDatabaseAPI
	Where(query any, args ...any) SQLDatabaseAPI
}

type SQLQueryOpt func(SQLDatabaseAPI) // Variadic func.

/* -~-~-~ Mongo Database ~-~-~- */

// Low-level API for our Mongo Database.
type MongoDatabaseAPI interface {
	Close(ctx context.Context)
	InsertOne(ctx context.Context, colName string, document any) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, colName string, filter any, limit, offset int) (*mongo.Cursor, error)
	FindOne(ctx context.Context, colName string, filter any) *mongo.SingleResult
	DeleteOne(ctx context.Context, colName string, filter any) (*mongo.DeleteResult, error)
	Count(ctx context.Context, colName string, filter any) (int64, error)
}

type MongoQueryOpt func(*bson.D) // Variadic func.

/* -~-~-~- Toolbox -~-~-~- */

// Use this to avoid importing the tools pkg.
// Our app.Tools fulfills this interface.
type Toolbox interface {
	GetAuthenticator() TokenAuthenticator
	GetRequestsValidator() RequestsValidator
	GetRateLimiter() RateLimiter
	GetPwdHasher() PwdHasher
	GetTLSTool() TLSTool

	TokenAuthenticator
	RequestsValidator
	RateLimiter
	PwdHasher
	TLSTool
}

type TokenAuthenticator interface {
	TokenGenerator
	TokenValidator
}

type TokenGenerator interface {
	GenerateToken(id int, username string, role Role) (string, error)
}

type PwdHasher interface {
	HashPassword(pwd string) string
	PasswordsMatch(plainPwd, hashedPwd string) bool
}

type TLSTool interface {
	GetServerCertificate() *x509.CertPool
	GetServerCreds() credentials.TransportCredentials
	GetClientCreds() credentials.TransportCredentials
}

// These below are grpc.UnaryInterceptor funcs.

type TokenValidator interface {
	ValidateToken(c context.Context, r any, i *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error)
}

type RequestsValidator interface {
	ValidateGRPC(c context.Context, r any, i *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error)
}

type RateLimiter interface {
	LimitGRPC(c context.Context, r any, i *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error)
}
