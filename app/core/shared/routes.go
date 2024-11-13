package shared

import (
	"context"
	"strconv"
	"strings"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ⭐ As we use grpc-gateway, for each endpoint we have 2 versions: GRPC and HTTP.

/* ———————————————————————————————— — — — ROUTES — — — ———————————————————————————————— */

// Routes are the representation of a single endpoint, and it holds all the things that
// both versions of said endpoint should share.
//
// This is the place to code behaviour that operates based on each route, like the auth level.
// We could have rate-limiting per route, a pool of connections, etc.
type Route struct {
	Name string
	Auth AuthMethod
}

// Map of routes.
//
// When we talk about a route, we usually talk about a pair of endpoints: one for GRPC and another one HTTP.
// The key of this map corresponds to the route name. It is derived from the last part of each GRPC method,
// the part after the last '/':
//
//	Method = /pbs.Service/Signup
//	Route  = Signup
//
// TODO — autogenerate this from the proto files.
var Routes = map[string]Route{

	// 🚑 Health Service
	"CheckHealth": {"CheckHealth", RouteAuthPublic},

	// 🔒 Auth Service
	"Signup": {"Signup", RouteAuthPublic},
	"Login":  {"Login", RouteAuthPublic},

	// 😎 Users Service
	"GetUser":     {"GetUser", RouteAuthSelf},
	"UpdateUser":  {"UpdateUser", RouteAuthSelf},
	"DeleteUser":  {"DeleteUser", RouteAuthSelf},
	"GetMyGroups": {"GetMyGroups", RouteAuthSelf},
	"GetUsers":    {"GetUsers", RouteAuthAdmin},

	// 👨‍👨‍👧‍👦 Groups Service
	"GetGroup":          {"GetGroup", RouteAuthUser},
	"CreateGroup":       {"CreateGroup", RouteAuthSelf},
	"InviteToGroup":     {"InviteToGroup", RouteAuthSelf},
	"AnswerGroupInvite": {"AnswerGroupInvite", RouteAuthSelf},

	// 🤖 GPT Service
	"NewGPTChat":     {"NewGPTChat", RouteAuthPublic},
	"ReplyToGPTChat": {"ReplyToGPTChat", RouteAuthPublic},
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

// 🔑 RouteAuthPublic can be accessed by anyone.
// 🔑 RouteAuthAdmin can only be accessed by users with the Admin role.
// 🔑 RouteAuthSelf can only be accessed by the user with the same ID as the one specified on the request URL.
// The PB auto-generated requests for these routes MUST include a UserId int32 field, to do that on
// the .proto request definition we just add
//
//	▶ int32 user_id = 1 [(buf.validate.field).int32.gt = 0, (google.api.field_behavior) = REQUIRED];
func (r Route) CanBeAccessed(claims *JWTClaims, req any) error {
	if r.Auth == RouteAuthPublic {
		return nil
	}

	if r.Auth == RouteAuthSelf {
		// Compare the UserID from the request URL with the one from the claims.
		// They should match.
		urlUserID := int(req.(PBReqWithUserID).GetUserId())
		if strconv.Itoa(urlUserID) != claims.Subject {
			return status.Errorf(codes.PermissionDenied, errs.AuthUserIDInvalid)
		}
		return nil
	}

	if r.Auth == RouteAuthAdmin {
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

var InvalidRoute = Route{"Invalid", RouteAuthInvalid}

/* ———————————————————————————————— — — — GET REQUEST'S ROUTE — — — ———————————————————————————————— */

// Our Routes are named by the last part of their GRPC Method.
// It's everything after the last slash.
//
//	Method = /pbs.Service/Signup
//	Route  = Signup
func GetRouteFromGRPCMethod(method string) Route {
	lastSlashIndex := strings.LastIndex(method, "/")
	if lastSlashIndex != -1 {
		return Routes[method[lastSlashIndex+1:]]
	}
	return InvalidRoute
}

// Gets the GRPC Method from the context and the Route from that Method.
func GetRouteFromCtx(ctx context.Context) Route {
	if method, ok := grpc.Method(ctx); ok {
		return GetRouteFromGRPCMethod(method)
	}
	return InvalidRoute
}
