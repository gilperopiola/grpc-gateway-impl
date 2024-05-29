package toolbox

import (
	"strconv"
	"strings"
	"time"

	"github.com/gilperopiola/god"
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
func (g jwtGenerator) GenerateToken(id int, username string, role core.Role) (string, error) {
	var (
		now    = time.Now()
		claims = &jwtClaims{
			jwt.RegisteredClaims{
				ID:        strconv.Itoa(id),
				IssuedAt:  jwt.NewNumericDate(now),
				ExpiresAt: g.expAtFn(now),
			},
			username,
			role,
		}
	)

	token, err := jwt.NewWithClaims(g.signingMethod, claims).SignedString([]byte(g.secret))
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
	ctxManager core.CtxManager
	keyFn      func(*jwt.Token) (any, error)
}

func NewJWTValidator(ctxManager core.CtxManager, secret string) core.TokenValidator {
	return &jwtValidator{
		ctxManager: ctxManager,
		keyFn:      defaultKeyFn(secret),
	}
}

// Returns a GRPC interceptor that validates a JWT token inside of the context.
func (v jwtValidator) ValidateToken(ctx god.Ctx, req any, grpcInfo *god.GRPCInfo, handler god.GRPCHandler) (any, error) {
	route := core.RouteNameFromGRPC(grpcInfo.FullMethod)

	bearer, err := v.getBearer(ctx)
	if err != nil {
		if core.AuthForRoute(route) == core.RouteAuthPublic {
			// If failed to get token but route is public: OK!
			return handler(ctx, req)
		}
		return nil, err
	}

	claims, err := v.getClaims(bearer)
	if err != nil {
		return nil, err
	}

	if err := core.CanAccessRoute(route, claims.ID, claims.Role, req); err != nil {
		return nil, err
	}

	ctx = v.ctxManager.AddUserInfo(ctx, claims.ID, claims.Username)

	return handler(ctx, req)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns the authorization field on the Metadata that lives the context.
func (v jwtValidator) getBearer(ctx god.Ctx) (string, error) {
	bearer, err := v.ctxManager.ExtractMetadata(ctx, "authorization")
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
func (v jwtValidator) getClaims(bearer string) (*jwtClaims, error) {
	token, err := jwt.ParseWithClaims(bearer, &jwtClaims{}, v.keyFn)
	if err == nil && token != nil && token.Valid {
		if claims, ok := token.Claims.(*jwtClaims); ok && claims.Valid() == nil {
			return claims, nil
		}
	}
	return nil, status.Errorf(codes.Unauthenticated, errs.AuthTokenInvalid)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

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
