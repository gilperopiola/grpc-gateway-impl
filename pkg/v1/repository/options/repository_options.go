package options

import (
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/db"
)

/* ----------------------------------- */
/*    - Repository Query Options -     */
/* ----------------------------------- */

// QueryOption defines a function which takes a *gorm.DB and modifies it.
// We use it to apply different options to our database queries.
type QueryOption func(db.GormAdapter)

/* ----------------------------------- */
/*         - General Options -         */
/* ----------------------------------- */

// WithField returns a QueryOption which filters by the given field and value.
func WithField(fieldName, fieldValue string) QueryOption {
	return func(db db.GormAdapter) {
		db.Where(fieldName+" = ?", fieldValue)
	}
}

// WithOr returns a QueryOption which filters by the given field and value using OR.
func WithOr(fieldName, fieldValue string) QueryOption {
	return func(db db.GormAdapter) {
		db.Or(fieldName+" = ?", fieldValue)
	}
}

// WithFilter returns a QueryOption which fuzzy-matches the given field with the given filter.
func WithFilter(fieldName, filter string) QueryOption {
	return func(db db.GormAdapter) {
		if filter != "" {
			db.Where(fieldName+" LIKE ?", "%"+filter+"%")
		}
	}
}

/* ----------------------------------- */
/*          - Users Options -          */
/* ----------------------------------- */

func WithUsername(username string) QueryOption {
	return WithField("username", username)
}

func WithUserID(userID int) QueryOption {
	return WithField("id", fmt.Sprint(userID))
}
