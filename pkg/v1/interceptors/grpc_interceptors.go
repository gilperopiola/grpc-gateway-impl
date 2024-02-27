package interceptors

import (
	"context"
	"log"
	"time"

	"github.com/bufbuild/protovalidate-go"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* ----------------------------------- */
/*        - gRPC Interceptors -        */
/* ----------------------------------- */

// GetAll returns all the gRPC interceptors as ServerOptions.
// If TLS is not enabled, return the default interceptors.
// If TLS is enabled, return the default interceptors + the TLS security interceptor.
func GetAll(logger *zap.Logger, validator *protovalidate.Validator, tlsEnabled bool, certPath, keyPath string) []grpc.ServerOption {
	if !tlsEnabled {
		return []grpc.ServerOption{
			newDefaultInterceptors(logger, validator),
		}
	}

	return []grpc.ServerOption{
		newGRPCTLSSecurityInterceptor(certPath, keyPath),
		newDefaultInterceptors(logger, validator),
	}
}

// newDefaultInterceptors returns the default gRPC interceptors.
func newDefaultInterceptors(logger *zap.Logger, validator *protovalidate.Validator) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		newGRPCLoggerInterceptor(logger),
		newGRPCValidatorInterceptor(validator),
		newGRPCRecoveryInterceptor(logger),
	)
}

// newGRPCRecoveryInterceptor returns a gRPC interceptor that recovers from panics.
func newGRPCRecoveryInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor(
		grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
			logger.Error("gRPC Panic!", zap.Any("panic", p))
			return status.Errorf(codes.Internal, "unexpected panic, something went wrong: %v", p)
		}),
	)
}

// newGRPCValidatorInterceptor takes a *protovalidate.Validator and returns a gRPC interceptor
// that enforces the validation rules written in the .proto files.
// It returns a gRPC InvalidArgument error if the validation fails.
func newGRPCValidatorInterceptor(protoValidator *protovalidate.Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := protoValidator.Validate(req.(protoreflect.ProtoMessage)); err != nil {
			return nil, fromValidationErrToGRPCInvalidArgErr(err)
		}
		return handler(ctx, req) // Call next handler.
	}
}

// newGRPCLoggerInterceptor returns a gRPC interceptor that logs every gRPC request that comes in through the gRPC server.
// It logs the full method of the request, and it runs before any validation.
func newGRPCLoggerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	sugar := logger.Sugar()

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			sugar.Errorw("gRPC Error",
				"method", info.FullMethod,
				"duration", duration,
				"error", err,
			)
		} else {
			sugar.Infow("gRPC Request",
				"method", info.FullMethod,
				"duration", duration,
			)
		}

		return resp, err
	}
}

// newGRPCTLSSecurityInterceptor returns a grpc.ServerOption that enables TLS communication.
// It loads the server's certificate and key from a file.
func newGRPCTLSSecurityInterceptor(certPath, keyPath string) grpc.ServerOption {
	creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		log.Fatalf(msgErrLoadingTLSCredentials_Fatal, err)
	}
	return grpc.Creds(creds)
}

const (
	msgErrLoadingTLSCredentials_Fatal = "Failed to load server TLS credentials: %v"
)
