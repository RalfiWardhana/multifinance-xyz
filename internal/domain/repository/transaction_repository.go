package repository

import (
	"context"
	"pt-xyz-multifinance/internal/domain/entity"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	GetByID(ctx context.Context, id uint64) (*entity.Transaction, error)
	GetByContractNumber(ctx context.Context, contractNumber string) (*entity.Transaction, error)
	GetByCustomerID(ctx context.Context, customerID uint64) ([]*entity.Transaction, error)
	Update(ctx context.Context, transaction *entity.Transaction) error
	Delete(ctx context.Context, id uint64) error
	GetAll(ctx context.Context, limit, offset int) ([]*entity.Transaction, error)
}
