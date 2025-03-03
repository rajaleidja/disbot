package repository

import (
    "database/sql"
    "discord-bot/internal/models"
)

type TeamRepository struct {
    db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
    return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(team *models.Team) error {
    err := r.db.QueryRow("INSERT INTO teams (name, created_by) VALUES ($1, $2) RETURNING id", team.Name, team.CreatedBy).Scan(&team.ID)
    return err
}

func (r *TeamRepository) AddMember(teamID int, userID string) error {
    _, err := r.db.Exec("INSERT INTO team_members (team_id, user_id) VALUES ($1, $2)", teamID, userID)
    return err
}

func (r *TeamRepository) GetMembers(teamID int) ([]string, error) {
    rows, err := r.db.Query("SELECT user_id FROM team_members WHERE team_id = $1", teamID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var members []string
    for rows.Next() {
        var userID string
        if err := rows.Scan(&userID); err != nil {
            return nil, err
        }
        members = append(members, userID)
    }

    return members, nil
}

func (r *TeamRepository) GetTeamsByUser(userID string) ([]models.Team, error) {
    rows, err := r.db.Query("SELECT t.id, t.name, t.created_by FROM teams t JOIN team_members tm ON t.id = tm.team_id WHERE tm.user_id = $1", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var teams []models.Team
    for rows.Next() {
        var team models.Team
        if err := rows.Scan(&team.ID, &team.Name, &team.CreatedBy); err != nil {
            return nil, err
        }
        teams = append(teams, team)
    }

    return teams, nil
}

func (r *TeamRepository) GetByName(name string) (*models.Team, error) {
    var team models.Team
    err := r.db.QueryRow("SELECT id, name, created_by FROM teams WHERE name = $1", name).Scan(&team.ID, &team.Name, &team.CreatedBy)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &team, nil
}