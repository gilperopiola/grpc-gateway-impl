package errs

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// This checks for gorm/mongo error types, not our custom ones.
// Adds an unnecesary but insignificant overhead sometimes as it checks for both DB types.
func IsDBNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, mongo.ErrNoDocuments)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - Database Errors -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type DBErr struct {
	Err     error
	Context string
}

func (dberr DBErr) Error() string {
	return fmt.Sprintf("%s -> %v", dberr.Context, dberr.Err)
}

func (dberr DBErr) Unwrap() error {
	if dberr.Err != nil {
		return dberr.Err
	}
	return fmt.Errorf(dberr.Context)
}
