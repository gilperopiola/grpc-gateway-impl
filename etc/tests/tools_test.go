package tests

import (
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools"

	"github.com/stretchr/testify/assert"
)

func TestSetupTools(t *testing.T) {
	cfg := core.LoadConfig()
	tools := tools.Setup(cfg, nil)
	assert.NotNil(t, tools)
}
