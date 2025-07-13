package usecase_test

import (
	"context"
	"errors"
	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, id uint64) (*entity.Transaction, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByContractNumber(ctx context.Context, contractNumber string) (*entity.Transaction, error) {
	args := m.Called(ctx, contractNumber)
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByCustomerID(ctx context.Context, customerID uint64) ([]*entity.Transaction, error) {
	args := m.Called(ctx, customerID)
	return args.Get(0).([]*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.Transaction, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*entity.Transaction), args.Error(1)
}

type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) Create(ctx context.Context, customer *entity.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockCustomerRepository) GetByID(ctx context.Context, id uint64) (*entity.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Customer), args.Error(1)
}

func (m *MockCustomerRepository) GetByNIK(ctx context.Context, nik string) (*entity.Customer, error) {
	args := m.Called(ctx, nik)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Customer), args.Error(1)
}

func (m *MockCustomerRepository) Update(ctx context.Context, customer *entity.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockCustomerRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCustomerRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.Customer, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*entity.Customer), args.Error(1)
}

type MockLimitRepository struct {
	mock.Mock
}

func (m *MockLimitRepository) Create(ctx context.Context, limit *entity.CustomerLimit) error {
	args := m.Called(ctx, limit)
	return args.Error(0)
}

func (m *MockLimitRepository) GetByCustomerAndTenor(ctx context.Context, customerID uint64, tenorMonths int) (*entity.CustomerLimit, error) {
	args := m.Called(ctx, customerID, tenorMonths)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CustomerLimit), args.Error(1)
}

func (m *MockLimitRepository) GetByCustomerID(ctx context.Context, customerID uint64) ([]*entity.CustomerLimit, error) {
	args := m.Called(ctx, customerID)
	return args.Get(0).([]*entity.CustomerLimit), args.Error(1)
}

func (m *MockLimitRepository) Update(ctx context.Context, limit *entity.CustomerLimit) error {
	args := m.Called(ctx, limit)
	return args.Error(0)
}

func (m *MockLimitRepository) UpdateUsedAmount(ctx context.Context, customerID uint64, tenorMonths int, amount float64) error {
	args := m.Called(ctx, customerID, tenorMonths, amount)
	return args.Error(0)
}

func (m *MockLimitRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTransactionUseCase_CreateTransaction_Success(t *testing.T) {
	// Setup
	mockTransactionRepo := new(MockTransactionRepository)
	mockCustomerRepo := new(MockCustomerRepository)
	mockLimitRepo := new(MockLimitRepository)

	transactionUC := usecase.NewTransactionUseCase(
		mockTransactionRepo,
		mockCustomerRepo,
		mockLimitRepo,
	)

	ctx := context.Background()
	customer := &entity.Customer{
		ID:       1,
		NIK:      "1234567890123456",
		FullName: "Test Customer",
	}

	limit := &entity.CustomerLimit{
		ID:              1,
		CustomerID:      1,
		TenorMonths:     3,
		LimitAmount:     1000000.0,
		UsedAmount:      0.0,
		AvailableAmount: 1000000.0,
	}

	transaction := &entity.Transaction{
		CustomerID:        1,
		TenorMonths:       3,
		OTRAmount:         500000.0,
		AdminFee:          50000.0,
		InterestAmount:    25000.0,
		AssetName:         "Honda Beat",
		AssetType:         entity.AssetMotor,
		TransactionSource: entity.SourceWeb,
	}

	// Mock expectations
	mockCustomerRepo.On("GetByID", ctx, uint64(1)).Return(customer, nil)
	mockLimitRepo.On("GetByCustomerAndTenor", ctx, uint64(1), 3).Return(limit, nil)
	mockTransactionRepo.On("Create", ctx, mock.AnythingOfType("*entity.Transaction")).Return(nil)
	mockLimitRepo.On("UpdateUsedAmount", ctx, uint64(1), 3, 500000.0).Return(nil)

	// Execute
	err := transactionUC.CreateTransaction(ctx, transaction)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, transaction.ContractNumber)
	assert.Equal(t, float64(191666.67), transaction.InstallmentAmount)

	// Verify all expectations were met
	mockTransactionRepo.AssertExpectations(t)
	mockCustomerRepo.AssertExpectations(t)
	mockLimitRepo.AssertExpectations(t)
}

func TestTransactionUseCase_CreateTransaction_CustomerNotFound(t *testing.T) {
	// Setup
	mockTransactionRepo := new(MockTransactionRepository)
	mockCustomerRepo := new(MockCustomerRepository)
	mockLimitRepo := new(MockLimitRepository)

	transactionUC := usecase.NewTransactionUseCase(
		mockTransactionRepo,
		mockCustomerRepo,
		mockLimitRepo,
	)

	ctx := context.Background()
	transaction := &entity.Transaction{
		CustomerID: 999,
	}

	// Mock expectations
	mockCustomerRepo.On("GetByID", ctx, uint64(999)).Return((*entity.Customer)(nil), errors.New("customer not found"))

	// Execute
	err := transactionUC.CreateTransaction(ctx, transaction)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "customer not found")

	mockCustomerRepo.AssertExpectations(t)
}

func TestTransactionUseCase_ValidateTransactionLimit_ExceedsLimit(t *testing.T) {
	// Setup
	mockTransactionRepo := new(MockTransactionRepository)
	mockCustomerRepo := new(MockCustomerRepository)
	mockLimitRepo := new(MockLimitRepository)

	transactionUC := usecase.NewTransactionUseCase(
		mockTransactionRepo,
		mockCustomerRepo,
		mockLimitRepo,
	)

	ctx := context.Background()
	limit := &entity.CustomerLimit{
		ID:              1,
		CustomerID:      1,
		TenorMonths:     3,
		LimitAmount:     1000000.0,
		UsedAmount:      800000.0,
		AvailableAmount: 200000.0,
	}

	// Mock expectations
	mockLimitRepo.On("GetByCustomerAndTenor", ctx, uint64(1), 3).Return(limit, nil)

	// Execute
	err := transactionUC.ValidateTransactionLimit(ctx, 1, 3, 300000.0)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transaction amount exceeds available limit")

	mockLimitRepo.AssertExpectations(t)
}

func TestTransactionUseCase_GetTransactionByID_Success(t *testing.T) {
	// Setup
	mockTransactionRepo := new(MockTransactionRepository)
	mockCustomerRepo := new(MockCustomerRepository)
	mockLimitRepo := new(MockLimitRepository)

	transactionUC := usecase.NewTransactionUseCase(
		mockTransactionRepo,
		mockCustomerRepo,
		mockLimitRepo,
	)

	ctx := context.Background()
	expectedTransaction := &entity.Transaction{
		ID:             1,
		ContractNumber: "XYZ123456",
		CustomerID:     1,
		OTRAmount:      500000.0,
	}

	// Mock expectations
	mockTransactionRepo.On("GetByID", ctx, uint64(1)).Return(expectedTransaction, nil)

	// Execute
	result, err := transactionUC.GetTransactionByID(ctx, 1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedTransaction, result)

	mockTransactionRepo.AssertExpectations(t)
}

func TestTransactionUseCase_UpdateTransactionStatus_Success(t *testing.T) {
	// Setup
	mockTransactionRepo := new(MockTransactionRepository)
	mockCustomerRepo := new(MockCustomerRepository)
	mockLimitRepo := new(MockLimitRepository)

	transactionUC := usecase.NewTransactionUseCase(
		mockTransactionRepo,
		mockCustomerRepo,
		mockLimitRepo,
	)

	ctx := context.Background()
	transaction := &entity.Transaction{
		ID:     1,
		Status: entity.StatusPending,
	}

	// Mock expectations
	mockTransactionRepo.On("GetByID", ctx, uint64(1)).Return(transaction, nil)
	mockTransactionRepo.On("Update", ctx, mock.AnythingOfType("*entity.Transaction")).Return(nil)

	// Execute
	err := transactionUC.UpdateTransactionStatus(ctx, 1, entity.StatusApproved)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, entity.StatusApproved, transaction.Status)

	mockTransactionRepo.AssertExpectations(t)
}
