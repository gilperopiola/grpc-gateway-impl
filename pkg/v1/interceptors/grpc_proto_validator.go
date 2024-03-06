package interceptors

import (
	"errors"
	"fmt"
	"log"
	"strings"

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

// validationErrToInvalidArgErr returns an InvalidArgument(3) gRPC error with its corresponding message.
// It gets translated to a 400 Bad Request on the HTTP error handler.
// Validation errors are always returned as InvalidArgument.
// This function is called from the validation interceptor.
func validationErrToInvalidArgErr(err error) error {
	message := fmt.Sprintf(v1.ErrMsgInValidationUnexpected, err)

	var validationErr *protovalidate.ValidationError
	if ok := errors.As(err, &validationErr); ok {
		brokenRules := validationErr.ToProto().GetViolations()
		message = fmt.Sprintf(v1.ErrMsgInValidation, parseBrokenRules(brokenRules))
	}

	var runtimeErr *protovalidate.RuntimeError
	if ok := errors.As(err, &runtimeErr); ok {
		message = fmt.Sprintf(v1.ErrMsgInValidationRuntime, runtimeErr)
	}

	return status.Error(codes.InvalidArgument, message)
}

// parseBrokenRules returns a string detailing the broken validation rules.
// It just concatenates the fieldname and the message of each broken rule.
//
// This is the default human-facing format in which validation errors translate.
func parseBrokenRules(brokenRules []*validate.Violation) string {
	out := ""
	for i, brokenRule := range brokenRules {
		out += parseBrokenRule(brokenRule)
		if i < len(brokenRules)-1 {
			out += ", "
		}
	}
	return out
}

// parseBrokenRule is the default human-facing format in which validation errors translate.
// Special cases: obfuscate the invalid regex errors and return a generic message.
func parseBrokenRule(v *validate.Violation) string {
	if strings.Contains(v.Message, "match regex pattern") {
		return fmt.Sprintf("%s value has an invalid format", v.FieldPath)
	}
	return fmt.Sprintf("%s %s", v.FieldPath, v.Message)
}
