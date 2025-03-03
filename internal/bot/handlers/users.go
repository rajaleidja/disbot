package handlers

import (
    "database/sql"
    "github.com/bwmarrin/discordgo"
    "log"
    "strconv"
    "strings"
)

func HandleUserRegister(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB) {
    if m.Author.Bot {
        return
    }

    if strings.HasPrefix(m.Content, "!register") {
        parts := strings.Split(m.Content, " ")
        if len(parts) != 3 {
            s.ChannelMessageSend(m.ChannelID, "Usage: !register <name> <age>")
            return
        }

        name := parts[1]
        age, err := strconv.Atoi(parts[2])
        if err != nil {
            s.ChannelMessageSend(m.ChannelID, "Invalid age. Please provide a number.")
            return
        }

        _, err = db.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", name, age)
        if err != nil {
            log.Println("Error registering user:", err)
            s.ChannelMessageSend(m.ChannelID, "Error registering user.")
            return
        }

        s.ChannelMessageSend(m.ChannelID, "User registered successfully!")
    }
}