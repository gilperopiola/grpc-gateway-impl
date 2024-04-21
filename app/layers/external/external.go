package external

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external/clients"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external/storage"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external/storage/sql"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - External Layer -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// The External Layer handles connections to external resources such as databases or other APIs.

type ExternalLayer struct {
	*storage.StorageLayer
	*clients.Clients
}

func SetupLayer(db sql.DB) *ExternalLayer {
	return &ExternalLayer{
		&storage.StorageLayer{db},
		&clients.Clients{},
	}
}
