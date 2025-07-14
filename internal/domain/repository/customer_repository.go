package repository

import (
	"context"
	"pt-xyz-multifinance/internal/domain/entity"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *entity.Customer) error
	GetByID(ctx context.Context, id uint64) (*entity.Customer, error)
	GetByUserID(ctx context.Context, userID uint64) (*entity.Customer, error)
	GetByNIK(ctx context.Context, nik string) (*entity.Customer, error)
	Update(ctx context.Context, customer *entity.Customer) error
	Delete(ctx context.Context, id uint64) error
	GetAll(ctx context.Context, limit, offset int) ([]*entity.Customer, error)
}
