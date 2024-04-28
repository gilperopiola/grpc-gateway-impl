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

/* -~-~-~-~- Service -~-~-~-~- */

// Embed all PB Services here! For now we only have 1.
// This interface is kind of our entire API. It has a method for each GRPC/HTTP endpoint we have.
type Service interface {
	pbs.UsersServiceServer
}

/* -~-~-~-~- Servers -~-~-~-~- */

type Servers interface {
	Run()
	Shutdown()
}

/* -~-~-~- Toolbox -~-~-~- */

// Use this to avoid importing the tools pkg.
// Our app.Toolbox fulfills this interface.
type Toolbox interface {
	APICaller
	DBTool
	PwdHasher
	RateLimiter
	RequestsValidator
	TLSTool
	TokenAuthenticator
}

// -> All our Tools have a Getter method, like GetDBTool for the DBTool.
// -> May seem redundant, but we need them to be able to get each particular Tool from the Toolbox,
// -> as the Toolbox interface abstracts them away.

type (
	APICaller interface {
		APICallerGetter
		// We still don't have any Clients.
	}
	DBTool interface { // Connects with a Database. Has implementations for SQL and Mongo.
		DBToolGetter
		GetDB() DB
		CreateUser(ctx context.Context, username, hashedPwd string) (*User, error)
		GetUser(ctx context.Context, opts ...any) (*User, error)
		GetUsers(ctx context.Context, page, pageSize int, opts ...any) (Users, int, error)
	}
	PwdHasher interface {
		PwdHasherGetter
		HashPassword(pwd string) string
		PasswordsMatch(plainPwd, hashedPwd string) bool
	}
	RateLimiter interface {
		RateLimiterGetter
		LimitGRPC(c context.Context, r any, i *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error) // grpc.UnaryServerInterceptor
	}
	RequestsValidator interface {
		RequestsValidatorGetter
		ValidateGRPC(c context.Context, r any, i *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error) // grpc.UnaryServerInterceptor
	}
	TLSTool interface {
		TLSToolGetter
		GetServerCertificate() *x509.CertPool
		GetServerCreds() credentials.TransportCredentials
		GetClientCreds() credentials.TransportCredentials
	}
	TokenAuthenticator interface {
		TokenAuthenticatorGetter
		GenerateToken(id int, username string, role Role) (string, error)
		ValidateToken(c context.Context, r any, i *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error) // grpc.UnaryServerInterceptor
	}
)

type (
	APICallerGetter          interface{ GetAPICaller() APICaller }
	DBToolGetter             interface{ GetDBTool() DBTool }
	PwdHasherGetter          interface{ GetPwdHasher() PwdHasher }
	RateLimiterGetter        interface{ GetRateLimiter() RateLimiter }
	RequestsValidatorGetter  interface{ GetRequestsValidator() RequestsValidator }
	TLSToolGetter            interface{ GetTLSTool() TLSTool }
	TokenAuthenticatorGetter interface{ GetTokenAuthenticator() TokenAuthenticator }
)

/* -~-~-~ Databases ~-~-~- */

// Used to kinda unify our SQL and Mongo DB Interfaces. Also lets us get the inner DB object which may be useful.
type DB interface {
	GetInnerDB() any
}

/* -~-~-~ SQL Database ~-~-~- */

// Low-level API for our SQL Database.
// It's an adapter for Gorm. Concrete types sql.sqlAdapter and mocks.Gorm implement this.
type SQLDB interface {
	DB
	AddError(err error) error
	AutoMigrate(dst ...any) error
	Close()
	Count(value *int64) SQLDB
	Create(value any) SQLDB
	Debug() SQLDB
	Delete(value any, where ...any) SQLDB
	Error() error
	Find(out any, where ...any) SQLDB
	First(out any, where ...any) SQLDB
	FirstOrCreate(out any, where ...any) SQLDB
	Group(query string) SQLDB
	InsertAdmin(hashedPwd string)
	Joins(query string, args ...any) SQLDB
	Limit(value int) SQLDB
	Model(value any) SQLDB
	Offset(value int) SQLDB
	Or(query any, args ...any) SQLDB
	Order(value string) SQLDB
	Pluck(column string, value any) SQLDB
	Raw(sql string, values ...any) SQLDB
	Rows() (*sql.Rows, error)
	RowsAffected() int64
	Row() *sql.Row
	Save(value any) SQLDB
	Scan(dest any) SQLDB
	Scopes(funcs ...func(SQLDB) SQLDB) SQLDB
	WithContext(ctx context.Context) SQLDB
	Where(query any, args ...any) SQLDB
}

type SQLDBOpt func(SQLDB) // Variadic func.

/* -~-~-~ Mongo Database ~-~-~- */

// Low-level API for our Mongo Database.
type MongoDB interface {
	DB
	Close(ctx context.Context)
	InsertOne(ctx context.Context, colName string, document any) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, colName string, filter any, limit, offset int) (*mongo.Cursor, error)
	FindOne(ctx context.Context, colName string, filter any) *mongo.SingleResult
	DeleteOne(ctx context.Context, colName string, filter any) (*mongo.DeleteResult, error)
	Count(ctx context.Context, colName string, filter any) (int64, error)
}

type MongoDBOpt func(*bson.D) // Variadic func.
