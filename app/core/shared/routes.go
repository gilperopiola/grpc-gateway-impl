package shared

import (
	"strings"

	"github.com/gilperopiola/god"
	"google.golang.org/grpc"
)

// 🔻🔻🔻 Routes 🔻🔻🔻

// ⭐ As we use grpc-gateway, for each endpoint we have 2 versions: GRPC and HTTP.

// Routes are the representation of a single endpoint, and it holds all the things that
// both versions of said endpoint should share.
//
// This is the place to code behaviour that operates based on each route, like the auth level.
// We could have rate-limiting per route, a pool of connections, etc.
type Route struct {
	Auth RouteAuth
}

// 🔻 All the routes in our app.
//
// So, when we talk about a route, we usually talk about the key string in this map.
// Each one of these keys correspond to the last part of each GRPC method, the part after the last '/'.
//
//	-> The route for the GRPC method '/pbs.SomeSvc/Login' is 'Login'.
var Routes = map[string]Route{

	// 🚑 Health Service
	"CheckHealth": {RouteAuthPublic},

	// 🔒 Auth Service
	"Signup": {RouteAuthPublic},
	"Login":  {RouteAuthPublic},

	// 😎 Users Service
	"GetUser":     {RouteAuthSelf},
	"UpdateUser":  {RouteAuthSelf},
	"DeleteUser":  {RouteAuthSelf},
	"GetMyGroups": {RouteAuthSelf},
	"GetUsers":    {RouteAuthAdmin},

	// 👨‍👨‍👧‍👦 Groups Service
	"GetGroup":          {RouteAuthUser},
	"CreateGroup":       {RouteAuthSelf},
	"InviteToGroup":     {RouteAuthSelf},
	"AnswerGroupInvite": {RouteAuthSelf},

	// 🤖 GPT Service
	"NewGPTChat":     {RouteAuthPublic},
	"ReplyToGPTChat": {RouteAuthPublic},
}

// TODO -> We should autogenerate this from the proto files.

// 🔻🔻🔻 Auth Required per Route 🔻🔻🔻

type RouteAuth string

const (
	RouteAuthPublic RouteAuth = "public"
	RouteAuthUser   RouteAuth = "user"
	RouteAuthSelf   RouteAuth = "self"
	RouteAuthAdmin  RouteAuth = "admin"
	RouteAPIKey     RouteAuth = "apiKey"
)

// Defaults to AuthAdmin for unknown routes.
func AuthForRoute(routeName string) RouteAuth {
	if route, ok := Routes[routeName]; ok {
		return route.Auth
	}
	return RouteAuthAdmin
}

// 🔻🔻🔻 Get Current Route Name 🔻🔻🔻

// Our Routes are named by the last part of their GRPC Method.
// It's everything after the last slash.
//
//	Method = /pbs.Service/Signup
//	Route  = Signup
func RouteNameFromGRPCMethod(method string) string {
	i := strings.LastIndex(method, "/")
	if i == -1 {
		return ""
	}
	return method[i+1:]
}

// Returns the route name from the context's data.
func RouteNameFromCtx(ctx god.Ctx) string {
	if method, ok := grpc.Method(ctx); ok {
		return RouteNameFromGRPCMethod(method)
	}
	return ""
}
