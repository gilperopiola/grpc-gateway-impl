package core

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

// All data that we need on a per-route basis lives here. For now it's just the Auth type.
type Route struct {
	Name string
	Auth RouteAuth
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - Routes -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// In GRPC, 'method' means the entire 'endpoint' name. In HTTP it's just GET, POST, etc.
//
// -> We're calling both of them 'routes'.
// -> Each route is just the last part of the corresponding GRPC method.
// -> So '/pbs.Service/Login' becomes 'Login' and that is the route for the Login endpoint.
// -> For HTTP, we need to get the route name from the HTTP Path, so we have a map for that.

var Routes = map[string]Route{
	"Signup": {
		Name: "Signup",
		Auth: RouteAuthPublic,
	},
	"Login": {
		Name: "Login",
		Auth: RouteAuthPublic,
	},
	"GetUser": {
		Name: "GetUser",
		Auth: RouteAuthSelf,
	},
	"GetUsers": {
		Name: "GetUsers",
		Auth: RouteAuthAdmin,
	},
	"CreateGroup": {
		Name: "CreateGroup",
		Auth: RouteAuthUser,
	},
	"GetGroup": {
		Name: "GetGroup",
		Auth: RouteAuthUser,
	},
}

// Use this to get the route name from the HTTP Path.
var RouteNamesFromHTTP = map[string]string{
	"POST /v1/signup": "Signup",
	"POST /v1/login":  "Login",
	"GET /v1/user":    "GetUser",
	"GET /v1/users":   "GetUsers",
	"POST /v1/groups": "CreateGroup",
	"GET /v1/group":   "GetGroup",
}

func AuthForRoute(routeName string) RouteAuth {
	if route, ok := Routes[routeName]; ok {
		return route.Auth
	}
	LogWeirdBehaviour("Route not found: " + routeName)
	return RouteAuthAdmin
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// The Route is the last part of the GRPC Method.
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

// The Route is not directly linked to the HTTP Path, so we need to get it from a map.
func RouteNameFromHTTP(req *http.Request) string {
	httpPath := req.Method + " " + req.URL.Path // -> e.g. 'GET /users/1' or 'POST /signup'

	lastSlashIndex := strings.LastIndex(httpPath, "/")
	if lastSlashIndex == -1 {
		LogWeirdBehaviour("No '/' found in HTTP: " + httpPath)
		return httpPath
	}

	httpPathLastPart := httpPath[lastSlashIndex+1:] // -> e.g. '1' or 'signup'

	// If the last part of the HTTP Path is not a number, then it's a full path and we can use it directly to get the
	// route name from the map. Otherwise, we need to remove the number from the end of the path, so '/users/1' becomes just '/users'
	// and then we can get the route name from the map.
	if _, err := strconv.Atoi(httpPathLastPart); err != nil {
		// -> httpPath = 'POST /signup'
		// -> httpPathLastPart = 'signup' -> Not a number.
		return RouteNamesFromHTTP[httpPath]
	}
	// -> httpPath = 'GET /users/1'
	// -> httpPathLastPart = '1' -> Number.
	// -> httpPathWithoutLastPart = 'GET /users'
	httpPathWithoutLastPart := httpPath[:lastSlashIndex-1]
	return RouteNamesFromHTTP[httpPathWithoutLastPart]
}

// Returns the route name from the context.
func RouteNameFromCtx(ctx context.Context) string {
	if method, ok := grpc.Method(ctx); ok {
		return RouteNameFromGRPC(method)
	}
	LogWeirdBehaviour("No GRPC Method found in context")
	return ""
}
