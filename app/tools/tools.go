package tools

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/logs"
)

var _ core.Tools = &Tools{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - Tools -              */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

// 🛠️ - Holds an instance of every single Tool on our app.
//
// A Tool lets us perform clear, explicit actions, to be used
// across our Service.
//
// Well defined Tools make the business logic easier to read and understand.
type Tools struct {
	core.TLSManager       // -> Holds and retrieves data for TLS communication.
	core.ContextManager   // -> Manages context.
	core.FileManager      // -> Creates folders and files.
	core.ImageLoader      // -> Loads images from different sources.
	core.ModelConverter   // -> Converts between models and PBs.
	core.PwdHasher        // -> Hashes and compares passwords.
	core.RateLimiter      // -> Limits rate of requests.
	core.RequestPaginator // -> Helps handling GRPC requests with pagination.
	core.RequestValidator // -> Validates GRPC requests.
	core.ShutdownJanitor  // -> Cleans up and frees resources on application shutdown.
	core.TokenGenerator   // -> Generates JWT Tokens.
	core.TokenValidator   // -> Validates JWT Tokens.
}

func Setup(cfg *core.Config) *Tools {
	tools := Tools{}

	// Security
	tools.TLSManager = NewTLSTool(&cfg.TLSCfg)

	// Request handling
	tools.ContextManager = NewCtxTool()
	tools.RequestPaginator = NewRequestsPaginator(1, 10) // TODO - Config
	tools.RequestValidator = NewProtoRequestValidator()

	// Auth -> JWT Tokens
	tools.TokenGenerator = NewJWTGenerator(cfg.JWTCfg.Secret, cfg.JWTCfg.SessionDays)
	tools.TokenValidator = NewJWTValidator(tools.ContextManager, cfg.JWTCfg.Secret, "TODOimproveapikey")

	// Other utilities
	tools.FileManager = NewFileManager("etc/data/")
	tools.ImageLoader = NewImageLoader()
	tools.PwdHasher = NewPwdHasher(cfg.PwdHasherCfg.Salt)
	tools.RateLimiter = NewRateLimiter(&cfg.RLimiterCfg)
	tools.ModelConverter = NewModelConverter()
	tools.ShutdownJanitor = NewShutdownJanitor()

	logs.InitModuleOK("Tools", "🛠️ ")
	return &tools
}
