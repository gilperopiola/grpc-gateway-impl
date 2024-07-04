package toolbox

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/toolbox/api_clients"
	"github.com/gilperopiola/grpc-gateway-impl/app/toolbox/db_tool/sqldb"
)

var _ core.Toolbox = (*Toolbox)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Toolbox -              */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

// 🛠️ Things that perform actions 🛠️
type Toolbox struct {
	core.APIs              // -> API Clients.
	core.DBTool            // -> Storage (DB, Cache, etc).
	core.TLSTool           // -> Holds and retrieves data for TLS communication.
	core.CtxManager        // -> Manages context.
	core.FileManager       // -> Creates folders and files.
	core.HealthChecker     // -> Checks health of own service.
	core.ModelConverter    // -> Converts between models and PBs.
	core.PwdHasher         // -> Hashes and compares passwords.
	core.RateLimiter       // -> Limits rate of requests.
	core.Retrier           // -> Executes a fn and retries if it fails.
	core.RequestsValidator // -> Validates GRPC requests.
	core.ShutdownJanitor   // -> Cleans up and frees resources on application shutdown.
	core.TokenGenerator    // -> Generates JWT Tokens.
	core.TokenValidator    // -> Validates JWT Tokens.
}

func Setup(cfg *core.Config, serviceFn ServiceFn) *Toolbox {
	toolbox := Toolbox{}

	toolbox.Retrier = NewRetrier(&cfg.RetrierCfg)
	toolbox.RateLimiter = NewRateLimiter(&cfg.RLimiterCfg)

	sqlDB := sqldb.NewSQLDB(&cfg.DBCfg, toolbox.Retrier)
	toolbox.DBTool = sqldb.NewDBTool(sqlDB)

	toolbox.APIs = api_clients.NewAPIClients()

	toolbox.CtxManager = NewCtxManager()
	toolbox.FileManager = NewFileManager("etc/data/")
	toolbox.HealthChecker = NewHealthChecker(serviceFn)

	toolbox.ModelConverter = NewModelConverter()
	toolbox.PwdHasher = NewPwdHasher(cfg.PwdHasherCfg.Salt)
	toolbox.RequestsValidator = NewRequestsValidator()

	toolbox.TLSTool = NewTLSTool(&cfg.TLSCfg)

	toolbox.TokenGenerator = NewJWTGenerator(cfg.JWTCfg.Secret, cfg.JWTCfg.SessionDays)
	toolbox.TokenValidator = NewJWTValidator(toolbox.CtxManager, cfg.JWTCfg.Secret)

	toolbox.ShutdownJanitor = NewShutdownJanitor()

	return &toolbox
}
