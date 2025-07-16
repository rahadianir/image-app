package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	// General
	AppName string
	Port    int

	// Database
	DBURL string
}

func InitConfig() *Config {
	return &Config{
		AppName: getEnvString("APP_NAME", "image-app"),
		Port:    getEnvInt("APP_PORT", 8080),
		DBURL:   getEnvString("DB_URL", "postgres://user:pass@localhost:5432/database?sslmode=disable"),
	}
}

func getEnvString(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}

func getEnvInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	valInt, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Errorf("failed to convert config key: %s, err: %v", key, err))
	}

	return valInt
}
