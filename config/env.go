package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// AppConfig holds all configuration parameters for the application
type AppConfig struct {
	ServerPort    int
	ServerHost    string
	MongoURI      string
	DatabaseName  string
	JWTSecret     string
	JWTExpiration time.Duration
	LogLevel      string
	Environment   string
	Debug         bool
}

var (
	defaultConfig = AppConfig{
		ServerPort:    8081,
		ServerHost:    "localhost",
		MongoURI:      "mongodb://localhost:27017",
		DatabaseName:  "default_database",
		JWTSecret:     "default_secret_key_change_in_production",
		JWTExpiration: 24 * time.Hour,
		LogLevel:      "info",
		Environment:   "development",
		Debug:         false,
	}
	Env *AppConfig
)

func LoadEnvironmentVariables(envFiles ...string) error {
	if len(envFiles) > 0 {
		for _, file := range envFiles {
			if err := godotenv.Load(file); err != nil {
				log.Printf("Warning: Could not load env file %s: %v", file, err)
			}
		}
	} else {
		if err := godotenv.Load(".env", ".env.local"); err != nil {
			log.Printf("Warning: Could not load default env files: %v", err)
		}
	}

	Env = &AppConfig{
		ServerPort:    getEnvAsInt("SERVER_PORT", defaultConfig.ServerPort),
		ServerHost:    getEnv("SERVER_HOST", defaultConfig.ServerHost),
		MongoURI:      getEnv("MONGO_URI", defaultConfig.MongoURI),
		DatabaseName:  getEnv("MONGO_DATABASE", defaultConfig.DatabaseName),
		JWTSecret:     getEnv("JWT_SECRET", defaultConfig.JWTSecret),
		JWTExpiration: time.Duration(getEnvAsInt64("JWT_EXPIRATION", int64(defaultConfig.JWTExpiration.Seconds()))) * time.Second,
		LogLevel:      getEnv("LOG_LEVEL", defaultConfig.LogLevel),
		Environment:   getEnv("APP_ENV", defaultConfig.Environment),
		Debug:         getEnvAsBool("DEBUG", defaultConfig.Debug),
	}

	if err := validateConfiguration(Env); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer for %s: %v. Using default: %d", key, err, defaultValue)
		return defaultValue
	}
	return value
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		log.Printf("Warning: Invalid 64-bit integer for %s: %v. Using default: %d", key, err, defaultValue)
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(strings.ToLower(valueStr))
	if err != nil {
		log.Printf("Warning: Invalid boolean for %s: %v. Using default: %v", key, err, defaultValue)
		return defaultValue
	}
	return value
}

func validateConfiguration(config *AppConfig) error {
	if config.ServerPort <= 0 || config.ServerPort > 65535 {
		return fmt.Errorf("invalid server port: %d", config.ServerPort)
	}

	if config.MongoURI == "" {
		return fmt.Errorf("mongodb URI cannot be empty")
	}

	if config.JWTSecret == defaultConfig.JWTSecret {
		log.Println("Warning: Using default JWT secret. Please change in production!")
	}

	return nil
}
