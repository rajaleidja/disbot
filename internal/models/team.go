package models

type Team struct {
    ID        int
    Name      string
    CreatedBy string
}

type TeamMember struct {
    TeamID int
    UserID string
}