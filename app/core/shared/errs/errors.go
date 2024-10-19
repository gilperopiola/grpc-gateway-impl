package errs

import (
	"fmt"
)

// Mimics errors.New, but allows for any type of message.
func New(msg any) error {
	return fmt.Errorf("%v", msg)
}
