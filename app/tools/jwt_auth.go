package tools

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

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

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// GenerateToken returns a JWT token with the given user id, username and role.
func (jwta jwtAuthenticator) GenerateToken(id int, username string, role core.Role) (string, error) {
	claims := jwta.newClaims(id, username, role, time.Now())

	token, err := jwt.NewWithClaims(jwta.signingMethod, claims).SignedString([]byte(jwta.secret))
	if err != nil {
		core.LogUnexpectedErr(err)
		return "", status.Errorf(codes.Internal, errs.AuthGeneratingToken, err)
	}

	return token, nil
}

// Returns a GRPC interceptor that validates a JWT token inside of the context.
func (jwta jwtAuthenticator) ValidateToken(ctx context.Context, req any, grpcInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	route := core.RouteNameFromGRPC(grpcInfo.FullMethod)

	bearer, err := jwta.getBearer(ctx)
	if err != nil {
		if core.AuthForRoute(route) == core.RouteAuthPublic {
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

	ctx = jwta.newContextWithUserInfo(ctx, claims.ID, claims.Username)

	return handler(ctx, req)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Checks if the user is allowed to access the route.
func (jwta jwtAuthenticator) canAccessRoute(route, userID string, role core.Role, req any) error {
	switch core.AuthForRoute(route) {

	case core.RouteAuthPublic:
		return nil

	case core.RouteAuthSelf:
		type PBReqWithUserID interface {
			GetUserId() int32
		}
		reqUserID := req.(PBReqWithUserID).GetUserId()
		if userID != strconv.Itoa(int(reqUserID)) {
			return status.Errorf(codes.PermissionDenied, errs.AuthUserIDInvalid)
		}

	case core.RouteAuthAdmin:
		if role != core.AdminRole {
			core.LogPotentialThreat("User " + userID + " tried to access admin route: " + route)
			return status.Errorf(codes.PermissionDenied, errs.AuthRoleInvalid)
		}

	default:
		core.LogWeirdBehaviour("Route unknown: " + route)
		return status.Errorf(codes.Unknown, errs.AuthRouteUnknown)
	}
	return nil
}

// Returns the authorization field on the Metadata that lives the context.
func (jwta jwtAuthenticator) getBearer(ctx context.Context) (string, error) {
	bearer, err := jwta.grpcMDMapFn(ctx).Get("authorization")
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, errs.AuthTokenNotFound)
	}
	if !strings.HasPrefix(bearer, "Bearer ") {
		core.LogWeirdBehaviour(errs.AuthTokenMalformed)
		return "", status.Errorf(codes.Unauthenticated, errs.AuthTokenMalformed)
	}
	return strings.TrimPrefix(bearer, "Bearer "), nil
}

// Parses the token object into *jwtClaims and validates said claims.
// Returns an error if claims are not valid.
func (jwta jwtAuthenticator) getClaims(bearer string) (*jwtClaims, error) {
	token, err := jwt.ParseWithClaims(bearer, &jwtClaims{}, jwta.keyFn)
	if err == nil && token != nil && token.Valid {
		if claims, ok := token.Claims.(*jwtClaims); ok && claims.Valid() == nil {
			return claims, nil
		}
	}
	return nil, status.Errorf(codes.Unauthenticated, errs.AuthTokenInvalid)
}

func (jwta jwtAuthenticator) newClaims(id int, username string, role core.Role, issuedAt time.Time) *jwtClaims {
	return &jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwta.expAtFn(issuedAt),
			ID:        strconv.Itoa(id),
		},
		Username: username,
		Role:     role,
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns a new context with the user ID and username inside.
func (jwta jwtAuthenticator) newContextWithUserInfo(ctx context.Context, userID, username string) context.Context {
	ctx = context.WithValue(ctx, &ContextKeyUserID{}, userID)
	ctx = context.WithValue(ctx, &ContextKeyUsername{}, username)
	return ctx
}

type (
	// These are used as keys to store the user ID and username in the context (it's a key-value store).
	ContextKeyUserID   struct{}
	ContextKeyUsername struct{}
)

var defaultGRPCMetadataMapFn = NewGRPCMetadataMap

// Gets the key to decrypt the token. Key = JWT Secret.
var defaultKeyFn = func(key string) func(*jwt.Token) (any, error) {
	return func(*jwt.Token) (any, error) { return []byte(key), nil }
}

var defaultExpiresAtFn = func(sessionDays int) func(time.Time) *jwt.NumericDate {
	return func(iAt time.Time) *jwt.NumericDate {
		duration := time.Hour * 24 * time.Duration(sessionDays)
		return jwt.NewNumericDate(iAt.Add(duration))
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
