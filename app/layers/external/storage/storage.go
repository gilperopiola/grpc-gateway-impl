package storage

import "github.com/gilperopiola/grpc-gateway-impl/app/layers/external/storage/sql"

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*     - External Layer: Storage -     */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// T0D0 this can probably be done in a way that ExternalLayer can hold many Storage(s),
// each of those being the same Storage struct but with a generic type that would be sql or MongoDB or so. DRY overload.
type Storage struct {
	DB sql.DB
}

func NewStorage(db sql.DB) *Storage {
	return &Storage{db}
}
