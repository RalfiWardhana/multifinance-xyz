package repository

import (
	"context"
	"fmt"
	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/domain/repository"

	"gorm.io/gorm"
)

type limitRepositoryImpl struct {
	db *gorm.DB
}

func NewLimitRepository(db *gorm.DB) repository.LimitRepository {
	return &limitRepositoryImpl{db: db}
}

func (r *limitRepositoryImpl) Create(ctx context.Context, limit *entity.CustomerLimit) error {
	if err := r.db.WithContext(ctx).Create(limit).Error; err != nil {
		return fmt.Errorf("failed to create customer limit: %w", err)
	}
	return nil
}

func (r *limitRepositoryImpl) GetByCustomerAndTenor(ctx context.Context, customerID uint64, tenorMonths int) (*entity.CustomerLimit, error) {
	var limit entity.CustomerLimit
	if err := r.db.WithContext(ctx).Where("customer_id = ? AND tenor_months = ?", customerID, tenorMonths).First(&limit).Error; err != nil {
		return nil, fmt.Errorf("failed to get customer limit: %w", err)
	}
	return &limit, nil
}

func (r *limitRepositoryImpl) GetByCustomerID(ctx context.Context, customerID uint64) ([]*entity.CustomerLimit, error) {
	var limits []*entity.CustomerLimit
	if err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&limits).Error; err != nil {
		return nil, fmt.Errorf("failed to get customer limits: %w", err)
	}
	return limits, nil
}

func (r *limitRepositoryImpl) Update(ctx context.Context, limit *entity.CustomerLimit) error {
	if err := r.db.WithContext(ctx).Save(limit).Error; err != nil {
		return fmt.Errorf("failed to update customer limit: %w", err)
	}
	return nil
}

func (r *limitRepositoryImpl) UpdateUsedAmount(ctx context.Context, customerID uint64, tenorMonths int, amount float64) error {

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var limit entity.CustomerLimit

		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("customer_id = ? AND tenor_months = ?", customerID, tenorMonths).
			First(&limit).Error; err != nil {
			return fmt.Errorf("failed to get customer limit for update: %w", err)
		}

		// Calculate new used amount
		newUsedAmount := limit.UsedAmount + amount

		// Validate that used amount doesn't go negative (for rollbacks)
		if newUsedAmount < 0 {
			newUsedAmount = 0
		}

		// Validate that used amount doesn't exceed limit amount
		if newUsedAmount > limit.LimitAmount {
			return fmt.Errorf("used amount %.2f would exceed limit amount %.2f", newUsedAmount, limit.LimitAmount)
		}

		// Update the used amount
		limit.UsedAmount = newUsedAmount

		if err := tx.Save(&limit).Error; err != nil {
			return fmt.Errorf("failed to update used amount: %w", err)
		}

		return nil
	})
}

func (r *limitRepositoryImpl) Delete(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Delete(&entity.CustomerLimit{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete customer limit: %w", err)
	}
	return nil
}
