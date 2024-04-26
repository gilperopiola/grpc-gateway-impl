package external

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external/clients"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external/storage/sql"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - External Layer -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Handles connections with external resources such as DBs, Files or other APIs.
type externalLayer struct {
	Storage core.StorageAPI
	Clients core.ClientsAPI
}

func SetupLayer(dbCfg *core.DBCfg) core.ExternalLayer {
	return &externalLayer{
		Storage: sql.SetupStorage(sql.NewGormDB(dbCfg)),
		Clients: clients.Setup(),
	}
}

func (e *externalLayer) GetStorage() core.StorageAPI {
	return e.Storage
}

func (e *externalLayer) GetClients() core.ClientsAPI {
	return e.Clients
}
