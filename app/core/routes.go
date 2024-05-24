package core

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

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

// T0D0 This isn't automatically updated when the .proto changes. Could be a good fit for code gen.
var Routes = map[string]Route{
	pbs.AuthService_ServiceDesc.Methods[0].MethodName: {
		Name: "Signup",
		Auth: RouteAuthPublic,
	},
	pbs.AuthService_ServiceDesc.Methods[1].MethodName: {
		Name: "Login",
		Auth: RouteAuthPublic,
	},

	pbs.UsersService_ServiceDesc.Methods[0].MethodName: {
		Name: "GetUser",
		Auth: RouteAuthSelf,
	},
	pbs.UsersService_ServiceDesc.Methods[1].MethodName: {
		Name: "GetUsers",
		Auth: RouteAuthAdmin,
	},
	pbs.UsersService_ServiceDesc.Methods[2].MethodName: {
		Name: "UpdateUser",
		Auth: RouteAuthSelf,
	},
	pbs.UsersService_ServiceDesc.Methods[3].MethodName: {
		Name: "DeleteUser",
		Auth: RouteAuthSelf,
	},
	pbs.UsersService_ServiceDesc.Methods[4].MethodName: {
		Name: "GetMyGroups",
		Auth: RouteAuthSelf,
	},

	pbs.GroupsService_ServiceDesc.Methods[0].MethodName: {
		Name: "CreateGroup",
		Auth: RouteAuthSelf,
	},
	pbs.GroupsService_ServiceDesc.Methods[1].MethodName: {
		Name: "GetGroup",
		Auth: RouteAuthUser,
	},
	pbs.GroupsService_ServiceDesc.Methods[2].MethodName: {
		Name: "InviteToGroup",
		Auth: RouteAuthSelf,
	},
	pbs.GroupsService_ServiceDesc.Methods[3].MethodName: {
		Name: "AnswerGroupInvite",
		Auth: RouteAuthSelf,
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
		LogUnexpectedErr(errors.New("No '/' found in GRPC: " + method))
		return method
	}
	return method[i+1:]
}

// The Route is not derived from the HTTP Path, we need to get it from a map.
func RouteNameFromHTTP(req *http.Request) string {
	httpPath := req.Method + " " + req.URL.Path // -> e.g. 'GET /users/1' or 'POST /signup'

	lastSlashIndex := strings.LastIndex(httpPath, "/")
	if lastSlashIndex == -1 {
		LogUnexpectedErr(errors.New("No '/' found in HTTP: " + httpPath))
		return httpPath
	}

	lastPart := httpPath[lastSlashIndex+1:] // -> e.g. '1' or 'signup'

	// If the last part of the HTTP Path is not a number, then it's a full path and we can use it directly to get the
	// route name from the map. Otherwise, we need to remove the number from the end of the path, so '/users/1' becomes just '/users'
	// and then we can get the route name from the map.
	if _, err := strconv.Atoi(lastPart); err != nil {
		// -> httpPath = 'POST /signup'
		// -> lastPart = 'signup' -> Not a number.
		return RouteNamesFromHTTP[httpPath]
	}
	// -> httpPath = 'GET /users/1'
	// -> lastPart = '1' -> Number.
	// -> pathWithoutLastPart = 'GET /users'
	pathWithoutLastPart := httpPath[:lastSlashIndex-1]
	return RouteNamesFromHTTP[pathWithoutLastPart]
}

// Returns the route name from the context.
func RouteNameFromCtx(ctx Ctx) string {
	if method, ok := grpc.Method(ctx); ok {
		return RouteNameFromGRPC(method)
	}
	LogWeirdBehaviour("No GRPC Method found in context")
	return ""
}
