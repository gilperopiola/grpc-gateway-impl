package options

import "gorm.io/gorm"

// QueryOption defines a function which takes a *gorm.DB and returns a *gorm.DB.
// We use it to apply different options to our database queries.
type QueryOption func(*gorm.DB) *gorm.DB

// WithField returns a QueryOption which filters by the given field and value.
func WithField(fieldName, fieldValue string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fieldName+" = ?", fieldValue)
	}
}

// WithOr returns a QueryOption which filters by the given field and value using OR.
func WithOr(fieldName, fieldValue string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Or(fieldName+" = ?", fieldValue)
	}
}

// WithLike returns a QueryOption which fuzzy-matches the given field with the given filter.
func WithLike(field, filter string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if filter == "" {
			return db
		}
		return db.Where(field+" LIKE ?", "%"+filter+"%")
	}
}
