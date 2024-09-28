package clients

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/clients/apis"
	"github.com/gilperopiola/grpc-gateway-impl/app/clients/dbs/sqldb"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

var _ core.Clients = &Clients{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - Clients -            */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

type Clients struct {
	core.APIs // -> API Clients
	core.DB   // -> High-level DB operations
}

func Setup(cfg *core.Config) *Clients {
	clients := Clients{}

	// APIs
	clients.APIs = apis.NewAPIs(cfg.APIsCfg)

	// DBs
	clients.DB = sqldb.NewSQLDBConnection(&cfg.DBCfg)

	return &clients
}
