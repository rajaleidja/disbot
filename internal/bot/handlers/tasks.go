package handlers

import (
	"database/sql"
	"discord-bot/internal/database/repository"
	"discord-bot/internal/models"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func HandleTaskCreate(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB) {
    if m.Author.Bot {
        return
    }

    if strings.HasPrefix(m.Content, "!task create") {
        parts := strings.Split(m.Content, " ")
        if len(parts) < 3 {
            s.ChannelMessageSend(m.ChannelID, "**Usage: !task create <description> [deadline]**")
            return
        }

        description := strings.Join(parts[2:len(parts)-1], " ")
        deadlineStr := parts[len(parts)-1]

        var deadline *time.Time
        if deadlineStr != "" {
            parsedDeadline, err := time.Parse("2006-01-02", deadlineStr)
            if err != nil {
                s.ChannelMessageSend(m.ChannelID, "Invalid deadline format. Use YYYY-MM-DD.")
                return
            }
            deadline = &parsedDeadline
        }

        taskRepo := repository.NewTaskRepository(db)
        err := taskRepo.Create(&models.Task{
            Description: description,
            Deadline: deadline,
            UserID:      m.Author.ID,
        })
        if err != nil {
            log.Println("Error creating task:", err)
            s.ChannelMessageSend(m.ChannelID, "Error creating task.")
            return
        }

        s.ChannelMessageSend(m.ChannelID, "Task created successfully!")
    }
}

func HandleTaskList(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB) {
    if m.Author.Bot {
        return
    }

    if m.Content == "!task list" {
        taskRepo := repository.NewTaskRepository(db)
        tasks, err := taskRepo.GetAll(m.Author.ID)
        if err != nil {
            log.Println("Error fetching tasks:", err)
            s.ChannelMessageSend(m.ChannelID, "Error fetching tasks.")
            return
        }

        if len(tasks) == 0 {
            s.ChannelMessageSend(m.ChannelID, "No tasks found.")
            return
        }

        var response string
        for _, task := range tasks {
            response += fmt.Sprintf("- Task #%d: \n> **Description:** %s \n> **Deadline:** %s", task.ID, task.Description, task.Deadline)
            response = response[:len(response)-21] + "\n"
        }

        s.ChannelMessageSend(m.ChannelID, response)
    }
}

func HandleTaskDelete(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB) {
    if m.Author.Bot {
        return
    }

    if strings.HasPrefix(m.Content, "!task delete") {
        parts := strings.Split(m.Content, " ")
        if len(parts) != 3 {
            s.ChannelMessageSend(m.ChannelID, "Usage: !task delete <id>")
            return
        }

        id, err := strconv.Atoi(parts[2])
        if err != nil {
            s.ChannelMessageSend(m.ChannelID, "Invalid task ID. Please provide a number.")
            return
        }

        taskRepo := repository.NewTaskRepository(db)
        err = taskRepo.Delete(id, m.Author.ID)
        if err != nil {
            log.Println("Error deleting task:", err)
            s.ChannelMessageSend(m.ChannelID, "Error deleting task.")
            return
        }

        s.ChannelMessageSend(m.ChannelID, "Task deleted successfully!")
    }
}

func HandleTaskUpdate(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB) {
    if m.Author.Bot {
        return
    }

    if strings.HasPrefix(m.Content, "!task update") {
        parts := strings.Split(m.Content, " ")
        if len(parts) < 4 {
            s.ChannelMessageSend(m.ChannelID, "Usage: !task update <id> <new_description>")
            return
        }

        id, err := strconv.Atoi(parts[2])
        if err != nil {
            s.ChannelMessageSend(m.ChannelID, "Invalid task ID. Please provide a number.")
            return
        }

        description := strings.Join(parts[3:], " ")

        taskRepo := repository.NewTaskRepository(db)
        err = taskRepo.Update(id, description, m.Author.ID)
        if err != nil {
            log.Println("Error updating task:", err)
            s.ChannelMessageSend(m.ChannelID, "Error updating task.")
            return
        }

        s.ChannelMessageSend(m.ChannelID, "Task updated successfully!")
    }
}

func HandleTaskFilter(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB) {
    if m.Author.Bot {
        return
    }

    if strings.HasPrefix(m.Content, "!task filter") {
        parts := strings.Split(m.Content, " ")
        if len(parts) < 3 {
            s.ChannelMessageSend(m.ChannelID, "Usage: !task filter <keyword>")
            return
        }

        keyword := strings.Join(parts[2:], " ")

        taskRepo := repository.NewTaskRepository(db)
        tasks, err := taskRepo.Filter(keyword, m.Author.ID)
        if err != nil {
            log.Println("Error filtering tasks:", err)
            s.ChannelMessageSend(m.ChannelID, "Error filtering tasks.")
            return
        }

        if len(tasks) == 0 {
            s.ChannelMessageSend(m.ChannelID, "No tasks found.")
            return
        }

        var response string
        for _, task := range tasks {
            response += fmt.Sprintf("Task #%d: %s\n", task.ID, task.Description)
        }

        s.ChannelMessageSend(m.ChannelID, response)
    }
}

func HandleTaskDeadline(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB) {
    if m.Author.Bot {
        return
    }

    if strings.HasPrefix(m.Content, "!task deadline") {
        parts := strings.Split(m.Content, " ")
        if len(parts) != 3 {
            s.ChannelMessageSend(m.ChannelID, "Usage: !task deadline <id>")
            return
        }

        id, err := strconv.Atoi(parts[2])
        if err != nil {
            s.ChannelMessageSend(m.ChannelID, "Invalid task ID. Please provide a number.")
            return
        }

        taskRepo := repository.NewTaskRepository(db)
        tasks, err := taskRepo.GetAll(m.Author.ID)
        if err != nil {
            log.Println("Error fetching tasks:", err)
            s.ChannelMessageSend(m.ChannelID, "Error fetching tasks.")
            return
        }

        var task *models.Task
        for _, t := range tasks {
            if t.ID == id {
                task = &t
                break
            }
        }

        if task == nil {
            s.ChannelMessageSend(m.ChannelID, "Task not found.")
            return
        }

        if task.Deadline == nil {
            s.ChannelMessageSend(m.ChannelID, "This task has no deadline.")
            return
        }

        remaining := time.Until(*task.Deadline)
        if remaining < 0 {
            s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The deadline for this task has passed(Was: %s)!", remaining.Round(time.Second)))
            return
        }

        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Time remaining for task #%d: %s", task.ID, remaining.Round(time.Second)))
    }
}

func HandleTaskCreateTeam(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB) {
    if m.Author.Bot {
        return
    }

    if strings.HasPrefix(m.Content, "!task create_team") {
        parts := strings.Split(m.Content, " ")
        if len(parts) < 5 {
            s.ChannelMessageSend(m.ChannelID, "Usage: !task create_team <team_name> <description> <assignee_id> [deadline]")
            return
        }

        teamName := parts[2]
        description := strings.Join(parts[3:len(parts)-2], " ")
        assigneeID := parts[len(parts)-3]
        fmt.Println(parts[3]) 
        fmt.Println(parts[4])
        fmt.Println(parts[len(parts)-3])
        deadlineStr := parts[len(parts)-1]

        var deadline *time.Time
        if deadlineStr != "" {
            parsedDeadline, err := time.Parse("2006-01-02", deadlineStr)
            if err != nil {
                s.ChannelMessageSend(m.ChannelID, "Invalid deadline format. Use YYYY-MM-DD.")
                return
            }
            deadline = &parsedDeadline
        }

        teamRepo := repository.NewTeamRepository(db)
        team, err := teamRepo.GetByName(teamName)
        if err != nil {
            log.Println("Error fetching team:", err)
            s.ChannelMessageSend(m.ChannelID, "Error fetching team.")
            return
        }

        taskRepo := repository.NewTaskRepository(db)
        err = taskRepo.Create(&models.Task{
            Description: description,
            Deadline:    deadline,
            UserID:      m.Author.ID,
            TeamID:      team.ID,
            AssigneeID:  assigneeID,
        })
        if err != nil {
            log.Println("Error creating task:", err)
            s.ChannelMessageSend(m.ChannelID, "Error creating task.")
            return
        }

        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Task created in team **%s**! \n**Assignee:** %s", teamName, assigneeID))
    }
}