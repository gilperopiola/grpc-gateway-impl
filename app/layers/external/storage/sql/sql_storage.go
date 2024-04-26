package sql

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

var _ core.StorageAPI = (*SQLStorage)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*     - External Layer: Storage -     */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type SQLStorage struct {
	DB core.SQLDatabaseAPI
}

func SetupStorage(db core.SQLDatabaseAPI) *SQLStorage {
	return &SQLStorage{db}
}

func (s *SQLStorage) CloseDB() {
	s.DB.Close()
}

// T0D0 this can probably be done in a way that ExternalLayer can hold many Storage(s),
// each of those being the same Storage struct but with a generic type that would be sql or MongoDB or so. DRY overload.
