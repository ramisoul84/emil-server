package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// Server configuration
type ServerConfig struct {
	Port         string        `env:"SERVER_PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"10s"`
}

// Database configuration
type DatabaseConfig struct {
	Host            string        `env:"DB_HOST" envDefault:"localhost"`
	Port            string        `env:"DB_PORT" envDefault:"5432"`
	User            string        `env:"DB_USER" envDefault:"postgres"`
	Password        string        `env:"DB_PASSWORD" envDefault:""`
	DB              string        `env:"DB_NAME" envDefault:"emil"`
	SSLMode         string        `env:"DB_SSL_MODE" envDefault:"disable"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"5m"`
	ConnMaxIdleTime time.Duration `env:"DB_CONN_MAX_IDLE_TIME" envDefault:"2m"`
	Timeout         time.Duration `env:"DB_TIMEOUT" envDefault:"5s"`
}

type JWTConfig struct {
	JWTSecret          string        `env:"JWT_SECRET" envDefault:""`
	AccessTokenExpiry  time.Duration `env:"ACCESS_TOKEN_EXPIRY" envDefault:"15m"`
	RefreshTokenExpiry time.Duration `env:"REFRESH_TOKEN_EXPIRY" envDefault:"7d"`
}

func Load(env string) (*Config, error) {
	var envFile string
	if env == "production" {
		envFile = ".env.prod"
	} else {
		envFile = ".env.dev"
	}

	// Load environment file (ignore error if file doesn't exist)
	_ = godotenv.Load(envFile)

	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			DB:              getEnv("DB_NAME", "emil"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getDurationEnv("DB_CONN_MAX_IDLE_TIME", 2*time.Minute),
			Timeout:         getDurationEnv("DB_TIMEOUT", 5*time.Second),
		},
		JWT: JWTConfig{
			JWTSecret:          getEnv("JWT_SECRET", "your-default-secret-key-change-in-production"),
			AccessTokenExpiry:  getDurationEnv("ACCESS_TOKEN_EXPIRY", 15*time.Minute),
			RefreshTokenExpiry: getDurationEnv("REFRESH_TOKEN_EXPIRY", 24*7*time.Hour),
		},
	}, nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
