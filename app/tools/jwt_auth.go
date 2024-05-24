package tools

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ core.TokenGenerator = (*jwtGenerator)(nil)
var _ core.TokenValidator = (*jwtValidator)(nil)

// Our JWT Tokens contain these Claims to identify the user
type jwtClaims struct {
	jwt.RegisteredClaims
	Username string    `json:"username"`
	Role     core.Role `json:"role"`
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - JWT Token Generator -       */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type jwtGenerator struct {
	secret        string
	sessionDays   int
	signingMethod jwt.SigningMethod
	expAtFn       func(iAt time.Time) *jwt.NumericDate
}

func NewJWTGenerator(secret string, sessionDays int) core.TokenGenerator {
	return &jwtGenerator{
		secret:        secret,
		sessionDays:   sessionDays,
		signingMethod: jwt.SigningMethodHS256,
		expAtFn:       defaultExpiresAtFn(sessionDays),
	}
}

// GenerateToken returns a JWT token with the given user id, username and role.
func (jwtgen jwtGenerator) GenerateToken(id int, username string, role core.Role) (string, error) {
	var (
		now    = time.Now()
		claims = &jwtClaims{
			jwt.RegisteredClaims{
				ID:        strconv.Itoa(id),
				IssuedAt:  jwt.NewNumericDate(now),
				ExpiresAt: jwtgen.expAtFn(now),
			},
			username,
			role,
		}
	)

	token, err := jwt.NewWithClaims(jwtgen.signingMethod, claims).SignedString([]byte(jwtgen.secret))
	if err != nil {
		core.LogUnexpectedErr(err)
		return "", status.Errorf(codes.Internal, errs.AuthGeneratingToken, err)
	}

	return token, nil
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - JWT Token Validator -      */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type jwtValidator struct {
	mdGetter           core.MetadataGetter
	routeAuthenticator core.RouteAuthenticator
	keyFn              func(*jwt.Token) (any, error)
}

func NewJWTValidator(mdGetter core.MetadataGetter, routeAuther core.RouteAuthenticator, secret string) core.TokenValidator {
	return &jwtValidator{
		mdGetter:           mdGetter,
		routeAuthenticator: routeAuther,
		keyFn:              defaultKeyFn(secret),
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns a GRPC interceptor that validates a JWT token inside of the context.
func (jwtval jwtValidator) ValidateToken(ctx core.Ctx, req any, grpcInfo *core.GRPCInfo, handler core.GRPCHandler) (any, error) {
	route := core.RouteNameFromGRPC(grpcInfo.FullMethod)

	bearer, err := jwtval.getBearer(ctx)
	if err != nil {
		if core.AuthForRoute(route) == core.RouteAuthPublic {
			// If failed to get token but route is public: OK!
			return handler(ctx, req)
		}
		return nil, err
	}

	claims, err := jwtval.getClaims(bearer)
	if err != nil {
		return nil, err
	}

	if err := jwtval.routeAuthenticator.CanAccessRoute(route, claims.ID, claims.Role, req); err != nil {
		return nil, err
	}

	ctx = jwtval.newContextWithUserInfo(ctx, claims.ID, claims.Username)

	return handler(ctx, req)
}

// Returns the authorization field on the Metadata that lives the context.
func (jwtval jwtValidator) getBearer(ctx core.Ctx) (string, error) {
	bearer, err := jwtval.mdGetter.GetMD(ctx, "authorization")
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
func (jwtval jwtValidator) getClaims(bearer string) (*jwtClaims, error) {
	token, err := jwt.ParseWithClaims(bearer, &jwtClaims{}, jwtval.keyFn)
	if err == nil && token != nil && token.Valid {
		if claims, ok := token.Claims.(*jwtClaims); ok && claims.Valid() == nil {
			return claims, nil
		}
	}
	return nil, status.Errorf(codes.Unauthenticated, errs.AuthTokenInvalid)
}

// Returns a new context with the user ID and username inside.
func (jwtval jwtValidator) newContextWithUserInfo(ctx core.Ctx, userID, username string) core.Ctx {
	ctx = context.WithValue(ctx, &ContextKeyUserID{}, userID)
	ctx = context.WithValue(ctx, &ContextKeyUsername{}, username)
	return ctx
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// These are used as keys to store the user ID and username in the context (it's a key-value store).
type (
	ContextKeyUserID   struct{}
	ContextKeyUsername struct{}
)

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
