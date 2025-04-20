package config

import (
	"os"
)

type EnvironmentType string

const (
	Production  EnvironmentType = "prod"
	Development EnvironmentType = "dev"
	Local       EnvironmentType = "local"
)

type Config struct {
	APP         string
	Environment EnvironmentType
	LogLevel    string
	AppURL      string

	Server struct {
		Host         string
		Port         string
		ReadTimeout  string
		WriteTimeout string
		IdleTimeout  string
	}

	Context struct {
		Timeout string
	}

	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		Sslmode  string
	}

	Redis struct {
		Host            string
		Port            string
		Password        string
		DB              string
		StorageDeadline string
	}

	Typesense struct {
		APIKey        string
		Host          string
		Port          string
		RetryCount    int
		RetryWaitTime string
		Timeout       string
	}

	OpenAI struct {
		APIKey  string
		Timeout string
	}

	Token struct {
		Secret string
	}

	Transliterator struct {
		Host    string
		Port    string
		Timeout string
	}

	Telegram struct {
		Token  string
		ChatID string
	}
}

func New() *Config {
	var config Config

	// general configuration
	config.APP = getEnv("APP", "")
	config.Environment = EnvironmentType(getEnv("ENVIRONMENT", "develop"))
	config.LogLevel = getEnv("LOG_LEVEL", "debug")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "5m")
	config.AppURL = getEnv("APP_URL", "")

	// server configuration
	config.Server.Host = getEnv("SERVER_HOST", "localhost")
	config.Server.Port = getEnv("SERVER_PORT", ":8000")
	config.Server.ReadTimeout = getEnv("SERVER_READ_TIMEOUT", "10s")
	config.Server.WriteTimeout = getEnv("SERVER_WRITE_TIMEOUT", "10s")
	config.Server.IdleTimeout = getEnv("SERVER_IDLE_TIMEOUT", "120s")

	// db configuration
	config.DB.Host = getEnv("POSTGRES_HOST", "localhost")
	config.DB.Port = getEnv("POSTGRES_PORT", "5432")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "whereismycity")
	config.DB.User = getEnv("POSTGRES_USER", "postgres")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "postgres")
	config.DB.Sslmode = getEnv("POSTGRES_SSLMODE", "disable")

	// redis configuration
	config.Redis.Host = getEnv("REDIS_HOST", "localhost")
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	config.Redis.DB = getEnv("REDIS_DB", "0")
	config.Redis.StorageDeadline = getEnv("REDIS_STORAGE_DEADLINE", "30m")

	// typesense configuration
	config.Typesense.APIKey = getEnv("TYPESENSE_API_KEY", "")
	config.Typesense.Host = getEnv("TYPESENSE_HOST", "localhost")
	config.Typesense.Port = getEnv("TYPESENSE_PORT", "8108")
	config.Typesense.RetryCount = 3
	config.Typesense.RetryWaitTime = "1s"
	config.Typesense.Timeout = getEnv("TYPESENSE_TIMEOUT", "30s")

	// embeddings configuration
	config.OpenAI.APIKey = getEnv("OPENAI_API_KEY", "")
	config.OpenAI.Timeout = getEnv("OPENAI_TIMEOUT", "30s")

	// token configuration
	config.Token.Secret = getEnv("TOKEN_SECRET", "secret")

	// transliterator configuration
	config.Transliterator.Host = getEnv("TRANSLITERATOR_HOST", "0.0.0.0")
	config.Transliterator.Port = getEnv("TRANSLITERATOR_PORT", "5005")
	config.Transliterator.Timeout = getEnv("TRANSLITERATOR_TIMEOUT", "30s")

	// telegram configuration
	config.Telegram.Token = getEnv("TELEGRAM_TOKEN", "")
	config.Telegram.ChatID = getEnv("TELEGRAM_CHAT_ID", "")

	return &config
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
