package tools

import (
	"context"
	"strings"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ core.TokenValidator = &jwtValidator{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - JWT Token Validator -      */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type jwtValidator struct {
	ctxTool core.CtxTool
	keyFn   func(*jwt.Token) (any, error)
}

func NewJWTValidator(ctxTool core.CtxTool, secret string) core.TokenValidator {
	return &jwtValidator{
		ctxTool: ctxTool,
		keyFn:   defaultKeyFn(secret),
	}
}

// Returns a GRPC interceptor that validates a JWT token inside of the context.
func (v jwtValidator) ValidateToken(ctx context.Context, req any, route string) (core.TokenClaims, error) {
	bearer, err := v.getBearer(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := v.getClaims(bearer)
	if err != nil {
		return nil, err
	}

	if err := core.CanAccessRoute(route, claims.ID, claims.Role, req); err != nil {
		return claims, err
	}

	return claims, nil
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns the authorization field on the Metadata that lives the context.
func (v jwtValidator) getBearer(ctx god.Ctx) (string, error) {
	bearer, err := v.ctxTool.GetMetadata(ctx, "authorization")
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, errs.AuthTokenNotFound)
	}

	if !strings.HasPrefix(bearer, "Bearer ") {
		core.LogWeirdBehaviour(errs.AuthTokenMalformed)
		return "", status.Errorf(codes.Unauthenticated, errs.AuthTokenMalformed)
	}

	return strings.TrimPrefix(bearer, "Bearer "), nil
}

// Parses the token object into *models.Claims and validates said claims.
// Returns an error if claims are not valid.
func (v jwtValidator) getClaims(bearer string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(bearer, &models.Claims{}, v.keyFn)
	if err == nil && token != nil && token.Valid {
		if claims, ok := token.Claims.(*models.Claims); ok && claims.Valid() == nil {
			return claims, nil
		}
	}
	return nil, status.Errorf(codes.Unauthenticated, errs.AuthTokenInvalid)
}

// Gets the key to decrypt the token.
// Key = JWT Secret.
var defaultKeyFn = func(key string) jwt.Keyfunc {
	return func(*jwt.Token) (any, error) { return []byte(key), nil }
}
