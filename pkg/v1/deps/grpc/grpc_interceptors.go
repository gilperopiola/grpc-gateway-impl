package grpc

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/deps"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"

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
func AllInterceptors(deps *deps.Deps, tlsEnabled bool) []grpc.ServerOption {
	out := make([]grpc.ServerOption, 0)
	if tlsEnabled {
		out = append(out, getGRPCTLSInterceptor(deps.ServerCreds)) // TLS interceptor.
	}
	out = append(out, getDefaultInterceptors(deps)) // Default interceptors.
	return out
}

// getDefaultInterceptors returns the default gRPC interceptors.
func getDefaultInterceptors(deps *deps.Deps) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		getGRPCRateLimiterInterceptor(deps.RateLimiter, deps.Logger),
		getGRPCLoggerInterceptor(deps.Logger),
		getGRPCJWTInterceptor(deps.Authenticator),
		getGRPCValidatorInterceptor(deps.Validator),
		getGRPCRecoveryInterceptor(deps.Logger),
	)
}

func getGRPCJWTInterceptor(tokenValidator deps.TokenValidator) grpc.UnaryServerInterceptor {
	return tokenValidator.Validate()
}

// getGRPCValidatorInterceptor takes a *Validator and returns a gRPC interceptor
// that enforces the validation rules written in the .proto files.
func getGRPCValidatorInterceptor(validator *deps.Validator) grpc.UnaryServerInterceptor {
	return validator.Validate()
}

// getGRPCLoggerInterceptor returns a gRPC interceptor that logs every gRPC request that comes in through the gRPC server.
func getGRPCLoggerInterceptor(logger *deps.Logger) grpc.UnaryServerInterceptor {
	sugar := logger.Sugar()
	return logger.LogGRPC(sugar)
}

// getGRPCTLSInterceptor returns a grpc.ServerOption that enables TLS communication.
// It loads the server's certificate and key from a file.
func getGRPCTLSInterceptor(serverCreds credentials.TransportCredentials) grpc.ServerOption {
	return grpc.Creds(serverCreds)
}

// getGRPCRecoveryInterceptor returns a gRPC interceptor that recovers from panics.
func getGRPCRecoveryInterceptor(logger *deps.Logger) grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor(
		grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
			logger.Error("gRPC Panic!", zap.Any("info", p))
			return status.Errorf(codes.Internal, errs.ErrMsgPanic)
		}),
	)
}

// getGRPCRateLimiterInterceptor returns a gRPC interceptor that limits the rate of requests that the server can process.
// Returns a gRPC ResourceExhausted error if the rate limit is exceeded.
func getGRPCRateLimiterInterceptor(limiter *rate.Limiter, logger *deps.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !limiter.Allow() {
			logger.Error("Rate limit exceeded!")
			return nil, status.Errorf(codes.ResourceExhausted, errs.ErrMsgRateLimitExceeded)
		}
		return handler(ctx, req)
	}
}
