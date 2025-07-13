package dto

type LoginRequest struct {
	NIK string `json:"nik" binding:"required,len=16"`
}

type LoginResponse struct {
	Token     string           `json:"token"`
	Customer  CustomerResponse `json:"customer"`
	ExpiresAt int64            `json:"expires_at"`
}
