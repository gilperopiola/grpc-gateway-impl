package sql

import (
	"fmt"
	"strconv"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*     - Database Query Options -      */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// WithField returns a  core.QueryOpt ion which filters by the given field and value.
func WithField(field, value string) core.DBQueryOpt {
	return func(db core.SQLDatabaseAPI) {
		if field == "" {
			core.LogWeirdBehaviour(fmt.Sprintf("field is empty. field: %s", field))
			return
		}
		db.Where(fmt.Sprintf("%s = ?", field), value)
	}
}

// WithOr returns a  core.QueryOpt ion which filters by the given field and value using OR.
func WithOr(field, value string) core.DBQueryOpt {
	return func(db core.SQLDatabaseAPI) {
		if field == "" {
			core.LogWeirdBehaviour(fmt.Sprintf("field is empty. field: %s", field))
			return
		}
		db.Or(fmt.Sprintf("%s = ?", field), value)
	}
}

// WithFilter returns a  core.QueryOpt ion which fuzzy-matches the given field with the given filter.
func WithFilter(field, filter string) core.DBQueryOpt {
	return func(db core.SQLDatabaseAPI) {
		if field == "" || filter == "" {
			core.LogWeirdBehaviour(fmt.Sprintf("field or filter empty. field: %s, filter: %s", field, filter))
			return
		}
		db.Where(fmt.Sprintf("%s LIKE ?", field), "%"+filter+"%")
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Users Options -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func WithUserID(userID int) core.DBQueryOpt {
	return WithField("id", strconv.Itoa(userID))
}

func WithUsername(username string) core.DBQueryOpt {
	return WithField("username", username)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns a slice of  core.QueryOpt s containing the given  core.QueryOpt .
func Slice(opt core.DBQueryOpt) []core.DBQueryOpt {
	if opt == nil {
		return nil
	}
	return []core.DBQueryOpt{opt}
}
