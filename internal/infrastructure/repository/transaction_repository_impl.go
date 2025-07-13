package repository

import (
	"context"
	"fmt"
	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/domain/repository"

	"gorm.io/gorm"
)

type transactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) repository.TransactionRepository {
	return &transactionRepositoryImpl{db: db}
}

func (r *transactionRepositoryImpl) Create(ctx context.Context, transaction *entity.Transaction) error {
	if err := r.db.WithContext(ctx).Create(transaction).Error; err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

func (r *transactionRepositoryImpl) GetByID(ctx context.Context, id uint64) (*entity.Transaction, error) {
	var transaction entity.Transaction
	if err := r.db.WithContext(ctx).Preload("Customer").First(&transaction, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get transaction by ID: %w", err)
	}
	return &transaction, nil
}

func (r *transactionRepositoryImpl) GetByContractNumber(ctx context.Context, contractNumber string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	if err := r.db.WithContext(ctx).Preload("Customer").Where("contract_number = ?", contractNumber).First(&transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to get transaction by contract number: %w", err)
	}
	return &transaction, nil
}

func (r *transactionRepositoryImpl) GetByCustomerID(ctx context.Context, customerID uint64) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	if err := r.db.WithContext(ctx).Preload("Customer").Where("customer_id = ?", customerID).Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactions by customer ID: %w", err)
	}
	return transactions, nil
}
func (r *transactionRepositoryImpl) Update(ctx context.Context, transaction *entity.Transaction) error {
	if err := r.db.WithContext(ctx).Save(transaction).Error; err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}
	return nil
}

func (r *transactionRepositoryImpl) Delete(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Transaction{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}
	return nil
}

func (r *transactionRepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	if err := r.db.WithContext(ctx).Preload("Customer").Limit(limit).Offset(offset).Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get all transactions: %w", err)
	}
	return transactions, nil
}
