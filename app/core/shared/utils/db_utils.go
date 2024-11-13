package utils

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// Works with gorm and mongo errors.
// Adds an unnecesary but insignificant overhead sometimes as it checks for both DB types.
func IsNotFound(err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return true
	}
	return false
}
