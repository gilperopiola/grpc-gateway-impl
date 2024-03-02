package interceptors

import (
	"errors"
	"fmt"
	"log"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* ----------------------------------- */
/*         - Proto Validator -         */
/* ----------------------------------- */

// NewProtoValidator returns a new instance of *protovalidate.Validator.
// It calls log.Fatalf if it fails to create the validator.
func NewProtoValidator() *protovalidate.Validator {
	protoValidator, err := protovalidate.New()
	if err != nil {
		log.Fatalf(v1.FatalErrMsgCreatingProtoValidator, err)
	}
	return protoValidator
}

// fromValidationErrToGRPCInvalidArgErr returns an InvalidArgument(3) gRPC error with its corresponding message.
// It gets translated to a 400 Bad Request on the error handler.
// Validation errors are always returned as InvalidArgument.
// This functions is called from the validation interceptor.
func fromValidationErrToGRPCInvalidArgErr(err error) error {
	outErrMsg := fmt.Sprintf(v1.ErrMsgInValidationUnexpected, err)

	var validationErr *protovalidate.ValidationError
	if ok := errors.As(err, &validationErr); ok {
		brokenRules := validationErr.ToProto().GetViolations()
		outErrMsg = fmt.Sprintf(v1.ErrMsgInValidation, makeStringFromBrokenValidationRules(brokenRules))
	}

	var runtimeErr *protovalidate.RuntimeError
	if ok := errors.As(err, &runtimeErr); ok {
		outErrMsg = fmt.Sprintf(v1.ErrMsgInValidationRuntime, runtimeErr)
	}

	return status.Error(codes.InvalidArgument, outErrMsg)
}

// makeStringFromBrokenValidationRules returns a string detailing the broken validation rules.
// The default concatenates the field path and the message of each broken rule.
// This is what the user will see as the error message:
// { "error": "username must be at least 3 characters long" } on a JSON 400 response.
func makeStringFromBrokenValidationRules(brokenRules []*validate.Violation) string {
	out := ""
	for i, brokenRule := range brokenRules {
		out += getErrMsgFromBrokenRule(brokenRule)
		if i < len(brokenRules)-1 {
			out += ", "
		}
	}
	return out
}

// getErrMsgFromBrokenRule is the default human-facing format in which validation errors translate.
func getErrMsgFromBrokenRule(v *validate.Violation) string {
	return fmt.Sprintf("%s %s", v.FieldPath, v.Message)
}
