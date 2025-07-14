package repository

import (
	"context"
	"pt-xyz-multifinance/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uint64) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint64) error
	GetAll(ctx context.Context, limit, offset int) ([]*entity.User, error)
	GetByRole(ctx context.Context, role entity.UserRole, limit, offset int) ([]*entity.User, error)
}
