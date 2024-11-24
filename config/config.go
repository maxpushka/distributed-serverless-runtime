package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	// Auth Config
	AuthJWTKey     string
	AuthJWTExpires time.Duration
	// Server Config
	ServerPort string
	// DB Config
	dbHost string
	dbPort string
	dbUser string
	dbPass string
	dbName string
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}
	expires := os.Getenv("AUTH_JWT_EXPIRES")
	if expires == "" {
		expires = "24h"
	}
	expiresDuration, err := time.ParseDuration(expires)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &Config{
		AuthJWTKey:     os.Getenv("AUTH_JWT_KEY"),
		AuthJWTExpires: expiresDuration,
		ServerPort:     os.Getenv("SERVER_PORT"),
		dbHost:         os.Getenv("DB_HOST"),
		dbPort:         os.Getenv("DB_PORT"),
		dbUser:         os.Getenv("DB_USER"),
		dbPass:         os.Getenv("DB_PASSWORD"),
		dbName:         os.Getenv("DB_NAME"),
	}, nil
}

func (c Config) DbUrl() string {
	return "host=" + c.dbHost + " port=" + c.dbPort + " user=" + c.dbUser + " password=" + c.dbPass + " dbname=" + c.dbName + " sslmode=disable"
}
