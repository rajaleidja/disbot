package bot

import (
	"database/sql"
	"discord-bot/internal/bot/handlers"

	"log"

	"github.com/bwmarrin/discordgo"
)

func Start(botToken string, db *sql.DB) {
    log.Println("Bot token:", botToken)

    discord, err := discordgo.New("Bot " + botToken)
    if err != nil {
        log.Fatal("Error creating Discord session: ", err)
    }

    discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
        handleMessage(s, m, db)
    })

    err = discord.Open()
    if err != nil {
        log.Fatal("Error opening connection: ", err)
    }
    defer discord.Close()

    log.Println("Bot is now running. Press CTRL+C to exit.")
    <-make(chan struct{})
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB) {
    log.Printf("Received message: %s from %s in channel %s\n", m.Content, m.Author.Username, m.ChannelID)

    if m.Author.Bot {
        log.Println("Ignoring message from bot")
        return
    }

    log.Printf("Message content: %s\n", m.Content)

    if m.Content == "!hello" {
        log.Println("Handling !hello command")
        _, err := s.ChannelMessageSend(m.ChannelID, "Hello, world!")
        if err != nil {
            log.Println("Error sending message: ", err)
        }
    }

    if m.Content == "!help" {
        _, err := s.ChannelMessageSend(m.ChannelID, "!task create\n!task list\n!task delete\n!task update\n!task filter\n!task deadline\n!task create_team\n!team create\n!team add")
        if err != nil {
            log.Println("Error sending message: ", err)
        }
    }

    if response := CheckEasterEgg(m.Content); response != "" {
        s.ChannelMessageSend(m.ChannelID, response)
        return
    }
    
    handlers.HandleUserRegister(s, m, db)

    handlers.HandleTaskCreate(s, m, db)

    handlers.HandleTaskList(s, m, db)

    handlers.HandleTaskDelete(s, m, db)

    handlers.HandleTaskUpdate(s, m, db)

    handlers.HandleTaskFilter(s, m, db)

    handlers.HandleTaskDeadline(s, m, db)

    handlers.HandleTeamCreate(s, m, db)

    handlers.HandleTeamAddMember(s, m, db)

    handlers.HandleTaskCreateTeam(s, m, db)
}