package tools

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/api_clients"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/db_tool/sqldb"
)

var _ core.Tools = (*Tools)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Tools -              */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

// 🛠️ Things that perform actions 🛠️
type Tools struct {
	core.ExternalAPIs     // -> API Clients.
	core.DBTool           // -> Storage (DB, Cache, etc).
	core.TLSTool          // -> Holds and retrieves data for TLS communication.
	core.CtxTool          // -> Manages context.
	core.FileManager      // -> Creates folders and files.
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

	tools.RateLimiter = NewRateLimiter(&cfg.RLimiterCfg)

	sqlDB := sqldb.NewSqlDB(&cfg.DBCfg)
	tools.DBTool = sqldb.NewDBTool(sqlDB)

	tools.ExternalAPIs = api_clients.NewAPIClients()

	tools.CtxTool = NewCtxTool()
	tools.FileManager = NewFileManager("etc/data/")

	tools.ModelConverter = NewModelConverter()
	tools.PwdHasher = NewPwdHasher(cfg.PwdHasherCfg.Salt)

	tools.RequestPaginator = NewRequestsPaginator(1, 10) // TODO - Config
	tools.RequestValidator = NewProtoRequestValidator()

	tools.TLSTool = NewTLSTool(&cfg.TLSCfg)

	tools.TokenGenerator = NewJWTGenerator(cfg.JWTCfg.Secret, cfg.JWTCfg.SessionDays)
	tools.TokenValidator = NewJWTValidator(tools.CtxTool, cfg.JWTCfg.Secret)

	tools.ShutdownJanitor = NewShutdownJanitor()

	return &tools
}
