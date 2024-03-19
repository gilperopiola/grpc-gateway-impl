package common

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/bufbuild/protovalidate-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* ----------------------------------- */
/*         - Input Validator -         */
/* ----------------------------------- */

// InputValidator is the interface that wraps the ValidateInput method.
// It is used to validate incoming gRPC requests. The rules are defined in the .proto files.
type InputValidator interface {
	ValidateInput(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

// protoValidator is a wrapper around the protovalidate.Validator.
type protoValidator struct {
	*protovalidate.Validator
}

// NewInputValidator returns a new instance of *protoValidator.
// It panics on failure.
func NewInputValidator() *protoValidator {
	validator, err := protovalidate.New()
	if err != nil {
		zap.S().Fatalf(errs.FatalErrMsgCreatingValidator, err)
	}
	return &protoValidator{Validator: validator}
}

// Validate returns an interceptor that validates incoming gRPC requests.
func (v *protoValidator) ValidateInput(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := v.Validate(req.(protoreflect.ProtoMessage)); err != nil {
		return nil, handleValidationErr(err)
	}
	return handler(ctx, req) // Next handler.
}

// handleValidationErr takes a ValidationError and returns an InvalidArgument(3) gRPC error with its corresponding message.
// Validation errors are always returned as InvalidArgument.
// They get translated to 400 Bad Request on the HTTP error handler.
func handleValidationErr(err error) error {
	var validationErr *protovalidate.ValidationError
	if ok := errors.As(err, &validationErr); ok {
		violations := validationErr.ToProto().GetViolations()
		return status.Error(codes.InvalidArgument, fmt.Sprintf(errs.ErrMsgInValidation, parseViolations(violations)))
	}

	var runtimeErr *protovalidate.RuntimeError
	if ok := errors.As(err, &runtimeErr); ok {
		return status.Error(codes.InvalidArgument, fmt.Sprintf(errs.ErrMsgInValidationRuntime, runtimeErr))
	}

	return status.Error(codes.InvalidArgument, fmt.Sprintf(errs.ErrMsgInValidationUnexpected, err))
}

// parseViolations returns a string detailing the validations violations.
// This is the human-facing format in which validation errors translate.
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

// parseViolation returns a string detailing the validation violations.
func parseViolation(v *validate.Violation) string {
	if strings.Contains(v.Message, "match regex pattern") {
		return fmt.Sprintf("%s value has an invalid format", v.FieldPath)
	}
	return fmt.Sprintf("%s %s", v.FieldPath, v.Message)
}
