package v1

import (
	"context"
	"fmt"
	"log"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// NewValidationInterceptor creates a new *protovalidate.Validator and returns a gRPC interceptor (also executed through HTTP calls)
// that enforces the validation rules written in the .proto files.
func NewValidationInterceptor() grpc.ServerOption {
	protoValidator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("Failed to create proto validator: %v", err)
	}

	fn := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := protoValidator.Validate(req.(protoreflect.ProtoMessage)); err != nil {
			return nil, fmt.Errorf("invalid request: %w", err)
		}
		return handler(ctx, req)
	}

	return grpc.UnaryInterceptor(fn)
}
