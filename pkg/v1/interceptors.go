package v1

import (
	"context"
	"errors"
	"log"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* ----------------------------------- */
/*          - Interceptors -           */
/* ----------------------------------- */

// NewValidationInterceptor creates a new *protovalidate.Validator and returns a gRPC interceptor (also executed through HTTP calls)
// that enforces the validation rules written in the .proto files.
func NewValidationInterceptor() grpc.UnaryServerInterceptor {
	protoValidator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("Failed to create proto validator: %v", err)
	}

	// This function is the interceptor that will be executed for every gRPC / HTTP call.
	fn := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := protoValidator.Validate(req.(protoreflect.ProtoMessage)); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return handler(ctx, req)
	}

	return fn
}

func NewHTTPErrorHandlerInterceptor() grpc.UnaryServerInterceptor {
	fn := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		got, err := handler(ctx, req)

		var validationErr *protovalidate.ValidationError
		if errors.As(err, &validationErr) {
			err = validationErr
		}
		return got, err
	}
	return fn
}
