package core

import "github.com/gilperopiola/grpc-gateway-impl/app/core/shared"

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
	Auth shared.RouteAuth
}

// This doesn't return any error on a not-found route, it just
// defaults to AuthAdmin.
func AuthForRoute(routeName string) shared.RouteAuth {
	if route, ok := Routes[routeName]; ok {
		return route.Auth
	}
	return shared.RouteAuthAdmin
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
	"CheckHealth": {shared.RouteAuthPublic},

	// Auth Service
	"Signup": {shared.RouteAuthPublic},
	"Login":  {shared.RouteAuthPublic},

	// Users Service
	"GetUser":     {shared.RouteAuthSelf},
	"UpdateUser":  {shared.RouteAuthSelf},
	"DeleteUser":  {shared.RouteAuthSelf},
	"GetMyGroups": {shared.RouteAuthSelf},
	"GetUsers":    {shared.RouteAuthAdmin},

	// Groups Service
	"GetGroup":          {shared.RouteAuthUser},
	"CreateGroup":       {shared.RouteAuthSelf},
	"InviteToGroup":     {shared.RouteAuthSelf},
	"AnswerGroupInvite": {shared.RouteAuthSelf},

	// GPT Service
	"NewGPTChat":     {shared.RouteAuthPublic},
	"ReplyToGPTChat": {shared.RouteAuthPublic},
}
