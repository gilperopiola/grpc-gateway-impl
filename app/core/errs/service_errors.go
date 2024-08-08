package errs

import (
	"errors"
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Returns a GRPC Status error with our custom ServiceErr inside.
func NewGRPCError(code codes.Code, err error, optionalInfo ...string) error {
	serviceErr := ServiceErr{err, code, optionalInfo}
	grpcStatus := status.New(code, serviceErr.Error())
	return grpcStatus.Err()
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - GRPC Service Errors -       */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Our Service Errors must be GRPC 'Status' errors, as our Service speaks GRPC.
type ServiceErr struct {
	Err    error
	Status codes.Code
	Info   []string
}

// Returns our ServiceErr as a string:
//
//	Example: "This is some additional information: actual error message"
func (serr ServiceErr) Error() string {
	errorMsg := utils.FirstOrDefault(serr.Info, serr.Status.String())
	return fmt.Sprintf("%v: %s", errorMsg, serr.Unwrap())
}

// Returns the inner error, or the status code as an error.
func (serr ServiceErr) Unwrap() error {
	if serr.Err != nil {
		return serr.Err
	}
	return fmt.Errorf(serr.Status.String())
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func GRPCNotFound[T int | string](resource string, identif T) error {
	return NewGRPCError(codes.NotFound, fmt.Errorf("%s %v not found", resource, identif))
}

// Translates to HTTP 409 Conflict Error.
func GRPCAlreadyExists(what string) error {
	return NewGRPCError(codes.AlreadyExists, errors.New(what+" already exists"))
}

// We return this on username or password mismatch on the Auth Service's Login.
func GRPCWrongLoginInfo() error {
	return NewGRPCError(codes.Unauthenticated, errors.New("wrong username or password"))
}

// We also return this from the Login, but after succesfully matching the credentials.
// Don't really know what could cause this, but the Login is kind of important so
// better be covered.
func GRPCGeneratingToken(err error) error {
	return NewGRPCError(codes.Unknown, err)
}

// We return this on unexpected errors coming from the DB Layer.
func GRPCFromDB(err error, route string) error {
	return NewGRPCError(codes.Internal, err, route)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
