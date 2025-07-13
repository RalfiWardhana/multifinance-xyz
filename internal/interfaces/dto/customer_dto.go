package dto

import (
	"time"
)

type CreateCustomerRequest struct {
	NIK             string               `json:"nik" binding:"required,len=16"`
	FullName        string               `json:"full_name" binding:"required,min=2,max=255"`
	LegalName       string               `json:"legal_name" binding:"required,min=2,max=255"`
	BirthPlace      string               `json:"birth_place" binding:"required,min=2,max=255"`
	BirthDate       string               `json:"birth_date" binding:"required"`
	Salary          float64              `json:"salary" binding:"required,min=0"`
	KTPPhotoPath    string               `json:"ktp_photo_path"`
	SelfiePhotoPath string               `json:"selfie_photo_path"`
	Limits          []CreateLimitRequest `json:"limits" binding:"required,min=1"`
}

type CreateLimitRequest struct {
	TenorMonths int     `json:"tenor_months" binding:"required,min=1"`
	LimitAmount float64 `json:"limit_amount" binding:"required,min=0"`
}

type CustomerResponse struct {
	ID              uint64    `json:"id"`
	NIK             string    `json:"nik"`
	FullName        string    `json:"full_name"`
	LegalName       string    `json:"legal_name"`
	BirthPlace      string    `json:"birth_place"`
	BirthDate       time.Time `json:"birth_date"`
	Salary          float64   `json:"salary"`
	KTPPhotoPath    string    `json:"ktp_photo_path"`
	SelfiePhotoPath string    `json:"selfie_photo_path"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CustomerLimitResponse struct {
	ID              uint64  `json:"id"`
	CustomerID      uint64  `json:"customer_id"`
	TenorMonths     int     `json:"tenor_months"`
	LimitAmount     float64 `json:"limit_amount"`
	UsedAmount      float64 `json:"used_amount"`
	AvailableAmount float64 `json:"available_amount"`
}
