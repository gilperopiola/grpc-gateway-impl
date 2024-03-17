package grpc

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/common"
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

/* Interceptors are used to intervene gRPC requests and responses.
/* They apply features across the entire API, as HTTP endpoints are converted to gRPC calls. */

// AllInterceptors returns all the gRPC interceptors as ServerOptions.
func AllInterceptors(components *components.Wrapper, tlsEnabled bool) []grpc.ServerOption {
	out := make([]grpc.ServerOption, 0)

	if tlsEnabled {
		out = append(out, tlsInterceptor(components.ServerCreds)) // TLS interceptor.
	}
	out = append(out, defaultInterceptors(components)) // Default interceptors.

	return out
}

// defaultInterceptors returns the default gRPC interceptors.
func defaultInterceptors(components *components.Wrapper) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		rateLimiterInterceptor(components.RateLimiter, components.Logger),
		loggerInterceptor(components.Logger),
		tokenValidationInterceptor(components.Authenticator),
		inputValidationInterceptor(components.InputValidator),
		recoveryInterceptor(components.Logger),
	)
}

// tokenValidationInterceptor returns a gRPC interceptor that validates if the user is allowed to access the endpoint.
func tokenValidationInterceptor(tokenValidator common.TokenValidator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, svInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return tokenValidator.Validate(ctx, req, svInfo, handler)
	}
}

// inputValidationInterceptor takes a *Validator and returns a gRPC interceptor
// that enforces the validation rules written in the .proto files.
func inputValidationInterceptor(validator common.InputValidator) grpc.UnaryServerInterceptor {
	return validator.ValidateInput()
}

// loggerInterceptor returns a gRPC interceptor that logs every gRPC request that comes in through the gRPC server.
func loggerInterceptor(logger *common.Logger) grpc.UnaryServerInterceptor {
	return logger.LogGRPC()
}

// rateLimiterInterceptor returns a gRPC interceptor that limits the rate of requests that the server can process.
// Returns a gRPC ResourceExhausted error if the rate limit is exceeded.
func rateLimiterInterceptor(limiter *rate.Limiter, logger *common.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !limiter.Allow() {
			logger.Error("Rate limit exceeded!")
			return nil, status.Errorf(codes.ResourceExhausted, errs.ErrMsgRateLimitExceeded)
		}
		return handler(ctx, req)
	}
}

// recoveryInterceptor returns a gRPC interceptor that recovers from panics.
func recoveryInterceptor(logger *common.Logger) grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor(
		grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
			logger.Error("gRPC Panic!", zap.Any("info", p))
			return status.Errorf(codes.Internal, errs.ErrMsgPanic)
		}),
	)
}

// tlsInterceptor returns a grpc.ServerOption that enables TLS communication.
// It loads the server's certificate and key from a file.
func tlsInterceptor(serverCreds credentials.TransportCredentials) grpc.ServerOption {
	return grpc.Creds(serverCreds)
}
