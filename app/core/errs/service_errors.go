package errs

import (
	"errors"
	"fmt"

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
	errorMsg := firstOrDefault(serr.Info, serr.Status.String())
	return fmt.Sprintf("%v: %s", errorMsg, serr.Unwrap())
}

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

func GRPCAlreadyExists(what string) error {
	return NewGRPCError(codes.AlreadyExists, errors.New(what+" already exists"))
}

// We return this on invalid password on Login - Auth Service.
// This can also mean a wrong username, but not a non-existing one.
func GRPCWrongLoginInfo() error {
	return NewGRPCError(codes.Unauthenticated, errors.New("wrong login information"))
}

// We also return this from the Login, but after checking the password.
// Don't really know what could cause this, but it may happen.
func GRPCGeneratingToken(err error) error {
	return NewGRPCError(codes.Unknown, err)
}

// We return this on unexpected DB errors from the Users Service.
func GRPCUsersDBCall(err error, route string) error {
	return NewGRPCError(codes.Internal, err, route)
}

// We return this on unexpected DB errors from the Groups Service.
func GRPCGroupsDBCall(err error, route string) error {
	return NewGRPCError(codes.Internal, err, route)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func firstOrDefault(slice []string, fallback string) string {
	if len(slice) > 0 {
		return slice[0]
	}
	return fallback
}
