package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Auth Config
	AuthJWTKey     []byte
	AuthJWTExpires time.Duration
	// Server Config
	ServerPort string
	// HotDuration specifies the duration
	// for which a runner is kept hot after loading a script.
	HotDuration time.Duration
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

	hot := os.Getenv("HOT_DURATION")
	if hot == "" {
		hot = "30m"
	}
	hotDuration, err := time.ParseDuration(hot)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Config{
		AuthJWTKey:     []byte(os.Getenv("AUTH_JWT_KEY")),
		AuthJWTExpires: expiresDuration,
		ServerPort:     os.Getenv("SERVER_PORT"),
		HotDuration:    hotDuration,
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
