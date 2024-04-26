package tools

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - Requests Validator -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var _ core.RequestsValidator = (*protoValidator)(nil)

// Validates GRPC requests based on rules written on .proto files.
// It uses the bufbuild/protovalidate library.
type protoValidator struct {
	*protovalidate.Validator
}

// New instance of *protoValidator. This panics on failure.
func NewProtoValidator() *protoValidator {
	validator, err := protovalidate.New()
	if err != nil {
		core.LogUnexpectedAndPanic(fmt.Errorf(errs.FailedToCreateProtoVal, err))
	}
	return &protoValidator{validator}
}

// Wraps the proto validation logic with a GRPC Interceptor.
func (v *protoValidator) ValidateGRPC(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if err := v.Validate(req.(protoreflect.ProtoMessage)); err != nil {
		return nil, validationErrToGRPC(err)
	}
	return handler(ctx, req)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Takes a *protovalidate.ValidationError and returns an InvalidArgument(3) GRPC error with its corresponding message.
// Validation errors are always returned as InvalidArgument.
// They get translated to 400 Bad Request on the HTTP error handler.
func validationErrToGRPC(err error) error {
	var validationErr *protovalidate.ValidationError
	if ok := errors.As(err, &validationErr); ok {
		violations := validationErr.ToProto().GetViolations()
		return status.Error(codes.InvalidArgument, fmt.Sprintf(errs.InReqValidation, parseViolations(violations)))
	}

	var runtimeErr *protovalidate.RuntimeError
	if ok := errors.As(err, &runtimeErr); ok {
		core.LogUnexpectedErr(runtimeErr)
		return status.Error(codes.InvalidArgument, fmt.Sprintf(errs.InReqValidationRuntime, runtimeErr))
	}

	core.LogUnexpectedErr(err)
	return status.Error(codes.InvalidArgument, fmt.Sprintf(errs.InReqValidationUnexpected, err))
}

// Returns a string detailing the validations violations.
// This string is the human-facing format in which validation errors translate.
func parseViolations(violations []*validate.Violation) string {
	out := ""
	for i, v := range violations {
		out += parseViolation(v)
		if i < len(violations)-1 {
			out += ", "
		}
	}
	return out
}

// For each broken rule in the validation process return a string explaining the cause.
// In human-facing format.
func parseViolation(v *validate.Violation) string {
	if strings.Contains(v.Message, "match regex pattern") {
		return fmt.Sprintf("%s value has an invalid format", v.FieldPath) // don't show regex pattern.
	}
	return fmt.Sprintf("%s %s", v.FieldPath, v.Message) // simple message.
}
