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

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}
	authConfig, err := AuthConfigFromEnv()
	if err != nil {
		return nil, err
	}
	serverConfig, err := ServerConfigFromEnv()
	if err != nil {
		return nil, err
	}
	dbConfig, err := DbConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return &Config{
		Auth:   *authConfig,
		Server: *serverConfig,
		Db:     *dbConfig,
	}, nil
}
