package tools

import (
	"context"
	"strings"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/logs"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ core.TokenValidator = &jwtValidator{}

/* ———————————————————————————————— — — — JWT TOKEN VALIDATOR — — — ———————————————————————————————— */

type jwtValidator struct {
	ctxTool core.ContextManager
	keyFn   func(*jwt.Token) (any, error)
	apiKey  string
}

func NewJWTValidator(ctxTool core.ContextManager, secret, apiKey string) core.TokenValidator {
	return &jwtValidator{
		ctxTool: ctxTool,
		keyFn: func(*jwt.Token) (any, error) {
			return []byte(secret), nil
		},
		apiKey: apiKey,
	}
}

// Validates a JWT Token against the Route to be accessed. Returns the Claims if valid, or a GRPC error if not.
// Errors returned can be Unauthenticated, PermissionDenied or NotFound.
// TODO — Change how this all works, it's breaking SRP.
func (v *jwtValidator) ValidateToken(ctx context.Context, req any, route shared.Route) (core.Claims, error) {
	if route.Auth == shared.RouteAuthAPIKey {
		apiKey, err := v.getAPIKeyFromCtx(ctx)
		if err != nil {
			return nil, err
		}

		if apiKey == v.apiKey { // TODO — Improve, hash, etc.
			return &shared.JWTClaims{}, nil
		}

		return nil, status.Errorf(codes.PermissionDenied, errs.AuthAPIKeyInvalid)
	}

	bearer, err := v.getBearerFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := v.getClaimsFromBearer(bearer)
	if err != nil {
		return nil, err
	}

	if err := route.CanBeAccessed(claims, req); err != nil {
		return nil, err
	}

	return claims, nil
}

// Returns the authorization field from the data that lives in the request's context.
func (v *jwtValidator) getBearerFromCtx(ctx context.Context) (string, error) {
	bearer, err := v.ctxTool.GetFromCtxMD(ctx, "authorization")
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, errs.AuthTokenNotFound)
	}

	if !strings.HasPrefix(bearer, "Bearer ") {
		logs.LogStrange(errs.AuthTokenMalformed)
		return "", status.Errorf(codes.Unauthenticated, errs.AuthTokenMalformed)
	}

	return strings.TrimPrefix(bearer, "Bearer "), nil
}

// Extracts the JWT Claims from the JWT token string.
// Returns an error if claims are not valid.
func (v *jwtValidator) getClaimsFromBearer(bearer string) (*shared.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(bearer, &shared.JWTClaims{}, v.keyFn)
	if err == nil && token != nil && token.Valid {
		if claims, ok := token.Claims.(*shared.JWTClaims); ok && claims.Valid() == nil {
			return claims, nil
		}
	}
	return nil, status.Errorf(codes.Unauthenticated, errs.AuthTokenInvalid)
}

// Returns the API Key from the data that lives in the request's context.
func (v *jwtValidator) getAPIKeyFromCtx(ctx context.Context) (string, error) {
	apiKey, err := v.ctxTool.GetFromCtxMD(ctx, "x-api-key")
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, errs.AuthAPIKeyNotFound)
	}
	return apiKey, nil
}
