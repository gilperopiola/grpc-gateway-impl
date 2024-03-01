package config

import (
	"os"
	"testing"

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

	config := LoadConfig()

	assert.True(t, config.IsProd)
	assert.Equal(t, ":9999", config.GRPCPort)
	assert.Equal(t, ":8888", config.HTTPPort)
	assert.False(t, config.TLS.Enabled)
	assert.Equal(t, "/path/to/cert", config.TLS.CertPath)
	assert.Equal(t, "/path/to/key", config.TLS.KeyPath)
}

func TestGetVar(t *testing.T) {
	const testKey = "TEST_VAR"
	const testValue = "test_value"

	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)

	assert.Equal(t, testValue, getVar(testKey, "fallback"))
	assert.Equal(t, "fallback", getVar("NON_EXISTING_VAR", "fallback"))
}

func TestGetVarBool(t *testing.T) {
	os.Setenv("TRUE_VAR", "true")
	os.Setenv("FALSE_VAR", "false")
	defer func() {
		os.Unsetenv("TRUE_VAR")
		os.Unsetenv("FALSE_VAR")
	}()

	assert.True(t, getVarBool("TRUE_VAR", false))
	assert.False(t, getVarBool("FALSE_VAR", true))
	assert.True(t, getVarBool("NON_EXISTING_VAR", true))
	assert.False(t, getVarBool("NON_EXISTING_VAR", false))
}
