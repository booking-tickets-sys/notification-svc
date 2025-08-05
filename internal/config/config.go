package config

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Redis   RedisConfig   `mapstructure:"redis"`
	Logging LoggingConfig `mapstructure:"logging"`
	Workers WorkersConfig `mapstructure:"workers"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// WorkersConfig holds worker pool configuration
type WorkersConfig struct {
	Count int `mapstructure:"count"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")

	// Set default values
	setDefaults()

	// Read environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, but that's okay - we'll use defaults and env vars
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Redis defaults
	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")

	// Workers defaults
	viper.SetDefault("workers.count", 5)
}

// SetupLogging configures the logger based on configuration
func SetupLogging(config *LoggingConfig) {
	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		logrus.Warnf("Invalid log level %s, using info", config.Level)
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// Set log format
	switch strings.ToLower(config.Format) {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	default:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	// Set output (for now, we only support stdout)
	if config.Output != "stdout" {
		logrus.Warnf("Output %s not supported, using stdout", config.Output)
	}
}
