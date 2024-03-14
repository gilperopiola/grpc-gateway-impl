package db

import "gorm.io/gorm"

type QueryOption func(*gorm.DB) *gorm.DB

func WithField(fieldName, fieldValue string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fieldName+" = ?", fieldValue)
	}
}

func WithOr(fieldName, fieldValue string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Or(fieldName+" = ?", fieldValue)
	}
}

func WithFilter(field, filter string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if filter == "" {
			return db
		}
		return db.Where(field+" LIKE ?", "%"+filter+"%")
	}
}
