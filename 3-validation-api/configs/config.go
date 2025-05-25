package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Email   EmailConfig
	Storage StorageConfig
}

type EmailConfig struct {
	Email       string
	Password    string
	Address     string
	SmtpAddress string
}

type StorageConfig struct {
	Path string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default config")
	}
	return &Config{
		Email: EmailConfig{
			Email:       os.Getenv("EMAIL"),
			Password:    os.Getenv("PASSWORD"),
			Address:     os.Getenv("ADDRESS"),
			SmtpAddress: os.Getenv("SMTP_ADDRESS"),
		},
		Storage: StorageConfig{
			Path: os.Getenv("STORAGE_PATH"),
		},
	}
}
