package config

import (
	"os"
    "log"
	"github.com/joho/godotenv"
)

type Config struct {
    BotToken    string
    DatabaseURL string
}

func Load() (*Config, error) {

    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    return &Config{
        BotToken:    os.Getenv("BOT_TOKEN"),
        DatabaseURL: os.Getenv("DATABASE_URL"),
    }, nil
}