package entity

import (
	"time"
)

type Customer struct {
	ID              uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	NIK             string     `json:"nik" gorm:"type:varchar(16);unique;not null;index"`
	FullName        string     `json:"full_name" gorm:"type:varchar(255);not null;index:idx_full_name,length:100"`
	LegalName       string     `json:"legal_name" gorm:"type:varchar(255);not null"`
	BirthPlace      string     `json:"birth_place" gorm:"type:varchar(255);not null"`
	BirthDate       time.Time  `json:"birth_date" gorm:"type:date;not null"`
	Salary          float64    `json:"salary" gorm:"type:decimal(15,2);not null"`
	KTPPhotoPath    string     `json:"ktp_photo_path" gorm:"type:varchar(500)"`
	SelfiePhotoPath string     `json:"selfie_photo_path" gorm:"type:varchar(500)"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       *time.Time `json:"deleted_at" gorm:"index"`
}

func (Customer) TableName() string {
	return "customers"
}
