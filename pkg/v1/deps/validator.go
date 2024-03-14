package deps

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* ----------------------------------- */
/*         - Proto Validator -         */
/* ----------------------------------- */

type Validator struct {
	*protovalidate.Validator
}

func (v *Validator) Validate() grpc.UnaryServerInterceptor {
	return validateRequest(v.Validator)
}

// NewValidator returns a new instance of *protovalidate.Validator.
// It panics on failure.
func NewValidator() *Validator {
	validator, err := protovalidate.New()
	if err != nil {
		log.Fatalf(errs.FatalErrMsgCreatingValidator, err)
	}
	return &Validator{Validator: validator}
}

func validateRequest(validator *protovalidate.Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := validator.Validate(req.(protoreflect.ProtoMessage)); err != nil {
			return nil, handleValidationErr(err)
		}
		return handler(ctx, req) // Next handler.
	}
}

// handleValidationErr takes a ValidationError and returns an InvalidArgument(3) gRPC error with its corresponding message.
// Validation errors are always returned as InvalidArgument.
// They get translated to 400 Bad Request on the HTTP error handler.
func handleValidationErr(err error) error {
	errMsg := fmt.Sprintf(errs.ErrMsgInValidationUnexpected, err)

	var validationErr *protovalidate.ValidationError
	if ok := errors.As(err, &validationErr); ok {
		violations := validationErr.ToProto().GetViolations()
		errMsg = fmt.Sprintf(errs.ErrMsgInValidation, parseViolations(violations))
	}

	var runtimeErr *protovalidate.RuntimeError
	if ok := errors.As(err, &runtimeErr); ok {
		errMsg = fmt.Sprintf(errs.ErrMsgInValidationRuntime, runtimeErr)
	}

	return status.Error(codes.InvalidArgument, errMsg)
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
