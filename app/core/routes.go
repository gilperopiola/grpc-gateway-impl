package core

import (
	"strconv"
	"strings"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// All settings we need on a per-route basis lives here. For now it's just the Auth type.
type Route struct {
	Auth models.RouteAuth
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - Routes -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// In GRPC, 'method' means the entire 'endpoint' name. In HTTP it's just GET, POST, etc.
//
// -> We're calling both of them 'routes'.
// -> Each route is just the last part of the corresponding GRPC method.
// -> So '/pbs.Service/Login' becomes 'Login' and that is the route for the Login endpoint.
// -> And HTTP calls GRPC, so we're covered.

// T0D0 generate this based on the .proto file.
var Routes = map[string]Route{

	// Auth Service
	"Signup": {models.RouteAuthPublic},
	"Login":  {models.RouteAuthPublic},

	// Users Service
	"GetUser":     {models.RouteAuthSelf},
	"GetUsers":    {models.RouteAuthAdmin},
	"UpdateUser":  {models.RouteAuthSelf},
	"DeleteUser":  {models.RouteAuthSelf},
	"GetMyGroups": {models.RouteAuthSelf},

	// Groups Service
	"CreateGroup":       {models.RouteAuthSelf},
	"GetGroup":          {models.RouteAuthUser},
	"InviteToGroup":     {models.RouteAuthSelf},
	"AnswerGroupInvite": {models.RouteAuthSelf},
}

func AuthForRoute(routeName string) models.RouteAuth {
	if route, ok := Routes[routeName]; ok {
		return route.Auth
	}
	LogWeirdBehaviour("Route not found: " + routeName)
	return models.RouteAuthAdmin
}

func CanAccessRoute(route, userID string, role models.Role, req any) error {
	switch AuthForRoute(route) {

	case models.RouteAuthPublic:
		return nil

	case models.RouteAuthSelf:
		type PBReqWithUserID interface {
			GetUserId() int32
		}
		reqUserID := req.(PBReqWithUserID).GetUserId()
		if userID != strconv.Itoa(int(reqUserID)) {
			return status.Errorf(codes.PermissionDenied, errs.AuthUserIDInvalid)
		}

	case models.RouteAuthAdmin:
		if role != models.AdminRole {
			LogPotentialThreat("User " + userID + " tried to access admin route: " + route)
			return status.Errorf(codes.PermissionDenied, errs.AuthRoleInvalid)
		}

	default:
		LogWeirdBehaviour("Route unknown: " + route)
		return status.Errorf(codes.Unknown, errs.AuthRouteUnknown)
	}

	return nil
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// A Route is the last part of a GRPC Method.
//
// -> Method = /pbs.Service/Signup
// -> Route  = Signup
func RouteNameFromGRPC(method string) string {
	i := strings.LastIndex(method, "/")
	if i == -1 {
		LogWeirdBehaviour("No '/' found in GRPC: " + method)
		return method
	}
	return method[i+1:]
}

// Returns the route name from the context.
func RouteNameFromCtx(ctx god.Ctx) string {
	if method, ok := grpc.Method(ctx); ok {
		return RouteNameFromGRPC(method)
	}
	LogWeirdBehaviour("No GRPC Method found in context")
	return ""
}
