package tests

import (
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/app/layers/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/modules"
	"github.com/gilperopiola/grpc-gateway-impl/etc/tests/mocks"

	"github.com/stretchr/testify/assert"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Tests: Service -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func TestNewService(t *testing.T) {
	mockRepo := &mocks.Storage{}
	jwtAuth := modules.NewJWTAuthenticator(jwtSecret, 1)
	pwdHasher := modules.NewPwdHasher(hashSalt)

	service := service.NewService(mockRepo, jwtAuth, pwdHasher)

	assert.NotNil(t, service)
}
