package tests

import (
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/app/external"
	"github.com/gilperopiola/grpc-gateway-impl/etc/tests/mocks"

	"github.com/stretchr/testify/assert"
)

/* ----------------------------------- */
/*         - Tests: Storage -       */
/* ----------------------------------- */

func TestNewStorage(t *testing.T) {
	external := external.NewExternalLayer(&mocks.Gorm{})
	assert.NotNil(t, external)
}
