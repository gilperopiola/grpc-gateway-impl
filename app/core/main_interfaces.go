package core

// ➤ Am I sure this is the best way to go about stitching together the whole
// project here as interfaces? — Not sure, however this allowed me to decouple
// the tools/service/servers/clients packages, each of one importing core instead.
//
// ➤ If you don't wanna import core, just copy n paste the interfaces.

import (
	"net/http"
	"time"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/models"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Interfaces -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

/* -~-~-~-~- Main Interfaces -~-~-~-~- */

// Each service is defined on a .proto file
type (
	ServiceLayer interface {
		RegisterGRPCEndpoints(grpcServer grpc.ServiceRegistrar)
		RegisterHTTPEndpoints(mux any, opts ...grpc.DialOption)
	}
	AuthSvc   = pbs.AuthServiceServer
	UsersSvc  = pbs.UsersSvcServer
	GroupsSvc = pbs.GroupsServiceServer
	GPTSvc    = pbs.GPTServiceServer
	HealthSvc = pbs.HealthServiceServer
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)

	LogGRPC(route string, duration time.Duration, err error)
	LogHTTPRequest(handler http.Handler) http.Handler
	LogDebug(msg string)
	LogUnexpected(err error) error
	LogIfErr(err error, optionalFmt ...string)
	LogFatal(err error)
	LogFatalIfErr(err error, optionalFmt ...string)
	WarnIfErr(err error, optionalFmt ...string)
	LogImportant(msg string)
	LogStrange(msg string, info ...any)
	LogThreat(msg string)
	LogResult(ofWhat string, err error)
	LogAPICall(url string, status int, body []byte)
	Sync() error
}

/* -~-~-~-~- Clients -~-~-~-~- */

type Clients interface {
	APIs
	DB
}

type (
	APIs interface {
		ChatGPTAPI
		WeatherAPI
	}

	DB interface {
		DBActions
		GetDB() any // *gorm.DB or *mongo.Client
		CloseDB()
	}

	// High-level DB operations
	DBActions interface {
		DBCreateUser(ctx god.Ctx, username, hashedPwd string) (*models.User, error)
		DBGetUser(ctx god.Ctx, opts ...any) (*models.User, error)
		DBGetUsers(ctx god.Ctx, page, pageSize int, opts ...any) ([]*models.User, int, error)

		DBCreateGroup(ctx god.Ctx, name string, ownerID int, invitedUserIDs []int) (*models.Group, error)
		DBGetGroup(ctx god.Ctx, opts ...any) (*models.Group, error)

		DBGetGPTChat(ctx god.Ctx, opts ...any) (*models.GPTChat, error)
		DBCreateGPTChat(ctx god.Ctx, title string) (*models.GPTChat, error)
		DBCreateGPTMessage(ctx god.Ctx, message *models.GPTMessage) (*models.GPTMessage, error)
	}
)

/* -~-~-~-~- Tools -~-~-~-~- */

// These are the interfaces to all of our tools.
// Each concrete tool lives on the tools pkg.
type Tools interface {
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

/* -~-~-~- Other types - GRPC and HTTP ~-~-~- */

// This isn't used like the other Tools, as it's instantiated per request.
// The implementation lives on the tools pkg.
// Used to wrap the http.ResponseWriter and then Type Assert it to this to get the extra methods.
type HTTPRespWriter interface {
	http.ResponseWriter

	GetWrittenBody() []byte
	GetWrittenStatus() int
}
