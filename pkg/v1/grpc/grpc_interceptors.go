package grpc

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/misc"

	"github.com/bufbuild/protovalidate-go"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

/* ----------------------------------- */
/*        - gRPC Interceptors -        */
/* ----------------------------------- */

// AllInterceptors returns all the gRPC interceptors as ServerOptions.
// If TLS is enabled, return the TLS security interceptor + the default interceptors.
// If TLS is not enabled, only return the default interceptors.
func AllInterceptors(c *cfg.Config, lg *zap.Logger, vldtr *protovalidate.Validator, svCreds credentials.TransportCredentials) []grpc.ServerOption {
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
func newDefaultInterceptors(logger *zap.Logger, validator *protovalidate.Validator, rlConfig *cfg.RateLimiterConfig) grpc.ServerOption {
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
	return misc.Validate(protoValidator)
}

// newGRPCLoggerInterceptor returns a gRPC interceptor that logs every gRPC request that comes in through the gRPC server.
// It logs the full method of the request, and it runs before any validation.
func newGRPCLoggerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return misc.LogGRPC(logger)
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
			return status.Errorf(codes.Internal, errs.ErrMsgPanic)
		}),
	)
}

// newGRPCRateLimiterInterceptor returns a gRPC interceptor that limits the rate of requests.
// It returns a gRPC ResourceExhausted error if the rate limit is exceeded.
func newGRPCRateLimiterInterceptor(cfg *cfg.RateLimiterConfig) grpc.UnaryServerInterceptor {
	limiter := rate.NewLimiter(rate.Limit(cfg.TokensPerSecond), cfg.MaxTokens)
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !limiter.Allow() {
			return nil, status.Errorf(codes.ResourceExhausted, errs.ErrMsgRateLimitExceeded)
		}
		return handler(ctx, req)
	}
}
