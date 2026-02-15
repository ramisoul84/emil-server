package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	switch strings.ToLower(valueStr) {
	case "true", "1", "yes", "on", "enabled":
		return true
	case "false", "0", "no", "off", "disabled":
		return false
	default:
		return defaultValue
	}
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, sep)
}

func validateConfig(cfg *Config) error {
	if cfg.Security.JWTSecret == "" {
		return fmt.Errorf("JWT Secret must be set")
	}
	if cfg.Security.Email == "" {
		return fmt.Errorf("Email must be set")
	}
	if cfg.Security.HashedPassword == "" {
		return fmt.Errorf("password must be set")
	}
	if cfg.Bot.Token == "" {
		return fmt.Errorf("Telegram bot token must be set")
	}
	if cfg.Database.Password == "" {
		return fmt.Errorf("Database password must be set")
	}

	return nil
}
