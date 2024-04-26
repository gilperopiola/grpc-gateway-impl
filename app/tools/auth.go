package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - JWT Auth -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var _ core.TokenAuthenticator = (*jwtAuthenticator)(nil)

type jwtAuthenticator struct {
	secret        string
	sessionDays   int
	signingMethod jwt.SigningMethod
	grpcMDMapFn   func(any) DataMap
	keyFn         func(*jwt.Token) (any, error)
	expAtFn       func(iAt time.Time) *jwt.NumericDate
}

// Our custom claims = the standard JWT RegisteredClaims + our own.
type jwtClaims struct {
	jwt.RegisteredClaims
	Username string    `json:"username"`
	Role     core.Role `json:"role"`
}

// New JWT authenticator with the given secret and duration.
func NewJWTAuthenticator(secret string, sessionDays int) *jwtAuthenticator {
	return &jwtAuthenticator{
		secret:        secret,
		sessionDays:   sessionDays,
		signingMethod: jwt.SigningMethodHS256,
		grpcMDMapFn:   defaultGRPCMetadataMapFn,
		keyFn:         defaultKeyFn(secret),
		expAtFn:       defaultExpiresAtFn(sessionDays),
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - Generate Token -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// GenerateToken returns a JWT token with the given user id, username and role.
func (a *jwtAuthenticator) GenerateToken(id int, username string, role core.Role) (string, error) {
	claims := a.NewClaims(id, username, role)
	tokenString, err := jwt.NewWithClaims(a.signingMethod, claims).SignedString([]byte(a.secret))
	if err != nil {
		core.LogUnexpectedErr(err)
		return "", status.Errorf(codes.Internal, "auth: error generating token: %v", err) // T0D0 move to errs.
	}
	return tokenString, nil
}

// newClaims have inside the RegisteredClaims (with ID and dates), as well as the Username and Role.
func (a *jwtAuthenticator) NewClaims(id int, username string, role core.Role) *jwtClaims {
	now := time.Now()
	return &jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        strconv.Itoa(id),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: a.expAtFn(now),
		},
		Username: username,
		Role:     role,
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - Validate Token -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns a GRPC interceptor that validates a JWT token inside of the context.
func (a *jwtAuthenticator) ValidateToken(ctx context.Context, req interface{}, svInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	route := core.GetRouteFromGRPC(svInfo.FullMethod)

	bearer, err := a.getBearer(ctx)
	if err != nil {
		if core.Routes[route].Auth == core.RouteAuthPublic {
			return handler(ctx, req)
		}
		return nil, err
	}

	claims, err := a.getClaims(bearer)
	if err != nil {
		return nil, err
	}

	if err := a.canAccessRoute(route, claims.ID, claims.Role, req); err != nil {
		return nil, err
	}

	ctx = a.contextWithUserInfo(ctx, claims.ID, claims.Username)

	return handler(ctx, req)
}

// Returns the authorization Metadata from the context.
func (a *jwtAuthenticator) getBearer(ctx context.Context) (string, error) {
	bearer, err := a.grpcMDMapFn(ctx).Get("authorization")
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "auth: token not found")
	}
	if !strings.HasPrefix(bearer, "Bearer ") {
		core.LogWeirdBehaviour(msgTokenMalformed())
		return "", status.Errorf(codes.Unauthenticated, "auth: token malformed")
	}
	return strings.TrimPrefix(bearer, "Bearer "), nil
}

// Parses the token object into *jwtClaims and validates them.
// Returns the claims if valid, or an error if not.
func (a *jwtAuthenticator) getClaims(bearer string) (*jwtClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(bearer, &jwtClaims{}, a.keyFn)
	if err == nil && jwtToken != nil && jwtToken.Valid {
		if claims, ok := jwtToken.Claims.(*jwtClaims); ok && claims.Valid() == nil {
			return claims, nil
		}
	}
	return nil, status.Errorf(codes.Unauthenticated, "auth: token invalid")
}

// Checks if the user is allowed to access the route.
func (a *jwtAuthenticator) canAccessRoute(route, userID string, role core.Role, req any) error {
	switch core.Routes[route].Auth {

	case core.RouteAuthPublic:
		return nil

	case core.RouteAuthSelf:
		type PBReqWithUserID interface{ GetUserId() int32 }
		if userID != fmt.Sprint(req.(PBReqWithUserID).GetUserId()) {
			return status.Errorf(codes.PermissionDenied, "auth: user id invalid")
		}

	case core.RouteAuthAdmin:
		if role != core.AdminRole {
			core.LogPotentialThreat(msgAdminUnauthorized(userID, route))
			return status.Errorf(codes.PermissionDenied, "auth: role invalid")
		}

	default:
		core.LogWeirdBehaviour(msgRouteUnknown(route))
		return status.Errorf(codes.Unknown, "auth: route unknown")
	}
	return nil
}

// Returns a new context with the user ID and username inside.
func (a *jwtAuthenticator) contextWithUserInfo(ctx context.Context, userID, username string) context.Context {
	ctx = context.WithValue(ctx, &ContextKeyUserID{}, userID)
	ctx = context.WithValue(ctx, &ContextKeyUsername{}, username)
	return ctx
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var defaultGRPCMetadataMapFn = NewGRPCMetadataMap

// Gets the key to decrypt the token. Key = JWT Secret.
var defaultKeyFn = func(key string) func(*jwt.Token) (any, error) {
	return func(*jwt.Token) (any, error) { return []byte(key), nil }
}

var defaultExpiresAtFn = func(days int) func(time.Time) *jwt.NumericDate {
	return func(iAt time.Time) *jwt.NumericDate {
		duration := time.Hour * 24 * time.Duration(days)
		return jwt.NewNumericDate(iAt.Add(duration))
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// These are used as keys to store the user ID and username in the context (it's a key-value store).
type (
	ContextKeyUserID   struct{}
	ContextKeyUsername struct{}
)

// Unexpected log messages.
var (
	msgTokenMalformed    = func() string { return "Token malformed" }
	msgRouteUnknown      = func(route string) string { return fmt.Sprintf("Route unknown: %s", route) }
	msgAdminUnauthorized = func(userID, route string) string {
		return fmt.Sprintf("User %s tried to access admin route: %s", userID, route)
	}
)
