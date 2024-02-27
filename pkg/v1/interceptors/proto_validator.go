package interceptors

import (
	"errors"
	"fmt"
	"log"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewProtoValidator returns a new instance of *protovalidate.Validator.
// It calls log.Fatalf if it fails to create the validator.
func NewProtoValidator() *protovalidate.Validator {
	protoValidator, err := protovalidate.New()
	if err != nil {
		log.Fatalf(msgNewProtoValidatorErr_Fatal, err)
	}
	return protoValidator
}

// fromValidationErrToGRPCInvalidArgErr returns an InvalidArgument(3) gRPC error with its corresponding message.
// It gets translated to a 400 Bad Request on the error handler.
// Validation errors are always returned as InvalidArgument.
// This functions is called from the validation interceptor.
func fromValidationErrToGRPCInvalidArgErr(err error) error {
	outErrorMsg := fmt.Sprintf(msgUnexpectedValidationErr, err)

	var validationErr *protovalidate.ValidationError
	if ok := errors.As(err, &validationErr); ok {
		brokenRules := validationErr.ToProto().GetViolations()
		outErrorMsg = fmt.Sprintf(msgErrInProtoValidation, makeStringFromBrokenValidationRules(brokenRules))
	}

	var runtimeErr *protovalidate.RuntimeError
	if ok := errors.As(err, &runtimeErr); ok {
		outErrorMsg = fmt.Sprintf(msgRuntimeErr, runtimeErr)
	}

	return status.Error(codes.InvalidArgument, outErrorMsg)
}

// makeStringFromBrokenValidationRules returns a string detailing the broken validation rules.
// The default concatenates the field path and the message of each broken rule.
// This is what the user will see as the error message:
// { "error": "username must be at least 3 characters long" } on a JSON 400 response.
func makeStringFromBrokenValidationRules(brokenRules []*validate.Violation) (out string) {
	for i, brokenRule := range brokenRules {
		out += getMsgFromBrokenRule(brokenRule)
		if i < len(brokenRules)-1 {
			out += ", "
		}
	}
	return
}

// getMsgFromBrokenRule is the default human-facing format in which validation errors translate.
func getMsgFromBrokenRule(v *validate.Violation) string {
	return fmt.Sprintf("%s %s", v.FieldPath, v.Message)
}

const (
	msgErrInProtoValidation    = "validation error: %v"
	msgRuntimeErr              = "unexpected runtime validation error: %v"
	msgUnexpectedValidationErr = "unexpected validation error: %v"

	msgNewProtoValidatorErr_Fatal = "Failed to create proto validator: %v" // Fatal error.
)
