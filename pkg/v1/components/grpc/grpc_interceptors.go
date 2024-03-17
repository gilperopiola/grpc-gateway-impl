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
	"google.golang.org/grpc/status"
)

/* ----------------------------------- */
/*        - gRPC Interceptors -        */
/* ----------------------------------- */

/* Interceptors are used to intervene gRPC requests and responses.
/* They apply features across the entire API, not only gRPC but also HTTP.
/* Even though we only use Unary Interceptors, Stream Interceptors also exist. */

// getUnaryInterceptors returns the gRPC Unary Interceptors.
// These Interceptors are then chained together and added to the gRPC Server as a ServerOption.
func getUnaryInterceptors(components *components.Wrapper) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		rateLimiterInterceptor(components.RateLimiter, components.Logger),
		loggerInterceptor(components.Logger),
		tokenValidationInterceptor(components.Authenticator),
		inputValidationInterceptor(components.InputValidator),
		contextCancelledInterceptor(),
		recoveryInterceptor(components.Logger),
	}
}

// tokenValidationInterceptor returns a gRPC interceptor that validates if the user is allowed to access the endpoint.
func tokenValidationInterceptor(tokenValidator common.TokenValidator) grpc.UnaryServerInterceptor {
	return tokenValidator.Validate
}

// inputValidationInterceptor takes a common.InputValidator and returns a gRPC interceptor
// that enforces the validation rules written in the .proto files.
func inputValidationInterceptor(validator common.InputValidator) grpc.UnaryServerInterceptor {
	return validator.ValidateInput
}

// loggerInterceptor returns a gRPC interceptor that logs every gRPC request that comes in through the gRPC server.
func loggerInterceptor(logger *common.Logger) grpc.UnaryServerInterceptor {
	return logger.LogGRPC
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

// contextCancelledInterceptor returns a gRPC interceptor that checks if the context has been cancelled before processing the request.
func contextCancelledInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if ctx.Err() != nil {
			return nil, ctx.Err()
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
