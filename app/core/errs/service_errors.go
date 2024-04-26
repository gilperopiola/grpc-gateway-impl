package errs

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Service Errors -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var (
	ErrSvcUserRelated = func(err error, methodName string) error {
		return NewGRPC(codes.Unknown, err, methodName)
	}
	ErrSvcOnTokenGeneration = func(err error) error {
		return NewGRPC(codes.Unknown, err)
	}
	ErrSvcUnauthenticated = func() error {
		return NewGRPC(codes.Unauthenticated, nil)
	}
	ErrSvcNotFound = func(entity string) error {
		err := fmt.Errorf("%s not found", entity)
		return NewGRPC(codes.NotFound, err)
	}
	ErrSvcAlreadyExists = func(entity string) error {
		err := fmt.Errorf("%s already exists", entity)
		return NewGRPC(codes.AlreadyExists, err)
	}
)

// NewGRPC returns a new ServiceError inside of a GRPC error.
func NewGRPC(code codes.Code, err error, messages ...string) error {
	msg := code.String()
	if len(messages) > 0 {
		msg = messages[0]
	}
	return status.Error(code, ServiceError{err, code, msg}.Error())
}

type ServiceError struct {
	Err  error
	Code codes.Code
	Msg  string
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
