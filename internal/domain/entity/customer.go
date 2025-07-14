package entity

import (
	"time"
)

type UserRole string

const (
	RoleAdmin    UserRole = "ADMIN"
	RoleCustomer UserRole = "CUSTOMER"
)

type User struct {
	ID        uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string     `json:"username" gorm:"type:varchar(100);unique;not null;index"`
	Email     string     `json:"email" gorm:"type:varchar(255);unique;not null;index"`
	Password  string     `json:"-" gorm:"type:varchar(255);not null"`
	Role      UserRole   `json:"role" gorm:"type:enum('ADMIN','CUSTOMER');not null;default:CUSTOMER;index"`
	IsActive  bool       `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsCustomer() bool {
	return u.Role == RoleCustomer
}

type Customer struct {
	ID              uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          uint64     `json:"user_id" gorm:"not null;unique;index"`
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
	User            User       `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (Customer) TableName() string {
	return "customers"
}
