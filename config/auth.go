package config

import (
	"log"
	"os"
	"time"
)

type AuthConfig struct {
	JWTKey     []byte
	JWTExpires time.Duration
}

func AuthConfigFromEnv() AuthConfig {
	key := os.Getenv("AUTH_JWT_KEY")
	if key == "" {
		log.Fatal("AUTH_JWT_KEY is required")
	}
	expires := os.Getenv("AUTH_JWT_EXPIRES")
	if expires == "" {
		expires = "24h"
	}
	expiresDuration, err := time.ParseDuration(expires)
	if err != nil {
		log.Fatal(err)
	}

	return AuthConfig{
		JWTKey:     []byte(key),
		JWTExpires: expiresDuration,
	}
}
