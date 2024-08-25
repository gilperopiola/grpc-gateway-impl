package tools

import (
	"context"
	"strconv"
	"strings"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/types/models"

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
		keyFn: func(*jwt.Token) (any, error) {
			return []byte(secret), nil
		},
	}
}

// Validates a JWT Token. Returns the Claims if valid, or a GRPC error if not.
// Errors returned can be Unauthenticated, PermissionDenied or Unknown.
func (v jwtValidator) ValidateToken(ctx context.Context, req any, route string) (models.TokenClaims, error) {
	bearer, err := v.getBearer(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := v.getClaims(bearer)
	if err != nil {
		return nil, err
	}

	if err := v.canAccessRoute(route, claims, req); err != nil {
		return claims, err
	}

	return claims, nil
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns the authorization field from the data that lives in the request's context.
func (v jwtValidator) getBearer(ctx god.Ctx) (string, error) {
	bearer, err := v.ctxTool.GetFromCtx(ctx, "authorization")
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, errs.AuthTokenNotFound)
	}

	if !strings.HasPrefix(bearer, "Bearer ") {
		logs.LogStrange(errs.AuthTokenMalformed)
		return "", status.Errorf(codes.Unauthenticated, errs.AuthTokenMalformed)
	}

	return strings.TrimPrefix(bearer, "Bearer "), nil
}

// Parses the token string into a *models.Claims.
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

// Determines if a set of Claims can access certain route with certain request.
func (v jwtValidator) canAccessRoute(route string, claims *models.Claims, req any) error {
	switch core.AuthForRoute(route) {

	// These routes only allow the user with the same ID as the one specified on the request to go through.
	case models.RouteAuthSelf:

		// Requests for routes with this Auth type must have an int32 UserID field.
		type PBReqWithUserID interface {
			GetUserId() int32
		}

		// Compare the UserID from the request with the one from the claims.
		// They should match.
		reqUserID := int(req.(PBReqWithUserID).GetUserId())
		if strconv.Itoa(reqUserID) != claims.ID {
			return status.Errorf(codes.PermissionDenied, errs.AuthUserIDInvalid)
		}

	// These routes only allow admin users to go through.
	case models.RouteAuthAdmin:
		if claims.Role != models.AdminRole {
			logs.LogThreat("User " + claims.ID + " tried to access admin route " + route)
			return status.Errorf(codes.PermissionDenied, errs.AuthRoleInvalid)
		}

	// Everyone can access these routes.
	// This is the last option because the GRPC Token Validation Interceptor already checks for it.
	case models.RouteAuthPublic:
		return nil

	default:
		logs.LogStrange("Auth for route " + route + " unhandled")
		return status.Errorf(codes.Unknown, errs.AuthRouteInvalid)
	}

	return nil
}
