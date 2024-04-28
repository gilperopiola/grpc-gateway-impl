package errs

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Service Errors -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// -> Service Errors have to be GRPC 'status' errors, as our Service speaks GRPC.
// -> But we want to have our own custom errors, so we first make a ServiceErr and then .Error() it as the GRPC Status message.

type ServiceErr struct {
	Err      error
	Code     codes.Code
	Metadata string // Holds relevant info like the route or sth.
}

// Returns a new ServiceError inside of a GRPC error 'status'.
func NewGRPCServiceErr(code codes.Code, err error, optionalMD ...string) error {
	md := code.String()
	if len(optionalMD) > 0 {
		md = optionalMD[0]
	}
	return status.Error(code, ServiceErr{err, code, md}.Unwrap().Error())
}

func (serr ServiceErr) Error() string {
	if serr.Metadata == "" {
		return serr.Unwrap().Error()
	}
	return fmt.Sprintf("%s -> %v", serr.Metadata, serr.Unwrap())
}

func (serr ServiceErr) Unwrap() error {
	if serr.Err != nil {
		return serr.Err
	}
	return fmt.Errorf(serr.Code.String())
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var (
	GRPCUsersDBCall = func(err error, route string) error {
		return NewGRPCServiceErr(codes.Unknown, err, route)
	}
	GRPCGeneratingToken = func(err error) error {
		return NewGRPCServiceErr(codes.Unknown, err)
	}
	GRPCUnauthenticated = func() error {
		return NewGRPCServiceErr(codes.Unauthenticated, nil)
	}
	GRPCNotFound = func(resource string) error {
		return NewGRPCServiceErr(codes.NotFound, fmt.Errorf("%s not found", resource))
	}
	GRPCAlreadyExists = func(resource string) error {
		return NewGRPCServiceErr(codes.AlreadyExists, fmt.Errorf("%s already exists", resource))
	}
)
