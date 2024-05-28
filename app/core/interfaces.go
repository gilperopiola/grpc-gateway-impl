package core

import (
	"crypto/x509"
	"database/sql"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/other"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// When a GRPC/HTTP Request arrives, our Servers pass it through Interceptors, and then through the Service.
// So: App -> Servers -> Interceptors -> Service.
//
// Our Service, assisted by our set of Tools (TokenGenerator, PwdHasher, etc), performs Actions (like GetUser or GenerateToken).
// These Actions sometimes let us communicate with external things, like a Database or the File System.
//
// To sum it all up:
// * App -> Servers -> Interceptors -> Service -> Actions -> External Resources (SQL Database, File System, etc).

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Interfaces -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

/* -~-~-~-~- Main Interfaces -~-~-~-~- */

type Servers interface {
	Run()
	Shutdown()
}

// This interface is kind of our entire API. It has a method for each GRPC/HTTP endpoint we have.
type Service interface {
	AuthSvc
	UsersSvc
	GroupsSvc

	RegisterGRPCServices(GRPCServiceRegistrar)
	RegisterHTTPServices(*HTTPMultiplexer, GRPCDialOptions)
}

// Used to kinda unify our SQL and Mongo DB Interfaces. Also lets us get the inner DB object which may be useful.
type DB interface {
	GetInnerDB() any
}

// With this you can avoid importing the tools pkg.
// Remember to add new tools on the app.go file as well.
type Toolbox interface {
	APIs
	DBTool
	TLSTool
	CtxManager
	FileManager
	TokenGenerator
	TokenValidator
	RequestsValidator
	ShutdownJanitor
	RateLimiter
	PwdHasher
	Retrier
}

/* -~-~-~-~- Toolbox: Tools -~-~-~-~- */

type (
	CtxManager interface {
		AddUserInfo(ctx Ctx, userID, username string) Ctx
		ExtractMetadata(ctx Ctx, key string) (string, error)
	}

	FileManager interface {
		CreateFolder(path string) error
		CreateFolders(paths ...string) error
	}

	PwdHasher interface {
		HashPassword(pwd string) string
		PasswordsMatch(plainPwd, hashedPwd string) bool
	}

	RateLimiter interface {
		LimitGRPC(c Ctx, r any, i *GRPCInfo, h GRPCHandler) (any, error) // grpc.UnaryServerInterceptor
	}

	ShutdownJanitor interface {
		AddCleanupFunc(fn func())
		AddCleanupFuncWithErr(fn func() error)
		Cleanup()
	}

	RequestsValidator interface {
		ValidateGRPC(c Ctx, r any, i *GRPCInfo, h GRPCHandler) (any, error) // grpc.UnaryServerInterceptor
	}

	Retrier interface {
		TryToConnectToDB(connectToDB func() (any, error), execOnFailure func()) (any, error)
	}

	TLSTool interface {
		GetServerCertificate() *x509.CertPool
		GetServerCreds() TLSCredentials
		GetClientCreds() TLSCredentials
	}

	TokenGenerator interface {
		GenerateToken(id int, username string, role Role) (string, error)
	}

	TokenValidator interface {
		ValidateToken(c Ctx, r any, i *GRPCInfo, h GRPCHandler) (any, error) // grpc.UnaryServerInterceptor
	}

	DBTool interface {
		GetDB() DB
		CloseDB()
		IsNotFound(err error) bool

		// Users
		CreateUser(ctx Ctx, username, hashedPwd string) (*User, error)
		GetUser(ctx Ctx, opts ...any) (*User, error)
		GetUsers(ctx Ctx, page, pageSize int, opts ...any) (Users, int, error)

		// Groups
		CreateGroup(ctx Ctx, name string, ownerID int, invitedUserIDs []int) (*Group, error)
		GetGroup(ctx Ctx, opts ...any) (*Group, error)
	}
)

type (
	APIs interface {
		InternalAPIs
		ExternalAPIs
	}

	InternalAPIs interface{}

	ExternalAPIs interface {
		WeatherAPI
	}

	WeatherAPI interface {
		GetCurrentWeather(ctx Ctx, lat, lon float64) (other.GetWeatherResponse, error)
	}
)

/* -~-~-~ SQL DB ~-~-~- */

// Low-level API for our SQL Database.
// It's an adapter for Gorm. Concrete types sql.sqlAdapter and mocks.Gorm implement this.
type SQLDB interface {
	DB
	AddError(err error) error
	AutoMigrate(dst ...any) error
	Association(column string) *gorm.Association
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
	WithContext(ctx Ctx) SQLDB
	Where(query any, args ...any) SQLDB
}

type SQLDBOpt func(SQLDB) // Variadic options

/* -~-~-~ Mongo DB ~-~-~- */

// Low-level API for our Mongo Database.
type MongoDB interface {
	DB
	Close(ctx Ctx)
	InsertOne(ctx Ctx, colName string, document any) (*mongo.InsertOneResult, error)
	Find(ctx Ctx, colName string, filter any, limit, offset int) (*mongo.Cursor, error)
	FindOne(ctx Ctx, colName string, filter any) *mongo.SingleResult
	DeleteOne(ctx Ctx, colName string, filter any) (*mongo.DeleteResult, error)
	Count(ctx Ctx, colName string, filter any) (int64, error)
}

type MongoDBOpt func(*bson.D) // Variadic options

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
