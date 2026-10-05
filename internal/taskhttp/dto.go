package taskhttp

import (
	"TaskAPI2/internal/task"
	"time"
)

type TaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}
type TaskAnswer struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ConvertRequest(request TaskRequest) task.Input {
	var t task.Input
	t.Title = request.Title
	t.Description = request.Description
	t.Done = request.Done
	return t
}

func ConvertAnswer(t task.Task) TaskAnswer {
	var a TaskAnswer
	a.ID = t.ID
	a.Title = t.Title
	a.Description = t.Description
	a.Done = t.Done
	a.CreatedAt = t.CreatedAt
	a.UpdatedAt = t.UpdatedAt
	return a
}
