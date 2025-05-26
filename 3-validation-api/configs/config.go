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
	Email      string
	Password   string
	SmtpHost   string
	SmtpPort   string
	SenderName string
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
			Email:      os.Getenv("EMAIL"),
			Password:   os.Getenv("PASSWORD"),
			SmtpHost:   os.Getenv("SMTP_HOST"),
			SmtpPort:   os.Getenv("SMTP_PORT"),
			SenderName: os.Getenv("SENDER_NAME"),
		},
		Storage: StorageConfig{
			Path: os.Getenv("STORAGE_PATH"),
		},
	}
}
