package dto

import "pt-xyz-multifinance/internal/domain/entity"

type RegisterRequest struct {
	Username        string                 `json:"username" binding:"required,min=3,max=50"`
	Email           string                 `json:"email" binding:"required,email"`
	Password        string                 `json:"password" binding:"required,min=8"`
	ConfirmPassword string                 `json:"confirm_password" binding:"required"`
	Role            string                 `json:"role" binding:"required,oneof=ADMIN CUSTOMER"`
	CustomerData    *CreateCustomerRequest `json:"customer_data,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string            `json:"token"`
	User      UserResponse      `json:"user"`
	Customer  *CustomerResponse `json:"customer,omitempty"`
	ExpiresAt int64             `json:"expires_at"`
}

type UserResponse struct {
	ID        uint64          `json:"id"`
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	Role      entity.UserRole `json:"role"`
	IsActive  bool            `json:"is_active"`
	CreatedAt string          `json:"created_at"`
}

type RegisterResponse struct {
	Message  string            `json:"message"`
	User     UserResponse      `json:"user"`
	Customer *CustomerResponse `json:"customer,omitempty"`
}
