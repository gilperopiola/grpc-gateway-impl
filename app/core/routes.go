package core

import "strings"

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - Routes -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// In GRPC, 'method' means the entire 'endpoint' name. In HTTP it's just GET, POST, etc.
//
// -> We're calling both of them 'routes'.
// -> Each route is just the last part of the corresponding GRPC method.
// -> So '/pbs.Service/Login' becomes 'Login' and that is the route for the Login endpoint.

var Routes = map[string]RouteInfo{
	"Signup":   {Auth: RouteAuthPublic},
	"Login":    {Auth: RouteAuthPublic},
	"GetUser":  {Auth: RouteAuthSelf},
	"GetUsers": {Auth: RouteAuthAdmin},
}

// All data that we need on a per-route basis lives here. For now it's just the Auth type.
type RouteInfo struct {
	Auth RouteAuth
}

// Route is the last part of the GRPC Method.
//
// -> Method = /pbs.Service/Signup
// -> Route  = Signup
//
func GetRouteFromGRPC(method string) string {
	i := strings.LastIndex(method, "/")
	if i == -1 {
		LogWeirdBehaviour("No '/' found in GRPC: " + method)
		return method
	}
	return method[i+1:]
}
