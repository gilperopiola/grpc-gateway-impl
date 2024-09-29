package tools

import (
	"strconv"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ core.TokenGenerator = &jwtGenerator{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - JWT Token Generator -       */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type jwtGenerator struct {
	secret          string
	sessionDuration time.Duration
	signingMethod   jwt.SigningMethod
}

func NewJWTGenerator(secret string, sessionDays int) core.TokenGenerator {
	return &jwtGenerator{
		secret:          secret,
		signingMethod:   jwt.SigningMethodHS256,
		sessionDuration: time.Hour * 24 * time.Duration(sessionDays),
	}
}

// GenerateToken returns a JWT token with the given user id, username and role.
func (g *jwtGenerator) GenerateToken(id int, username string, role shared.Role) (string, error) {
	claims := g.newClaims(id, username, role)

	token, err := jwt.NewWithClaims(g.signingMethod, claims).SignedString([]byte(g.secret))
	if err != nil {
		logs.LogUnexpected(err)
		return "", status.Errorf(codes.Internal, errs.AuthGeneratingToken, err)
	}

	return token, nil
}

func (g *jwtGenerator) newClaims(id int, username string, role shared.Role) *shared.JWTClaims {
	now := time.Now()
	return &shared.JWTClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(id),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(g.sessionDuration)),
		},
	}
}
