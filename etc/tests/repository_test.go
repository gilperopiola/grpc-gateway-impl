package tests

import (
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/app/storage"
	"github.com/gilperopiola/grpc-gateway-impl/etc/tests/mocks"

	"github.com/stretchr/testify/assert"
)

/* ----------------------------------- */
/*         - Tests: Storage -       */
/* ----------------------------------- */

func TestNewStorage(t *testing.T) {
	storage := storage.NewStorage(&mocks.Gorm{})
	assert.NotNil(t, storage)
}
