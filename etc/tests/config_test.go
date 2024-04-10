package tests

import (
	"os"
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("IS_PROD", "true")
	os.Setenv("GRPC_PORT", ":9999")
	os.Setenv("HTTP_PORT", ":8888")
	os.Setenv("TLS_ENABLED", "false")
	os.Setenv("TLS_CERT_PATH", "/path/to/cert")
	os.Setenv("TLS_KEY_PATH", "/path/to/key")

	defer func() {
		os.Unsetenv("IS_PROD")
		os.Unsetenv("GRPC_PORT")
		os.Unsetenv("HTTP_PORT")
		os.Unsetenv("TLS_ENABLED")
		os.Unsetenv("TLS_CERT_PATH")
		os.Unsetenv("TLS_KEY_PATH")
	}()

	config := core.LoadConfig()

	assert.True(t, core.IsProd)
	assert.Equal(t, ":9999", core.GRPCPort)
	assert.Equal(t, ":8888", core.HTTPPort)
	assert.False(t, config.TLSCfg.Enabled)
	assert.Equal(t, "/path/to/cert", config.TLSCfg.CertPath)
	assert.Equal(t, "/path/to/key", config.TLSCfg.KeyPath)
}

func TestGetVar(t *testing.T) {
	const testKey = "TEST_VAR"
	const testValue = "test_value"

	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)

	assert.Equal(t, testValue, os.Getenv(testKey))
}
