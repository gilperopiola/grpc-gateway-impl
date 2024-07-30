package core

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"

	"google.golang.org/grpc"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - Routes -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// As we use grpc-gateway, for each endpoint we have 2 versions: GRPC and HTTP.

// Routes are the representation of a single endpoint, and it holds all the things that
// both versions of said endpoint should share.
//
// This is the place to code behaviour that operates based on each route, like the auth level.
// We could have rate-limiting per route, a pool of connections, etc.
type Route struct {
	Auth models.RouteAuth
}

// And this is the map of all the routes in our app.
//
// So, when we talk about a route, we usually talk about the key string in this map.
// Each one of these keys correspond to the last part of each GRPC method, the part after the last '/'.
// For example, the route for the method /pbs.Svc/Login is "Login".
//
// This needs manual updating when the .proto routes change. TODO - Generate this based on the proto.
// We used to have a way of mapping the HTTP route to this, but I don't know why I scraped it. Think I want it now :(
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

// In GRPC, method means the entire route name.
// In HTTP it's just GET, POST, etc.

func AuthForRoute(routeName string) models.RouteAuth {
	if route, ok := Routes[routeName]; ok {
		return route.Auth
	}
	LogWeirdBehaviour("Route not found: " + routeName)
	return models.RouteAuthAdmin
}

func CanAccessRoute(route, claimsUserID string, claimsRole models.Role, req any) error {
	switch AuthForRoute(route) {

	case models.RouteAuthPublic:
		return nil

	case models.RouteAuthAdmin:
		if claimsRole != models.AdminRole {
			LogPotentialThreat("User " + claimsUserID + " tried to access admin route: " + route)
			return errors.New(errs.AuthRoleInvalid)
		}

	case models.RouteAuthSelf:
		// Requests for routes with this Auth type must have an int32 UserID field.
		type PBReqWithUserID interface {
			GetUserId() int32
		}

		// Compare the UserID from the request with the one from the claims.
		// They should match.
		reqUserID := int(req.(PBReqWithUserID).GetUserId())
		if strconv.Itoa(reqUserID) != claimsUserID {
			return errors.New(errs.AuthUserIDInvalid)
		}

	default:
		LogWeirdBehaviour("Route unknown: " + route)
		return errors.New(errs.AuthRouteUnknown)
	}

	return nil
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// A Route is the last part of a GRPC Method, the part after the last slash.
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
