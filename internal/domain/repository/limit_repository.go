package repository

import (
	"context"
	"pt-xyz-multifinance/internal/domain/entity"
)

type LimitRepository interface {
	Create(ctx context.Context, limit *entity.CustomerLimit) error
	GetByCustomerAndTenor(ctx context.Context, customerID uint64, tenorMonths int) (*entity.CustomerLimit, error)
	GetByCustomerID(ctx context.Context, customerID uint64) ([]*entity.CustomerLimit, error)
	Update(ctx context.Context, limit *entity.CustomerLimit) error
	UpdateUsedAmount(ctx context.Context, customerID uint64, tenorMonths int, amount float64) error
	Delete(ctx context.Context, id uint64) error
}
