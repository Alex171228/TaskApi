package memory

import (
	"TaskAPI2/internal/task"
	"context"
	"sync"
)

type MemoryRepository struct {
	tasks []task.Task
	id    int
	mu    sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		id: 1,
	}
}

func (m *MemoryRepository) Create(ctx context.Context, t task.Task) (task.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t.ID = m.id
	m.tasks = append(m.tasks, t)
	m.id++
	return t, nil
}

func (m *MemoryRepository) Get(ctx context.Context) ([]task.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]task.Task, len(m.tasks))
	copy(result, m.tasks)
	return result, nil
}

func (m *MemoryRepository) Delete(ctx context.Context, id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, t := range m.tasks {
		if t.ID == id {
			m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
			return nil
		}
	}
	return task.ErrNotFound
}

func (m *MemoryRepository) GetbyID(ctx context.Context, id int) (task.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, t := range m.tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return task.Task{}, task.ErrNotFound
}

func (m *MemoryRepository) Put(ctx context.Context, t task.Task, id int) (task.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, tk := range m.tasks {
		if tk.ID == id {
			t.ID = tk.ID
			t.CreatedAt = tk.CreatedAt
			m.tasks[i] = t
			return t, nil
		}
	}
	return task.Task{}, task.ErrNotFound
}
