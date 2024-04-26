package servers

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - GRPC Interceptors -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Interceptors are used to intervene GRPC Requests and Responses.
// Even though we just use Unary Interceptors, Stream Interceptors are also available.
//
// -> grpc.UnaryServerInterceptor = func(ctx, req, info, handler) (any, error)

// Returns the GRPC Unary Interceptors.
// These Interceptors are then chained together and added to the GRPC Server as a grpc.ServerOption.
func defaultGRPCInterceptors(toolbox core.Toolbox) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		toolbox.LimitGRPC,
		core.LogGRPCRequest,
		toolbox.ValidateToken,
		toolbox.ValidateGRPC,
		contextCancelledInterceptor(),
		panicRecoveryInterceptor(),
	}
}

// contextCancelledInterceptor returns a GRPC interceptor that checks if the context has been cancelled before processing the request.
func contextCancelledInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if ctx.Err() != nil {
			core.LogWeirdBehaviour("Context cancelled early", req)
			return nil, ctx.Err()
		}
		return handler(ctx, req) // T0D0 Should we also check after the service call?
	}
}

// panicRecoveryInterceptor returns a GRPC interceptor that recovers from panics.
func panicRecoveryInterceptor() grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor(
		grpc_recovery.WithRecoveryHandler(func(panicInfo any) error {
			zap.S().Error("GRPC Panic!", zap.Any("info", panicInfo))
			return status.Errorf(codes.Internal, errs.PanicMsg)
		}),
	)
}
