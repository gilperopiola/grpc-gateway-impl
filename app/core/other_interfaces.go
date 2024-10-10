package core

import (
	"context"
	"crypto/x509"
	"image"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/clients/apis/apimodels"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared"
	"google.golang.org/grpc/credentials"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*      - Interfaces | Clients -       */
/* -~-~-~-~-~-~-~-~-~-~ DBs and APIs - */

type (
	// Completions API
	ChatGPTAPI interface {
		SendToGPT(ctx context.Context, prompt string, prevMsgs ...apimodels.GPTMessage) (string, error)
	}

	// WeatherMap API client
	WeatherAPI interface {
		GetCurrentWeather(ctx god.Ctx, lat, lon float64) (*apimodels.GetWeatherResponse, error)
	}
)

// With this you can avoid importing the tools pkg.
type (

	/* -~-~-~- Tools: Security -~-~-~- */

	// Used to manage TLS certificates and credentials.
	TLSTool interface {
		GetServerCertificate() *x509.CertPool
		GetServerCreds() credentials.TransportCredentials
		GetClientCreds() credentials.TransportCredentials
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
		GenerateToken(id int, username string, role shared.Role) (string, error)
	}

	// Validates authorization tokens.
	// Current implementation uses JWT.
	TokenValidator interface {
		ValidateToken(ctx god.Ctx, req any, route string) (Claims, error)
	}

	Claims interface {
		GetUserInfo() (id, username string)
	}

	/* -~-~-~- Tools: Other -~-~-~- */

	// Used to add and get values from a request's context or headers.
	CtxTool interface {
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

type (
	// Low-level BaseSQLDB DB interface
	// It's an adapter for Gorm
	BaseSQLDB interface {
		Close()
		AddError(err error) error
		AutoMigrate(dst ...any) error
		Association(column string) SqlDBAssoc
		Count(value *int64) BaseSQLDB
		Create(value any) BaseSQLDB
		Debug() BaseSQLDB
		Delete(value any, where ...any) BaseSQLDB
		Error() error
		Find(out any, where ...any) BaseSQLDB
		First(out any, where ...any) BaseSQLDB
		FirstOrCreate(out any, where ...any) BaseSQLDB
		Group(query string) BaseSQLDB
		InsertAdmin(hashedPwd string)
		Joins(query string, args ...any) BaseSQLDB
		Limit(value int) BaseSQLDB
		Model(value any) BaseSQLDB
		Offset(value int) BaseSQLDB
		Or(query any, args ...any) BaseSQLDB
		Order(value string) BaseSQLDB
		Pluck(column string, value any) BaseSQLDB
		Preload(query string, args ...any) BaseSQLDB
		Raw(sql string, values ...any) BaseSQLDB
		Row() SqlRow
		Rows() (SqlRows, error)
		RowsAffected() int64
		Save(value any) BaseSQLDB
		Scan(dest any) BaseSQLDB
		Scopes(funcs ...func(BaseSQLDB) BaseSQLDB) BaseSQLDB
		Unscoped() BaseSQLDB
		WithContext(ctx god.Ctx) BaseSQLDB
		Where(query any, args ...any) BaseSQLDB
	}

	// Low-level BaseMongoDB DB interface
	BaseMongoDB interface {
		Close(ctx god.Ctx)
		Count(ctx god.Ctx, colName string, filter any) (int64, error)
		Find(ctx god.Ctx, colName string, filter any, limit, offset int) (*mongo.Cursor, error)
		FindOne(ctx god.Ctx, colName string, filter any) *mongo.SingleResult
		InsertOne(ctx god.Ctx, colName string, document any) (*mongo.InsertOneResult, error)
		DeleteOne(ctx god.Ctx, colName string, filter any) (*mongo.DeleteResult, error)
	}
)

type (
	SqlDBOpt   func(BaseSQLDB) // Optional functions to apply to a query.
	MongoDBOpt func(*bson.D)   // Optional functions to apply to a query.
	SqlRow     interface{ Scan(dest ...any) error }
	SqlRows    interface {
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
