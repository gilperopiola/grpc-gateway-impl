package errs

import (
	"fmt"
)

/* ----------------------------------- */
/*         - Database Errors -         */
/* ----------------------------------- */

type DBError struct {
	Err     error
	Context string
}

func (e DBError) Error() string {
	return fmt.Sprintf("%s -> %v", e.Context, e.Err)
}

func (e DBError) Unwrap() error {
	if e.Err != nil {
		return e.Err
	}
	return fmt.Errorf(e.Context)
}
