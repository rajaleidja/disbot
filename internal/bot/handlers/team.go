package handlers

import (
    "database/sql"
    "discord-bot/internal/database/repository"
    "discord-bot/internal/models"
    "fmt"
    "log"
    "strings"

    "github.com/bwmarrin/discordgo"
)

func HandleTeamCreate(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB) {
    if m.Author.Bot {
        return
    }

    if strings.HasPrefix(m.Content, "!team create") {
        parts := strings.Split(m.Content, " ")
        if len(parts) < 3 {
            s.ChannelMessageSend(m.ChannelID, "Usage: !team create <name>")
            return
        }

        name := strings.Join(parts[2:], " ")

        teamRepo := repository.NewTeamRepository(db)
        err := teamRepo.Create(&models.Team{
            Name:      name,
            CreatedBy: m.Author.ID,
        })
        if err != nil {
            log.Println("Error creating team:", err)
            s.ChannelMessageSend(m.ChannelID, "Error creating team.")
            return
        }

        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Team '%s' created successfully!", name))
    }
}

func HandleTeamAddMember(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB) {
    if m.Author.Bot {
        return
    }

    if strings.HasPrefix(m.Content, "!team add") {
        parts := strings.Split(m.Content, " ")
        if len(parts) != 4 {
            s.ChannelMessageSend(m.ChannelID, "Usage: !team add <team_name> <user_id>")
            return
        }

        teamName := parts[2]
        userID := parts[3]

        teamRepo := repository.NewTeamRepository(db)
        team, err := teamRepo.GetByName(teamName)
        if err != nil {
            log.Println("Error fetching team:", err)
            s.ChannelMessageSend(m.ChannelID, "Error fetching team.")
            return
        }

		if team == nil {
            s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Team '%s' not found.", teamName))
            return
        }

        err = teamRepo.AddMember(team.ID, userID)
        if err != nil {
            log.Println("Error adding member:", err)
            s.ChannelMessageSend(m.ChannelID, "Error adding member.")
            return
        }

        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User %s added to team '%s'!", userID, teamName))
    }
}