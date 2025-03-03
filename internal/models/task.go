package models

import "time"

type Task struct {
    ID          int 
    Description string
    Deadline    *time.Time
    UserID      string 
    TeamID      int 
    AssigneeID  string
}