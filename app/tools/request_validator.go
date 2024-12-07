package tools

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var _ core.RequestValidator = &protoRequestValidator{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - Request Validator -       */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Validates GRPC requests based on rules written on .proto files.
// It uses the bufbuild/protovalidate library.
type protoRequestValidator struct {
	*protovalidate.Validator
}

func NewProtoRequestValidator() core.RequestValidator {
	validator, err := protovalidate.New()
	logs.LogFatalIfErr(err, errs.FailedToCreateProtoVal)
	return &protoRequestValidator{validator}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns a non-nil error if the request doesn't comply with
// our validation rules defined on the protofiles.
func (prv protoRequestValidator) ValidateRequest(req any) error {
	if err := prv.Validate(req.(protoreflect.ProtoMessage)); err != nil {
		return validationErrToGRPCErr(err)
	}
	return nil
}

// Takes a *protovalidate.ValidationError and returns an InvalidArgument(3) GRPC error with its corresponding message.
// Validation errors are always returned as InvalidArgument.
// They get translated to 400 Bad Request on the HTTP error handler.
func validationErrToGRPCErr(err error) error {
	var validationErr *protovalidate.ValidationError
	if ok := errors.As(err, &validationErr); ok {
		humanFacingMsg := parseValidationErr(validationErr)
		return status.Error(codes.InvalidArgument, fmt.Sprintf(errs.ValidatingRequest, humanFacingMsg))
	}

	var runtimeErr *protovalidate.RuntimeError
	if ok := errors.As(err, &runtimeErr); ok {
		logs.LogUnexpected(runtimeErr)
		return status.Error(codes.InvalidArgument, fmt.Sprintf(errs.ValidatingRequestRuntime, runtimeErr))
	}

	logs.LogUnexpected(err)
	return status.Error(codes.InvalidArgument, fmt.Sprintf(errs.ValidatingRequestUnexpected, err))
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type failedValidation validate.Violation

// Returns a string with the info on what validations did a request fail.
// It's the human-facing format of all GRPC Invalid Argument and HTTP Bad Request errors.
func parseValidationErr(vErr *protovalidate.ValidationError) string {
	out := ""
	for i, failedV := range vErr.Violations {
		out += (*failedValidation)(failedV).HumanFacing()
		if i < len(vErr.Violations)-1 {
			out += (", ")
		}
	}
	return out
}

// Pointer receiver because validate.Violation has a mutex.
func (v *failedValidation) HumanFacing() string {
	if strings.Contains(v.Message, "match regex pattern") {
		return fmt.Sprintf("%s value has an invalid format", v.FieldPath) // Don't display the regex patterns we use.
	}
	return fmt.Sprintf("%s %s", v.FieldPath, v.Message) // Example -> 'username field cannot be empty'
}
