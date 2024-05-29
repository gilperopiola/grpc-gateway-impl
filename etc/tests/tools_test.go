package tests

import (
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/toolbox"

	"github.com/stretchr/testify/assert"
)

func TestSetupTools(t *testing.T) {
	cfg := core.LoadConfig()
	toolbox := toolbox.Setup(cfg)
	assert.NotNil(t, toolbox)
}
