package app

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/api_clients"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/db_tool/sqldb"
)

var _ core.Toolbox = (*Tools)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - App Tools -            */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

// 🛠️ Things that perform actions 🛠️
type Tools struct {
	core.APIs               // -> API Clients.
	core.DBTool             // -> Storage (DB, Cache, etc).
	core.FileManager        // -> Creates folders and files.
	core.MetadataGetter     // -> Gets metadata from GRPC requests.
	core.PwdHasher          // -> Hashes and compares passwords.
	core.RateLimiter        // -> Limits rate of requests.
	core.Retrier            // -> Executes a fn and retries if it fails.
	core.RequestsValidator  // -> Validates GRPC requests.
	core.ShutdownJanitor    // -> Cleans up and frees resources on application shutdown.
	core.RouteAuthenticator // -> Allows or denies user access per route.
	core.TLSTool            // -> Holds and retrieves data for TLS communication.
	core.TokenGenerator     // -> Generates JWT Tokens.
	core.TokenValidator     // -> Validates JWT Tokens.
}

func (app *App) SetupTools() {
	cfg := app.Config

	app.Tools.APIs = api_clients.NewAPIClients()

	app.Tools.Retrier = tools.NewRetrier(&cfg.RetrierCfg)

	sqlDB := sqldb.NewSQLDB(&cfg.DBCfg, app.Tools.Retrier)
	app.Tools.DBTool = sqldb.NewDBTool(sqlDB)

	app.Tools.FileManager = tools.NewFileManager("data/")

	app.Tools.MetadataGetter = tools.NewMetadataGetter()

	app.Tools.PwdHasher = tools.NewPwdHasher(cfg.PwdHasherCfg.Salt)

	app.Tools.RateLimiter = tools.NewRateLimiter(&cfg.RLimiterCfg)

	app.Tools.RequestsValidator = tools.NewRequestsValidator()

	app.Tools.ShutdownJanitor = tools.NewShutdownJanitor()

	app.Tools.RouteAuthenticator = tools.NewRouteAuthenticator(core.AuthForRoute)

	app.Tools.TLSTool = tools.NewTLSTool(&cfg.TLSCfg)

	app.Tools.TokenGenerator = tools.NewJWTGenerator(cfg.JWTCfg.Secret, cfg.JWTCfg.SessionDays)

	app.Tools.TokenValidator = tools.NewJWTValidator(app.MetadataGetter, app.RouteAuthenticator, app.JWTCfg.Secret)
}
