package entity

import (
	"time"
)

type CustomerLimit struct {
	ID              uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	CustomerID      uint64    `json:"customer_id" gorm:"not null"`
	TenorMonths     int       `json:"tenor_months" gorm:"not null"`
	LimitAmount     float64   `json:"limit_amount" gorm:"not null"`
	UsedAmount      float64   `json:"used_amount" gorm:"default:0"`
	AvailableAmount float64   `json:"available_amount" gorm:"type:decimal(15,2) GENERATED ALWAYS AS (limit_amount - used_amount) STORED"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Customer Customer `json:"customer" gorm:"foreignKey:CustomerID"`
}

func (CustomerLimit) TableName() string {
	return "customer_limits"
}
