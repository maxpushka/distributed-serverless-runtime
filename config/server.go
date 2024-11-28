package config

import (
	"log"
	"os"
	"path/filepath"
)

type ServerConfig struct {
	Host            string
	Port            string
	FileStoragePath string
}

func ServerConfigFromEnv() ServerConfig {
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	fileStoragePath := os.Getenv("SERVER_FILE_STORAGE_PATH")
	if fileStoragePath == "" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		fileStoragePath = dir
	}
	err := os.MkdirAll(fileStoragePath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return ServerConfig{
		Host:            host,
		Port:            port,
		FileStoragePath: fileStoragePath,
	}
}

func (c ServerConfig) ConnectionString() string {
	return c.Host + ":" + c.Port
}

func (c ServerConfig) ConfigDir() (string, error) {
	filePath := filepath.Join(c.FileStoragePath, "configs")
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		log.Print(err)
		return "", err
	}
	return filePath, nil
}

func (c ServerConfig) ExecutableDir() (string, error) {
	filePath := filepath.Join(c.FileStoragePath, "executables")
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		log.Print(err)
		return "", err
	}
	return filePath, nil
}
