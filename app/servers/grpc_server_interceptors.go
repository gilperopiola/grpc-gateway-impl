package servers

import (
	"context"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/modules"

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

// Interceptors are used to intervene GRPC Requests and Responses.
// Even though we only use Unary Interceptors, Stream Interceptors are also available.

// getUnaryInterceptors returns the gRPC Unary Interceptors.
// These Interceptors are then chained together and added to the gRPC Server as a ServerOption.
func getUnaryInterceptors(allModules *modules.All) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		rateLimiterInterceptor(allModules.RateLimiter),
		loggerInterceptor(),
		tokenValidationInterceptor(allModules.Authenticator),
		inputValidationInterceptor(allModules.InputValidator),
		contextCancelledInterceptor(),
		recoveryInterceptor(),
	}
}

// tokenValidationInterceptor returns a gRPC interceptor that validates if the user is allowed to access the endpoint.
func tokenValidationInterceptor(tokenValidator modules.TokenValidator) grpc.UnaryServerInterceptor {
	return tokenValidator.Validate
}

// inputValidationInterceptor takes a modules.InputValidator and returns a gRPC interceptor
// that enforces the validation rules written in the .proto files.
func inputValidationInterceptor(validator modules.InputValidator) grpc.UnaryServerInterceptor {
	return validator.ValidateInput
}

// loggerInterceptor returns a gRPC interceptor that logs every gRPC request that comes in through the gRPC server.
func loggerInterceptor() grpc.UnaryServerInterceptor {
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

// recoveryInterceptor returns a gRPC interceptor that recovers from panics.
func recoveryInterceptor() grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor(
		grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
			zap.S().Error("gRPC Panic!", zap.Any("info", p))
			return status.Errorf(codes.Internal, errs.ErrMsgPanic)
		}),
	)
}
