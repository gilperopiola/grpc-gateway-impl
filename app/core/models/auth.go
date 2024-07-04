package models

type Role string

const (
	DefaultRole Role = "default"
	AdminRole   Role = "admin"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type RouteAuth string

const (
	RouteAuthPublic RouteAuth = "public"
	RouteAuthUser   RouteAuth = "user"
	RouteAuthSelf   RouteAuth = "self"
	RouteAuthAdmin  RouteAuth = "admin"
)
