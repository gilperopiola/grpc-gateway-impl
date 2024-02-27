package interceptors

import (
	"context"
	"time"

	"github.com/bufbuild/protovalidate-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* ----------------------------------- */
/*        - gRPC Interceptors -        */
/* ----------------------------------- */

func GetAll(logger *zap.Logger, validator *protovalidate.Validator) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		newGRPCLoggerInterceptor(logger),
		newGRPCValidatorInterceptor(validator),
	)
}

// newGRPCValidatorInterceptor takes a *protovalidate.Validator and returns a gRPC interceptor
// that enforces the validation rules written in the .proto files.
// It returns a gRPC InvalidArgument error if the validation fails.
func newGRPCValidatorInterceptor(protoValidator *protovalidate.Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := protoValidator.Validate(req.(protoreflect.ProtoMessage)); err != nil {
			return nil, fromValidationErrToGRPCInvalidArgErr(err)
		}
		return handler(ctx, req) // Call next handler.
	}
}

// newGRPCLoggerInterceptor returns a gRPC interceptor that logs every gRPC request that comes in through the gRPC server.
// It logs the full method of the request, and it runs before any validation.
func newGRPCLoggerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	sugar := logger.Sugar()

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			sugar.Errorw("gRPC Error",
				"method", info.FullMethod,
				"duration", duration,
				"error", err,
			)
		} else {
			sugar.Infow("gRPC Request",
				"method", info.FullMethod,
				"duration", duration,
			)
		}

		return resp, err
	}
}
