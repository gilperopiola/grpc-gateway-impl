package errs

import (
	"fmt"
)

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
