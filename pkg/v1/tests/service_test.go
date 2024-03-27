package tests

import (
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/common"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/mocks"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"

	"github.com/stretchr/testify/assert"
)

/* ----------------------------------- */
/*          - Tests: Service -         */
/* ----------------------------------- */

func TestNewService(t *testing.T) {
	mockRepo := &mocks.Repository{}
	jwtAuth := common.NewJWTAuthenticator(jwtSecret, 1)
	pwdHasher := common.NewPwdHasher(hashSalt)

	service := service.NewService(mockRepo, jwtAuth, pwdHasher)

	assert.NotNil(t, service)
}
