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

type External interface {
	GetStorage() *storage.Storage
	GetClients() *clients.Clients

	GetDB() sql.DB
}

type externalLayer struct {
	storage.Storage
	clients.Clients
}

func NewExternalLayer(db sql.DB) External {
	return &externalLayer{
		Storage: storage.Storage{DB: db},
		Clients: clients.Clients{},
	}
}

func (e *externalLayer) GetStorage() *storage.Storage { return &e.Storage }
func (e *externalLayer) GetClients() *clients.Clients { return &e.Clients }
func (e *externalLayer) GetDB() sql.DB                { return e.Storage.DB }
