package common

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// PublicMethods are the methods that do not require authentication.
var PublicMethods = []string{
	"/users.UsersService/Signup",
	"/users.UsersService/Login",
}

// SelfAuthMethods are the methods that only allow the user to access their own data.
var SelfAuthMethods = []string{
	"/users.UsersService/GetUser",
}

// AdminMethods are the methods that only allow admins to access them.
var AdminMethods = []string{
	"/users.UsersService/GetUsers",
}

/* ----------------------------------- */
/*            - JWT Auth -             */
/* ----------------------------------- */

// Claims are our custom claims that encompass the standard JWT claims and also our own.
type Claims struct {
	Username string      `json:"username"`
	Role     models.Role `json:"role"`
	jwt.RegisteredClaims
}

// Authenticator is the interface that wraps the TokenGenerator and TokenValidator interfaces.
type Authenticator interface {
	TokenGenerator
	TokenValidator
}

// TokenGenerator is an interface with a method to generate an authentication token.
type TokenGenerator interface {
	Generate(id int, username string, role models.Role) (string, error)
}

// TokenValidator is an interface with a method to validate if a user is allowed to access certain endpoint.
type TokenValidator interface {
	Validate(ctx context.Context, req interface{}, svInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

// jwtAuthenticator implements both TokenGenerator and TokenValidator.
type jwtAuthenticator struct {
	secret      string
	sessionDays int
}

// NewJWTAuthenticator returns a new JWT authenticator with the given secret and session days.
func NewJWTAuthenticator(secret string, sessionDays int) *jwtAuthenticator {
	return &jwtAuthenticator{secret: secret, sessionDays: sessionDays}
}

/* ----------------------------------- */
/*         - Generate Token -          */
/* ----------------------------------- */

// Generate returns a JWT token with the given user id, username and role.
func (a *jwtAuthenticator) Generate(id int, username string, role models.Role) (string, error) {

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

// RequestWithUserID is an interface that lets us abstract .pb request types that have a GetUserId method.
type RequestWithUserID interface {
	GetUserId() int32
}

// Validate returns a gRPC interceptor that validates the JWT token from the context.
func (a *jwtAuthenticator) Validate(ctx context.Context, req interface{}, svInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	grpcMethod := svInfo.FullMethod

	// If the method does not require authentication, skip validation.
	if slices.Contains(PublicMethods, grpcMethod) {
		return handler(ctx, req)
	}

	// Get the token from the gRPC metadata.
	bearer, err := a.getBearer(ctx)
	if err != nil {
		return nil, err
	}

	// Get the claims from the token and validate them.
	claims, err := a.getClaims(bearer)
	if err != nil {
		return nil, err
	}

	if err := a.checkPermissions(claims.ID, claims.Role, grpcMethod, req); err != nil {
		return nil, err
	}

	// Add the user info to the context.
	ctx = a.addInfoToCtx(ctx, claims)

	return handler(ctx, req)
}

// UserIDKey and UsernameKey are used to store the user id and username in the context.
type UserIDKey struct{}
type UsernameKey struct{}

/* ----------------------------------- */
/*        - Helper Functions -         */
/* ----------------------------------- */

// getBearer returns the token from the gRPC metadata from the context.
func (a *jwtAuthenticator) getBearer(ctx context.Context) (string, error) {
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

// getClaims parses the token into Claims and then validates them.
func (a *jwtAuthenticator) getClaims(bearer string) (*Claims, error) {
	jwtToken, err := jwt.ParseWithClaims(bearer, &Claims{}, a.keyFunc)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "auth: token invalid")
	}
	claims, ok := jwtToken.Claims.(*Claims)
	if !ok || !jwtToken.Valid || claims.Valid() != nil {
		return nil, status.Errorf(codes.Unauthenticated, "auth: token invalid")
	}
	return claims, nil
}

// checkPermissions checks if the user is allowed to access the method.
func (a *jwtAuthenticator) checkPermissions(userID string, userRole models.Role, grpcMethod string, req interface{}) error {
	checkUserID := slices.Contains(SelfAuthMethods, grpcMethod)
	checkIsAdmin := slices.Contains(AdminMethods, grpcMethod)

	// If the method only allows the user to access their own data, check if the JWT User ID is the same as the one on the request.
	if checkUserID && fmt.Sprint(req.(RequestWithUserID).GetUserId()) != userID {
		return status.Errorf(codes.PermissionDenied, "auth: user id invalid")
	}

	// If the method only allows admins, check if the user is an admin.
	if checkIsAdmin && userRole != models.AdminRole {
		return status.Errorf(codes.PermissionDenied, "auth: role invalid")
	}

	return nil
}

// addInfoToCtx adds the user id and username to the context.
func (a *jwtAuthenticator) addInfoToCtx(c context.Context, claims *Claims) context.Context {
	c = context.WithValue(c, &UserIDKey{}, claims.ID)
	return context.WithValue(c, &UsernameKey{}, claims.Username)
}

// keyFunc returns the key for validating JWT tokens.
func (a *jwtAuthenticator) keyFunc(_ *jwt.Token) (interface{}, error) {
	return []byte(a.secret), nil
}
