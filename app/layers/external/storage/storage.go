package storage

import "github.com/gilperopiola/grpc-gateway-impl/app/core/special_types"

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*     - External Layer: Storage -     */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// T0D0 this can probably be done in a way that ExternalLayer can hold many Storage(s),
// each of those being the same Storage struct but with a generic type that would be SQLDB or MongoDB or so. DRY overload.
type Storage struct {
	DB special_types.SQLDB
}

func NewStorage(db special_types.SQLDB) *Storage {
	return &Storage{db}
}
