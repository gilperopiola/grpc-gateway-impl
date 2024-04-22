package external

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external/clients"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external/storage"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external/storage/sql"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - External Layer -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Handles connections with external resources such as DBs, Files or other APIs.
type ExternalLayer struct {
	*storage.Storage
	*clients.Clients
}

func SetupLayer(dbCfg *core.DBCfg) *ExternalLayer {
	return &ExternalLayer{
		setupStorage(sql.NewGormDB(dbCfg)),
		setupClients(),
	}
}

func setupStorage(db core.SQLDatabaseAPI) *storage.Storage {
	return &storage.Storage{db}
}

func setupClients() *clients.Clients {
	return &clients.Clients{}
}
