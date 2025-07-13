package entity

import (
	"time"
)

type Customer struct {
	ID              uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	NIK             string     `json:"nik" gorm:"unique;not null"`
	FullName        string     `json:"full_name" gorm:"not null"`
	LegalName       string     `json:"legal_name" gorm:"not null"`
	BirthPlace      string     `json:"birth_place" gorm:"not null"`
	BirthDate       time.Time  `json:"birth_date" gorm:"not null"`
	Salary          float64    `json:"salary" gorm:"not null"`
	KTPPhotoPath    string     `json:"ktp_photo_path"`
	SelfiePhotoPath string     `json:"selfie_photo_path"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       *time.Time `json:"deleted_at" gorm:"index"`
}

func (Customer) TableName() string {
	return "customers"
}
