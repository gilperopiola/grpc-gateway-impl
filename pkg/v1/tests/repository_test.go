package tests

import (
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/mocks"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository"

	"github.com/stretchr/testify/assert"
)

/* ----------------------------------- */
/*         - Tests: Repository -       */
/* ----------------------------------- */

func TestNewRepository(t *testing.T) {
	repository := repository.NewRepository(&mocks.Gorm{})
	assert.NotNil(t, repository)
}
