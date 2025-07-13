package dto

import (
	"pt-xyz-multifinance/internal/domain/entity"
	"time"
)

type CreateTransactionRequest struct {
	CustomerID        uint64                   `json:"customer_id" binding:"required"`
	TenorMonths       int                      `json:"tenor_months" binding:"required,min=1"`
	OTRAmount         float64                  `json:"otr_amount" binding:"required,min=0"`
	AdminFee          float64                  `json:"admin_fee" binding:"required,min=0"`
	InterestAmount    float64                  `json:"interest_amount" binding:"required,min=0"`
	AssetName         string                   `json:"asset_name" binding:"required,min=2"`
	AssetType         entity.AssetType         `json:"asset_type" binding:"required"`
	TransactionSource entity.TransactionSource `json:"transaction_source" binding:"required"`
}

type TransactionResponse struct {
	ID                uint64                   `json:"id"`
	ContractNumber    string                   `json:"contract_number"`
	CustomerID        uint64                   `json:"customer_id"`
	TenorMonths       int                      `json:"tenor_months"`
	OTRAmount         float64                  `json:"otr_amount"`
	AdminFee          float64                  `json:"admin_fee"`
	InstallmentAmount float64                  `json:"installment_amount"`
	InterestAmount    float64                  `json:"interest_amount"`
	AssetName         string                   `json:"asset_name"`
	AssetType         entity.AssetType         `json:"asset_type"`
	Status            entity.TransactionStatus `json:"status"`
	TransactionSource entity.TransactionSource `json:"transaction_source"`
	Customer          CustomerResponse         `json:"customer"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
}

type UpdateTransactionStatusRequest struct {
	Status entity.TransactionStatus `json:"status" binding:"required"`
}
