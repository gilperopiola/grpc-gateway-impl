package interceptors

import (
	"context"
	"log"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* ----------------------------------- */
/*        - gRPC Interceptors -        */
/* ----------------------------------- */

// NewGRPCLogger returns a gRPC interceptor that logs every gRPC request.
func NewGRPCLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Printf("gRPC request: %s\n", info.FullMethod)
		return handler(ctx, req)
	}
}

// NewGRPCValidator takes a *protovalidate.Validator and returns a gRPC interceptor
// that enforces the validation rules written in the .proto files.
func NewGRPCValidator(protoValidator *protovalidate.Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := protoValidator.Validate(req.(protoreflect.ProtoMessage)); err != nil {
			return nil, v1.GetGRPCErrFromValidationErr(err)
		}
		return handler(ctx, req) // Call next handler.
	}
}
