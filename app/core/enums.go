package core

type Role string

const (
	DefaultRole Role = "default"
	AdminRole   Role = "admin"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type RouteAuth string

const (
	RouteAuthPublic RouteAuth = "public"
	RouteAuthSelf   RouteAuth = "self"
	RouteAuthAdmin  RouteAuth = "admin"
)
