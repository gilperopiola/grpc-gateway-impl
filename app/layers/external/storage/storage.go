package storage

import "github.com/gilperopiola/grpc-gateway-impl/app/layers/external/storage/sql"

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*     - External Layer: Storage -     */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// T0D0 this can probably be done in a way that ExternalLayer can hold many StorageLayer(s),
// each of those being the same StorageLayer struct but with a generic type that would be sql or MongoDB or so. DRY overload.
type StorageLayer struct {
	DB sql.DB
}

func NewStorageLayer(db sql.DB) *StorageLayer {
	return &StorageLayer{db}
}
