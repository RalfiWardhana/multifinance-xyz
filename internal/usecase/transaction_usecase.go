package usecase

import (
	"context"
	"fmt"
	"math"
	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/domain/repository"
	"pt-xyz-multifinance/pkg/logger"
	"pt-xyz-multifinance/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, transaction *entity.Transaction) error
	GetTransactionByID(ctx context.Context, id uint64) (*entity.Transaction, error)
	GetTransactionByContractNumber(ctx context.Context, contractNumber string) (*entity.Transaction, error)
	GetTransactionsByCustomerID(ctx context.Context, customerID uint64) ([]*entity.Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id uint64, status entity.TransactionStatus) error
	GetAllTransactions(ctx context.Context, limit, offset int) ([]*entity.Transaction, error)
	ValidateTransactionLimit(ctx context.Context, customerID uint64, tenorMonths int, amount float64) error
	ApproveTransaction(ctx context.Context, id uint64) error
	RejectTransaction(ctx context.Context, id uint64, reason string) error
}

type transactionUseCase struct {
	transactionRepo repository.TransactionRepository
	customerRepo    repository.CustomerRepository
	limitRepo       repository.LimitRepository
	db              *gorm.DB
}

func NewTransactionUseCase(
	transactionRepo repository.TransactionRepository,
	customerRepo repository.CustomerRepository,
	limitRepo repository.LimitRepository,
	db *gorm.DB,
) TransactionUseCase {
	return &transactionUseCase{
		transactionRepo: transactionRepo,
		customerRepo:    customerRepo,
		limitRepo:       limitRepo,
		db:              db,
	}
}

func (uc *transactionUseCase) CreateTransaction(ctx context.Context, transaction *entity.Transaction) error {
	if err := utils.ValidateTenor(transaction.TenorMonths); err != nil {
		return err
	}

	return uc.db.Transaction(func(tx *gorm.DB) error {

		_, err := uc.customerRepo.GetByID(ctx, transaction.CustomerID)
		if err != nil {
			return fmt.Errorf("customer not found: %w", err)
		}

		if err := uc.ValidateTransactionLimit(ctx, transaction.CustomerID, transaction.TenorMonths, transaction.OTRAmount); err != nil {
			return err
		}

		transaction.ContractNumber = uc.generateContractNumber()

		if transaction.InstallmentAmount == 0 {
			transaction.InstallmentAmount = uc.calculateInstallmentAmount(
				transaction.OTRAmount,
				transaction.AdminFee,
				transaction.InterestAmount,
				transaction.TenorMonths,
			)
		}

		if err := uc.transactionRepo.Create(ctx, transaction); err != nil {
			logger.Error("Failed to create transaction", "error", err)
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		if err := uc.limitRepo.UpdateUsedAmount(ctx, transaction.CustomerID, transaction.TenorMonths, transaction.OTRAmount); err != nil {
			logger.Error("Failed to update used amount", "error", err)
			return fmt.Errorf("failed to update used amount: %w", err)
		}

		logger.Info("Transaction created successfully with transaction", "transactionID", transaction.ID, "contractNumber", transaction.ContractNumber)
		return nil
	})
}

func (uc *transactionUseCase) ApproveTransaction(ctx context.Context, id uint64) error {
	return uc.db.Transaction(func(tx *gorm.DB) error {
		transaction, err := uc.transactionRepo.GetByID(ctx, id)
		if err != nil {
			return fmt.Errorf("transaction not found: %w", err)
		}

		if transaction.Status != entity.StatusPending {
			return fmt.Errorf("transaction cannot be approved, current status: %s", transaction.Status)
		}

		oldStatus := transaction.Status
		transaction.Status = entity.StatusApproved
		if err := uc.transactionRepo.Update(ctx, transaction); err != nil {
			return fmt.Errorf("failed to update transaction status: %w", err)
		}

		logger.Info("Transaction approved",
			"transactionID", id,
			"oldStatus", oldStatus,
			"newStatus", transaction.Status,
			"contractNumber", transaction.ContractNumber)

		return nil
	})
}

func (uc *transactionUseCase) RejectTransaction(ctx context.Context, id uint64, reason string) error {
	return uc.db.Transaction(func(tx *gorm.DB) error {
		transaction, err := uc.transactionRepo.GetByID(ctx, id)
		if err != nil {
			return fmt.Errorf("transaction not found: %w", err)
		}

		if transaction.Status != entity.StatusPending {
			return fmt.Errorf("transaction cannot be rejected, current status: %s", transaction.Status)
		}

		oldStatus := transaction.Status
		transaction.Status = entity.StatusRejected
		if err := uc.transactionRepo.Update(ctx, transaction); err != nil {
			return fmt.Errorf("failed to update transaction status: %w", err)
		}

		rollbackAmount := -transaction.OTRAmount
		if err := uc.limitRepo.UpdateUsedAmount(ctx, transaction.CustomerID, transaction.TenorMonths, rollbackAmount); err != nil {
			return fmt.Errorf("failed to rollback used amount: %w", err)
		}

		logger.Info("Transaction rejected with limit rollback",
			"transactionID", id,
			"oldStatus", oldStatus,
			"newStatus", transaction.Status,
			"reason", reason,
			"rollbackAmount", transaction.OTRAmount)

		return nil
	})
}

func (uc *transactionUseCase) UpdateTransactionStatus(ctx context.Context, id uint64, status entity.TransactionStatus) error {
	transaction, err := uc.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	transaction.Status = status
	return uc.transactionRepo.Update(ctx, transaction)
}

func (uc *transactionUseCase) ValidateTransactionLimit(ctx context.Context, customerID uint64, tenorMonths int, amount float64) error {
	limit, err := uc.limitRepo.GetByCustomerAndTenor(ctx, customerID, tenorMonths)
	if err != nil {
		return fmt.Errorf("customer limit not found for tenor %d months", tenorMonths)
	}

	availableAmount := limit.AvailableAmount()
	if amount > availableAmount {
		return fmt.Errorf("transaction amount exceeds available limit. Available: %.2f, Requested: %.2f",
			availableAmount, amount)
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

func (uc *transactionUseCase) GetAllTransactions(ctx context.Context, limit, offset int) ([]*entity.Transaction, error) {
	return uc.transactionRepo.GetAll(ctx, limit, offset)
}
