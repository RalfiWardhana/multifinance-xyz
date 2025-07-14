package usecase_test

import (
	"context"
	"testing"
	"time"

	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/interfaces/dto"
	"pt-xyz-multifinance/internal/usecase"
	"pt-xyz-multifinance/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UseCaseTestSuite struct {
	suite.Suite
	authUseCase        usecase.AuthUseCase
	customerUseCase    usecase.CustomerUseCase
	transactionUseCase usecase.TransactionUseCase
	userRepo           *mocks.MockUserRepository
	customerRepo       *mocks.MockCustomerRepository
	limitRepo          *mocks.MockLimitRepository
	transactionRepo    *mocks.MockTransactionRepository
	db                 *gorm.DB
}

func (suite *UseCaseTestSuite) SetupTest() {
	suite.userRepo = new(mocks.MockUserRepository)
	suite.customerRepo = new(mocks.MockCustomerRepository)
	suite.limitRepo = new(mocks.MockLimitRepository)
	suite.transactionRepo = new(mocks.MockTransactionRepository)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)
	suite.db = db

	suite.authUseCase = usecase.NewAuthUseCase(
		suite.userRepo, suite.customerRepo, suite.limitRepo, suite.db)
	suite.customerUseCase = usecase.NewCustomerUseCase(
		suite.customerRepo, suite.limitRepo, suite.db)
	suite.transactionUseCase = usecase.NewTransactionUseCase(
		suite.transactionRepo, suite.customerRepo, suite.limitRepo, suite.db)
}

func (suite *UseCaseTestSuite) TestAuthUseCase_RegisterSuccess() {
	ctx := context.Background()
	req := &dto.RegisterRequest{
		Username:        "testuser",
		Email:           "test@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
		Role:            "CUSTOMER",
		CustomerData: &dto.CreateCustomerRequest{
			NIK:        "1234567890123456",
			FullName:   "John Doe",
			LegalName:  "John Doe",
			BirthPlace: "Jakarta",
			BirthDate:  "1990-01-01",
			Salary:     5000000,
			Limits: []dto.CreateLimitRequest{
				{TenorMonths: 1, LimitAmount: 100000},
				{TenorMonths: 2, LimitAmount: 200000},
				{TenorMonths: 3, LimitAmount: 300000},
				{TenorMonths: 4, LimitAmount: 400000},
			},
		},
	}

	suite.userRepo.On("GetByUsername", mock.Anything, "testuser").Return(nil, gorm.ErrRecordNotFound)
	suite.userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, gorm.ErrRecordNotFound)
	suite.customerRepo.On("GetByNIK", mock.Anything, "1234567890123456").Return(nil, gorm.ErrRecordNotFound)

	suite.userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(1).(*entity.User)
		user.ID = 1
	})

	suite.customerRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Customer")).Return(nil).Run(func(args mock.Arguments) {
		customer := args.Get(1).(*entity.Customer)
		customer.ID = 1
	})

	suite.limitRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.CustomerLimit")).Return(nil).Times(4)

	user, customer, err := suite.authUseCase.Register(ctx, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.NotNil(suite.T(), customer)
	assert.Equal(suite.T(), "testuser", user.Username)
	assert.Equal(suite.T(), entity.RoleCustomer, user.Role)
}

func (suite *UseCaseTestSuite) TestAuthUseCase_LoginSuccess() {
	ctx := context.Background()
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &entity.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role:     entity.RoleCustomer,
		IsActive: true,
	}

	customer := &entity.Customer{
		ID:       1,
		UserID:   1,
		NIK:      "1234567890123456",
		FullName: "John Doe",
	}

	suite.userRepo.On("GetByUsername", mock.Anything, "testuser").Return(user, nil)
	suite.customerRepo.On("GetAll", mock.Anything, 1000, 0).Return([]*entity.Customer{customer}, nil)

	token, returnedUser, returnedCustomer, err := suite.authUseCase.Login(ctx, "testuser", password)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
	assert.Equal(suite.T(), user.ID, returnedUser.ID)
	assert.Equal(suite.T(), customer.ID, returnedCustomer.ID)
}

func (suite *UseCaseTestSuite) TestCustomerUseCase_CreateCustomerSuccess() {
	ctx := context.Background()
	customer := &entity.Customer{
		NIK:        "1234567890123456",
		FullName:   "John Doe",
		LegalName:  "John Doe",
		BirthPlace: "Jakarta",
		BirthDate:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Salary:     5000000,
	}

	limits := []*entity.CustomerLimit{
		{TenorMonths: 1, LimitAmount: 100000},
		{TenorMonths: 2, LimitAmount: 200000},
		{TenorMonths: 3, LimitAmount: 300000},
		{TenorMonths: 4, LimitAmount: 400000},
	}

	suite.customerRepo.On("GetByNIK", mock.Anything, "1234567890123456").Return(nil, gorm.ErrRecordNotFound)
	suite.customerRepo.On("Create", mock.Anything, customer).Return(nil).Run(func(args mock.Arguments) {
		customer.ID = 1
	})
	suite.limitRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.CustomerLimit")).Return(nil).Times(4)

	err := suite.customerUseCase.CreateCustomer(ctx, customer, limits)

	assert.NoError(suite.T(), err)
}

func (suite *UseCaseTestSuite) TestTransactionUseCase_CreateTransactionSuccess() {
	ctx := context.Background()

	customer := &entity.Customer{
		ID:       1,
		NIK:      "1234567890123456",
		FullName: "John Doe",
	}

	limit := &entity.CustomerLimit{
		ID:          1,
		CustomerID:  1,
		TenorMonths: 1,
		LimitAmount: 1000000,
		UsedAmount:  0,
	}

	transaction := &entity.Transaction{
		CustomerID:        1,
		TenorMonths:       1,
		OTRAmount:         500000,
		AdminFee:          50000,
		InterestAmount:    0,
		AssetName:         "Smartphone",
		AssetType:         entity.AssetWhiteGoods,
		TransactionSource: entity.SourceEcommerce,
		Status:            entity.StatusPending,
	}

	suite.customerRepo.On("GetByID", mock.Anything, uint64(1)).Return(customer, nil)
	suite.limitRepo.On("GetByCustomerAndTenor", mock.Anything, uint64(1), 1).Return(limit, nil)
	suite.transactionRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Transaction")).Return(nil).Run(func(args mock.Arguments) {
		tx := args.Get(1).(*entity.Transaction)
		tx.ID = 1
		tx.ContractNumber = "XYZ001"
	})
	suite.limitRepo.On("UpdateUsedAmount", mock.Anything, uint64(1), 1, float64(500000)).Return(nil)

	err := suite.transactionUseCase.CreateTransaction(ctx, transaction)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), transaction.ContractNumber)
}

func (suite *UseCaseTestSuite) TestTransactionUseCase_CreateTransactionInsufficientLimit() {
	ctx := context.Background()

	customer := &entity.Customer{
		ID:       1,
		NIK:      "1234567890123456",
		FullName: "John Doe",
	}

	limit := &entity.CustomerLimit{
		ID:          1,
		CustomerID:  1,
		TenorMonths: 1,
		LimitAmount: 1000000,
		UsedAmount:  800000,
	}

	transaction := &entity.Transaction{
		CustomerID:        1,
		TenorMonths:       1,
		OTRAmount:         500000,
		AdminFee:          50000,
		InterestAmount:    0,
		AssetName:         "Smartphone",
		AssetType:         entity.AssetWhiteGoods,
		TransactionSource: entity.SourceEcommerce,
		Status:            entity.StatusPending,
	}

	suite.customerRepo.On("GetByID", mock.Anything, uint64(1)).Return(customer, nil)
	suite.limitRepo.On("GetByCustomerAndTenor", mock.Anything, uint64(1), 1).Return(limit, nil)

	err := suite.transactionUseCase.CreateTransaction(ctx, transaction)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "exceeds available limit")
}

func TestUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UseCaseTestSuite))
}
