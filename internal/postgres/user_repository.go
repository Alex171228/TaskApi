package postgres

import (
	"TaskAPI2/internal/user"
	"context"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}
func (u *UserRepository) Create(ctx context.Context, email string, passwordHash string) (user.User, error) {
	const querry = `
INSERT INTO users (email,password_hash) VALUES ($1,$2)
RETURNING id,email
`
	var account user.User
	rowscan := u.db.QueryRowContext(ctx, querry, email, passwordHash)
	if err := rowscan.Scan(&account.ID, &account.Email); err != nil {
		return user.User{}, err
	}
	return account, nil
}

func (u *UserRepository) FindbyEmail(ctx context.Context, email string) (user.User, error) {
	const querry = `
SELECT id,email, password_hash FROM users
WHERE email = $1
`
	var account user.User
	rowscan := u.db.QueryRowContext(ctx, querry, email)
	if err := rowscan.Scan(&account.ID, &account.Email, &account.PasswordHash); err != nil {
		return user.User{}, err
	}
	return account, nil
}
