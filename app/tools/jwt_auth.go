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

var _ core.TokenAuthenticator = (*jwtAuthenticator)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - JWT Auth -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type jwtAuthenticator struct {
	secret        string
	sessionDays   int
	signingMethod jwt.SigningMethod
	grpcMDMapFn   func(any) DataMap
	keyFn         func(*jwt.Token) (any, error)
	expAtFn       func(iAt time.Time) *jwt.NumericDate
}

// Our jwtClaims -> Standard JWT RegisteredClaims + Username + Role.
type jwtClaims struct {
	jwt.RegisteredClaims
	Username string    `json:"username"`
	Role     core.Role `json:"role"`
}

func NewJWTAuthenticator(secret string, sessionDays int) core.TokenAuthenticator {
	return &jwtAuthenticator{
		secret:        secret,
		sessionDays:   sessionDays,
		signingMethod: jwt.SigningMethodHS256,
		grpcMDMapFn:   defaultGRPCMetadataMapFn,
		keyFn:         defaultKeyFn(secret),
		expAtFn:       defaultExpiresAtFn(sessionDays),
	}
}

func (jwta jwtAuthenticator) GetTokenAuthenticator() core.TokenAuthenticator {
	return jwta
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// GenerateToken returns a JWT token with the given user id, username and role.
func (jwta jwtAuthenticator) GenerateToken(id int, username string, role core.Role) (string, error) {
	claims := jwta.newClaims(id, username, role, time.Now())

	token, err := jwt.NewWithClaims(jwta.signingMethod, claims).SignedString([]byte(jwta.secret))
	if err != nil {
		core.LogUnexpectedErr(err)
		return "", status.Errorf(codes.Internal, "auth: error generating token: %v", err) // T0D0 move to errs.
	}

	return token, nil
}

// Returns a GRPC interceptor that validates a JWT token inside of the context.
func (jwta jwtAuthenticator) ValidateToken(ctx context.Context, req interface{}, svInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	route := core.RouteNameFromGRPC(svInfo.FullMethod)

	bearer, err := jwta.getBearer(ctx)
	if err != nil {
		if core.Routes[route].Auth == core.RouteAuthPublic {
			return handler(ctx, req)
		}
		return nil, err
	}

	claims, err := jwta.getClaims(bearer)
	if err != nil {
		return nil, err
	}

	if err := jwta.canAccessRoute(route, claims.ID, claims.Role, req); err != nil {
		return nil, err
	}

	ctx = jwta.contextWithUserInfo(ctx, claims.ID, claims.Username)

	return handler(ctx, req)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Checks if the user is allowed to access the route.
func (jwta jwtAuthenticator) canAccessRoute(route, userID string, role core.Role, req any) error {
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

// Returns the authorization Metadata from the context.
func (jwta jwtAuthenticator) getBearer(ctx context.Context) (string, error) {
	bearer, err := jwta.grpcMDMapFn(ctx).Get("authorization")
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
func (jwta jwtAuthenticator) getClaims(bearer string) (*jwtClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(bearer, &jwtClaims{}, jwta.keyFn)
	if err == nil && jwtToken != nil && jwtToken.Valid {
		if claims, ok := jwtToken.Claims.(*jwtClaims); ok && claims.Valid() == nil {
			return claims, nil
		}
	}
	return nil, status.Errorf(codes.Unauthenticated, "auth: token invalid")
}

func (jwta jwtAuthenticator) newClaims(id int, username string, role core.Role, issuedAt time.Time) *jwtClaims {
	return &jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        strconv.Itoa(id),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwta.expAtFn(issuedAt),
		},
		Username: username,
		Role:     role,
	}
}

// Returns a new context with the user ID and username inside.
func (jwta jwtAuthenticator) contextWithUserInfo(ctx context.Context, userID, username string) context.Context {
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
