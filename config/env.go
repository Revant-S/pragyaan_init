package config

import (
	"log"
	"os"
	"strconv"
)

type AppConfig struct {
	ServerPort    int
	ServerHost    string
	MongoURI      string
	DatabaseName  string
	JWTSecret     string
	JWTExpiration int64
	LogLevel      string
	Environment   string
}

var (
	Env *AppConfig
)

func LoadEnvironmentVariables() error {

	Env = &AppConfig{
		ServerPort:    getEnvAsInt("SERVER_PORT", 8080),
		ServerHost:    getEnv("SERVER_HOST", "localhost"),
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DatabaseName:  getEnv("MONGO_DATABASE", "default_database"),
		JWTSecret:     getEnv("JWT_SECRET", generateRandomSecret()),
		JWTExpiration: getEnvAsInt64("JWT_EXPIRATION", 24*60*60),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		Environment:   getEnv("APP_ENV", "development"),
	}

	if Env.MongoURI == "" {
		log.Println("Warning: MongoDB URI is not set")
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
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func generateRandomSecret() string {
	return "default_secret_key_change_in_production"
}

func IsProduction() bool {
	return Env.Environment == "production"
}

func IsDevelopment() bool {
	return Env.Environment == "development"
}
