package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Redis  RedisConfig
	Server ServerConfig
	Asynq  AsynqConfig
	Email  EmailConfig
	Log    LogConfig
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type ServerConfig struct {
	Port string
	Host string
}

type AsynqConfig struct {
	Concurrency int
	Queues      map[string]int
}

type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	From     string
	Password string
}

type LogConfig struct {
	Level string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load("config.env")

	concurrency, _ := strconv.Atoi(getEnv("ASYNQ_CONCURRENCY", "10"))
	critical, _ := strconv.Atoi(getEnv("ASYNQ_QUEUE_CRITICAL", "6"))
	defaultQueue, _ := strconv.Atoi(getEnv("ASYNQ_QUEUE_DEFAULT", "3"))
	low, _ := strconv.Atoi(getEnv("ASYNQ_QUEUE_LOW", "1"))
	smtpPort, _ := strconv.Atoi(getEnv("EMAIL_SMTP_PORT", "587"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	config := &Config{
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Asynq: AsynqConfig{
			Concurrency: concurrency,
			Queues: map[string]int{
				"critical": critical,
				"default":  defaultQueue,
				"low":      low,
			},
		},
		Email: EmailConfig{
			SMTPHost: getEnv("EMAIL_SMTP_HOST", "smtp.gmail.com"),
			SMTPPort: smtpPort,
			From:     getEnv("EMAIL_FROM", "noreply@example.com"),
			Password: getEnv("EMAIL_PASSWORD", ""),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
