package repository

import (
    "database/sql"
    "discord-bot/internal/models"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
    _, err := r.db.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", user.Name, user.Age)
    return err
}

func (r *UserRepository) FindByID(id int) (*models.User, error) {
    var user models.User
    err := r.db.QueryRow("SELECT id, name, age FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Age)
    if err != nil {
        return nil, err
    }
    return &user, nil
}