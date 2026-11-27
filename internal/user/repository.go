package user

import "context"

type Repository interface {
	Create(ctx context.Context, email string, passwordHash string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
}
