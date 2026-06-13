package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	defaultEnvironment    = "development"
	defaultPort           = "8080"
	defaultDBHost         = "db"
	defaultDBPort         = "5432"
	defaultDBUser         = "postgres"
	defaultDBPassword     = "postgres"
	defaultDBName         = "carro_ideal"
	defaultDBMaxOpenConns = 25
	defaultDBMaxIdleConns = 5
	defaultOpenAIModel    = "gpt-4o-mini"
	defaultOpenAITimeout  = 15
)

type Config struct {
	Environment       string
	Port              string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	SessionSecret     string
	LogLevel          string
	AllowedOrigins    []string
	RateLimit         int
	RateWindow        int
	CacheTTL          int
	TLSCertFile       string
	TLSKeyFile        string
	OpenAIAPIKey      string
	OpenAIModel       string
	OpenAITimeoutSecs int
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load .env file: %w", err)
	}

	port := getEnv("PORT", defaultPort)
	cfg := &Config{
		Environment:       getEnv("ENVIRONMENT", defaultEnvironment),
		Port:              port,
		DBHost:            getEnv("DB_HOST", defaultDBHost),
		DBPort:            getEnv("DB_PORT", defaultDBPort),
		DBUser:            getEnv("DB_USER", defaultDBUser),
		DBPassword:        getEnv("DB_PASSWORD", defaultDBPassword),
		DBName:            getEnv("DB_NAME", defaultDBName),
		DBMaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", defaultDBMaxOpenConns),
		DBMaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", defaultDBMaxIdleConns),
		SessionSecret:     getEnv("SESSION_SECRET", ""),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		AllowedOrigins:    splitEnv("ALLOWED_ORIGINS"),
		RateLimit:         getEnvInt("RATE_LIMIT_REQUESTS", 120),
		RateWindow:        getEnvInt("RATE_LIMIT_WINDOW_SECONDS", 60),
		CacheTTL:          getEnvInt("CACHE_TTL_SECONDS", 60),
		TLSCertFile:       getEnv("TLS_CERT_FILE", ""),
		TLSKeyFile:        getEnv("TLS_KEY_FILE", ""),
		OpenAIAPIKey:      getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:       getEnv("OPENAI_MODEL", defaultOpenAIModel),
		OpenAITimeoutSecs: getEnvInt("OPENAI_TIMEOUT_SECONDS", defaultOpenAITimeout),
	}

	if cfg.SessionSecret == "" {
		return nil, fmt.Errorf("missing required environment variable: SESSION_SECRET")
	}

	if cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" || cfg.DBHost == "" || cfg.DBPort == "" {
		return nil, fmt.Errorf("database environment variables must be configured")
	}
	if (cfg.TLSCertFile == "") != (cfg.TLSKeyFile == "") {
		return nil, fmt.Errorf("TLS_CERT_FILE and TLS_KEY_FILE must be configured together")
	}
	if cfg.RateLimit < 1 || cfg.RateWindow < 1 || cfg.CacheTTL < 1 {
		return nil, fmt.Errorf("rate limit, rate window, and cache TTL must be positive")
	}

	return cfg, nil
}

func splitEnv(key string) []string {
	value := getEnv(key, "")
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
