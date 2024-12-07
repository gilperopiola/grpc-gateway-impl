package core

import (
	"crypto/x509"
	"image"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/credentials"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*      - Interfaces | Clients -       */
/* -~-~-~-~-~-~-~-~-~-~ APIs and DBs - */

// With this you can avoid importing the tools pkg.
type (

	/* -~-~-~- Tools: Security -~-~-~- */

	// Used to manage TLS certificates and credentials.
	TLSManager interface {
		GetServerCertificate() *x509.CertPool
		GetServerCreds() credentials.TransportCredentials
		GetClientCreds() credentials.TransportCredentials
	}

	/* -~-~-~- Tools: Request handling -~-~-~- */

	// For pagination, we add a page and pageSize fields to all incoming 'GET Many' requests in the .protos:
	//
	//	▶ optional int32 page = 1 		[json_name = "page"];
	//	▶ optional int32 page_size = 3 	[json_name = "page_size"];
	//
	// get a list of many resources Used to obtain page and pageSize from paginated requests and also to compose
	// the *pbs.PaginationInfo struct for the corresponding paginated response.
	// Designed to work with GRPC, usually on the 'GetMany' methods.
	RequestPaginator interface {
		PaginatedRequest(req PaginatedRequest) (page int, pageSize int)
		PaginatedResponse(currentPage, pageSize, totalRecords int) *pbs.PaginationInfo
	}

	// This isn't a tool, but a type used by the RequestsPaginator tool.
	// The Paginator only works on protobuf autogenerated structs that have GetPage() and GetPageSize() methods.
	PaginatedRequest interface {
		GetPage() int32
		GetPageSize() int32
		// A .proto example would be a message with these 2 fields:
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
		GenerateToken(id int, username string, role models.UserRole) (string, error)
	}

	// Validates authorization tokens.
	// Current implementation uses JWT.
	TokenValidator interface {
		ValidateToken(ctx god.Ctx, req any, route Route) (Claims, error)
	}

	Claims interface {
		GetUserInfo() (id, username string)
	}

	/* -~-~-~- Tools: Other -~-~-~- */

	// Used to add and get values from a request's context or headers.
	ContextManager interface {
		AddToCtx(ctx god.Ctx, key, value string) god.Ctx
		GetFromCtx(ctx god.Ctx, key string) (string, error)
		GetFromCtxMD(ctx god.Ctx, key string) (string, error)

		AddUserInfoToCtx(ctx god.Ctx, userID, username string) god.Ctx
		GetUserIDFromCtx(ctx god.Ctx) string
		GetUsernameFromCtx(ctx god.Ctx) string
	}

	// File system operations.
	FileManager interface {
		CreateFolder(path string) error
		CreateFolders(paths ...string) error
	}

	ImageLoader interface {
		LoadImgFromFile(path string) (image.Image, error)
		LoadImgFromURL(url string) (image.Image, error)
		LoadImgFromBytes(b []byte) (image.Image, error)
		LoadImgFromBase64(b64 string) (image.Image, error)
	}

	// Encode and decode models.
	ModelConverter interface {
		UserToUserInfoPB(*models.User) *pbs.UserInfo
		UsersToUsersInfoPB([]*models.User) []*pbs.UserInfo

		GroupToGroupInfoPB(*models.Group) *pbs.GroupInfo
		GroupsToGroupsInfoPB([]*models.Group) []*pbs.GroupInfo
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

type (
	// Low-level InnerSqlDB DB interface
	// It's an adapter for Gorm
	InnerSqlDB interface {
		Close()
		AddError(err error) error
		AutoMigrate(dst ...any) error
		Association(column string) SqlDBAssoc
		Count(value *int64) InnerSqlDB
		Create(value any) InnerSqlDB
		Debug() InnerSqlDB
		Delete(value any, where ...any) InnerSqlDB
		Error() error
		Find(out any, where ...any) InnerSqlDB
		First(out any, where ...any) InnerSqlDB
		FirstOrCreate(out any, where ...any) InnerSqlDB
		Group(query string) InnerSqlDB
		InsertAdmin(hashedPwd string)
		Joins(query string, args ...any) InnerSqlDB
		Limit(value int) InnerSqlDB
		Model(value any) InnerSqlDB
		Offset(value int) InnerSqlDB
		Or(query any, args ...any) InnerSqlDB
		Order(value string) InnerSqlDB
		Pluck(column string, value any) InnerSqlDB
		Preload(query string, args ...any) InnerSqlDB
		Raw(sql string, values ...any) InnerSqlDB
		RowsAffected() int64
		Save(value any) InnerSqlDB
		Scan(dest any) InnerSqlDB
		Scopes(funcs ...func(InnerSqlDB) InnerSqlDB) InnerSqlDB
		Unscoped() InnerSqlDB
		WithContext(ctx god.Ctx) InnerSqlDB
		Where(query any, args ...any) InnerSqlDB
	}

	// Low-level InnerMongoDB DB interface
	InnerMongoDB interface {
		Close(ctx god.Ctx)
		Count(ctx god.Ctx, colName string, filter any) (int64, error)
		Find(ctx god.Ctx, colName string, filter any, limit, offset int) (*mongo.Cursor, error)
		FindOne(ctx god.Ctx, colName string, filter any) *mongo.SingleResult
		InsertOne(ctx god.Ctx, colName string, document any) (*mongo.InsertOneResult, error)
		DeleteOne(ctx god.Ctx, colName string, filter any) (*mongo.DeleteResult, error)
	}
)

type (
	SqlDBOpt   func(InnerSqlDB) // Optional functions to apply to a query.
	MongoDBOpt func(*bson.D)    // Optional functions to apply to a query.
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
		ContextManager
	}

	// Unimplemented.
	RedisKVTool interface {
		KVTool
		// RedisTool?
	}
)
