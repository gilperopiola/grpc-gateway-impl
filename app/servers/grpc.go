package servers

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - GRPC Server -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// -> Server Options are used to configure the GRPC Server.
// -> Our interceptors are actually added here, chained together as a single ServerOption.

// Returns the GRPC Server Options, interceptors included.
func getGRPCServerOptions(toolbox core.Toolbox, tlsEnabled bool) god.GRPCServerOpts {
	serverOpts := []grpc.ServerOption{}

	if tlsEnabled {
		serverOpts = append(serverOpts, grpc.Creds(toolbox.GetServerCreds()))
	}

	// Chain all Unary Interceptors into a single ServerOption and add it as a Server Opt.
	defaultInterceptors := getGRPCInterceptors(toolbox)
	serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(defaultInterceptors...))

	return serverOpts
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Dial Options are used by the HTTP Gateway when connecting to the GRPC Server.
func getGRPCDialOptions(tlsClientCreds god.TLSCreds) god.GRPCDialOpts {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(tlsClientCreds),
		grpc.WithUserAgent(customUserAgent),
	}
}

const customUserAgent = "by @gilperopiola"

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - GRPC Interceptors -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Interceptors are used to intervene GRPC Requests and Responses.
// Even though we just use Unary Interceptors, Stream Interceptors are also available.
//
// grpc.UnaryServerInterceptor -> func(ctx, req, info, handler) (any, error)

// Returns the GRPC Unary Interceptors.
// These Interceptors are then chained together and added to the GRPC Server as a grpc.ServerOption.
func getGRPCInterceptors(toolbox core.Toolbox) god.GRPCInterceptors {
	return []grpc.UnaryServerInterceptor{
		toolbox.LimitGRPC,
		handlePanicsAndRecover,
		core.LogGRPCRequest,
		toolbox.ValidateToken,
		toolbox.ValidateGRPC,
		handleContextCancellation,
	}
}

// Returns a GRPC Interceptor that checks if the context has been cancelled before processing the request.
func handleContextCancellation(ctx god.Ctx, req any, _ *god.GRPCInfo, handler god.GRPCHandler) (any, error) {
	if ctx.Err() != nil {
		core.LogWeirdBehaviour("Context cancelled early", req)
		return nil, ctx.Err()
	}
	return handler(ctx, req) // T0D0 Should we also check after the service call?
}

// Returns a GRPC Interceptor that recovers from panics and logs them
// -> Adapted from github.com/grpc-ecosystem/go-grpc-middleware/recovery
func handlePanicsAndRecover(ctx god.Ctx, req any, _ *god.GRPCInfo, handler god.GRPCHandler) (resp any, err error) {
	handlerCalled := false

	defer func() {
		if panicInfo := recover(); panicInfo != nil || !handlerCalled {
			zap.S().Error("GRPC Panic!", zap.Any("info", panicInfo), zap.Any("context", ctx))
			err = status.Errorf(codes.Internal, errs.PanicMsg)
		}
	}()

	resp, err = handler(ctx, req) // <- Panics happen here.

	// This var is checked on the defer, if it panicked then this place is never reached and handlerCalled = false.
	handlerCalled = true

	return resp, err
}
