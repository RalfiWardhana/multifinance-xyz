package repository_test

import (
	"context"
	"testing"
	"time"

	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/domain/repository"
	repoImpl "pt-xyz-multifinance/internal/infrastructure/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type RepositoryTestSuite struct {
	suite.Suite
	db              *gorm.DB
	userRepo        repository.UserRepository
	customerRepo    repository.CustomerRepository
	limitRepo       repository.LimitRepository
	transactionRepo repository.TransactionRepository
}

func (suite *RepositoryTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	err = db.AutoMigrate(&entity.User{}, &entity.Customer{}, &entity.CustomerLimit{}, &entity.Transaction{})
	suite.Require().NoError(err)

	suite.db = db
	suite.userRepo = repoImpl.NewUserRepository(db)
	suite.customerRepo = repoImpl.NewCustomerRepository(db)
	suite.limitRepo = repoImpl.NewLimitRepository(db)
	suite.transactionRepo = repoImpl.NewTransactionRepository(db)
}

func (suite *RepositoryTestSuite) TestUserRepository() {
	ctx := context.Background()

	user := &entity.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     entity.RoleCustomer,
		IsActive: true,
	}

	err := suite.userRepo.Create(ctx, user)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), user.ID)

	retrievedUser, err := suite.userRepo.GetByID(ctx, user.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Username, retrievedUser.Username)

	userByUsername, err := suite.userRepo.GetByUsername(ctx, "testuser")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.ID, userByUsername.ID)

	user.Email = "updated@example.com"
	err = suite.userRepo.Update(ctx, user)
	assert.NoError(suite.T(), err)

	updatedUser, err := suite.userRepo.GetByID(ctx, user.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "updated@example.com", updatedUser.Email)
}

func (suite *RepositoryTestSuite) TestCustomerRepository() {
	ctx := context.Background()

	user := &entity.User{
		Username: "customer",
		Email:    "customer@example.com",
		Password: "hashedpassword",
		Role:     entity.RoleCustomer,
		IsActive: true,
	}
	err := suite.userRepo.Create(ctx, user)
	suite.Require().NoError(err)

	customer := &entity.Customer{
		UserID:     user.ID,
		NIK:        "1234567890123456",
		FullName:   "John Doe",
		LegalName:  "John Doe",
		BirthPlace: "Jakarta",
		BirthDate:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Salary:     5000000,
	}

	err = suite.customerRepo.Create(ctx, customer)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), customer.ID)

	retrievedCustomer, err := suite.customerRepo.GetByID(ctx, customer.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), customer.NIK, retrievedCustomer.NIK)

	customerByNIK, err := suite.customerRepo.GetByNIK(ctx, "1234567890123456")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), customer.ID, customerByNIK.ID)

	customerByUserID, err := suite.customerRepo.GetByUserID(ctx, user.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), customer.ID, customerByUserID.ID)
}

func (suite *RepositoryTestSuite) TestLimitRepository() {
	ctx := context.Background()

	user := &entity.User{
		Username: "limituser",
		Email:    "limit@example.com",
		Password: "hashedpassword",
		Role:     entity.RoleCustomer,
		IsActive: true,
	}
	err := suite.userRepo.Create(ctx, user)
	suite.Require().NoError(err)

	customer := &entity.Customer{
		UserID:     user.ID,
		NIK:        "1234567890123457",
		FullName:   "Jane Doe",
		LegalName:  "Jane Doe",
		BirthPlace: "Bandung",
		BirthDate:  time.Date(1992, 5, 15, 0, 0, 0, 0, time.UTC),
		Salary:     6000000,
	}
	err = suite.customerRepo.Create(ctx, customer)
	suite.Require().NoError(err)

	limit := &entity.CustomerLimit{
		CustomerID:  customer.ID,
		TenorMonths: 1,
		LimitAmount: 1000000,
		UsedAmount:  0,
	}

	err = suite.limitRepo.Create(ctx, limit)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), limit.ID)

	retrievedLimit, err := suite.limitRepo.GetByCustomerAndTenor(ctx, customer.ID, 1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), limit.LimitAmount, retrievedLimit.LimitAmount)

	err = suite.limitRepo.UpdateUsedAmount(ctx, customer.ID, 1, 500000)
	assert.NoError(suite.T(), err)

	updatedLimit, err := suite.limitRepo.GetByCustomerAndTenor(ctx, customer.ID, 1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), float64(500000), updatedLimit.UsedAmount)
}

func (suite *RepositoryTestSuite) TestTransactionRepository() {
	ctx := context.Background()

	user := &entity.User{
		Username: "transuser",
		Email:    "trans@example.com",
		Password: "hashedpassword",
		Role:     entity.RoleCustomer,
		IsActive: true,
	}
	err := suite.userRepo.Create(ctx, user)
	suite.Require().NoError(err)

	customer := &entity.Customer{
		UserID:     user.ID,
		NIK:        "1234567890123458",
		FullName:   "Bob Smith",
		LegalName:  "Bob Smith",
		BirthPlace: "Surabaya",
		BirthDate:  time.Date(1988, 8, 20, 0, 0, 0, 0, time.UTC),
		Salary:     7000000,
	}
	err = suite.customerRepo.Create(ctx, customer)
	suite.Require().NoError(err)

	transaction := &entity.Transaction{
		ContractNumber:    "XYZ001",
		CustomerID:        customer.ID,
		TenorMonths:       1,
		OTRAmount:         500000,
		AdminFee:          50000,
		InstallmentAmount: 550000,
		InterestAmount:    0,
		AssetName:         "Smartphone",
		AssetType:         entity.AssetWhiteGoods,
		Status:            entity.StatusPending,
		TransactionSource: entity.SourceEcommerce,
	}

	err = suite.transactionRepo.Create(ctx, transaction)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), transaction.ID)

	retrievedTransaction, err := suite.transactionRepo.GetByID(ctx, transaction.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), transaction.ContractNumber, retrievedTransaction.ContractNumber)

	transactionByContract, err := suite.transactionRepo.GetByContractNumber(ctx, "XYZ001")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), transaction.ID, transactionByContract.ID)

	transaction.Status = entity.StatusApproved
	err = suite.transactionRepo.Update(ctx, transaction)
	assert.NoError(suite.T(), err)

	updatedTransaction, err := suite.transactionRepo.GetByID(ctx, transaction.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), entity.StatusApproved, updatedTransaction.Status)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
