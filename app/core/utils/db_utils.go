package utils

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// Works with gorm and mongo errors.
// Adds a tiny slight insignificant overhead.
func IsNotFound(err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return true
	}
	return false
}
