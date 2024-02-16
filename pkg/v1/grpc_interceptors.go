package v1

import (
	"context"
	"errors"
	"fmt"
	"log"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* ----------------------------------- */
/*        - gRPC Interceptors -        */
/* ----------------------------------- */

// GetInterceptorsAsServerOption returns a gRPC server option that chains all interceptors together.
// These may be gRPC interceptors, but they are also executed through HTTP calls.
func GetInterceptorsAsServerOption() grpc.ServerOption {
	protoValidator := newProtoValidator()

	return grpc.ChainUnaryInterceptor(
		NewValidationInterceptor(protoValidator),
	)
}

// NewValidationInterceptor instantiates a new *protovalidate.Validator and returns a gRPC interceptor
// that enforces the validation rules written in the .proto files.
func NewValidationInterceptor(protoValidator *protovalidate.Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := protoValidator.Validate(req.(protoreflect.ProtoMessage)); err != nil {
			return nil, getGRPCErrFromValidationErr(err)
		}
		return handler(ctx, req) // Call next handler.
	}
}

// getGRPCErrFromValidationErr returns an InvalidArgument gRPC error with its respective message.
// It translates to a 400 Bad Request HTTP status code.
func getGRPCErrFromValidationErr(err error) error {
	var validationErr *protovalidate.ValidationError
	if ok := errors.As(err, &validationErr); ok {
		return status.Error(codes.InvalidArgument, getErrorMsgFromViolations(validationErr.ToProto()))
	}
	return status.Error(codes.InvalidArgument, fmt.Sprintf(errMsgUnexpectedValidation, err))
}

// newProtoValidator returns a new instance of *protovalidate.Validator.
// It calls log.Fatalf if it fails to create the validator.
func newProtoValidator() *protovalidate.Validator {
	protoValidator, err := protovalidate.New()
	if err != nil {
		log.Fatalf(errMsgProtoValidator, err)
	}
	return protoValidator
}

// getErrorMsgFromViolations returns a formatted error message based on the validate violations.
func getErrorMsgFromViolations(violations *validate.Violations) string {
	out := ""
	for i, v := range violations.Violations {
		out += fmt.Sprintf("%s %s", v.FieldPath, v.Message)
		if i < len(violations.Violations)-1 {
			out += ", "
		}
	}
	return out
}

var (
	errMsgProtoValidator       = "Failed to create proto validator: %v"
	errMsgUnexpectedValidation = "unexpected validation error: %v"
)
