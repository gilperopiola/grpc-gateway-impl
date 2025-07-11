package core

// ➤ Am I sure this is the best way to go about stitching together the whole
// project here as interfaces? — Not sure, however this allowed me to decouple
// the tools/service/servers/clients packages, each of one importing core instead.
//
// ➤ If you don't wanna import core, just copy n paste the interfaces.

import (
	"context"
	"net/http"
	"time"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/apimodels"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/google/uuid"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

/* -~-~-~-~- Main Interfaces -~-~-~-~- */

// Each service is defined in a .proto file
type ServiceLayer interface {
	RegisterGRPCEndpoints(grpcServer grpc.ServiceRegistrar)
	RegisterHTTPEndpoints(mux any, opts ...grpc.DialOption)
}

/* -~-~-~-~- Clients -~-~-~-~- */

type (
	Clients interface {
		// Database access
		GetDB() any
		CloseDB() error

		// Repositories
		UserRepository() UserRepository
		GroupRepository() GroupRepository
		GPTChatRepository() GPTChatRepository

		// API clients
		APIClients
	}

	MongoDB interface {
		MongoRepository
		GetDB() any // *gorm.DB or *mongo.Client
		CloseDB()
	}

	APIClients interface {
		GPTAPI
		WeatherAPI
	}
)

type (
	MongoRepository interface { // High-level DB operations
		MongoUsersRepository
		MongoGroupsRepository
		MongoGPTRepository
	}
	MongoUsersRepository interface {
		DBCreateUser(ctx god.Ctx, username, hashedPwd string) (*models.User, error)
		DBGetUser(ctx god.Ctx, opts ...any) (*models.User, error)
		DBGetUsers(ctx god.Ctx, page, pageSize int, opts ...any) ([]*models.User, int, error)
	}
	MongoGroupsRepository interface {
		DBCreateGroup(ctx god.Ctx, name string, ownerID int, invitedUserIDs []int) (*models.Group, error)
		DBGetGroup(ctx god.Ctx, opts ...any) (*models.Group, error)
	}
	MongoGPTRepository interface {
		DBGetGPTChat(ctx god.Ctx, opts ...any) (*models.GPTChat, error)
		DBCreateGPTChat(ctx god.Ctx, title string) (*models.GPTChat, error)
		DBCreateGPTMessage(ctx god.Ctx, message *models.GPTMessage) (*models.GPTMessage, error)
	}
)

type (
	GPTAPI interface {
		SendRequestToGPT(ctx context.Context, prompt string, prevMsgs ...apimodels.GPTChatMsg) (string, error)
		SendRequestToDallE(ctx context.Context, prompt string, size pbs.GPTImageSize) (apimodels.GPTImageMsg, error)
	}
	WeatherAPI interface {
		GetCurrentWeather(ctx god.Ctx, lat, lon float64) (*apimodels.GetWeatherResponse, error)
	}
)

/* -~-~-~-~- Tools -~-~-~-~- */

// These are the interfaces to all of our tools.
// Each concrete tool lives on the tools pkg.
type Tools interface {
	TokenGenerator
	TokenValidator
	RequestPaginator
	RequestValidator
	ShutdownJanitor
	RateLimiter
	PwdHasher
	TLSManager
	FileManager
	ContextManager
	ModelConverter
	ImageLoader
	IDGenerator[string]
	FileDownloader
}

/* -~-~-~-~- Other -~-~-~-~- */

type IDType interface {
	string | int | int64 | uuid.UUID
}

// This isn't used like the other tools, as it's instantiated per request.
// The implementation lives on the tools pkg.
// Used to wrap the http.ResponseWriter and then Type Assert it to this to get the extra methods.
type HTTPRespWriter interface {
	http.ResponseWriter

	GetWrittenBody() []byte
	GetWrittenStatus() int
}

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
