package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the API Gateway
type Config struct {
	App      AppConfig
	Logging  LoggingConfig
	Server   ServerConfig
	Bot      TelegramBotConfig
	Database DatabaseConfig
	Security SecurityConfig
}

// AppConfig holds application metadata
type AppConfig struct {
	Environment string
	Name        string
	Version     string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Output string
	File   string
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host               string
	Port               string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	ShutdownTimeout    time.Duration
	EnableCORS         bool
	CORSAllowedOrigins []string
	CORSAllowedMethods []string
	CORSAllowedHeaders []string
	ExposeHeaders      []string
	AllowCredentials   bool
	MaxAge             int
}

// TelegramBotConfig holds telegram bot configuration
type TelegramBotConfig struct {
	Token          string
	AdminUsernames []string
}

// DatabaseConfig holds postgres configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	QueryTimeout    time.Duration
}

// SecurityConfig holds Security configuration
type SecurityConfig struct {
	JWTSecret            string
	AccessTokenExpiresIn time.Duration
	Email                string
	HashedPassword       string
}

func Load(env string) (*Config, error) {
	var envFile string
	switch strings.ToLower(env) {
	case "production":
		envFile = ".env.prod"
	default:
		envFile = ".env.dev"
	}

	_ = godotenv.Load(envFile)

	app := AppConfig{
		Environment: getEnv("APP_ENV", "development"),
		Name:        getEnv("APP_NAME", "emil-server"),
		Version:     getEnv("APP_VERSION", "1.0.0"),
	}

	logging := LoggingConfig{
		Level:  getEnv("LOG_LEVEL", "debug"),
		Output: getEnv("LOG_OUTPUT", "stdout"),
		File:   getEnv("LOG_FILE", ""),
	}

	server := ServerConfig{
		Host:               getEnv("SERVER_HOST", "127.0.0.1"),
		Port:               getEnv("SERVER_PORT", "8080"),
		ReadTimeout:        getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:       getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:        getEnvAsDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
		ShutdownTimeout:    getEnvAsDuration("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
		EnableCORS:         getEnvAsBool("SERVER_ENABLE_CORS", true),
		CORSAllowedOrigins: getEnvAsSlice("SERVER_CORS_ALLOWED_ORIGINS", []string{"*"}, ","),
		CORSAllowedMethods: getEnvAsSlice("SERVER_CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, ","),
		CORSAllowedHeaders: getEnvAsSlice("SERVER_CORS_ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization"}, ","),
		ExposeHeaders:      getEnvAsSlice("SERVER_CORS_EXPOSE_HEADERS", []string{"Content-Length,Set-Cookie,X-Total-Count"}, ","),
		AllowCredentials:   getEnvAsBool("SERVER_CORS_ALLOW_CREDENTIALS", true),
		MaxAge:             getEnvAsInt("SERVER_CORS_MAX_AGE", 86400),
	}

	bot := TelegramBotConfig{
		Token:          getEnv("BOT_TOKEN", ""),
		AdminUsernames: getEnvAsSlice("BOT_ADMINS", []string{}, ","),
	}

	database := DatabaseConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnv("DB_PORT", "5432"),
		User:            getEnv("DB_USER", "ramisoul"),
		Password:        getEnv("DB_PASSWORD", ""),
		Name:            getEnv("DB_NAME", "dating"),
		SSLMode:         getEnv("DB_SSL_MODE", "disable"),
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		ConnMaxIdleTime: getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 1*time.Minute),
		QueryTimeout:    getEnvAsDuration("DB_QUERY_TIMEOUT", 5*time.Second),
	}

	security := SecurityConfig{
		JWTSecret:            getEnv("JWT_SECRET", ""),
		AccessTokenExpiresIn: getEnvAsDuration("ACCESS_TOKEN_EXPIRES_IN", 15*time.Minute),
		Email:                getEnv("EMAIL", ""),
		HashedPassword:       getEnv("HASHED_PASSWORD", ""),
	}

	cfg := &Config{
		App:      app,
		Logging:  logging,
		Server:   server,
		Bot:      bot,
		Database: database,
		Security: security,
	}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}
