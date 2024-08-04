package models

import "github.com/golang-jwt/jwt/v4"

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

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// These are the claims that live encrypted on our JWT Tokens.
// A JWT Token String, when decoded, returns one of this.
type Claims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	Role     Role   `json:"role"`
}

type TokenClaims interface {
	GetUserInfo() (id, username string)
}

func (c *Claims) GetUserInfo() (string, string) {
	return c.ID, c.Username
}
