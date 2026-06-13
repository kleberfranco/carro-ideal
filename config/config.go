package config

import (
	"fmt"
	"os"
	"strconv"

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
)

type Config struct {
	Environment    string
	Port           string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBMaxOpenConns int
	DBMaxIdleConns int
	SessionSecret  string
	LogLevel       string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load .env file: %w", err)
	}

	port := getEnv("PORT", defaultPort)
	cfg := &Config{
		Environment:    getEnv("ENVIRONMENT", defaultEnvironment),
		Port:           port,
		DBHost:         getEnv("DB_HOST", defaultDBHost),
		DBPort:         getEnv("DB_PORT", defaultDBPort),
		DBUser:         getEnv("DB_USER", defaultDBUser),
		DBPassword:     getEnv("DB_PASSWORD", defaultDBPassword),
		DBName:         getEnv("DB_NAME", defaultDBName),
		DBMaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", defaultDBMaxOpenConns),
		DBMaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", defaultDBMaxIdleConns),
		SessionSecret:  getEnv("SESSION_SECRET", ""),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}

	if cfg.SessionSecret == "" {
		return nil, fmt.Errorf("missing required environment variable: SESSION_SECRET")
	}

	if cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" || cfg.DBHost == "" || cfg.DBPort == "" {
		return nil, fmt.Errorf("database environment variables must be configured")
	}

	return cfg, nil
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
