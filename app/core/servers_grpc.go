package core

import (
	"context"
	"net"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - gRPC Server -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func newGRPCServer(usersSvc pbs.UsersServiceServer, serverOpts []grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(serverOpts...)
	pbs.RegisterUsersServiceServer(grpcServer, usersSvc)
	return grpcServer
}

func runGRPC(grpcServer *grpc.Server) {
	zap.S().Infof("GRPC Gateway Implementation | GRPC Port %s 🚀", GRPCPort)

	lis, err := net.Listen("tcp", GRPCPort)
	if err != nil {
		LogUnexpectedAndPanic(err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			LogUnexpectedAndPanic(err)
		}
	}()
}

func shutdownGRPC(grpcServer *grpc.Server) {
	zap.S().Info("Shutting down gRPC server...")
	grpcServer.GracefulStop()
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - gRPC Interceptors -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Interceptors are used to intervene GRPC Requests and Responses.
// Even though we just use Unary Interceptors, Stream Interceptors are also available.
// grpc.UnaryServerInterceptor = func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (any, error)

// Returns the gRPC Unary Interceptors.
// These Interceptors are then chained together and added to the gRPC Server as a ServerOption.
func defaultGRPCInterceptors(tools Toolbox) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		rateLimiterInterceptor(tools.GetRateLimiter()),
		requestsLoggerInterceptor(),
		tokenAuthInterceptor(tools.GetAuthenticator()),
		requestsValidationInterceptor(tools.GetRequestsValidator()),
		contextCancelledInterceptor(),
		panicRecoveryInterceptor(),
	}
}

// Wraps a TokenValidator in an grpc.UnaryServerInterceptor. Enforces authentication rules.
func tokenAuthInterceptor(tokenValidator TokenValidator) grpc.UnaryServerInterceptor {
	return tokenValidator.ValidateToken
}

// Wraps an InputValidator in an grpc.UnaryServerInterceptor. Enforces request validation rules.
func requestsValidationInterceptor(reqsValidator RequestsValidator) grpc.UnaryServerInterceptor {
	return reqsValidator.ValidateRequest
}

// requestsLoggerInterceptor returns a gRPC interceptor that logs every gRPC request that comes in through the gRPC server.
func requestsLoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			zap.S().Errorw("gRPC Error", ZapRoute(info.FullMethod), ZapDuration(duration), ZapError(err))
		} else {
			zap.S().Infow("gRPC Request", ZapRoute(info.FullMethod), ZapDuration(duration))
		}

		return resp, err
	}
}

// rateLimiterInterceptor returns a gRPC interceptor that limits the rate of requests that the server can process.
// Returns a gRPC ResourceExhausted error if the rate limit is exceeded.
func rateLimiterInterceptor(limiter *rate.Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if !limiter.Allow() {
			zap.S().Error("Rate limit exceeded!")
			return nil, status.Errorf(codes.ResourceExhausted, errs.ErrMsgRateLimitExceeded)
		}
		return handler(ctx, req)
	}
}

// contextCancelledInterceptor returns a gRPC interceptor that checks if the context has been cancelled before processing the request.
func contextCancelledInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if ctx.Err() != nil {
			LogWeirdBehaviour("Context cancelled before processing request.", req)
			return nil, ctx.Err()
		}
		return handler(ctx, req)
	}
}

// panicRecoveryInterceptor returns a gRPC interceptor that recovers from panics.
func panicRecoveryInterceptor() grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor(
		grpc_recovery.WithRecoveryHandler(func(p any) error {
			zap.S().Error("gRPC Panic!", zap.Any("info", p))
			return status.Errorf(codes.Internal, errs.ErrMsgPanic)
		}),
	)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - gRPC Server Options -       */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Server Options are used to configure the gRPC Server.
// Our interceptors are actually added here, chained together as a ServerOption.

// Returns the gRPC Server Options, interceptors included.
func defaultGRPCServerOpts(tools Toolbox, tlsEnabled bool) []grpc.ServerOption {
	serverOpts := []grpc.ServerOption{}

	// Add TLS Option if enabled.
	if tlsEnabled {
		serverOpts = append(serverOpts, grpc.Creds(tools.GetTLSServerCreds()))
	}

	// Chain all Unary Interceptors into a single ServerOption and add it to the slice.
	interceptorsOpt := grpc.ChainUnaryInterceptor(defaultGRPCInterceptors(tools)...)
	serverOpts = append(serverOpts, interceptorsOpt)

	return serverOpts
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - gRPC Dial Options -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Dial Options are used by the HTTP Gateway when connecting to the gRPC Server.
func defaultGRPCDialOpts(tlsClientCreds credentials.TransportCredentials) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(tlsClientCreds),
		grpc.WithUserAgent(customUserAgent),
	}
}

const customUserAgent = "by @gilperopiola"
