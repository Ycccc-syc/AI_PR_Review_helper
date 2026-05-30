package config

import (
	"bufio"
	"os"
	"strings"
)

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

	JWTSecret string
}

func Load() *Config {
	loadDotEnv()

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

		JWTSecret: getEnv("JWT_SECRET", "change-me-in-production"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// WarnDefaults 在 Load 后调用，对不安全默认值输出警告
func (c *Config) WarnDefaults() {
	if c.JWTSecret == "change-me-in-production" {
		println("[WARN] JWT_SECRET 使用默认值，生产环境请务必更换为强随机字符串")
	}
	if c.GitHubToken == "" {
		println("[WARN] GITHUB_TOKEN 未设置，公开仓库限制 60次/小时，私仓无法访问")
	}
	if c.DeepSeekAPIKey == "" {
		println("[WARN] DEEPSEEK_API_KEY 未设置，AI 分析功能不可用")
	}
}

// loadDotEnv 从项目根目录读取 .env 文件，注入环境变量（不覆盖已有值）
func loadDotEnv() {
	paths := []string{".env", "../.env", "../../.env"}
	var file *os.File
	var err error
	for _, p := range paths {
		file, err = os.Open(p)
		if err == nil {
			break
		}
	}
	if file == nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, "=")
		if idx == -1 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}
