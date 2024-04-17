package modules

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

// jwtClaims are our custom claims that encompass the standard JWT RegisteredClaims and also our own.
type jwtClaims struct {
	jwt.RegisteredClaims
	Username string      `json:"username"`
	Role     models.Role `json:"role"`
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - JWT Auth -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// jwtAuthenticator implements TokenAuthenticator.
type jwtAuthenticator struct {
	secret          string
	sessionDays     int
	signingMethod   jwt.SigningMethod
	keyFn           func(*jwt.Token) (any, error)
	expiresAtFn     func(issuedAt time.Time) *jwt.NumericDate
	headersAccessor func(any) KeyValStoreAccessor
}

// NewJWTAuthenticator returns a new JWT authenticator with the given secret and session days.
func NewJWTAuthenticator(secret string, sessDays int) *jwtAuthenticator {
	return &jwtAuthenticator{
		secret:          secret,
		sessionDays:     sessDays,
		signingMethod:   jwt.SigningMethodHS256,
		keyFn:           defaultKeyFn(secret),
		expiresAtFn:     defaultExpiresAtFn(sessDays),
		headersAccessor: NewGRPCMetadataAccessor,
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
			ExpiresAt: a.expiresAtFn(time.Now()),
		},
		Username: username,
		Role:     role,
	}
}

var defaultExpiresAtFn = func(sessDays int) func(issuedAt time.Time) *jwt.NumericDate {
	return func(issuedAt time.Time) *jwt.NumericDate {
		duration := time.Hour * 24 * time.Duration(sessDays)
		return jwt.NewNumericDate(issuedAt.Add(duration))
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
	authHeader, err := a.headersAccessor(ctx).Get("authorization")
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "auth: token not found")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", status.Errorf(codes.Unauthenticated, "auth: token malformed")
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

// GetClaims parses the token into Claims and then validates them.
func (a *jwtAuthenticator) GetClaims(bearer string) (*jwtClaims, error) {
	if jwtToken, err := jwt.ParseWithClaims(bearer, &jwtClaims{}, a.keyFn); err == nil && jwtToken != nil && jwtToken.Valid {
		if claims, ok := jwtToken.Claims.(*jwtClaims); ok && claims.Valid() == nil {
			return claims, nil
		}
	}
	return nil, status.Errorf(codes.Unauthenticated, "auth: token invalid")
}

// CanAccessRoute checks if the user is allowed to access the gRPC method.
func (a *jwtAuthenticator) CanAccessRoute(route, userID string, role models.Role, req interface{}) error {
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
			return status.Errorf(codes.PermissionDenied, "auth: role invalid")
		}
	default:
		return status.Errorf(codes.Unknown, "auth: route unknown")
	}
	return nil
}

// ContextWithUserInfo adds the user id and username to the context.
func (a *jwtAuthenticator) CtxWithUserInfo(c context.Context, userID, username string) context.Context {
	c = context.WithValue(c, &CtxKeyUserID{}, userID)
	c = context.WithValue(c, &CtxKeyUsername{}, username)
	return c
}

// CtxKeyUserID and CtxKeyUsername are used to store the user ID and username in the context.
type CtxKeyUserID struct{}
type CtxKeyUsername struct{}

var defaultKeyFn = func(secret string) func(_ *jwt.Token) (any, error) {
	return func(_ *jwt.Token) (any, error) { return []byte(secret), nil }
}
