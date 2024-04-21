package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - JWT Auth -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var _ core.TokenAuthenticator = (*jwtAuthenticator)(nil) // jwtAuthenticator implements TokenAuthenticator.

type jwtAuthenticator struct {
	secret        string
	sessionDays   int
	signingMethod jwt.SigningMethod
	grpcMDMapFn   func(any) DataMap
	keyFn         func(*jwt.Token) (any, error)
	expAtFn       func(iAt time.Time) *jwt.NumericDate
}

// Our custom claims, the standard JWT RegisteredClaims + our own.
type jwtClaims struct {
	jwt.RegisteredClaims
	Username string      `json:"username"`
	Role     models.Role `json:"role"`
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

// Generate returns a JWT token with the given user id, username and role.
func (a *jwtAuthenticator) Generate(id int, username string, role models.Role) (string, error) {
	claims := a.NewClaims(id, username, role)
	tokenString, err := jwt.NewWithClaims(a.signingMethod, claims).SignedString([]byte(a.secret))
	if err != nil {
		go core.LogUnexpected(err)
		return "", status.Errorf(codes.Internal, "auth: error generating token: %v", err) // T0D0 move to errs.
	}
	return tokenString, nil
}

// newClaims have inside the RegisteredClaims (with ID and dates), as well as the Username and Role.
func (a *jwtAuthenticator) NewClaims(id int, username string, role models.Role) *jwtClaims {
	return &jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        fmt.Sprint(id),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: a.expAtFn(time.Now()),
		},
		Username: username,
		Role:     role,
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - Validate Token -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Validate returns a gRPC interceptor that validates the JWT token from the context.
func (a *jwtAuthenticator) Validate(ctx context.Context, req interface{}, svInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	route := core.GetRouteFromGRPC(svInfo.FullMethod)

	bearer, err := a.GetBearer(ctx)
	if err != nil {
		if core.Routes[route].Auth == core.RouteAuthPublic {
			return handler(ctx, req)
		}
		return nil, err
	}

	claims, err := a.GetClaims(bearer)
	if err != nil {
		return nil, err
	}

	if err := a.CanAccessRoute(route, claims.ID, claims.Role, req); err != nil {
		return nil, err
	}

	ctx = a.CtxWithUserInfo(ctx, claims.ID, claims.Username)

	return handler(ctx, req)
}

func (a *jwtAuthenticator) GetBearer(ctx context.Context) (string, error) { // From ctx -> get md -> get bearer
	bearer, err := a.grpcMDMapFn(ctx).Get("authorization")
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "auth: token not found")
	}
	if !strings.HasPrefix(bearer, "Bearer ") {
		core.LogWeirdBehaviour("token bearer malformed: " + bearer)
		return "", status.Errorf(codes.Unauthenticated, "auth: token malformed")
	}
	return strings.TrimPrefix(bearer, "Bearer "), nil
}

// Parses the token object into *jwtClaims and validates them.
// Returns the claims if valid, or an error if not.
func (a *jwtAuthenticator) GetClaims(bearer string) (*jwtClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(bearer, &jwtClaims{}, a.keyFn)
	if err == nil && jwtToken != nil && jwtToken.Valid {
		if claims, ok := jwtToken.Claims.(*jwtClaims); ok && claims.Valid() == nil {
			return claims, nil
		}
	}
	core.LogUnexpected(err)
	return nil, status.Errorf(codes.Unauthenticated, "auth: token invalid")
}

// Checks if the user is allowed to access the gRPC method.
func (a *jwtAuthenticator) CanAccessRoute(route, userID string, role models.Role, req any) error {
	switch core.Routes[route].Auth {

	case core.RouteAuthPublic:
		return nil

	case core.RouteAuthSelf:
		type PBReqWithUserID interface{ GetUserId() int32 }
		if userID != fmt.Sprint(req.(PBReqWithUserID).GetUserId()) {
			return status.Errorf(codes.PermissionDenied, "auth: user id invalid")
		}

	case core.RouteAuthAdmin:
		if role != models.AdminRole {
			core.LogPotentialThreat(fmt.Sprintf("User %s tried to access admin route %s", userID, route))
			return status.Errorf(codes.PermissionDenied, "auth: role invalid")
		}

	default:
		core.LogWeirdBehaviour(fmt.Sprintf("Route unknown: %s", route))
		return status.Errorf(codes.Unknown, "auth: route unknown")
	}
	return nil
}

// Returns a new context with a user ID and username added.
func (a *jwtAuthenticator) CtxWithUserInfo(ctx context.Context, userID, username string) context.Context {
	ctx = context.WithValue(ctx, &CtxKeyUserID{}, userID)
	ctx = context.WithValue(ctx, &CtxKeyUsername{}, username)
	return ctx
}

// These are used as keys to store the user ID and username in the context (it's a key-value store).
type CtxKeyUserID struct{}
type CtxKeyUsername struct{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var defaultGRPCMetadataMapFn = NewGRPCMetadataMap

// Gets the key to decrypt the token. Key = JWT Secret.
var defaultKeyFn = func(key string) func(_ *jwt.Token) (any, error) {
	return func(_ *jwt.Token) (any, error) { return []byte(key), nil }
}

var defaultExpiresAtFn = func(sessDays int) func(time.Time) *jwt.NumericDate {
	return func(issuedAt time.Time) *jwt.NumericDate {
		duration := time.Hour * 24 * time.Duration(sessDays)
		return jwt.NewNumericDate(issuedAt.Add(duration))
	}
}
