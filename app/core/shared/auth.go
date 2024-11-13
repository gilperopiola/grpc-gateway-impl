package shared

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/models"
	"github.com/golang-jwt/jwt/v4"
)

// These are the claims that live encrypted on our JWT Tokens.
// A JWT Token String, when decoded, returns one of this. And we
// create a new one each time we generate a token.
type JWTClaims struct {
	jwt.RegisteredClaims
	Username string          `json:"username"`
	Role     models.UserRole `json:"role"`
}

func (c *JWTClaims) GetUserInfo() (string, string) {
	return c.Subject, c.Username
}
