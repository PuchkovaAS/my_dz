package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db     DbConfig
	Logger LoggerConfig
	Auth   AuthConfig
}

type LoggerConfig struct {
	LogFile string
}

type DbConfig struct {
	Dsn string
}

type AuthConfig struct {
	Secret string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default config")
	}
	return &Config{
		Db: DbConfig{
			Dsn: fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
				os.Getenv("DB_HOST"),
				os.Getenv("DB_USER"),
				os.Getenv("DB_PASSWORD"),
				os.Getenv("DB_NAME"),
				os.Getenv("DB_PORT"),
				os.Getenv("DB_SSLMODE"),
			),
		},
		Logger: LoggerConfig{
			LogFile: os.Getenv("LOG_FILE"),
		},

		Auth: AuthConfig{
			Secret: os.Getenv("SECRET"),
		},
	}
}
