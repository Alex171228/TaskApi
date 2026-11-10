package postgres

import (
	"TaskAPI2/internal/task"
	"context"
	"database/sql"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db}
}

//	type Repository interface {
//		Create(ctx context.Context, t Task) (Task, error)
//		Get(ctx context.Context) ([]Task, error)
//		Delete(ctx context.Context, id int) error
//		GetbyID(ctx context.Context, id int) (Task, error)
//		Put(ctx context.Context, t Task, id int) (Task, error)
//	}
func (tr *TaskRepository) Create(ctx context.Context, t task.Task) (task.Task, error) {
	//	1 Создать запрос
	//	2 выполнить запрос
	//	3 вернуть результат запроса
	const query = `
	INSERT INTO tasks (title, description, done, created_at, updated_at)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id,title, description, done, created_at, updated_at
`
	rowscan := tr.db.QueryRowContext(ctx, query, t.Title, t.Description, t.CreatedAt, t.UpdatedAt)
	var res task.Task
	if err := rowscan.Scan(&res.ID, &res.Title, &res.Description, &res.Done, &res.CreatedAt, &res.UpdatedAt); err != nil {
		return task.Task{}, err
	}
	return res, nil
}
func (tr *TaskRepository) Get(ctx context.Context) ([]task.Task, error) {

}
func (tr *TaskRepository) Delete(ctx context.Context, id int) error {

}
func (tr *TaskRepository) GetbyID(ctx context.Context, id int) (task.Task, error) {

}
func (tr *TaskRepository) Put(ctx context.Context, t task.Task, id int) (task.Task, error) {

}
