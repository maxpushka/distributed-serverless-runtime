package config

import "os"

type ServerConfig struct {
	Host string
	Port string
}

func ServerConfigFromEnv() (*ServerConfig, error) {
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	return &ServerConfig{
		Host: host,
		Port: port,
	}, nil
}

func (c ServerConfig) ConnectionString() string {
	return c.Host + ":" + c.Port
}