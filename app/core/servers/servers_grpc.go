package servers

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - GRPC Server -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func newGRPCServer(service core.Service, serverOpts []grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(serverOpts...)
	pbs.RegisterUsersServiceServer(grpcServer, service)
	return grpcServer
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// -> Server Options are used to configure the GRPC Server.
// -> Our interceptors are actually added here, chained together as a single ServerOption.

// Returns the GRPC Server Options, interceptors included.
func defaultGRPCServerOpts(toolbox core.Toolbox, tlsEnabled bool) []grpc.ServerOption {
	serverOpts := []grpc.ServerOption{}

	// Add TLS Option if enabled.
	if tlsEnabled {
		serverOpts = append(serverOpts, grpc.Creds(toolbox.GetTLSTool().GetServerCreds()))
	}

	// Chain all Unary Interceptors into a single ServerOption and add it to the slice.
	defaultInterceptors := defaultGRPCInterceptors(toolbox)
	serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(defaultInterceptors...))

	return serverOpts
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Dial Options are used by the HTTP Gateway when connecting to the GRPC Server.
func defaultGRPCDialOpts(tlsClientCreds credentials.TransportCredentials) []grpc.DialOption {
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
// -> grpc.UnaryServerInterceptor = func(ctx, req, info, handler) (any, error)

type interceptorFunc grpc.UnaryServerInterceptor

// Returns the GRPC Unary Interceptors.
// These Interceptors are then chained together and added to the GRPC Server as a grpc.ServerOption.
func defaultGRPCInterceptors(toolbox core.Toolbox) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		toolbox.LimitGRPC,
		core.LogGRPCRequest,
		toolbox.ValidateToken,
		toolbox.ValidateGRPC,
		grpc.UnaryServerInterceptor(handleContextCancellation()),
		grpc.UnaryServerInterceptor(recoverFromPanic()),
	}
}

// Returns a GRPC interceptor that checks if the context has been cancelled before processing the request.
func handleContextCancellation() interceptorFunc {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if ctx.Err() != nil {
			core.LogWeirdBehaviour("Context cancelled early", req)
			return nil, ctx.Err()
		}
		return handler(ctx, req) // T0D0 Should we also check after the service call?
	}
}

// Returns a GRPC interceptor that recovers from panics.
func recoverFromPanic() interceptorFunc {
	return interceptorFunc(grpc_recovery.UnaryServerInterceptor(
		grpc_recovery.WithRecoveryHandler(func(panicInfo any) error {
			zap.S().Error("GRPC Panic!", zap.Any("info", panicInfo))
			return status.Errorf(codes.Internal, errs.PanicMsg)
		}),
	))
}
