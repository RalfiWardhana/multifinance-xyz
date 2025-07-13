package repository

import (
	"context"
	"fmt"
	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/domain/repository"

	"gorm.io/gorm"
)

type customerRepositoryImpl struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) repository.CustomerRepository {
	return &customerRepositoryImpl{db: db}
}

func (r *customerRepositoryImpl) Create(ctx context.Context, customer *entity.Customer) error {
	if err := r.db.WithContext(ctx).Create(customer).Error; err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}
	return nil
}

func (r *customerRepositoryImpl) GetByID(ctx context.Context, id uint64) (*entity.Customer, error) {
	var customer entity.Customer
	if err := r.db.WithContext(ctx).First(&customer, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get customer by ID: %w", err)
	}
	return &customer, nil
}

func (r *customerRepositoryImpl) GetByNIK(ctx context.Context, nik string) (*entity.Customer, error) {
	var customer entity.Customer
	if err := r.db.WithContext(ctx).Where("nik = ?", nik).First(&customer).Error; err != nil {
		return nil, fmt.Errorf("failed to get customer by NIK: %w", err)
	}
	return &customer, nil
}

func (r *customerRepositoryImpl) Update(ctx context.Context, customer *entity.Customer) error {
	if err := r.db.WithContext(ctx).Save(customer).Error; err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}
	return nil
}

func (r *customerRepositoryImpl) Delete(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Customer{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	return nil
}

func (r *customerRepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]*entity.Customer, error) {
	var customers []*entity.Customer
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&customers).Error; err != nil {
		return nil, fmt.Errorf("failed to get all customers: %w", err)
	}
	return customers, nil
}
