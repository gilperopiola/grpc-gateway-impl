package modules

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

/* ----------------------------------- */
/*        - gRPC Methods Auth -        */
/* ----------------------------------- */

type EndpointAuthType string

const (
	AuthTypePublic EndpointAuthType = "public"
	AuthTypeSelf   EndpointAuthType = "self"
	AuthTypeAdmin  EndpointAuthType = "admin"
)

var authTypePerGRPCMethod = map[string]EndpointAuthType{
	"Signup": AuthTypePublic,
	"Login":  AuthTypePublic,

	"GetUser":  AuthTypeSelf,
	"GetUsers": AuthTypeAdmin,
}

/* ----------------------------------- */
/*            - JWT Auth -             */
/* ----------------------------------- */

// TokenAuthenticator is the interface that wraps the TokenGenerator and TokenValidator interfaces.
type TokenAuthenticator interface {
	TokenGenerator
	TokenValidator
}

// TokenGenerator is an interface with a method to generate an authentication token.
type TokenGenerator interface {
	Generate(id int, username string, role models.Role) (string, error)
}

// TokenValidator is an interface with a method to validate if a user is allowed to access certain endpoint.
type TokenValidator interface {
	Validate(ctx context.Context, req interface{}, svI *grpc.UnaryServerInfo, h grpc.UnaryHandler) (interface{}, error)
}

// jwtAuthenticator implements both TokenGenerator and TokenValidator.
type jwtAuthenticator struct {
	jwtSecret string
	jwtDays   int

	jwtSignMethod    jwt.SigningMethod
	jwtKeyFunc       func(_ *jwt.Token) (interface{}, error)
	jwtExpiresAtFunc func(issuedAt time.Time) *jwt.NumericDate
}

// jwtClaims are our custom claims that encompass the standard JWT RegisteredClaims and also our own.
type jwtClaims struct {
	jwt.RegisteredClaims
	Username string      `json:"username"`
	Role     models.Role `json:"role"`
}

// NewJWTAuthenticator returns a new JWT authenticator with the given secret and session days.
func NewJWTAuthenticator(secret string, sessionDays int) *jwtAuthenticator {
	var (
		signingMethod = jwt.SigningMethodHS256
		keyFunc       = func(_ *jwt.Token) (interface{}, error) { return []byte(secret), nil }
		expiresAtFunc = func(issuedAt time.Time) *jwt.NumericDate {
			tokenDuration := time.Hour * 24 * time.Duration(sessionDays)
			return jwt.NewNumericDate(issuedAt.Add(tokenDuration))
		}
	)

	return &jwtAuthenticator{
		secret, sessionDays, signingMethod,
		keyFunc, expiresAtFunc,
	}
}

/* ----------------------------------- */
/*         - Generate Token -          */
/* ----------------------------------- */

// Generate returns a JWT token with the given user id, username and role.
func (a *jwtAuthenticator) Generate(id int, username string, role models.Role) (string, error) {
	claims := a.NewClaims(id, username, role).(*jwtClaims)
	tokenString, err := jwt.NewWithClaims(a.jwtSignMethod, claims).SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", status.Errorf(codes.Internal, "auth: error generating token: %v", err) // T0D0 move to errs.
	}
	return tokenString, nil
}

/* ----------------------------------- */
/*         - Validate Token -          */
/* ----------------------------------- */

// Validate returns a gRPC interceptor that validates the JWT token from the context.
func (a *jwtAuthenticator) Validate(ctx context.Context, req interface{}, svInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	grpcMethod := svInfo.FullMethod

	bearer, err := a.GetBearer(ctx) // Get token from metadata
	if errors.Is(err, errTokenMalformed) || (errors.Is(err, errTokenNotFound) && authTypePerGRPCMethod[grpcMethod] != AuthTypePublic) {
		return nil, err
	}

	untypedClaims, err := a.GetClaims(bearer) // Get claims from token
	if err != nil {
		return nil, err
	}

	claims := untypedClaims.(*jwtClaims)
	if err := a.IsMethodAllowed(grpcMethod, claims.ID, claims.Role, req); err != nil {
		return nil, err
	}

	ctx = a.CtxWithUserInfo(ctx, claims.ID, claims.Username)

	return handler(ctx, req)
}

/* ----------------------------------- */
/*        - Helper Functions -         */
/* ----------------------------------- */

// IsMethodAllowed checks if the user is allowed to access the gRPC method.
func (a *jwtAuthenticator) IsMethodAllowed(grpcMethod, userID string, userRole models.Role, req interface{}) error {
	type PBReqWithUserID interface{ GetUserId() int32 }

	switch authTypePerGRPCMethod[grpcMethod] {
	case AuthTypePublic:
		return nil
	case AuthTypeSelf:
		if userID != fmt.Sprint(req.(PBReqWithUserID).GetUserId()) {
			return status.Errorf(codes.PermissionDenied, "auth: user id invalid")
		}
	case AuthTypeAdmin:
		if userRole != models.AdminRole {
			return status.Errorf(codes.PermissionDenied, "auth: role invalid")
		}
	default:
		return status.Errorf(codes.Unknown, "auth: method unknown")
	}
	return nil
}

// GetBearer returns the token from the gRPC metadata from the context.
func (a *jwtAuthenticator) GetBearer(ctx context.Context) (string, error) {
	authMD := metadata.ValueFromIncomingContext(ctx, "authorization")
	if len(authMD) == 0 || authMD[0] == "" {
		return "", errTokenNotFound
	}
	if !strings.HasPrefix(authMD[0], "Bearer ") {
		return "", errTokenMalformed
	}
	return strings.TrimPrefix(authMD[0], "Bearer "), nil
}

// GetClaims parses the token into Claims and then validates them.
func (a *jwtAuthenticator) GetClaims(bearer string) (any, error) {
	if jwtToken, err := jwt.ParseWithClaims(bearer, &jwtClaims{}, a.jwtKeyFunc); err == nil && jwtToken != nil && jwtToken.Valid {
		if claims, ok := jwtToken.Claims.(*jwtClaims); ok && claims.Valid() == nil {
			return claims, nil
		}
	}
	return nil, status.Errorf(codes.Unauthenticated, "auth: token invalid")
}

// newClaims have inside the RegisteredClaims (with ID and dates), as well as the Username and Role.
func (a *jwtAuthenticator) NewClaims(id int, username string, role models.Role) any {
	return &jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        fmt.Sprint(id),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: a.jwtExpiresAtFunc(time.Now()),
		},
		Username: username,
		Role:     role,
	}
}

// ContextWithUserInfo adds the user id and username to the context.
func (a *jwtAuthenticator) CtxWithUserInfo(c context.Context, userID string, username string) context.Context {
	c = context.WithValue(c, &CtxKeyUserID{}, userID)
	return context.WithValue(c, &CtxKeyUsername{}, username)
}

// CtxKeyUserID and CtxKeyUsername are used to store the user ID and username in the context.
type CtxKeyUserID struct{}
type CtxKeyUsername struct{}

var errTokenNotFound = fmt.Errorf("auth: token not found")
var errTokenMalformed = fmt.Errorf("auth: token malformed")
