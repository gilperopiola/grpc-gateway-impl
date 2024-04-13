package core

import "strings"

/* ----------------------------------- */
/*              - Routes -             */
/* ----------------------------------- */

// Ugh we have this horrible conflict between gRPC and HTTP's naming conventions.
// In gRPC, Method means the entire endpoint name. In HTTP it's just GET, POST, etc.
//
// We're using 'route' for both to avoid confusion.

var Routes = map[string]RouteInfo{
	"Signup": {RouteAuthPublic},
	"Login":  {RouteAuthPublic},

	"GetUser":  {RouteAuthSelf},
	"GetUsers": {RouteAuthAdmin},
}

type RouteInfo struct {
	Auth RouteAuth
}

type RouteAuth string

const (
	RouteAuthPublic RouteAuth = "public"
	RouteAuthSelf   RouteAuth = "self"
	RouteAuthAdmin  RouteAuth = "admin"
)

// Route equals the last part of the gRPC Method.
//
// - Method = /pbs.Service/Login
// - Route  = Login
//
func GetRouteFromGRPC(method string) string {
	index := strings.LastIndex(method, "/")
	if index == -1 {
		return method
	}
	return method[index+1:]
}
