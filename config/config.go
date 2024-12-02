package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Auth     Auth     `env-prefix:"AUTH"`
	Server   Server   `env-prefix:"SERVER"`
	Database Database `env-prefix:"DATABASE"`
	Executor Executor `env-prefix:"EXECUTOR"`
}

type Auth struct {
	JWTKey     string        `env:"JWT_KEY"`
	JWTExpires time.Duration `env:"JWT_EXPIRES" env-default:"24h"`
}

type Server struct {
	Host            string `env:"HOST" env-default:"0.0.0.0"`
	Port            string `env:"PORT" env-default:"8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" env-default:"."`
}

func (c Server) ConnectionString() string {
	return c.Host + ":" + c.Port
}

func (c Server) ConfigDir() (string, error) {
	filePath := filepath.Join(c.FileStoragePath, "configs")
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		log.Print(err)
		return "", err
	}
	return filePath, nil
}

func (c Server) ExecutableDir() (string, error) {
	filePath := filepath.Join(c.FileStoragePath, "executables")
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		log.Print(err)
		return "", err
	}
	return filePath, nil
}

type Database struct {
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Name     string `env:"NAME"`
}

func (c Database) ConnectionString() string {
	return "host=" + c.Host + " port=" + c.Port + " user=" + c.User + " password=" + c.Password + " dbname=" + c.Name + " sslmode=disable"
}

type Executor struct {
	HotDuration time.Duration `env:"HOT_DURATION" env-default:"30m"`
}

func New() (c Config, err error) {
	err = cleanenv.ReadEnv(&c)
	return c, err
}
