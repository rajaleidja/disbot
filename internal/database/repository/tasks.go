package repository

import (
    "database/sql"
    "discord-bot/internal/models"
)

type TaskRepository struct {
    db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
    return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *models.Task) error {
    _, err := r.db.Exec("INSERT INTO tasks (description, deadline, user_id) VALUES ($1, $2, $3)", task.Description, task.Deadline, task.UserID)
    return err
}

func (r *TaskRepository) GetAll(userID string) ([]models.Task, error) {
    rows, err := r.db.Query("SELECT id, description, deadline FROM tasks WHERE user_id = $1", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []models.Task
    for rows.Next() {
        var task models.Task
        var deadline sql.NullTime
        if err := rows.Scan(&task.ID, &task.Description, &deadline); err != nil {
            return nil, err
        }
        if deadline.Valid {
            task.Deadline = &deadline.Time
        }
        task.UserID = userID
        tasks = append(tasks, task)
    }

    return tasks, nil
}

func (r *TaskRepository) Delete(id int, userID string) error {
    _, err := r.db.Exec("DELETE FROM tasks WHERE id = $1 AND user_id = $2", id, userID)
    return err
}

func (r *TaskRepository) Update(id int, description string, userID string) error {
    _, err := r.db.Exec("UPDATE tasks SET description = $1 WHERE id = $2 AND user_id = $3", description, id, userID)
    return err
}

func (r *TaskRepository) Filter(keyword string, userID string) ([]models.Task, error) {
    rows, err := r.db.Query("SELECT id, description, deadline FROM tasks WHERE description LIKE $1 AND user_id = $2", "%"+keyword+"%", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []models.Task
    for rows.Next() {
        var task models.Task
        var deadline sql.NullTime
        if err := rows.Scan(&task.ID, &task.Description, &deadline); err != nil {
            return nil, err
        }
        if deadline.Valid {
            task.Deadline = &deadline.Time
        }
        task.UserID = userID
        tasks = append(tasks, task)
    }

    return tasks, nil
}