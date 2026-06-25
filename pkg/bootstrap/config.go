package bootstrap

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var (
	config   Config
	configMu sync.RWMutex
)

type Config struct {
	Environment string
	ServiceName string
	Port        int

	DBHost    string
	DBPort    int
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string

	GRPCTLSCertPath string
	GRPCTLSKeyPath  string

	LocationDefaultCountryCode string
	LocationLimitDefault       int32
	LocationLimitMax           int32
}

func LoadConfig(dotenvPath string) (Config, error) {
	if dotenvPath != "" {
		_ = godotenv.Load(dotenvPath)
	}

	cfg := Config{
		Environment:                getEnv("ENVIRONMENT", "development"),
		ServiceName:                getEnv("SERVICE_NAME", "neuraclinic-location"),
		Port:                       getEnvInt("PORT", 8000),
		DBHost:                     getEnv("DB_HOST", ""),
		DBPort:                     getEnvInt("DB_PORT", 5432),
		DBUser:                     getEnv("DB_USER", ""),
		DBPass:                     getEnv("DB_PASS", ""),
		DBName:                     getEnv("DB_NAME", ""),
		DBSSLMode:                  getEnv("DB_SSLMODE", "disable"),
		GRPCTLSCertPath:            getEnv("GRPC_TLS_CERT_PATH", ""),
		GRPCTLSKeyPath:             getEnv("GRPC_TLS_KEY_PATH", ""),
		LocationDefaultCountryCode: strings.ToUpper(getEnv("LOCATION_DEFAULT_COUNTRY_CODE", "MX")),
		LocationLimitDefault:       int32(getEnvInt("LOCATION_LIMIT_DEFAULT", 20)),
		LocationLimitMax:           int32(getEnvInt("LOCATION_LIMIT_MAX", 100)),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	setConfig(cfg)
	return cfg, nil
}

func GetConfig() Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return config
}

func setConfig(cfg Config) {
	configMu.Lock()
	config = cfg
	configMu.Unlock()
}

func (c Config) Validate() error {
	required := map[string]string{
		"DB_HOST":                       c.DBHost,
		"DB_USER":                       c.DBUser,
		"DB_PASS":                       c.DBPass,
		"DB_NAME":                       c.DBName,
		"GRPC_TLS_CERT_PATH":            c.GRPCTLSCertPath,
		"GRPC_TLS_KEY_PATH":             c.GRPCTLSKeyPath,
		"LOCATION_DEFAULT_COUNTRY_CODE": c.LocationDefaultCountryCode,
	}

	for key, value := range required {
		if value == "" {
			return fmt.Errorf("missing required config key: %s", key)
		}
	}

	if c.Port <= 0 {
		return fmt.Errorf("PORT must be greater than zero")
	}
	if c.DBPort <= 0 {
		return fmt.Errorf("DB_PORT must be greater than zero")
	}
	if len(c.LocationDefaultCountryCode) != 2 {
		return fmt.Errorf("LOCATION_DEFAULT_COUNTRY_CODE must be ISO 3166-1 alpha-2")
	}
	if c.LocationLimitDefault <= 0 {
		return fmt.Errorf("LOCATION_LIMIT_DEFAULT must be greater than zero")
	}
	if c.LocationLimitMax < c.LocationLimitDefault {
		return fmt.Errorf("LOCATION_LIMIT_MAX must be greater than or equal to LOCATION_LIMIT_DEFAULT")
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
