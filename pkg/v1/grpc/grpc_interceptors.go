package grpc

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/dependencies"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"
	"github.com/golang-jwt/jwt/v4"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

/* ----------------------------------- */
/*        - gRPC Interceptors -        */
/* ----------------------------------- */

// AllInterceptors returns all the gRPC interceptors as ServerOptions.
func AllInterceptors(deps *dependencies.Dependencies, tlsEnabled bool) []grpc.ServerOption {
	out := make([]grpc.ServerOption, 0)
	if tlsEnabled {
		out = append(out, getGRPCTLSInterceptor(deps.ServerCreds)) // TLS interceptor.
	}
	out = append(out, getDefaultInterceptors(deps)) // Default interceptors.
	return out
}

// getDefaultInterceptors returns the default gRPC interceptors.
func getDefaultInterceptors(deps *dependencies.Dependencies) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		getGRPCRateLimiterInterceptor(deps.RateLimiter, deps.Logger),
		getGRPCLoggerInterceptor(deps.Logger),
		getGRPCJWTInterceptor(AnyRole, false, "asdasfjoasidf"),
		getGRPCValidatorInterceptor(deps.Validator),
		getGRPCRecoveryInterceptor(deps.Logger),
	)
}

type Role string

const (
	AnyRole   Role = "any"
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

type CustomClaims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     Role   `json:"role"`
	jwt.RegisteredClaims
}

func getGRPCJWTInterceptor(role Role, shouldMatchUserID bool, secret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, svInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		if svInfo.FullMethod == "/users.UsersService/Signup" || svInfo.FullMethod == "/users.UsersService/Login" {
			return handler(ctx, req) // Next handler.
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Internal, "metadata not found")
		}

		authHeader := md["authorization"]
		if len(authHeader) == 0 || len(authHeader[0]) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "token not found")
		}

		authHeaderValue := authHeader[0]

		if !strings.HasPrefix(authHeaderValue, "Bearer ") {
			return nil, status.Errorf(codes.Unauthenticated, "token invalid format")
		}

		authHeaderValue = strings.TrimPrefix(authHeaderValue, "Bearer ")

		keyFunc := func(token *jwt.Token) (interface{}, error) { return []byte("some"), nil }
		jwtToken, err := jwt.ParseWithClaims(authHeaderValue, &CustomClaims{}, keyFunc)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "token invalid")
		}

		claims, ok := jwtToken.Claims.(*CustomClaims)
		if !ok || !jwtToken.Valid || claims.Valid() != nil || (role != AnyRole && claims.Role != role) {
			return nil, status.Errorf(codes.Unauthenticated, "token or role invalid")
		}

		if shouldMatchUserID {

			// Get user ID from svInfo and compare with token's ID
			str := strings.Split(svInfo.FullMethod, "/")
			if len(str) < 3 {
				return nil, status.Errorf(codes.Internal, "error getting user id from method")
			}

			// Get user ID from method
			// Example: /v1/users/1 -> 1
			urlUserID, err := strconv.Atoi(str[len(str)-1])
			if err != nil || claims.ID != fmt.Sprint(urlUserID) {
				return nil, status.Errorf(codes.Unauthenticated, "user id mismatch")
			}
		}

		addUserInfoToContext(ctx, claims)

		return handler(ctx, req) // Next handler.
	}
}

func addUserInfoToContext(c context.Context, claims *CustomClaims) {
	userID, _ := strconv.Atoi(claims.ID)
	c = context.WithValue(c, "UserID", userID)
	c = context.WithValue(c, "Username", claims.Username)
	c = context.WithValue(c, "Email", claims.Email)
}

func getTokenFromAuthorizationHeader(c context.Context, secret string) (*jwt.Token, error) {
	// Get token string from headers
	tokenString, ok := c.Value("Authorization").(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "user id not found in context")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Decode string into actual *jwt.Token
	return decodeTokenString(tokenString, secret)
}

// decodeTokenString decodes a JWT token string into a *jwt.Token
func decodeTokenString(tokenString, secret string) (*jwt.Token, error) {

	// Check length
	if len(tokenString) < 40 {
		return nil, status.Errorf(codes.Unauthenticated, "asdasds")
	}

	// Make key function and return parsed token
	keyFunc := func(token *jwt.Token) (interface{}, error) { return []byte(secret), nil }
	return jwt.ParseWithClaims(tokenString, &CustomClaims{}, keyFunc)
}

// getGRPCValidatorInterceptor takes a *Validator and returns a gRPC interceptor
// that enforces the validation rules written in the .proto files.
func getGRPCValidatorInterceptor(validator *dependencies.Validator) grpc.UnaryServerInterceptor {
	return validator.Validate()
}

// getGRPCLoggerInterceptor returns a gRPC interceptor that logs every gRPC request that comes in through the gRPC server.
func getGRPCLoggerInterceptor(logger *dependencies.Logger) grpc.UnaryServerInterceptor {
	sugar := logger.Sugar()
	return logger.LogGRPC(sugar)
}

// getGRPCTLSInterceptor returns a grpc.ServerOption that enables TLS communication.
// It loads the server's certificate and key from a file.
func getGRPCTLSInterceptor(serverCreds credentials.TransportCredentials) grpc.ServerOption {
	return grpc.Creds(serverCreds)
}

// getGRPCRecoveryInterceptor returns a gRPC interceptor that recovers from panics.
func getGRPCRecoveryInterceptor(logger *dependencies.Logger) grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor(
		grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
			logger.Error("gRPC Panic!", zap.Any("info", p))
			return status.Errorf(codes.Internal, errs.ErrMsgPanic)
		}),
	)
}

// getGRPCRateLimiterInterceptor returns a gRPC interceptor that limits the rate of requests that the server can process.
// Returns a gRPC ResourceExhausted error if the rate limit is exceeded.
func getGRPCRateLimiterInterceptor(limiter *rate.Limiter, logger *dependencies.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !limiter.Allow() {
			logger.Error("Rate limit exceeded!")
			return nil, status.Errorf(codes.ResourceExhausted, errs.ErrMsgRateLimitExceeded)
		}
		return handler(ctx, req)
	}
}
