package shared

import "github.com/golang-jwt/jwt/v4"

/* -~-~-~-~ Auth Roles ~-~-~-~- */

type Role string

const (
	DefaultRole Role = "default"
	AdminRole   Role = "admin"
)

/* -~-~-~-~ Auth required per Route ~-~-~-~- */

type RouteAuth string

const (
	RouteAuthPublic RouteAuth = "public"
	RouteAuthUser   RouteAuth = "user"
	RouteAuthSelf   RouteAuth = "self"
	RouteAuthAdmin  RouteAuth = "admin"
)

/* -~-~-~-~ Auth Claims ~-~-~-~- */

type Claims interface {
	GetUserInfo() (id, username string)
}

// These are the claims that live encrypted on our JWT Tokens.
// A JWT Token String, when decoded, returns one of this. And we
// create a new one each time we generate a token.
type JWTClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	Role     Role   `json:"role"`
}

func (c *JWTClaims) GetUserInfo() (string, string) {
	return c.Subject, c.Username
}
