package models

type Todo struct {
	ID          string `json:"id" db:"id"`
	TaskName    string `json:"taskName" db: "taskName"`
	Description string `json:"Description" db:"Description"`
}
