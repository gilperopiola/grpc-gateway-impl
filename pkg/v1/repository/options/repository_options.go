package options

import (
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/db"
)

/* ----------------------------------- */
/*    - Repository Query Options -     */
/* ----------------------------------- */

// QueryOpt defines a function which takes a *gorm.DB and modifies it.
// We use it to apply different options to our database queries.
type QueryOpt func(db.GormAdapter)

// Slice returns a slice of QueryOpts.
func Slice(opt QueryOpt) []QueryOpt {
	if opt == nil {
		return nil
	}
	return []QueryOpt{opt}
}

/* ----------------------------------- */
/*         - General Options -         */
/* ----------------------------------- */

// WithField returns a QueryOption which filters by the given field and value.
func WithField(fieldName, fieldValue string) QueryOpt {
	return func(db db.GormAdapter) {
		db.Where(fmt.Sprintf("%s = ?", fieldName), fieldValue)
	}
}

// WithOr returns a QueryOption which filters by the given field and value using OR.
func WithOr(fieldName, fieldValue string) QueryOpt {
	return func(db db.GormAdapter) {
		db.Or(fmt.Sprintf("%s = ?", fieldName), fieldValue)
	}
}

// WithFilter returns a QueryOption which fuzzy-matches the given field with the given filter.
func WithFilter(fieldName, filter string) QueryOpt {
	return func(db db.GormAdapter) {
		if filter != "" {
			db.Where(fmt.Sprintf("%s LIKE ?", fieldName), "%"+filter+"%")
		}
	}
}

/* ----------------------------------- */
/*          - Users Options -          */
/* ----------------------------------- */

func WithUsername(username string) QueryOpt {
	return WithField("username", username)
}

func WithUserID(userID int) QueryOpt {
	return WithField("id", fmt.Sprint(userID))
}
