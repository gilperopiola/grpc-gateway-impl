package toolbox

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/toolbox/api_clients"
	"github.com/gilperopiola/grpc-gateway-impl/app/toolbox/db_tool/sqldb"
)

var _ core.Toolbox = (*Toolbox)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - Toolbox -              */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

// 🛠️ Things that perform actions 🛠️
type Toolbox struct {
	core.APIs              // -> API Clients.
	core.DBTool            // -> Storage (DB, Cache, etc).
	core.TLSTool           // -> Holds and retrieves data for TLS communication.
	core.CtxManager        // -> Manages context.
	core.FileManager       // -> Creates folders and files.
	core.PwdHasher         // -> Hashes and compares passwords.
	core.RateLimiter       // -> Limits rate of requests.
	core.Retrier           // -> Executes a fn and retries if it fails.
	core.RequestsValidator // -> Validates GRPC requests.
	core.ShutdownJanitor   // -> Cleans up and frees resources on application shutdown.
	core.TokenGenerator    // -> Generates JWT Tokens.
	core.TokenValidator    // -> Validates JWT Tokens.
}

func Setup(cfg *core.Config) core.Toolbox {
	t := &Toolbox{}

	t.Retrier = NewRetrier(&cfg.RetrierCfg)
	t.RateLimiter = NewRateLimiter(&cfg.RLimiterCfg)

	sqlDB := sqldb.NewSQLDB(&cfg.DBCfg, t.Retrier)
	t.DBTool = sqldb.NewDBTool(sqlDB)

	t.APIs = api_clients.NewAPIClients()

	t.CtxManager = NewCtxManager()
	t.FileManager = NewFileManager("etc/data/")

	t.PwdHasher = NewPwdHasher(cfg.PwdHasherCfg.Salt)
	t.RequestsValidator = NewRequestsValidator()

	t.TLSTool = NewTLSTool(&cfg.TLSCfg)

	t.TokenGenerator = NewJWTGenerator(cfg.JWTCfg.Secret, cfg.JWTCfg.SessionDays)
	t.TokenValidator = NewJWTValidator(t.CtxManager, cfg.JWTCfg.Secret)

	t.ShutdownJanitor = NewShutdownJanitor()

	return t
}
