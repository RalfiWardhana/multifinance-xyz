package usecase

import (
	"context"
	"fmt"
	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/domain/repository"
	"pt-xyz-multifinance/pkg/logger"
	"pt-xyz-multifinance/pkg/utils"

	"gorm.io/gorm"
)

type CustomerUseCase interface {
	CreateCustomer(ctx context.Context, customer *entity.Customer, limits []*entity.CustomerLimit) error
	GetCustomerByID(ctx context.Context, id uint64) (*entity.Customer, error)
	GetCustomerByNIK(ctx context.Context, nik string) (*entity.Customer, error)
	UpdateCustomer(ctx context.Context, customer *entity.Customer) error
	DeleteCustomer(ctx context.Context, id uint64) error
	GetAllCustomers(ctx context.Context, limit, offset int) ([]*entity.Customer, error)
	GetCustomerLimits(ctx context.Context, customerID uint64) ([]*entity.CustomerLimit, error)
}

type customerUseCase struct {
	customerRepo repository.CustomerRepository
	limitRepo    repository.LimitRepository
	db           *gorm.DB
}

func NewCustomerUseCase(customerRepo repository.CustomerRepository, limitRepo repository.LimitRepository, db *gorm.DB) CustomerUseCase {
	return &customerUseCase{
		customerRepo: customerRepo,
		limitRepo:    limitRepo,
		db:           db,
	}
}

func (uc *customerUseCase) CreateCustomer(ctx context.Context, customer *entity.Customer, limits []*entity.CustomerLimit) error {

	tenorLimits := make([]utils.TenorLimit, len(limits))
	for i, limit := range limits {
		tenorLimits[i] = utils.TenorLimit{
			Tenor:       limit.TenorMonths,
			LimitAmount: limit.LimitAmount,
		}
	}

	if err := utils.ValidateTenorLimits(tenorLimits); err != nil {
		return fmt.Errorf("tenor validation failed: %w", err)
	}

	for _, limit := range limits {

		if err := utils.ValidateTenor(limit.TenorMonths); err != nil {
			return err
		}

		if limit.LimitAmount <= 0 {
			return fmt.Errorf("limit amount must be greater than 0 for tenor %d months", limit.TenorMonths)
		}
	}

	return uc.db.Transaction(func(tx *gorm.DB) error {

		existingCustomer, err := uc.customerRepo.GetByNIK(ctx, customer.NIK)
		if err == nil && existingCustomer != nil {
			return fmt.Errorf("customer with NIK %s already exists", customer.NIK)
		}

		if err := uc.customerRepo.Create(ctx, customer); err != nil {
			logger.Error("Failed to create customer", "error", err)
			return fmt.Errorf("failed to create customer: %w", err)
		}

		for _, limit := range limits {
			limit.CustomerID = customer.ID
			if err := uc.limitRepo.Create(ctx, limit); err != nil {
				logger.Error("Failed to create customer limit", "error", err)
				return fmt.Errorf("failed to create customer limit: %w", err)
			}
		}

		logger.Info("Customer created successfully with all required tenors", "customerID", customer.ID, "tenorsCount", len(limits))
		return nil
	})
}

func (uc *customerUseCase) DeleteCustomer(ctx context.Context, id uint64) error {
	return uc.db.Transaction(func(tx *gorm.DB) error {

		limits, err := uc.limitRepo.GetByCustomerID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to get customer limits: %w", err)
		}

		for _, limit := range limits {
			if err := uc.limitRepo.Delete(ctx, limit.ID); err != nil {
				return fmt.Errorf("failed to delete customer limit: %w", err)
			}
		}

		if err := uc.customerRepo.Delete(ctx, id); err != nil {
			return fmt.Errorf("failed to delete customer: %w", err)
		}

		logger.Info("Customer deleted successfully with transaction", "customerID", id)
		return nil
	})
}

func (uc *customerUseCase) GetCustomerByID(ctx context.Context, id uint64) (*entity.Customer, error) {
	return uc.customerRepo.GetByID(ctx, id)
}

func (uc *customerUseCase) GetCustomerByNIK(ctx context.Context, nik string) (*entity.Customer, error) {
	return uc.customerRepo.GetByNIK(ctx, nik)
}

func (uc *customerUseCase) UpdateCustomer(ctx context.Context, customer *entity.Customer) error {
	return uc.customerRepo.Update(ctx, customer)
}

func (uc *customerUseCase) GetAllCustomers(ctx context.Context, limit, offset int) ([]*entity.Customer, error) {
	return uc.customerRepo.GetAll(ctx, limit, offset)
}

func (uc *customerUseCase) GetCustomerLimits(ctx context.Context, customerID uint64) ([]*entity.CustomerLimit, error) {
	return uc.limitRepo.GetByCustomerID(ctx, customerID)
}
