package servers

import (
	"context"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"

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
		newTokenValidationInterceptor(tools),
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

// Returns a GRPC Interceptor that validates JWT tokens.
// It adds the UserID and Username to the request's context.
func newTokenValidationInterceptor(tools core.Tools) grpc.UnaryServerInterceptor {

	return func(c context.Context, req any, i *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
		route := core.RouteNameFromGRPC(i.FullMethod)

		if core.AuthForRoute(route) != models.RouteAuthPublic {
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
		core.LogGRPC(i.FullMethod, duration, err)
		return resp, err
	}
}

// Returns a GRPC Interceptor that checks if the context has been cancelled before processing the request.
func newCtxCancelledInterceptor() grpc.UnaryServerInterceptor {
	return func(c context.Context, req any, _ *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
		if err := c.Err(); err != nil {
			return nil, status.Errorf(codes.Canceled, err.Error())
		}

		// Call next handler.
		return next(c, req)
	}
}

// Returns a GRPC Interceptor that sets a rate limit on the server.
func newRateLimitingInterceptor(tools core.Tools) grpc.UnaryServerInterceptor {
	return func(c context.Context, req any, _ *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
		if ok := tools.AllowRate(); !ok {
			err := status.Errorf(codes.ResourceExhausted, errs.RateLimitedMsg)
			core.LogUnexpected(err)
			return nil, err
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
				zap.L().Error("GRPC Panic", zap.Any("error", err), zap.Any("context", c))
				err = status.Errorf(codes.Internal, errs.PanicMsg)
			}
		}()

		// Call next handler.
		resp, err = next(c, req) // <- Panics happen here.

		handlerFinishedOK = true
		return resp, err
	}
}
