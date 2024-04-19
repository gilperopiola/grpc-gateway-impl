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

func NewGRPCServer(usersService pbs.UsersServiceServer, svOptions []grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(svOptions...)
	pbs.RegisterUsersServiceServer(grpcServer, usersService)
	return grpcServer
}

func RunGRPCServer(grpcServer *grpc.Server) {
	zap.S().Infof("Running gRPC on port %s!\n", GRPCPort)

	lis, err := net.Listen("tcp", GRPCPort)
	if err != nil {
		zap.S().Fatalf(errs.FatalErrMsgStartingGRPC, err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			zap.S().Fatalf(errs.FatalErrMsgServingGRPC, err)
		}
	}()
}

func ShutdownGRPCServer(grpcServer *grpc.Server) {
	zap.S().Info("Shutting down gRPC server...")
	grpcServer.GracefulStop()
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - gRPC Interceptors -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Interceptors are used to intervene GRPC Requests and Responses.
// Even though we just use Unary Interceptors, Stream Interceptors are also available.
// grpc.UnaryServerInterceptor = func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (any, error)

// getInterceptors returns the gRPC Unary Interceptors.
// These Interceptors are then chained together and added to the gRPC Server as a ServerOption.
func getInterceptors(tools ToolsAccessor) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		rateLimiterInterceptor(tools.GetRateLimiter()),
		requestsLoggerInterceptor(),
		tokenValidationInterceptor(tools.GetAuthenticator()),
		inputValidationInterceptor(tools.GetInputValidator()),
		contextCancelledInterceptor(),
		panicRecoveryInterceptor(),
	}
}

// Wraps a TokenValidator in an grpc.UnaryServerInterceptor. Enforces authentication rules.
func tokenValidationInterceptor(tokenValidator TokenValidator) grpc.UnaryServerInterceptor {
	return tokenValidator.Validate
}

// Wraps an InputValidator in an grpc.UnaryServerInterceptor. Enforces request validation rules.
func inputValidationInterceptor(inputValidator InputValidator) grpc.UnaryServerInterceptor {
	return inputValidator.ValidateInput
}

// requestsLoggerInterceptor returns a gRPC interceptor that logs every gRPC request that comes in through the gRPC server.
func requestsLoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			zap.S().Errorw("gRPC Error", ZapEndpoint(info.FullMethod), ZapDuration(duration), ZapError(err))
		} else {
			zap.S().Infow("gRPC Request", ZapEndpoint(info.FullMethod), ZapDuration(duration))
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
			LogWeirdBehaviour("Context cancelled before processing request.")
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

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - gRPC Server Options -       */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Server Options are used to configure the gRPC Server.
// Our interceptors are actually added here, chained together as a ServerOption.

// AllServerOptions returns the gRPC Server Options.
func AllServerOptions(tools ToolsAccessor, tlsEnabled bool) []grpc.ServerOption {
	serverOptions := []grpc.ServerOption{}

	// Add TLS Option if enabled.
	if tlsEnabled {
		serverOptions = append(serverOptions, grpc.Creds(tools.GetTLSServerCreds()))
	}

	// Chain all Unary Interceptors into a single ServerOption and add it to the slice.
	interceptorsOption := grpc.ChainUnaryInterceptor(getInterceptors(tools)...)
	serverOptions = append(serverOptions, interceptorsOption)

	return serverOptions
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - gRPC Dial Options -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Dial Options are used by the HTTP Gateway when connecting to the gRPC Server.

const (
	customUserAgent = "by @gilperopiola"
)

// AllDialOptions returns the gRPC Dial Options.
func AllDialOptions(tlsClientCreds credentials.TransportCredentials) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(tlsClientCreds),
		grpc.WithUserAgent(customUserAgent),
	}
}
