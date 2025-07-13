// File: internal/domain/entity/limit.go (Simple fix - no generated column)
package entity

import (
	"time"
)

type CustomerLimit struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	CustomerID  uint64    `json:"customer_id" gorm:"not null;index"`
	TenorMonths int       `json:"tenor_months" gorm:"not null;index"`
	LimitAmount float64   `json:"limit_amount" gorm:"type:decimal(15,2);not null"`
	UsedAmount  float64   `json:"used_amount" gorm:"type:decimal(15,2);default:0"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Customer Customer `json:"customer" gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
}

// AvailableAmount calculates available amount dynamically
func (cl *CustomerLimit) AvailableAmount() float64 {
	return cl.LimitAmount - cl.UsedAmount
}

func (CustomerLimit) TableName() string {
	return "customer_limits"
}
