package handler

import (
	"fmt"
	"net/http"
	"pt-xyz-multifinance/internal/interfaces/dto"
	"pt-xyz-multifinance/internal/usecase"
	"pt-xyz-multifinance/pkg/response"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
}

func NewAuthHandler(authUseCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMsg := h.getValidationErrorMessage(validationErrors)
			response.Error(c, http.StatusBadRequest, "Validation Failed", errorMsg)
			return
		}
		response.Error(c, http.StatusBadRequest, "Invalid JSON format", err.Error())
		return
	}

	fmt.Printf("DEBUG Register Request: %+v\n", req)
	if req.CustomerData != nil {
		fmt.Printf("DEBUG Customer Data: %+v\n", *req.CustomerData)
		fmt.Printf("DEBUG KTP Photo Path: %s\n", req.CustomerData.KTPPhotoPath)
		fmt.Printf("DEBUG Selfie Photo Path: %s\n", req.CustomerData.SelfiePhotoPath)
	}

	if validationErr := h.validateRegisterRequest(&req); validationErr != "" {
		response.Error(c, http.StatusBadRequest, "Business Rule Validation Failed", validationErr)
		return
	}

	user, customer, err := h.authUseCase.Register(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Registration failed", err.Error())
		return
	}

	registerResponse := dto.RegisterResponse{
		Message: "Registration successful",
		User: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	if customer != nil {
		fmt.Printf("DEBUG Created Customer: %+v\n", *customer)
		fmt.Printf("DEBUG Customer KTP Photo: %s\n", customer.KTPPhotoPath)
		fmt.Printf("DEBUG Customer Selfie Photo: %s\n", customer.SelfiePhotoPath)

		registerResponse.Customer = &dto.CustomerResponse{
			ID:              customer.ID,
			UserID:          customer.UserID,
			NIK:             customer.NIK,
			FullName:        customer.FullName,
			LegalName:       customer.LegalName,
			BirthPlace:      customer.BirthPlace,
			BirthDate:       customer.BirthDate,
			Salary:          customer.Salary,
			KTPPhotoPath:    customer.KTPPhotoPath,
			SelfiePhotoPath: customer.SelfiePhotoPath,
			CreatedAt:       customer.CreatedAt,
			UpdatedAt:       customer.UpdatedAt,
		}
	}

	response.Success(c, http.StatusCreated, "Registration successful", registerResponse)
}

func (h *AuthHandler) validateRegisterRequest(req *dto.RegisterRequest) string {

	if req.Role != "ADMIN" && req.Role != "CUSTOMER" {
		return "Role must be either ADMIN or CUSTOMER"
	}

	if req.Role == "CUSTOMER" {
		if req.CustomerData == nil {
			return "Customer data is required for CUSTOMER role"
		}

		if len(req.CustomerData.Limits) != 4 {
			return "Customer must have exactly 4 tenor limits (1, 2, 3, 4 months)"
		}

		tenorMap := make(map[int]bool)
		for _, limit := range req.CustomerData.Limits {
			if limit.TenorMonths < 1 || limit.TenorMonths > 4 {
				return "Invalid tenor months. Only 1, 2, 3, 4 months allowed"
			}
			if tenorMap[limit.TenorMonths] {
				return "Duplicate tenor detected. Each tenor (1,2,3,4) must appear exactly once"
			}
			tenorMap[limit.TenorMonths] = true
			if limit.LimitAmount <= 0 {
				return "Limit amount must be greater than 0"
			}
		}

		requiredTenors := []int{1, 2, 3, 4}
		for _, required := range requiredTenors {
			if !tenorMap[required] {
				return "Missing required tenor: " + string(rune(required)) + " months"
			}
		}
	}

	return ""
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	token, user, customer, err := h.authUseCase.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	loginResponse := dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		},
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	if customer != nil {
		loginResponse.Customer = &dto.CustomerResponse{
			ID:              customer.ID,
			UserID:          customer.UserID,
			NIK:             customer.NIK,
			FullName:        customer.FullName,
			LegalName:       customer.LegalName,
			BirthPlace:      customer.BirthPlace,
			BirthDate:       customer.BirthDate,
			Salary:          customer.Salary,
			KTPPhotoPath:    customer.KTPPhotoPath,
			SelfiePhotoPath: customer.SelfiePhotoPath,
			CreatedAt:       customer.CreatedAt,
			UpdatedAt:       customer.UpdatedAt,
		}
	}

	response.Success(c, http.StatusOK, "Login successful", loginResponse)
}

func (h *AuthHandler) getValidationErrorMessage(validationErrors validator.ValidationErrors) string {
	var messages []string

	for _, err := range validationErrors {
		switch err.Tag() {
		case "required":
			messages = append(messages, err.Field()+" is required")
		case "email":
			messages = append(messages, "Email format is invalid")
		case "min":
			messages = append(messages, err.Field()+" must be at least "+err.Param()+" characters")
		case "max":
			messages = append(messages, err.Field()+" cannot exceed "+err.Param()+" characters")
		case "oneof":
			if err.Field() == "Role" {
				messages = append(messages, "Role must be either ADMIN or CUSTOMER")
			} else {
				messages = append(messages, err.Field()+" must be one of the allowed values")
			}
		default:
			messages = append(messages, err.Field()+" is invalid")
		}
	}

	return strings.Join(messages, "; ")
}
