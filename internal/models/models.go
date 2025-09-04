package models

import "time"

type Todo struct {
	ID          string    `json:"id" db:"id"`
	TaskName    string    `json:"taskName" db:"taskName"`
	Description *string   `json:"description" db:"description"`
	Completed   bool      `json:"completed" db:"completed"`
	CreatedAt   time.Time `json:"createdAt" db:"createdAt"`
}

type UpdateTodoRequest struct {
	TaskName    *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
}
