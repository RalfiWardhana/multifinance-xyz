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
	ContractNumber    string            `json:"contract_number" gorm:"type:varchar(50);unique;not null;index"`
	CustomerID        uint64            `json:"customer_id" gorm:"not null;index"`
	TenorMonths       int               `json:"tenor_months" gorm:"not null"`
	OTRAmount         float64           `json:"otr_amount" gorm:"type:decimal(15,2);not null"`
	AdminFee          float64           `json:"admin_fee" gorm:"type:decimal(15,2);not null"`
	InstallmentAmount float64           `json:"installment_amount" gorm:"type:decimal(15,2);not null"`
	InterestAmount    float64           `json:"interest_amount" gorm:"type:decimal(15,2);not null"`
	AssetName         string            `json:"asset_name" gorm:"type:varchar(255);not null"`
	AssetType         AssetType         `json:"asset_type" gorm:"type:enum('WHITE_GOODS','MOTOR','MOBIL');not null;index"`
	Status            TransactionStatus `json:"status" gorm:"type:enum('PENDING','APPROVED','REJECTED','ACTIVE','COMPLETED','DEFAULTED');default:PENDING;index"`
	TransactionSource TransactionSource `json:"transaction_source" gorm:"type:enum('ECOMMERCE','WEB','DEALER');not null"`
	CreatedAt         time.Time         `json:"created_at" gorm:"autoCreateTime;index"`
	UpdatedAt         time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         *time.Time        `json:"deleted_at" gorm:"index"`
	Customer          Customer          `json:"customer" gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
}

func (Transaction) TableName() string {
	return "transactions"
}
