package servers

import (
	"context"
	"runtime"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Interceptors are a chain of handlers that wrap around our service's handler.

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - GRPC Interceptors -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// RateLimiter + PanicRecoverer + GRPCLogger + TokenValidator + RequestValidator + CtxCancelled.
func getInterceptors(tools core.Tools) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		newRateLimitingInterceptor(tools),
		newPanicRecovererInterceptor(),
		newLoggingInterceptor(),
		newAuthValidationInterceptor(tools),
		newRequestValidationInterceptor(tools),
		newCtxCancelledInterceptor(),
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns a GRPC Interceptor that validates requests.
func newRequestValidationInterceptor(tools core.Tools) grpc.UnaryServerInterceptor {
	return func(c context.Context, req any, _ *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
		if err := tools.ValidateRequest(req); err != nil {
			return nil, err
		}

		// Call next handler.
		return next(c, req)
	}
}

// Returns a GRPC Interceptor that validates the auth to access the desired Route is OK.
// Public Routes need no authorization.
// JWT
// It adds the UserID and Username to the request's context.
func newAuthValidationInterceptor(tools core.Tools) grpc.UnaryServerInterceptor {
	return func(c context.Context, req any, i *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {

		route := shared.RouteNameFromGRPCMethod(i.FullMethod)
		authForRoute := shared.AuthForRoute(route)

		// If it's a public endpoint, just go ahead.
		// Note that the user's info is not added to the context.
		if authForRoute == shared.RouteAuthPublic {
			return next(c, req)
		}

		// Public routes skip token validation.
		// That means the context will not have the user's info.
		if shared.AuthForRoute(route) != shared.RouteAuthPublic {
			claims, err := tools.ValidateToken(c, req, route)
			if err != nil {
				return nil, err
			}

			// Gets user info from claims and adds it to the request's context.
			userID, username := claims.GetUserInfo()
			c = tools.AddUserInfoToCtx(c, userID, username)
		}

		// Call next handler.
		return next(c, req)
	}
}

// Returns a GRPC Interceptor that logs GRPC requests.
func newLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(c context.Context, req any, i *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
		start := time.Now()

		// Call next handler.
		resp, err := next(c, req)

		duration := time.Since(start)
		logs.LogGRPC(i.FullMethod, duration, err)
		return resp, err
	}
}

// Returns a GRPC Interceptor that checks if the context has been cancelled before processing the request.
func newCtxCancelledInterceptor() grpc.UnaryServerInterceptor {
	return func(c context.Context, req any, _ *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
		if err := c.Err(); err != nil {
			return nil, status.Error(codes.Canceled, err.Error())
		}

		// Call next handler.
		return next(c, req)
	}
}

// Returns a GRPC Interceptor that limits the request-processing-rate of the server.
func newRateLimitingInterceptor(tools core.Tools) grpc.UnaryServerInterceptor {
	return func(c context.Context, req any, _ *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
		if ok := tools.AllowRate(); !ok {
			logs.LogStrange("rate limit exceeded")
			return nil, status.Error(codes.ResourceExhausted, errs.RateLimitedMsg)
		}

		// Call next handler.
		return next(c, req)
	}
}

// Returns a GRPC Interceptor that recovers from panics and logs them.
func newPanicRecovererInterceptor() grpc.UnaryServerInterceptor {
	return func(c context.Context, req any, _ *grpc.UnaryServerInfo, next grpc.UnaryHandler) (resp any, err error) {

		// This is set to true after the call to the next handler gets back.
		handlerFinishedOK := false

		// And this gets executed when a panic happens or after this func finishes.
		// Panics will recover, logging the error and returning to the user a standard panic response.
		defer func() {
			if err := recover(); err != nil || !handlerFinishedOK {
				stackBuf := make([]byte, 2048)
				stackBuf = stackBuf[:runtime.Stack(stackBuf, false)]
				zap.L().Error("GRPC Panic", zap.Any("error", err), zap.Any("stack", stackBuf))
				err = status.Error(codes.Internal, errs.PanicMsg)
			}
		}()

		// Call next handler.
		resp, err = next(c, req) // <- Panics happen here.

		handlerFinishedOK = true
		return resp, err
	}
}
