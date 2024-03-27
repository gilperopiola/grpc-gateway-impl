package errs

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* ----------------------------------- */
/*          - Service Errors -         */
/* ----------------------------------- */

type ServiceError struct {
	Err  error
	Code codes.Code
	Msg  string
}

// NewGRPC returns a new ServiceError inside of a gRPC error.
func NewGRPC(err error, code codes.Code, messages ...string) error {
	msg := code.String()
	if len(messages) > 0 {
		msg = messages[0]
	}
	return status.Error(code, ServiceError{err, code, msg}.Error())
}

func (e ServiceError) Error() string {
	if e.Msg == "" {
		return e.Unwrap().Error()
	}
	return fmt.Sprintf("%s -> %v", e.Msg, e.Unwrap())
}

func (e ServiceError) Unwrap() error {
	if e.Err != nil {
		return e.Err
	}
	return fmt.Errorf(e.Code.String())
}
