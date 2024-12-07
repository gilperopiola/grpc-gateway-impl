package core

import (
	"strconv"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 🔑 RouteAuthPublic can be accessed by anyone.
// 🔑 RouteAuthSelf can only be accessed by the user with the same ID as the one specified on the request URL.
// The PB auto-generated requests for these routes MUST include a UserId int32 field, to do that on
// the .proto request definition we just add
//
//	▶ int32 user_id = 1 [(buf.validate.field).int32.gt = 0, (google.api.field_behavior) = REQUIRED];
//
// 🔑 RouteAuthAdmin can only be accessed by users with the Admin role.
// 🔑 RouteAuthAPIKey can only be accessed by valid API key holders.
func AccessRoute(route Route, claims *JWTClaims, req any) error {
	authNeeded := route.Auth

	if authNeeded == RouteAuthPublic {
		return nil
	}

	if authNeeded == RouteAuthSelf {
		// Compare the UserID from the request URL with the one from the claims.
		// They should match.
		urlUserID := int(req.(PBReqWithUserID).GetUserId())
		if strconv.Itoa(urlUserID) != claims.Subject {
			return status.Errorf(codes.PermissionDenied, errs.AuthUserIDInvalid)
		}
		return nil
	}

	if authNeeded == RouteAuthAdmin {
		if claims.Role != models.AdminRole {
			// logs.LogThreat("User " + claims.Subject + " tried to access admin route " + r.Name)
			return status.Errorf(codes.PermissionDenied, errs.AuthRoleInvalid)
		}
		return nil
	}

	// logs.LogStrange("Auth for route " + r.Name + " unhandled")
	return status.Errorf(codes.NotFound, errs.AuthRouteInvalid)
}

// All Protobuf requests with a userID on the URL should implement this.
type PBReqWithUserID interface {
	GetUserId() int32
}

/* ———————————————————————————————— — — — AUTH REQUIRED PER ROUTE — — — ———————————————————————————————— */

type AuthMethod string

const (
	RouteAuthInvalid AuthMethod = "invalid"
	RouteAuthPublic  AuthMethod = "public"
	RouteAuthUser    AuthMethod = "user"
	RouteAuthSelf    AuthMethod = "self"
	RouteAuthAdmin   AuthMethod = "admin"
	RouteAuthAPIKey  AuthMethod = "key"
)

/* ———————————————————————————————— — — — JWT CLAIMS — — — ———————————————————————————————— */

// These are the claims that live encrypted on our JWT Tokens.
// A JWT Token String, when decoded, returns one of this.
type JWTClaims struct {
	jwt.RegisteredClaims
	Username string          `json:"username"`
	Role     models.UserRole `json:"role"`
}

func (c *JWTClaims) GetUserInfo() (string, string) {
	return c.Subject, c.Username
}
