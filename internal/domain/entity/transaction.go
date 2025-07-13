package entity

import (
	"time"
)

type TransactionStatus string
type AssetType string
type TransactionSource string

const (
	StatusPending   TransactionStatus = "PENDING"
	StatusApproved  TransactionStatus = "APPROVED"
	StatusRejected  TransactionStatus = "REJECTED"
	StatusActive    TransactionStatus = "ACTIVE"
	StatusCompleted TransactionStatus = "COMPLETED"
	StatusDefaulted TransactionStatus = "DEFAULTED"
)

const (
	AssetWhiteGoods AssetType = "WHITE_GOODS"
	AssetMotor      AssetType = "MOTOR"
	AssetMobil      AssetType = "MOBIL"
)

const (
	SourceEcommerce TransactionSource = "ECOMMERCE"
	SourceWeb       TransactionSource = "WEB"
	SourceDealer    TransactionSource = "DEALER"
)

type Transaction struct {
	ID                uint64            `json:"id" gorm:"primaryKey;autoIncrement"`
	ContractNumber    string            `json:"contract_number" gorm:"unique;not null"`
	CustomerID        uint64            `json:"customer_id" gorm:"not null"`
	TenorMonths       int               `json:"tenor_months" gorm:"not null"`
	OTRAmount         float64           `json:"otr_amount" gorm:"not null"`
	AdminFee          float64           `json:"admin_fee" gorm:"not null"`
	InstallmentAmount float64           `json:"installment_amount" gorm:"not null"`
	InterestAmount    float64           `json:"interest_amount" gorm:"not null"`
	AssetName         string            `json:"asset_name" gorm:"not null"`
	AssetType         AssetType         `json:"asset_type" gorm:"not null"`
	Status            TransactionStatus `json:"status" gorm:"default:PENDING"`
	TransactionSource TransactionSource `json:"transaction_source" gorm:"not null"`
	CreatedAt         time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         *time.Time        `json:"deleted_at" gorm:"index"`

	// Relations
	Customer Customer `json:"customer" gorm:"foreignKey:CustomerID"`
}

func (Transaction) TableName() string {
	return "transactions"
}
