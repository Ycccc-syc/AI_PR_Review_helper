package config

import "os"

type Config struct {
	Port    string
	GinMode string
	LogLevel string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RedisAddr     string
	RedisPassword string

	GitHubToken    string
	DeepSeekAPIKey string
	DeepSeekModel  string
}

func Load() *Config {
	return &Config{
		Port:     getEnv("PORT", "8080"),
		GinMode:  getEnv("GIN_MODE", "debug"),
		LogLevel: getEnv("LOG_LEVEL", "debug"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "qiniudemo"),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		GitHubToken:    getEnv("GITHUB_TOKEN", ""),
		DeepSeekAPIKey: getEnv("DEEPSEEK_API_KEY", ""),
		DeepSeekModel:  getEnv("DEEPSEEK_MODEL", "deepseek-chat"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
