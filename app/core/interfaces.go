package core

import (
	"crypto/x509"
	"net/http"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Interfaces -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

/* -~-~-~-~- Main Interfaces -~-~-~-~- */

// These are our different Services, as defined on the .proto files.
type (
	AuthSvc   = pbs.AuthServiceServer
	UsersSvc  = pbs.UsersServiceServer
	GroupsSvc = pbs.GroupsServiceServer
	HealthSvc = pbs.HealthServiceServer
)

// With this you can avoid importing the tools pkg.
// Remember to add new tools on the app.go file as well.
type Tools interface {
	ExternalAPIs
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

// Unifies our SQL and Mongo interfaces.
type AnyDB interface {
	GetInnerDB() any // Implementations return *gorm.DB or *mongo.Client
}

/* -~-~-~-~- Tools: Tools -~-~-~-~- */

// These are the interfaces to all of our tools.
// Each concrete tool lives on the tools pkg.

type (
	CtxTool interface {
		AddUserInfoToCtx(ctx god.Ctx, userID, username string) god.Ctx
		GetFromCtx(ctx god.Ctx, key string) (string, error)
	}

	FileManager interface {
		CreateFolder(path string) error
		CreateFolders(paths ...string) error
	}

	ModelConverter interface {
		UserToUserInfoPB(*models.User) *pbs.UserInfo
		UsersToUsersInfoPB(models.Users) []*pbs.UserInfo

		GroupToGroupInfoPB(*models.Group) *pbs.GroupInfo
		GroupsToGroupsInfoPB(models.Groups) []*pbs.GroupInfo
	}

	PwdHasher interface {
		HashPassword(pwd string) string
		PasswordsMatch(plainPwd, hashedPwd string) bool
	}

	RateLimiter interface {
		AllowRate() bool
	}

	ShutdownJanitor interface {
		AddCleanupFunc(fn func())
		AddCleanupFuncWithErr(fn func() error)
		Cleanup()
	}

	// Used to obtain page and pageSize from paginated requests and also to compose
	// the *pbs.PaginationInfo struct for the corresponding paginated response.
	// Designed to work with GRPC, usually on the 'GetMany' methods.
	RequestPaginator interface {
		PaginatedRequest(req PaginatedRequest) (page int, pageSize int)
		PaginatedResponse(currentPage, pageSize, totalRecords int) *pbs.PaginationInfo
	}

	// Used to validate that incoming requests fields actually comply with
	// our predefined rules and formats.
	// Our only implementation is for GRPC.
	RequestValidator interface {
		ValidateRequest(req any) error
	}

	TLSTool interface {
		GetServerCertificate() *x509.CertPool
		GetServerCreds() god.TLSCreds
		GetClientCreds() god.TLSCreds
	}

	TokenGenerator interface {
		GenerateToken(id int, username string, role models.Role) (string, error)
	}

	TokenValidator interface {
		ValidateToken(ctx god.Ctx, req any, route string) (TokenClaims, error)
	}

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
)

type (
	InternalAPIs interface{}

	ExternalAPIs interface {
		OpenWeatherAPI
	}

	OpenWeatherAPI interface {
		GetCurrentWeather(ctx god.Ctx, lat, lon float64) (APIResponse, error)
	}

	APIResponse any
)

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

type TokenClaims interface {
	GetUserInfo() (id, username string)
}

/* -~-~-~- SQL DB ~-~-~- */

// Low-level API for our SQL Database.
// It's an adapter for Gorm. Concrete types sql.sqlAdapter and mocks.Gorm implement this.
type (
	SqlDB interface {
		AnyDB
		AddError(err error) error
		AutoMigrate(dst ...any) error
		Association(column string) SqlDBAssociation
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

	SqlDBOpt func(SqlDB) // Optional functions to apply to a query
)

// Used to avoid importing the gorm and sql pkgs: *gorm.Association, *sql.Row, *sql.Rows.
type (
	SqlDBAssociation interface {
		Append(values ...interface{}) error
	}
	SqlRow interface {
		Scan(dest ...any) error
	}
	SqlRows interface {
		Next() bool
		Scan(dest ...any) error
		Close() error
	}
)

/* -~-~-~- Mongo DB ~-~-~- */

// Low-level API for our Mongo Database.
type MongoDB interface {
	AnyDB
	Close(ctx god.Ctx)

	Count(ctx god.Ctx, colName string, filter any) (int64, error)
	Find(ctx god.Ctx, colName string, filter any, limit, offset int) (*mongo.Cursor, error)
	FindOne(ctx god.Ctx, colName string, filter any) *mongo.SingleResult
	InsertOne(ctx god.Ctx, colName string, document any) (*mongo.InsertOneResult, error)
	DeleteOne(ctx god.Ctx, colName string, filter any) (*mongo.DeleteResult, error)
}

type MongoDBOpt func(*bson.D) // Optional functions to apply to a query
