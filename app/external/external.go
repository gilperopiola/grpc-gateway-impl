package external

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core/special_types"
	"github.com/gilperopiola/grpc-gateway-impl/app/external/clients"
	"github.com/gilperopiola/grpc-gateway-impl/app/external/storage"
	"github.com/gilperopiola/grpc-gateway-impl/app/external/storage/sqldb"
)

/* ----------------------------------- */
/*         - External Layer -          */
/* ----------------------------------- */

// The External Layer handles connections to external resources such as databases or other APIs.

type ExternalLayer interface {
	GetStorage() *storage.Storage
	GetDB() special_types.SQLDB
	GetClients() *clients.Clients
}

type externalLayer struct {
	storage.Storage
	clients.Clients
}

func NewExternalLayer(db sqldb.Database) ExternalLayer {
	return &externalLayer{
		Storage: storage.Storage{DB: db},
		Clients: clients.Clients{},
	}
}

func (e *externalLayer) GetStorage() *storage.Storage {
	return &e.Storage
}

func (e *externalLayer) GetDB() special_types.SQLDB {
	return e.Storage.DB
}

func (e *externalLayer) GetClients() *clients.Clients {
	return &e.Clients
}
