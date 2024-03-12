package dependencies

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

/* ----------------------------------- */
/*            - JWT Auth -             */
/* ----------------------------------- */

type Role string

const (
	AnyRole     Role = "any"
	DefaultRole Role = "default"
	AdminRole   Role = "admin"
)

type Claims struct {
	Username string `json:"username"`
	Role     Role   `json:"role"`
	jwt.RegisteredClaims
}

type TokenGenerator interface {
	Generate(id int, username string, role Role) (string, error)
}

type TokenValidator interface {
	Validate(expectedRole Role) grpc.UnaryServerInterceptor

	getBearerFromCtx(ctx context.Context) (string, error)
	getClaimsFromBearer(bearer string, expectedRole Role) (*Claims, error)
	addClaimsInfoToCtx(c context.Context, claims *Claims) context.Context
	keyFunc(_ *jwt.Token) (interface{}, error)
}

type jwtAuthenticator struct {
	secret      string
	sessionDays int
}

func NewJWTAuthenticator(secret string, sessionDays int) *jwtAuthenticator {
	return &jwtAuthenticator{secret: secret, sessionDays: sessionDays}
}

/* ----------------------------------- */
/*         - Generate Token -          */
/* ----------------------------------- */

func (a *jwtAuthenticator) Generate(id int, username string, role Role) (string, error) {

	// New claims with Username, Role and ID.
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        fmt.Sprint(id),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(a.sessionDays))),
		},
		Username: username,
		Role:     role,
	}

	// Generate JWT token and get signed string from it.
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(a.secret))
	if err != nil {
		return "", status.Errorf(codes.Internal, "auth: error generating token: %v", err)
	}

	return tokenString, nil
}

/* ----------------------------------- */
/*         - Validate Token -          */
/* ----------------------------------- */

// RequestWithUserID is an interface that lets us use protobuf request types that have a GetUserId method.
type RequestWithUserID interface {
	GetUserId() int32
}

var methodsWithoutAuth = []string{
	"/users.UsersService/Signup",
	"/users.UsersService/Login",
}

var methodsWithSelfAuth = []string{
	"/users.UsersService/GetUser",
}

func (a *jwtAuthenticator) Validate(expectedRole Role) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, svInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		grpcMethod := svInfo.FullMethod

		// If the method does not require authentication, skip validation.
		if slices.Contains(methodsWithoutAuth, grpcMethod) {
			return handler(ctx, req)
		}

		// Get the token from the gRPC metadata.
		bearer, err := a.getBearerFromCtx(ctx)
		if err != nil {
			return nil, err
		}

		// Get the claims from the token and validate them.
		claims, err := a.getClaimsFromBearer(bearer, expectedRole)
		if err != nil {
			return nil, err
		}

		// If the method only allows the user to access their own data, check if the JWT User ID is the same as the one on the request.
		if slices.Contains(methodsWithSelfAuth, grpcMethod) && fmt.Sprint(req.(RequestWithUserID).GetUserId()) != claims.ID {
			return nil, status.Errorf(codes.PermissionDenied, "auth: user id invalid")
		}

		// Add the user info to the context.
		ctx = a.addClaimsInfoToCtx(ctx, claims)

		return handler(ctx, req)
	}
}

// getBearerFromCtx returns the token from the gRPC metadata from the context.
func (a *jwtAuthenticator) getBearerFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Internal, "auth: metadata not found")
	}

	authHeaders := md["authorization"]
	if len(authHeaders) == 0 || len(authHeaders[0]) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "auth: token not found")
	}

	if !strings.HasPrefix(authHeaders[0], "Bearer ") {
		return "", status.Errorf(codes.Unauthenticated, "auth: token invalid format")
	}

	return strings.TrimPrefix(authHeaders[0], "Bearer "), nil
}

// getClaimsFromBearer parses the token into Claims and then validates them.
func (a *jwtAuthenticator) getClaimsFromBearer(bearer string, expectedRole Role) (*Claims, error) {
	jwtToken, err := jwt.ParseWithClaims(bearer, &Claims{}, a.keyFunc)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "auth: token invalid")
	}
	claims, ok := jwtToken.Claims.(*Claims)
	if !ok || !jwtToken.Valid || claims.Valid() != nil {
		return nil, status.Errorf(codes.Unauthenticated, "auth: token invalid")
	}
	if expectedRole != AnyRole && claims.Role != expectedRole {
		return nil, status.Errorf(codes.Unauthenticated, "auth: role invalid")
	}
	return claims, nil
}

type UserIDKey struct{}
type UsernameKey struct{}

// addClaimsInfoToCtx adds the user id and username to the context.
func (a *jwtAuthenticator) addClaimsInfoToCtx(c context.Context, claims *Claims) context.Context {
	c = context.WithValue(c, &UserIDKey{}, claims.ID)
	return context.WithValue(c, &UsernameKey{}, claims.Username)
}

// keyFunc returns the key for validating JWT tokens.
func (a *jwtAuthenticator) keyFunc(_ *jwt.Token) (interface{}, error) {
	return []byte(a.secret), nil
}
