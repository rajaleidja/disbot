package main

import (
    "discord-bot/internal/bot"
    "discord-bot/internal/config"
    "discord-bot/internal/database"
    "log"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Error loading config: ", err)
    }

    db, err := database.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Error connecting to database: ", err)
    }
    defer db.Close()

    bot.Start(cfg.BotToken, db)
}