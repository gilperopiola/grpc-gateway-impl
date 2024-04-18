package servers

import (
	"context"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/interfaces"
	"github.com/gilperopiola/grpc-gateway-impl/app/modules"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - gRPC Interceptors -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Interceptors are used to intervene GRPC Requests and Responses.
// Even though we just use Unary Interceptors, Stream Interceptors are also available.
// grpc.UnaryServerInterceptor = func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (any, error)

// getInterceptors returns the gRPC Unary Interceptors.
// These Interceptors are then chained together and added to the gRPC Server as a ServerOption.
func getInterceptors(modules *modules.Active) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		rateLimiterInterceptor(modules.RateLimiter),
		requestsLoggerInterceptor(),
		tokenValidationInterceptor(modules.Authenticator),
		inputValidationInterceptor(modules.InputValidator),
		contextCancelledInterceptor(),
		panicRecoveryInterceptor(),
	}
}

// Wraps a TokenValidator in an grpc.UnaryServerInterceptor. Enforces authentication rules.
func tokenValidationInterceptor(tokenValidator interfaces.TokenValidator) grpc.UnaryServerInterceptor {
	return tokenValidator.Validate
}

// Wraps an InputValidator in an grpc.UnaryServerInterceptor. Enforces request validation rules.
func inputValidationInterceptor(inputValidator interfaces.InputValidator) grpc.UnaryServerInterceptor {
	return inputValidator.ValidateInput
}

// requestsLoggerInterceptor returns a gRPC interceptor that logs every gRPC request that comes in through the gRPC server.
func requestsLoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			zap.S().Errorw("gRPC Error", core.ZapEndpoint(info.FullMethod), core.ZapDuration(duration), core.ZapError(err))
		} else {
			zap.S().Infow("gRPC Request", core.ZapEndpoint(info.FullMethod), core.ZapDuration(duration))
		}

		return resp, err
	}
}

// rateLimiterInterceptor returns a gRPC interceptor that limits the rate of requests that the server can process.
// Returns a gRPC ResourceExhausted error if the rate limit is exceeded.
func rateLimiterInterceptor(limiter *rate.Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !limiter.Allow() {
			zap.S().Error("Rate limit exceeded!")
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

// panicRecoveryInterceptor returns a gRPC interceptor that recovers from panics.
func panicRecoveryInterceptor() grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor(
		grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
			zap.S().Error("gRPC Panic!", zap.Any("info", p))
			return status.Errorf(codes.Internal, errs.ErrMsgPanic)
		}),
	)
}
