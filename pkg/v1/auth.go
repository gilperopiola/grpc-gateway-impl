package v1

import (
	"fmt"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/grpc"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(id int, username, email string, role grpc.Role, secret string, sessionDays int) (string, error) {

	var (
		issuedAt  = time.Now()
		expiresAt = time.Now().Add(time.Hour * 24 * time.Duration(sessionDays))
	)

	// Generate claims containing Username, Email, Role and ID
	claims := &grpc.CustomClaims{
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        fmt.Sprint(id),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	// Generate token (object)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate token (string)
	return token.SignedString([]byte(secret))
}

var (
	pathUserIDKey = "user_id"
)
