package models

import "time"

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Task struct {
	TaskID      int       `json:"task_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      int       `json:"user_id"`
}
