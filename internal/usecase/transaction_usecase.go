package usecase

import (
	"context"
	"fmt"
	"math"
	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/domain/repository"
	"pt-xyz-multifinance/pkg/logger"
	"time"
)

type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, transaction *entity.Transaction) error
	GetTransactionByID(ctx context.Context, id uint64) (*entity.Transaction, error)
	GetTransactionByContractNumber(ctx context.Context, contractNumber string) (*entity.Transaction, error)
	GetTransactionsByCustomerID(ctx context.Context, customerID uint64) ([]*entity.Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id uint64, status entity.TransactionStatus) error
	GetAllTransactions(ctx context.Context, limit, offset int) ([]*entity.Transaction, error)
	ValidateTransactionLimit(ctx context.Context, customerID uint64, tenorMonths int, amount float64) error
}

type transactionUseCase struct {
	transactionRepo repository.TransactionRepository
	customerRepo    repository.CustomerRepository
	limitRepo       repository.LimitRepository
}

func NewTransactionUseCase(
	transactionRepo repository.TransactionRepository,
	customerRepo repository.CustomerRepository,
	limitRepo repository.LimitRepository,
) TransactionUseCase {
	return &transactionUseCase{
		transactionRepo: transactionRepo,
		customerRepo:    customerRepo,
		limitRepo:       limitRepo,
	}
}

func (uc *transactionUseCase) CreateTransaction(ctx context.Context, transaction *entity.Transaction) error {
	// Validate customer exists
	_, err := uc.customerRepo.GetByID(ctx, transaction.CustomerID)
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}

	// Validate transaction limit
	if err := uc.ValidateTransactionLimit(ctx, transaction.CustomerID, transaction.TenorMonths, transaction.OTRAmount); err != nil {
		return err
	}

	// Generate contract number
	transaction.ContractNumber = uc.generateContractNumber()

	// Calculate installment amount if not provided
	if transaction.InstallmentAmount == 0 {
		transaction.InstallmentAmount = uc.calculateInstallmentAmount(
			transaction.OTRAmount,
			transaction.AdminFee,
			transaction.InterestAmount,
			transaction.TenorMonths,
		)
	}

	// Create transaction
	if err := uc.transactionRepo.Create(ctx, transaction); err != nil {
		logger.Error("Failed to create transaction", "error", err)
		return err
	}

	// Update used limit
	if err := uc.limitRepo.UpdateUsedAmount(ctx, transaction.CustomerID, transaction.TenorMonths, transaction.OTRAmount); err != nil {
		logger.Error("Failed to update used amount", "error", err)
		return err
	}

	logger.Info("Transaction created successfully", "transactionID", transaction.ID, "contractNumber", transaction.ContractNumber)
	return nil
}

func (uc *transactionUseCase) ValidateTransactionLimit(ctx context.Context, customerID uint64, tenorMonths int, amount float64) error {
	limit, err := uc.limitRepo.GetByCustomerAndTenor(ctx, customerID, tenorMonths)
	if err != nil {
		return fmt.Errorf("customer limit not found for tenor %d months", tenorMonths)
	}

	if limit.UsedAmount+amount > limit.LimitAmount {
		return fmt.Errorf("transaction amount exceeds available limit. Available: %.2f, Requested: %.2f",
			limit.LimitAmount-limit.UsedAmount, amount)
	}

	return nil
}

func (uc *transactionUseCase) generateContractNumber() string {
	return fmt.Sprintf("XYZ%d", time.Now().Unix())
}

func (uc *transactionUseCase) calculateInstallmentAmount(otr, adminFee, interest float64, tenor int) float64 {
	totalAmount := otr + adminFee + interest
	return math.Round((totalAmount/float64(tenor))*100) / 100
}

func (uc *transactionUseCase) GetTransactionByID(ctx context.Context, id uint64) (*entity.Transaction, error) {
	return uc.transactionRepo.GetByID(ctx, id)
}

func (uc *transactionUseCase) GetTransactionByContractNumber(ctx context.Context, contractNumber string) (*entity.Transaction, error) {
	return uc.transactionRepo.GetByContractNumber(ctx, contractNumber)
}

func (uc *transactionUseCase) GetTransactionsByCustomerID(ctx context.Context, customerID uint64) ([]*entity.Transaction, error) {
	return uc.transactionRepo.GetByCustomerID(ctx, customerID)
}

func (uc *transactionUseCase) UpdateTransactionStatus(ctx context.Context, id uint64, status entity.TransactionStatus) error {
	transaction, err := uc.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	transaction.Status = status
	return uc.transactionRepo.Update(ctx, transaction)
}

func (uc *transactionUseCase) GetAllTransactions(ctx context.Context, limit, offset int) ([]*entity.Transaction, error) {
	return uc.transactionRepo.GetAll(ctx, limit, offset)
}
