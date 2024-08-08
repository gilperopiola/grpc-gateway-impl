package core

import (
	"crypto/x509"
	"net/http"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/apis/apimodels"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Interfaces -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

/* -~-~-~-~- Services -~-~-~-~- */

// These are our different Services, as defined on the .proto files.
type (
	AuthSvc   = pbs.AuthServiceServer
	UsersSvc  = pbs.UsersServiceServer
	GroupsSvc = pbs.GroupsServiceServer
	HealthSvc = pbs.HealthServiceServer
)

/* -~-~-~-~- Tools -~-~-~-~- */

// These are the interfaces to all of our tools.
// Each concrete tool lives on the tools pkg.

type (

	// With this you can avoid importing the tools pkg.
	// Remember to add new tools on the app.go file as well.
	Tools interface {
		APIs
		DBTool
		TLSTool
		CtxTool
		FileManager
		ModelConverter
		TokenGenerator
		TokenValidator
		RequestPaginator
		RequestValidator
		ShutdownJanitor
		RateLimiter
		PwdHasher
	}

	/* -~-~-~- Tools: APIs -~-~-~- */

	// Contains all external and internal API clients.
	APIs interface {
		OpenWeatherAPI
	}

	// OpenWeatherMap API client.
	OpenWeatherAPI interface {
		GetCurrentWeather(ctx god.Ctx, lat, lon float64) (*apimodels.GetCurrentWeatherResponse, error)
	}

	/* -~-~-~- Tools: DBs -~-~-~- */

	// High-level interactions for any DB.
	// It connects the Service with the DB, calling the low-level db-type-specific methods.
	DBTool interface {
		GetDB() AnyDB
		CloseDB()
		IsNotFound(err error) bool

		// Users
		CreateUser(ctx god.Ctx, username, hashedPwd string) (*models.User, error)
		GetUser(ctx god.Ctx, opts ...any) (*models.User, error)
		GetUsers(ctx god.Ctx, page, pageSize int, opts ...any) (models.Users, int, error)

		// Groups
		CreateGroup(ctx god.Ctx, name string, ownerID int, invitedUserIDs []int) (*models.Group, error)
		GetGroup(ctx god.Ctx, opts ...any) (*models.Group, error)
	}

	// Low-level interactions for our SQL DB. It's an adapter for Gorm.
	// Both concrete types sqldb.sqlAdapter and mocks.Gorm implement this.
	SqlDB interface {
		AnyDB
		AddError(err error) error
		AutoMigrate(dst ...any) error
		Association(column string) SqlDBAssoc
		Close()
		Count(value *int64) SqlDB
		Create(value any) SqlDB
		Debug() SqlDB
		Delete(value any, where ...any) SqlDB
		Error() error
		Find(out any, where ...any) SqlDB
		First(out any, where ...any) SqlDB
		FirstOrCreate(out any, where ...any) SqlDB
		Group(query string) SqlDB
		InsertAdmin(hashedPwd string)
		Joins(query string, args ...any) SqlDB
		Limit(value int) SqlDB
		Model(value any) SqlDB
		Offset(value int) SqlDB
		Or(query any, args ...any) SqlDB
		Order(value string) SqlDB
		Pluck(column string, value any) SqlDB
		Raw(sql string, values ...any) SqlDB
		Row() SqlRow
		Rows() (SqlRows, error)
		RowsAffected() int64
		Save(value any) SqlDB
		Scan(dest any) SqlDB
		Scopes(funcs ...func(SqlDB) SqlDB) SqlDB
		WithContext(ctx god.Ctx) SqlDB
		Where(query any, args ...any) SqlDB
	}

	// Low-level interactions for our Mongo DB.
	MongoDB interface {
		AnyDB
		Close(ctx god.Ctx)

		Count(ctx god.Ctx, colName string, filter any) (int64, error)
		Find(ctx god.Ctx, colName string, filter any, limit, offset int) (*mongo.Cursor, error)
		FindOne(ctx god.Ctx, colName string, filter any) *mongo.SingleResult
		InsertOne(ctx god.Ctx, colName string, document any) (*mongo.InsertOneResult, error)
		DeleteOne(ctx god.Ctx, colName string, filter any) (*mongo.DeleteResult, error)
	}

	// Unifies our SQL and Mongo interfaces.
	AnyDB interface {
		GetInnerDB() any // Implementations return *gorm.DB or *mongo.Client
	}

	/* -~-~-~- Tools: Security -~-~-~- */

	// Used to manage TLS certificates and credentials.
	TLSTool interface {
		GetServerCertificate() *x509.CertPool
		GetServerCreds() god.TLSCreds
		GetClientCreds() god.TLSCreds
	}

	/* -~-~-~- Tools: Request handling -~-~-~- */

	// Used to obtain page and pageSize from paginated requests and also to compose
	// the *pbs.PaginationInfo struct for the corresponding paginated response.
	// Designed to work with GRPC, usually on the 'GetMany' methods.
	RequestPaginator interface {
		PaginatedRequest(req PaginatedRequest) (page int, pageSize int)
		PaginatedResponse(currentPage, pageSize, totalRecords int) *pbs.PaginationInfo
	}

	// Used to validate that incoming requests' fields follow our predefined rules and formats.
	// GRPC Interceptor.
	RequestValidator interface {
		ValidateRequest(req any) error
	}

	/* -~-~-~- Tools: Auth -~-~-~- */

	// Generates authorization tokens.
	// Current implementation uses JWT.
	TokenGenerator interface {
		GenerateToken(id int, username string, role models.Role) (string, error)
	}

	// Validates authorization tokens.
	// Current implementation uses JWT.
	TokenValidator interface {
		ValidateToken(ctx god.Ctx, req any, route string) (models.TokenClaims, error)
	}

	/* -~-~-~- Tools: Other -~-~-~- */

	// Used to add and get values from a request's context or headers.
	CtxTool interface {
		AddToCtx(ctx god.Ctx, key, value string) god.Ctx
		GetFromCtx(ctx god.Ctx, key string) (string, error)

		AddUserInfoToCtx(ctx god.Ctx, userID, username string) god.Ctx
		GetUserIDFromCtx(ctx god.Ctx) string
		GetUsernameFromCtx(ctx god.Ctx) string
	}

	// File system operations.
	FileManager interface {
		CreateFolder(path string) error
		CreateFolders(paths ...string) error
	}

	// Encode and decode models.
	ModelConverter interface {
		UserToUserInfoPB(*models.User) *pbs.UserInfo
		UsersToUsersInfoPB(models.Users) []*pbs.UserInfo

		GroupToGroupInfoPB(*models.Group) *pbs.GroupInfo
		GroupsToGroupsInfoPB(models.Groups) []*pbs.GroupInfo
	}

	// Hashes and compares passwords.
	PwdHasher interface {
		HashPassword(pwd string) string
		PasswordsMatch(plainPwd, hashedPwd string) bool
	}

	// Used to limit the rate of incoming requests.
	// GRPC Interceptor.
	RateLimiter interface {
		AllowRate() bool
	}

	// Used to cleanup resources before exiting.
	ShutdownJanitor interface {
		AddCleanupFunc(fn func())
		AddCleanupFuncWithErr(fn func() error)
		Cleanup()
	}
)

/* -~-~-~- Other types - SQL and Mongo ~-~-~- */

// This isn't used like the other Tools, as it's instantiated per request.
// The implementation lives on the tools pkg.
// Used to wrap the http.ResponseWriter and then Type Assert it to this to get the extra methods.
type HTTPRespWriter interface {
	http.ResponseWriter

	GetWrittenBody() []byte
	GetWrittenStatus() int
}

// This isn't a tool, but a type used by the RequestsPaginator tool.
// The Paginator only works on protobuf autogenerated structs that have GetPage() and GetPageSize() methods.
type PaginatedRequest interface {
	GetPage() int32
	GetPageSize() int32

	// A .proto example would be a message that contained these 2 fields:
	//	optional int32 page = 1 		[json_name = "page"];
	//	optional int32 page_size = 3 	[json_name = "page_size"];
}

/* -~-~-~- Other types - SQL and Mongo ~-~-~- */

type (
	SqlDBOpt   func(SqlDB)   // Optional functions to apply to a query.
	MongoDBOpt func(*bson.D) // Optional functions to apply to a query.

	SqlRow  interface{ Scan(dest ...any) error }
	SqlRows interface {
		Next() bool
		Scan(...any) error
		Close() error
	}
	SqlDBAssoc interface{ Append(...interface{}) error }
)

/* -~-~-~- Other types - Unimplemented Tools ~-~-~- */

type (

	// Unimplemented. Idea for a general key-value store interface,
	// with just Get and Set methods.
	KVTool interface {
		Get(key string) (string, error)
		Set(key, value string) error
	}

	// Unimplemented. The idea is to use the same Get and Set signature for all KV Tool interfaces,
	// and each implementation can also have its own methods.
	// Then we can write functions that accept a KVTool and works for any implementation.
	CtxKVTool interface {
		KVTool
		CtxTool
	}

	// Unimplemented.
	RedisKVTool interface {
		KVTool
		// RedisTool?
	}
)
