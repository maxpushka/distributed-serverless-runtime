package config

import (
	"errors"
	"log"
	"os"
)

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func DbConfigFromEnv() (*DbConfig, error) {
	for _, envVar := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"} {
		if os.Getenv(envVar) == "" {
			log.Fatal(envVar + " is required")
			return nil, errors.New(envVar + " is required")
		}
	}
	return &DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}, nil
}

func (c DbConfig) DbUrl() string {
	return "host=" + c.Host + " port=" + c.Port + " user=" + c.User + " password=" + c.Password + " dbname=" + c.Name + " sslmode=disable"
}
