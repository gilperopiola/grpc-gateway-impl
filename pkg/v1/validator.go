package v1

import (
	"errors"
	"fmt"
	"log"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	errMsgProtoValidator = "Failed to create proto validator: %v" // Fatal error.

	errMsgValidation           = "validation error: %v"
	errMsgUnexpectedValidation = "unexpected validation error: %v"
)

// NewProtoValidator returns a new instance of *protovalidate.Validator.
// It calls log.Fatalf if it fails to create the validator.
func NewProtoValidator() *protovalidate.Validator {
	protoValidator, err := protovalidate.New()
	if err != nil {
		log.Fatalf(errMsgProtoValidator, err)
	}
	return protoValidator
}

// GetGRPCErrFromValidationErr returns an InvalidArgument(3) gRPC error with its corresponding message.
// It gets translated to a 400 Bad Request on the error handler.
func GetGRPCErrFromValidationErr(err error) error {
	var validationErr *protovalidate.ValidationError
	if ok := errors.As(err, &validationErr); ok {
		return status.Error(codes.InvalidArgument, fmt.Sprintf(errMsgValidation, GetErrorMsgFromValidateViolations(validationErr.ToProto())))
	}
	return status.Error(codes.InvalidArgument, fmt.Sprintf(errMsgUnexpectedValidation, err))
}

// GetErrorMsgFromValidateViolations returns a formatted error based on the validation rules that were broken.
func GetErrorMsgFromValidateViolations(violations *validate.Violations) error {
	out := ""
	for i, v := range violations.Violations {
		out += fmt.Sprintf("%s %s", v.FieldPath, v.Message)
		if i < len(violations.Violations)-1 {
			out += ", "
		}
	}
	return errors.New(out)
}
