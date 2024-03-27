package cfg

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

/* ----------------------------------- */
/*             - Config -              */
/* ----------------------------------- */

// Config holds the configuration of the entire API.
type Config struct {
	MainCfg
	DBCfg
	JWTCfg
	TLSCfg
	LoggerCfg
	RLimiterCfg
}

// Load sets up the configuration from the environment variables.
func Load() *Config {

	// The project is either run from the root folder or the /cmd folder.
	// If it's run from /cmd, we add a '..' prefix to the filesystem paths to move them back to the root folder.
	// Otherwise, we just add a '.', staying on the root.
	projectName := envString("PROJECT_NAME", "grpc-gateway-impl")
	wdPrefix := getWdPrefix(projectName)

	developmentDBAdminPass := "n8zAyv96oAtfQoNof-_ulH4pS0Dqf61VThTZbbOLXCU=" // T0D0 remove.

	return &Config{
		MainCfg: MainCfg{
			ProjectName: projectName,
			IsProd:      envBool__("IS_PROD", false),
			IsDev:       !envBool__("IS_PROD", false),
			GRPCPort:    envString("GRPC_PORT", ":50053"),
			HTTPPort:    envString("HTTP_PORT", ":8083"),
			HashSalt:    envString("HASH_SALT", "s0m3_s4l7"), // Used to hash passwords.
		},
		DBCfg: DBCfg{
			Username: envString("DB_USERNAME", "root"),
			Password: envString("DB_PASSWORD", ""),
			Hostname: envString("DB_HOSTNAME", "localhost"),
			Port:     envString("DB_PORT", "3306"),
			Schema:   envString("DB_SCHEMA", "grpc-gateway-impl"),
			Params:   envString("DB_PARAMS", "?charset=utf8&parseTime=True&loc=Local"),

			GormLogLevel: gormLevelMap[envString("DB_LOG_LEVEL", "error")],
			Migrate:      envBool__("DB_MIGRATE", true),
			InsertAdmin:  envBool__("DB_INSERT_ADMIN", true),
			AdminPwd:     envString("DB_ADMIN_PASSWORD", developmentDBAdminPass), // Should already be hashed with our salt.
		},
		JWTCfg: JWTCfg{
			Secret:      envString("JWT_SECRET", "please_set_the_env_var"),
			SessionDays: envInt___("JWT_SESSION_DAYS", 7),
		},
		TLSCfg: TLSCfg{
			Enabled:  envBool__("TLS_ENABLED", false),
			CertPath: envString("TLS_CERT_PATH", wdPrefix+"/server.crt"),
			KeyPath:  envString("TLS_KEY_PATH", wdPrefix+"/server.key"),
		},
		LoggerCfg: LoggerCfg{
			Level:           zapLevelMap[envString("LOG_LEVEL", "info")],
			StacktraceLevel: zapLevelMap[envString("LOG_STACKTRACE_LEVEL", "dpanic")],
			LogCaller:       envBool__("LOG_CALLER", false),
		},
		RLimiterCfg: RLimiterCfg{
			MaxTokens:       envInt___("RATE_LIMITER_MAX_TOKENS", 40),
			TokensPerSecond: envInt___("RATE_LIMITER_TOKENS_PER_SECOND", 10),
		},
	}
}

// Project configuration
type MainCfg struct {
	ProjectName string
	IsProd      bool
	IsDev       bool
	GRPCPort    string
	HTTPPort    string
	HashSalt    string
}

// DB configuration
type DBCfg struct {
	Username string
	Password string
	Hostname string
	Port     string
	Schema   string
	Params   string

	GormLogLevel int // options are in the gormLevelMap var
	Migrate      bool
	InsertAdmin  bool
	AdminPwd     string // Should already be hashed with our salt
}

var gormLevelMap = map[string]int{
	"silent": int(logger.Silent),
	"error":  int(logger.Error),
	"warn":   int(logger.Warn),
	"info":   int(logger.Info),
}

// JWT Auth configuration.
type JWTCfg struct {
	Secret      string
	SessionDays int
}

// TLS configuration.
type TLSCfg struct {
	Enabled  bool // If enabled, use TLS between HTTP and gRPC.
	CertPath string
	KeyPath  string
}

// Logger configuration.
type LoggerCfg struct {
	Level           int // options are in the zapLevelMap var.
	StacktraceLevel int // options are in the zapLevelMap var.
	LogCaller       bool
}

var zapLevelMap = map[string]int{
	"debug":  int(zap.DebugLevel),
	"info":   int(zap.InfoLevel),
	"warn":   int(zap.WarnLevel),
	"error":  int(zap.ErrorLevel),
	"dpanic": int(zap.DPanicLevel),
	"panic":  int(zap.PanicLevel),
	"fatal":  int(zap.FatalLevel),
}

// Rate Limiter configuration.
type RLimiterCfg struct {
	MaxTokens       int // Max tokens the bucket can hold.
	TokensPerSecond int // Tokens reloaded per second.
}

/* ----------------------------------- */
/*            - Helpers -              */
/* ----------------------------------- */

// envString returns the value of an env var or a fallback value if it doesn't exist.
func envString(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// envBool__ returns the value of an env var as a boolean or a fallback value if it doesn't exist.
func envBool__(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true" || value == "TRUE" || value == "1"
	}
	return fallback
}

// envInt___ returns the value of an env var as an int or a fallback value if it doesn't exist.
func envInt___(key string, fallback int) int {
	if value, err := strconv.Atoi(envString(key, "")); err == nil {
		return value
	}
	return fallback
}

// getWdPrefix returns the prefix that needs to be added to the default paths to start at the root folder.
func getWdPrefix(projectName string) string {
	if isWorkingDirRootFolder(projectName) {
		return "."
	}
	return ".."
}

// isWorkingDirRootFolder returns true if the working directory is the root folder.
// The project name is used to determine where we are.
func isWorkingDirRootFolder(projectName string) bool {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf(errs.FatalErrMsgGettingWorkingDir, err) // Our logger is not initialized yet.
	}
	return strings.HasSuffix(workingDir, projectName)
}
