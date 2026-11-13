package postgres

import (
	"TaskAPI2/internal/task"
	"context"
	"database/sql"
	"errors"
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
	rowscan := tr.db.QueryRowContext(ctx, query, t.Title, t.Description, t.Done, t.CreatedAt, t.UpdatedAt)
	var res task.Task
	if err := rowscan.Scan(&res.ID, &res.Title, &res.Description, &res.Done, &res.CreatedAt, &res.UpdatedAt); err != nil {
		return task.Task{}, err
	}
	return res, nil
}
func (tr *TaskRepository) Get(ctx context.Context) ([]task.Task, error) {
	const query = `
	SELECT  id, title, description, done, created_at, updated_at
	FROM  tasks
	ORDER BY created_at DESC, id DESC
`
	rows, err := tr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]task.Task, 0)
	for rows.Next() {
		var t task.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Done, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
func (tr *TaskRepository) Delete(ctx context.Context, id int) error {
	const query = `
	DELETE
	FROM  tasks
	WHERE id = $1
	`
	result, err := tr.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return task.ErrNotFound
	}
	return nil
}
func (tr *TaskRepository) GetbyID(ctx context.Context, id int) (task.Task, error) {
	const query = `
	SELECT  id, title, description, done, created_at, updated_at
	FROM  tasks
	WHERE id = $1
`
	rowscan := tr.db.QueryRowContext(ctx, query, id)
	var res task.Task
	if err := rowscan.Scan(&res.ID, &res.Title, &res.Description, &res.Done, &res.CreatedAt, &res.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return task.Task{}, task.ErrNotFound
		}
		return task.Task{}, err
	}
	return res, nil

}
func (tr *TaskRepository) Put(ctx context.Context, t task.Task, id int) (task.Task, error) {
	const query = `
	UPDATE tasks
	SET (title, description, done, updated_at) = ($1, $2, $3, $4)
	WHERE id = $5
	RETURNING id, title, description, done, created_at, updated_at
`
	rowscan := tr.db.QueryRowContext(ctx, query, t.Title, t.Description, t.Done, t.UpdatedAt, id)
	var res task.Task
	if err := rowscan.Scan(&res.ID, &res.Title, &res.Description, &res.Done, &res.CreatedAt, &res.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return task.Task{}, task.ErrNotFound
		}
		return task.Task{}, err
	}
	return res, nil
}
