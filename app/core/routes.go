package core

import (
	"context"
	"strings"

	"google.golang.org/grpc"
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

func (r Route) CanBeAccessed(claims *JWTClaims, req any) error {
	return AccessRoute(r, claims, req)
}

func (r Route) GetName() string {
	return r.Name
}

func (r Route) GetAuth() AuthMethod {
	return r.Auth
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

var InvalidRoute = Route{"Invalid", RouteAuthInvalid}
