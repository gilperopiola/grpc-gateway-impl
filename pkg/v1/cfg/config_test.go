package cfg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetTestingConfig() *Config {
	return &Config{
		MainConfig: &MainConfig{
			ProjectName: "grpc-gateway-impl-test",
			IsProd:      false,
			GRPCPort:    ":50053",
			HTTPPort:    ":8083",
			HashSalt:    "t3s7_s4l7",
		},
		DBConfig: &DBConfig{
			Username:    "",
			Password:    "",
			Hostname:    "",
			Port:        "",
			Schema:      "",
			Params:      "",
			Migrate:     true,
			InsertAdmin: true,
			AdminPwd:    "test_admin_pwd",
		},
		JWTConfig: &JWTConfig{
			Secret:      "test_jwt_secret",
			SessionDays: 1,
		},
		TLSConfig: &TLSConfig{
			Enabled:  true,
			CertPath: "",
			KeyPath:  "",
		},
		RateLimiterConfig: &RateLimiterConfig{
			MaxTokens:       5,
			TokensPerSecond: 2,
		},
	}
}

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

	config := Load()

	assert.True(t, config.IsProd)
	assert.Equal(t, ":9999", config.GRPCPort)
	assert.Equal(t, ":8888", config.HTTPPort)
	assert.False(t, config.TLSConfig.Enabled)
	assert.Equal(t, "/path/to/cert", config.TLSConfig.CertPath)
	assert.Equal(t, "/path/to/key", config.TLSConfig.KeyPath)
}

func TestGetVar(t *testing.T) {
	const testKey = "TEST_VAR"
	const testValue = "test_value"

	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)

	assert.Equal(t, testValue, envString(testKey, "fallback"))
	assert.Equal(t, "fallback", envString("NON_EXISTING_VAR", "fallback"))
}

func TestGetVarBool(t *testing.T) {
	os.Setenv("TRUE_VAR", "true")
	os.Setenv("FALSE_VAR", "false")
	defer func() {
		os.Unsetenv("TRUE_VAR")
		os.Unsetenv("FALSE_VAR")
	}()

	assert.True(t, envBulean("TRUE_VAR", false))
	assert.False(t, envBulean("FALSE_VAR", true))
	assert.True(t, envBulean("NON_EXISTING_VAR", true))
	assert.False(t, envBulean("NON_EXISTING_VAR", false))
}
