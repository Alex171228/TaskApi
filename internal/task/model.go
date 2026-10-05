package task

import (
	"errors"
	"time"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Done        bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
type Input struct {
	Title       string
	Description string
	Done        bool
}

var ErrNotFound = errors.New("Task not found")
