package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Auth   AuthConfig
	Server ServerConfig
	Db     DbConfig
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	authConfig := AuthConfigFromEnv()
	serverConfig := ServerConfigFromEnv()
	dbConfig := DbConfigFromEnv()
	return &Config{
		Auth:   authConfig,
		Server: serverConfig,
		Db:     dbConfig,
	}
}
