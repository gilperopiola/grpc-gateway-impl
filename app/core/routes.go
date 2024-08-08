package core

import (
	"strings"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"

	"google.golang.org/grpc"
)

// Our Routes need manual updating when a .proto route changes.
// TODO - Generate this based on the proto.

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
// We used to have a way of mapping the HTTP route to this,
// I don't know why I scraped it. Think I want it now :(
var Routes = map[string]Route{

	// Health Service
	"CheckHealth": {models.RouteAuthPublic},

	// Auth Service
	"Signup": {models.RouteAuthPublic},
	"Login":  {models.RouteAuthPublic},

	// Users Service
	"GetUser":     {models.RouteAuthSelf},
	"UpdateUser":  {models.RouteAuthSelf},
	"DeleteUser":  {models.RouteAuthSelf},
	"GetMyGroups": {models.RouteAuthSelf},

	"GetUsers": {models.RouteAuthAdmin},

	// Groups Service
	"GetGroup": {models.RouteAuthUser},

	"CreateGroup":       {models.RouteAuthSelf},
	"InviteToGroup":     {models.RouteAuthSelf},
	"AnswerGroupInvite": {models.RouteAuthSelf},
}

// Simple enough.
func RouteExists(routeName string) bool {
	if _, ok := Routes[routeName]; ok {
		return true
	}
	LogStrange("Route not found: " + routeName)
	return false
}

// This doesn't return any error on a not-found route, it just
// creates a log and defaults to AuthAdmin.
func AuthForRoute(routeName string) models.RouteAuth {
	if route, ok := Routes[routeName]; ok {
		return route.Auth
	}
	LogStrange("Route with Auth not found: " + routeName)
	return models.RouteAuthAdmin
}

// Our Routes are named by the last part of their GRPC Method.
// It's everything after the last slash.
//
//	Method = /pbs.Service/Signup
//	Route  = Signup
func RouteNameFromGRPC(method string) string {
	i := strings.LastIndex(method, "/")
	if i == -1 {
		LogStrange("No '/' found in GRPC Method " + method)
		return ""
	}
	return method[i+1:]
}

// Returns the route name from the context's data.
func RouteNameFromCtx(ctx god.Ctx) string {
	if method, ok := grpc.Method(ctx); ok {
		return RouteNameFromGRPC(method)
	}
	LogStrange("No GRPC Method found in context")
	return ""
}
