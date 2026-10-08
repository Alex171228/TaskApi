package task

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, t Task) (Task, error)
	Get(ctx context.Context) ([]Task, error)
	Delete(ctx context.Context, id int) error
	GetbyID(ctx context.Context, id int) (Task, error)
	Put(ctx context.Context, t Task, id int) (Task, error)
}
