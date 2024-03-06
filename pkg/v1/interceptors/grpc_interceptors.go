package interceptors

import (
	"context"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	"github.com/gilperopiola/grpc-gateway-impl/server/config"

	"github.com/bufbuild/protovalidate-go"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* ----------------------------------- */
/*        - gRPC Interceptors -        */
/* ----------------------------------- */

// GetAll returns all the gRPC interceptors as ServerOptions.
// If TLS is enabled, return the TLS security interceptor + the default interceptors.
// If TLS is not enabled, only return the default interceptors.
func GetAll(c *config.Config, lg *zap.Logger, vldtr *protovalidate.Validator, svCreds credentials.TransportCredentials) []grpc.ServerOption {
	out := make([]grpc.ServerOption, 0)

	// Add TLS interceptor.
	if c.TLS.Enabled {
		out = append(out, newGRPCTLSInterceptor(svCreds))
	}

	// Add default interceptors.
	out = append(out, newDefaultInterceptors(lg, vldtr, c.RateLimiter))

	return out
}

// newDefaultInterceptors returns the default gRPC interceptors.
func newDefaultInterceptors(logger *zap.Logger, validator *protovalidate.Validator, rlConfig config.RateLimiterConfig) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		newGRPCRateLimiterInterceptor(rlConfig),
		newGRPCLoggerInterceptor(logger),
		newGRPCValidatorInterceptor(validator),
		newGRPCRecoveryInterceptor(logger),
	)
}

// newGRPCValidatorInterceptor takes a *protovalidate.Validator and returns a gRPC interceptor
// that enforces the validation rules written in the .proto files.
// It returns a gRPC InvalidArgument error if the validation fails.
func newGRPCValidatorInterceptor(protoValidator *protovalidate.Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := protoValidator.Validate(req.(protoreflect.ProtoMessage)); err != nil {
			return nil, validationErrToInvalidArgErr(err)
		}
		return handler(ctx, req) // Call next handler.
	}
}

// newGRPCLoggerInterceptor returns a gRPC interceptor that logs every gRPC request that comes in through the gRPC server.
// It logs the full method of the request, and it runs before any validation.
func newGRPCLoggerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return v1.LogGRPC(logger)
}

// newGRPCTLSInterceptor returns a grpc.ServerOption that enables TLS communication.
// It loads the server's certificate and key from a file.
func newGRPCTLSInterceptor(serverCredentials credentials.TransportCredentials) grpc.ServerOption {
	return grpc.Creds(serverCredentials)
}

// newGRPCRecoveryInterceptor returns a gRPC interceptor that recovers from panics.
func newGRPCRecoveryInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor(
		grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
			logger.Error("gRPC Panic!", zap.Any("info", p))
			return status.Errorf(codes.Internal, v1.ErrMsgPanic)
		}),
	)
}

// newGRPCRateLimiterInterceptor returns a gRPC interceptor that limits the rate of requests.
// It returns a gRPC ResourceExhausted error if the rate limit is exceeded.
func newGRPCRateLimiterInterceptor(config config.RateLimiterConfig) grpc.UnaryServerInterceptor {
	limiter := rate.NewLimiter(rate.Limit(config.TokensPerSecond), config.MaxTokens)
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !limiter.Allow() {
			return nil, status.Errorf(codes.ResourceExhausted, v1.ErrMsgRateLimitExceeded)
		}
		return handler(ctx, req)
	}
}
