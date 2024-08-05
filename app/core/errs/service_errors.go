package errs

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Service Errors -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Our Service Errors must be GRPC 'Status' errors, as our Service speaks GRPC.
//
// But we also want to have our own custom errors,
// so we first make a ServiceErr and then .Error() it as the GRPC Status message.
type ServiceErr struct {
	Err    error
	Status codes.Code
	Info   string
}

// Returns an error with a GRPC Status code.
func NewGRPCError[ErrT any](code codes.Code, err ErrT, optionalInfo ...string) error {
	info := firstOrDefault(optionalInfo, code.String())
	return omitError(status.New(code, ServiceErr{err, code, info}.Unwrap().Error()).WithDetails()).Err()
}

func omitError(status *status.Status, err error) *status.Status {
	return status
}

func firstOrDefault(slice []string, fallback string) string {
	if len(slice) > 0 {
		return slice[0]
	}
	return fallback
}

func (serr ServiceErr) Error() string {
	if serr.Info == "" {
		return serr.Unwrap().Error()
	}
	return fmt.Sprintf("%s -> %v", serr.Info, serr.Unwrap())
}

func (serr ServiceErr) Unwrap() error {
	if serr.Err != nil {
		return serr.Err
	}
	return fmt.Errorf(serr.Status.String())
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func GRPCNotFound[T int | string](resource string, identif T) error {
	return NewGRPCError(codes.NotFound, fmt.Errorf("%v %s not found", identif, resource))
}

func GRPCUsersDBCall(err error, route string) error {
	return NewGRPCError(codes.Unknown, err, route, "users")
}
func GRPCGroupsDBCall(err error, route string) error {
	return NewGRPCError(codes.Unknown, err, route, "groups")
}
func GRPCGeneratingToken(err error) error {
	return NewGRPCError(codes.Unknown, err)
}
func GRPCUnauthenticated() error {
	return NewGRPCError(codes.Unauthenticated, nil)
}
func GRPCAlreadyExists(resource string) error {
	return NewGRPCError(codes.AlreadyExists, fmt.Errorf("%s already exists", resource))
}
func GRPCInternal(err error) error {
	return NewGRPCError(codes.Internal, err)
}
