package task

import (
	"context"
	"errors"
	"strings"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) Create(ctx context.Context, input Input) (Task, error) {
	title := strings.TrimSpace(input.Title)
	if len(title) < 3 || len(input.Description) < 3 {
		return Task{}, ErrInvalidInput
	}
	now := time.Now()
	t := Task{
		Title:       title,
		Description: input.Description,
		Done:        input.Done,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return s.repo.Create(ctx, t)
}

func (s *Service) List(ctx context.Context) ([]Task, error) {
	return s.repo.Get(ctx)
}

func (s *Service) Get(ctx context.Context, id int) (Task, error) {
	if id < 0 {
		return Task{}, errors.New("Not valid id")
	}
	return s.repo.GetbyID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id < 0 {
		return errors.New("Not valid id")
	}
	return s.repo.Delete(ctx, id)
}
func (s *Service) Update(ctx context.Context, id int, input Input) (Task, error) {
	title := strings.TrimSpace(input.Title)
	if len(title) < 3 || len(input.Description) < 3 {
		return Task{}, ErrInvalidInput
	}
	task, err := s.repo.GetbyID(ctx, id)
	if err != nil {
		return Task{}, err
	}
	task.Title = title
	task.Description = input.Description
	task.Done = input.Done
	task.UpdatedAt = time.Now()
	return s.repo.Put(ctx, task, id)
}
