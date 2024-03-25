package tests

import (
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/mocks"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/db"

	"github.com/stretchr/testify/assert"
)

/* ----------------------------------- */
/*         - Tests: Repository -       */
/* ----------------------------------- */

func TestNewRepository(t *testing.T) {
	repository := repository.NewRepository(db.NewDB(nil, &mocks.Gorm{}))
	assert.NotNil(t, repository)
}

func OnCallTo(m *mocks.Gorm, methodName string, arguments ...interface{}) {
	m.On(methodName, arguments...).Return(m).Once()
}

func OnGetError(m *mocks.Gorm, expectedError error) {
	m.On("Error").Return(expectedError).Once()
}

const (
	int64PtrTypeName = "*int64"
)
