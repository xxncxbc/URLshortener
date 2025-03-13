package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type DbConfig struct {
	Dsn string
}

type AuthConfig struct {
	AccessSecret  string
	RefreshSecret string
}

type Config struct {
	Db   DbConfig
	Auth AuthConfig
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file, using default config")
	}
	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			AccessSecret:  os.Getenv("ACCESS_SECRET"),
			RefreshSecret: os.Getenv("REFRESH_SECRET"),
		},
	}
}
