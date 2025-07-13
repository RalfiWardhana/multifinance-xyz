package handler

import (
	"net/http"
	"pt-xyz-multifinance/internal/interfaces/dto"
	"pt-xyz-multifinance/internal/usecase"
	"pt-xyz-multifinance/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
}

func NewAuthHandler(authUseCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	token, err := h.authUseCase.Login(c.Request.Context(), req.NIK)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	customer, err := h.authUseCase.GetCustomerFromToken(c.Request.Context(), token)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get customer info", err.Error())
		return
	}

	loginResponse := dto.LoginResponse{
		Token: token,
		Customer: dto.CustomerResponse{
			ID:        customer.ID,
			NIK:       customer.NIK,
			FullName:  customer.FullName,
			LegalName: customer.LegalName,
		},
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	response.Success(c, http.StatusOK, "Login successful", loginResponse)
}
