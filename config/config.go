package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	RateLimitIPRequests      int
	RateLimitIPWindow        time.Duration
	RateLimitIPBlockDuration time.Duration

	RateLimitTokenRequests      int
	RateLimitTokenWindow        time.Duration
	RateLimitTokenBlockDuration time.Duration

	ServerPort string
}

func LoadConfig() *Config {
	return &Config{
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		RateLimitIPRequests:      getEnvAsInt("RATE_LIMIT_IP_REQUESTS", 5),
		RateLimitIPWindow:        getEnvAsDuration("RATE_LIMIT_IP_WINDOW", time.Second),
		RateLimitIPBlockDuration: getEnvAsDuration("RATE_LIMIT_IP_BLOCK_DURATION", 5*time.Minute),

		RateLimitTokenRequests:      getEnvAsInt("RATE_LIMIT_TOKEN_REQUESTS", 10),
		RateLimitTokenWindow:        getEnvAsDuration("RATE_LIMIT_TOKEN_WINDOW", time.Second),
		RateLimitTokenBlockDuration: getEnvAsDuration("RATE_LIMIT_TOKEN_BLOCK_DURATION", 5*time.Minute),

		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
