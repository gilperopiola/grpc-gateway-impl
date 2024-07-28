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

// -> Server Options are used to configure the GRPC Server.
// -> Our interceptors are actually added here, chained together as a single ServerOption.

// Returns the GRPC Server Options, interceptors included.
func getGRPCServerOpts(tools core.Tools, tls bool) god.GRPCServerOpts {
	serverOpts := []grpc.ServerOption{}

	if tls {
		tlsOpt := grpc.Creds(tools.GetServerCreds())
		serverOpts = append(serverOpts, tlsOpt)
	}

	// Chain all Unary Interceptors into a single ServerOption and add it.
	defaultInterceptors := getGRPCInterceptors(tools)
	serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(defaultInterceptors...))

	return serverOpts
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Dial Options are used by the HTTP Gateway when connecting to the GRPC Server.
func getGRPCDialOpts(tlsClientCreds god.TLSCreds) []grpc.DialOption {
	const customUserAgent = "by @gilperopiola"
	return []grpc.DialOption{
		grpc.WithTransportCredentials(tlsClientCreds),
		grpc.WithUserAgent(customUserAgent),
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - GRPC Interceptors -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Interceptors are used to intervene GRPC Requests and Responses.
// Even though we just use Unary Interceptors, Stream Interceptors are also available.
//
// grpc.UnaryServerInterceptor -> func(ctx, req, info, handler) (any, error)

// Returns the GRPC Unary Interceptors.
// These Interceptors are then chained together and added to the GRPC Server as a grpc.ServerOption.
func getGRPCInterceptors(tools core.Tools) god.GRPCInterceptors {
	return []grpc.UnaryServerInterceptor{
		tools.LimitGRPC,
		handlePanicsAndRecover,
		core.LogGRPCRequest,
		tools.ValidateToken,
		tools.ValidateGRPC,
		handleCtxCancel,
	}
}

// Returns a GRPC Interceptor that checks if the context has been cancelled before processing the request.
func handleCtxCancel(ctx god.Ctx, req any, _ *god.GRPCInfo, handler god.GRPCHandler) (any, error) {
	if err := ctx.Err(); err != nil {
		core.LogWeirdBehaviour("Context cancelled early", req, err)
		return nil, err
	}
	return handler(ctx, req) // - Should we also check after the service call?
}

// Returns a GRPC Interceptor that recovers from panics and logs them.
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

	// This is checked on the defer, if it panicked then this place is never reached and handlerCalled = false
	// which means it panicked.
	handlerCalled = true

	return resp, err
}
