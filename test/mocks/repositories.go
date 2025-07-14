package mocks

import (
	"context"
	"pt-xyz-multifinance/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint64) (*entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByRole(ctx context.Context, role entity.UserRole, limit, offset int) ([]*entity.User, error) {
	args := m.Called(ctx, role, limit, offset)
	return args.Get(0).([]*entity.User), args.Error(1)
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

func (m *MockCustomerRepository) GetByUserID(ctx context.Context, userID uint64) (*entity.Customer, error) {
	args := m.Called(ctx, userID)
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

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, id uint64) (*entity.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByContractNumber(ctx context.Context, contractNumber string) (*entity.Transaction, error) {
	args := m.Called(ctx, contractNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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
