package storage

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*     - External Layer: Storage -     */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var _ core.StorageAPI = (*Storage)(nil)

type Storage struct {
	DB core.SQLDatabaseAPI
}

func Setup(db core.SQLDatabaseAPI) *Storage {
	return &Storage{db}
}

// T0D0 this can probably be done in a way that ExternalLayer can hold many Storage(s),
// each of those being the same Storage struct but with a generic type that would be sql or MongoDB or so. DRY overload.
