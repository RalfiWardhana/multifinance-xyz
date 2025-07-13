package usecase

import (
	"context"
	"fmt"
	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/domain/repository"
	"pt-xyz-multifinance/pkg/logger"
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
}

func NewCustomerUseCase(customerRepo repository.CustomerRepository, limitRepo repository.LimitRepository) CustomerUseCase {
	return &customerUseCase{
		customerRepo: customerRepo,
		limitRepo:    limitRepo,
	}
}

func (uc *customerUseCase) CreateCustomer(ctx context.Context, customer *entity.Customer, limits []*entity.CustomerLimit) error {
	// Check if NIK already exists
	existingCustomer, err := uc.customerRepo.GetByNIK(ctx, customer.NIK)
	if err == nil && existingCustomer != nil {
		return fmt.Errorf("customer with NIK %s already exists", customer.NIK)
	}

	// Create customer
	if err := uc.customerRepo.Create(ctx, customer); err != nil {
		logger.Error("Failed to create customer", "error", err)
		return err
	}

	// Create customer limits
	for _, limit := range limits {
		limit.CustomerID = customer.ID
		if err := uc.limitRepo.Create(ctx, limit); err != nil {
			logger.Error("Failed to create customer limit", "error", err)
			return err
		}
	}

	logger.Info("Customer created successfully", "customerID", customer.ID)
	return nil
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

func (uc *customerUseCase) DeleteCustomer(ctx context.Context, id uint64) error {
	return uc.customerRepo.Delete(ctx, id)
}

func (uc *customerUseCase) GetAllCustomers(ctx context.Context, limit, offset int) ([]*entity.Customer, error) {
	return uc.customerRepo.GetAll(ctx, limit, offset)
}

func (uc *customerUseCase) GetCustomerLimits(ctx context.Context, customerID uint64) ([]*entity.CustomerLimit, error) {
	return uc.limitRepo.GetByCustomerID(ctx, customerID)
}
