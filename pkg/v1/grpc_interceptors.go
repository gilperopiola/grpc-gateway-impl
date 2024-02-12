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

// GetInterceptors returns a gRPC server option that chains all interceptors together.
// These may be gRPC interceptors, but they are also executed through HTTP calls.
func GetInterceptors() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		NewValidationInterceptor(),
	)
}

// NewValidationInterceptor creates a new *protovalidate.Validator and returns a gRPC interceptor
// that enforces the validation rules written in the .proto files.
func NewValidationInterceptor() grpc.UnaryServerInterceptor {
	protoValidator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("Failed to create proto validator: %v", err)
	}

	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		err := protoValidator.Validate(req.(protoreflect.ProtoMessage))
		if err == nil {
			return handler(ctx, req) // If there's no validation error, we call the next handler.
		}

		// If there was an error, we check if it's a ValidationError. If it's not, we return a generic error.
		var validationErr *protovalidate.ValidationError
		if ok := errors.As(err, &validationErr); !ok {
			return nil, status.Error(codes.InvalidArgument, "validation error")
		}

		// If it is, we go through each violation and format the error message accordingly.
		return nil, status.Error(codes.InvalidArgument, getErrorMsgFromViolations(validationErr.ToProto()))
	}
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
