package clients

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/clients/apis"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/db"
)

var _ core.Clients = &Clients{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - Clients -            */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

type Clients struct {
	core.APIClients // -> API Clients
	core.DB         // -> High-level DB Client
}

func Setup(cfg *core.Config, tools core.Tools) *Clients {
	clients := Clients{
		APIClients: apis.NewAPIs(&cfg.APIsCfg),
		DB:         db.NewSQLDBConn(&cfg.DBCfg, tools.HashPassword),
	}
	logs.InitModuleOK("Clients", "🔱")
	return &clients
}
