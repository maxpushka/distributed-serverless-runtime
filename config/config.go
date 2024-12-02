package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Auth     AuthConfig
	Server   ServerConfig
	Db       DbConfig
	Executor ExecutorConfig
}

type ExecutorConfig struct {
	HotDuration time.Duration
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	authConfig := AuthConfigFromEnv()
	serverConfig := ServerConfigFromEnv()
	dbConfig := DbConfigFromEnv()

	hot := os.Getenv("HOT_DURATION")
	if hot == "" {
		hot = "30m"
	}
	hotDuration, err := time.ParseDuration(hot)
	if err != nil {
		log.Fatal(err)
	}
	executorConfig := ExecutorConfig{HotDuration: hotDuration}

	return &Config{
		Auth:     authConfig,
		Server:   serverConfig,
		Db:       dbConfig,
		Executor: executorConfig,
	}
}
