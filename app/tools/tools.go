package tools

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/api_clients"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/db_tool/sqldb"
)

var _ core.Tools = (*Tools)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Tools -              */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

// 🛠️ Things that perform actions 🛠️
type Tools struct {
	core.ExternalAPIs      // -> API Clients.
	core.DBTool            // -> Storage (DB, Cache, etc).
	core.TLSTool           // -> Holds and retrieves data for TLS communication.
	core.CtxTool           // -> Manages context.
	core.FileManager       // -> Creates folders and files.
	core.HealthChecker     // -> Checks health of own service.
	core.ModelConverter    // -> Converts between models and PBs.
	core.PwdHasher         // -> Hashes and compares passwords.
	core.RateLimiter       // -> Limits rate of requests.
	core.RequestsPaginator // -> Helps handling GRPC requests with pagination.
	core.RequestsValidator // -> Validates GRPC requests.
	core.ShutdownJanitor   // -> Cleans up and frees resources on application shutdown.
	core.TokenGenerator    // -> Generates JWT Tokens.
	core.TokenValidator    // -> Validates JWT Tokens.
}

func Setup(cfg *core.Config, serviceFn ServiceFunc) *Tools {
	tools := Tools{}

	tools.RateLimiter = NewRateLimiter(&cfg.RLimiterCfg)

	sqlDB := sqldb.NewSqlDB(&cfg.DBCfg)
	tools.DBTool = sqldb.NewDBTool(sqlDB)

	tools.ExternalAPIs = api_clients.NewAPIClients()

	tools.CtxTool = NewCtxTool()
	tools.FileManager = NewFileManager("etc/data/")
	tools.HealthChecker = NewHealthChecker(serviceFn)

	tools.ModelConverter = NewModelConverter()
	tools.PwdHasher = NewPwdHasher(cfg.PwdHasherCfg.Salt)

	tools.RequestsPaginator = NewRequestsPaginator(1, 10) // TODO - Config
	tools.RequestsValidator = NewRequestsValidator()

	tools.TLSTool = NewTLSTool(&cfg.TLSCfg)

	tools.TokenGenerator = NewJWTGenerator(cfg.JWTCfg.Secret, cfg.JWTCfg.SessionDays)
	tools.TokenValidator = NewJWTValidator(tools.CtxTool, cfg.JWTCfg.Secret)

	tools.ShutdownJanitor = NewShutdownJanitor()

	return &tools
}

type ServiceFunc func(context.Context, *pbs.AnswerGroupInviteRequest) (*pbs.AnswerGroupInviteResponse, error)
