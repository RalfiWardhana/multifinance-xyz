package handler

import (
	"fmt"
	"net/http"
	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/interfaces/dto"
	"pt-xyz-multifinance/internal/usecase"
	"pt-xyz-multifinance/pkg/response"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CustomerHandler struct {
	customerUseCase usecase.CustomerUseCase
}

func NewCustomerHandler(customerUseCase usecase.CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{
		customerUseCase: customerUseCase,
	}
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req dto.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMsg := h.getValidationErrorMessage(validationErrors)
			response.Error(c, http.StatusBadRequest, "Validation Failed", errorMsg)
			return
		}
		response.Error(c, http.StatusBadRequest, "Invalid JSON format", err.Error())
		return
	}

	if validationErr := h.validateCustomerRequest(&req); validationErr != "" {
		response.Error(c, http.StatusBadRequest, "Business Rule Validation Failed", validationErr)
		return
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid birth date format", "Birth date must be in YYYY-MM-DD format (example: 1990-05-15)")
		return
	}

	customer := &entity.Customer{
		NIK:             req.NIK,
		FullName:        req.FullName,
		LegalName:       req.LegalName,
		BirthPlace:      req.BirthPlace,
		BirthDate:       birthDate,
		Salary:          req.Salary,
		KTPPhotoPath:    req.KTPPhotoPath,
		SelfiePhotoPath: req.SelfiePhotoPath,
	}

	var limits []*entity.CustomerLimit
	for _, limitReq := range req.Limits {
		limits = append(limits, &entity.CustomerLimit{
			TenorMonths: limitReq.TenorMonths,
			LimitAmount: limitReq.LimitAmount,
		})
	}

	if err := h.customerUseCase.CreateCustomer(c.Request.Context(), customer, limits); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create customer", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Customer created successfully with all required tenors", h.toCustomerResponse(customer))
}

func (h *CustomerHandler) validateCustomerRequest(req *dto.CreateCustomerRequest) string {
	// 1. Check if exactly 4 limits are provided
	if len(req.Limits) < 4 {
		missingCount := 4 - len(req.Limits)
		return fmt.Sprintf("Incomplete tenor limits. You provided %d tenor(s), but PT XYZ requires exactly 4 tenors (1, 2, 3, 4 months). Missing %d tenor(s).",
			len(req.Limits), missingCount)
	}

	if len(req.Limits) > 4 {
		extraCount := len(req.Limits) - 4
		return fmt.Sprintf("Too many tenor limits. You provided %d tenor(s), but PT XYZ only accepts exactly 4 tenors (1, 2, 3, 4 months). Please remove %d tenor(s).",
			len(req.Limits), extraCount)
	}

	// 2. Check for valid tenor values and duplicates
	tenorMap := make(map[int]bool)
	providedTenors := make([]int, 0)

	for i, limit := range req.Limits {
		// Check valid tenor range
		if limit.TenorMonths < 1 || limit.TenorMonths > 4 {
			return fmt.Sprintf("Invalid tenor at position %d: %d months is not allowed. PT XYZ only accepts tenors: 1, 2, 3, 4 months.",
				i+1, limit.TenorMonths)
		}

		// Check for duplicates
		if tenorMap[limit.TenorMonths] {
			return fmt.Sprintf("Duplicate tenor detected: %d months appears more than once. Each tenor (1, 2, 3, 4) must appear exactly once.",
				limit.TenorMonths)
		}

		tenorMap[limit.TenorMonths] = true
		providedTenors = append(providedTenors, limit.TenorMonths)

		// Check limit amount
		if limit.LimitAmount <= 0 {
			return fmt.Sprintf("Invalid limit amount for tenor %d months: %.2f. Limit amount must be greater than 0.",
				limit.TenorMonths, limit.LimitAmount)
		}
	}

	// 3. Check if all required tenors (1,2,3,4) are present
	requiredTenors := []int{1, 2, 3, 4}
	missingTenors := make([]int, 0)

	for _, required := range requiredTenors {
		if !tenorMap[required] {
			missingTenors = append(missingTenors, required)
		}
	}

	if len(missingTenors) > 0 {
		return fmt.Sprintf("Missing required tenors: %v months. PT XYZ requires ALL tenors (1, 2, 3, 4 months) to be provided. You provided: %v months.",
			missingTenors, providedTenors)
	}

	return ""
}

func (h *CustomerHandler) getValidationErrorMessage(validationErrors validator.ValidationErrors) string {
	var messages []string

	for _, err := range validationErrors {
		switch err.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("%s is required", getFieldDisplayName(err.Field())))
		case "len":
			if err.Field() == "NIK" {
				messages = append(messages, "NIK must be exactly 16 digits")
			} else if err.Field() == "Limits" {
				messages = append(messages, "PT XYZ requires exactly 4 tenor limits (1, 2, 3, 4 months)")
			} else {
				messages = append(messages, fmt.Sprintf("%s must be exactly %s characters", getFieldDisplayName(err.Field()), err.Param()))
			}
		case "min":
			if err.Field() == "Limits" {
				messages = append(messages, "At least 1 tenor limit is required")
			} else if err.Field() == "LimitAmount" {
				messages = append(messages, "Limit amount must be greater than 0")
			} else {
				messages = append(messages, fmt.Sprintf("%s must be at least %s", getFieldDisplayName(err.Field()), err.Param()))
			}
		case "max":
			if err.Field() == "Limits" {
				messages = append(messages, "Maximum 4 tenor limits allowed (PT XYZ policy)")
			} else if err.Field() == "TenorMonths" {
				messages = append(messages, "Tenor months cannot exceed 4 (PT XYZ only offers 1, 2, 3, 4 months)")
			} else {
				messages = append(messages, fmt.Sprintf("%s cannot exceed %s", getFieldDisplayName(err.Field()), err.Param()))
			}
		case "oneof":
			if err.Field() == "TenorMonths" {
				messages = append(messages, "Tenor months must be one of: 1, 2, 3, 4 months only")
			} else {
				messages = append(messages, fmt.Sprintf("%s must be one of the allowed values", getFieldDisplayName(err.Field())))
			}
		default:
			messages = append(messages, fmt.Sprintf("%s is invalid: %s", getFieldDisplayName(err.Field()), err.Tag()))
		}
	}

	return strings.Join(messages, "; ")
}

func getFieldDisplayName(field string) string {
	switch field {
	case "NIK":
		return "NIK (Identity Number)"
	case "FullName":
		return "Full Name"
	case "LegalName":
		return "Legal Name"
	case "BirthPlace":
		return "Birth Place"
	case "BirthDate":
		return "Birth Date"
	case "Salary":
		return "Salary"
	case "Limits":
		return "Tenor Limits"
	case "TenorMonths":
		return "Tenor Months"
	case "LimitAmount":
		return "Limit Amount"
	default:
		return field
	}
}

func (h *CustomerHandler) GetCustomerByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid customer ID", "Customer ID must be a valid number")
		return
	}

	customer, err := h.customerUseCase.GetCustomerByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Customer not found", fmt.Sprintf("No customer found with ID: %d", id))
		return
	}

	response.Success(c, http.StatusOK, "Customer retrieved successfully", h.toCustomerResponse(customer))
}

func (h *CustomerHandler) GetCustomerLimits(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid customer ID", "Customer ID must be a valid number")
		return
	}

	limits, err := h.customerUseCase.GetCustomerLimits(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Customer limits not found", fmt.Sprintf("No limits found for customer ID: %d", id))
		return
	}

	var limitResponses []dto.CustomerLimitResponse
	for _, limit := range limits {
		limitResponses = append(limitResponses, dto.CustomerLimitResponse{
			ID:              limit.ID,
			CustomerID:      limit.CustomerID,
			TenorMonths:     limit.TenorMonths,
			LimitAmount:     limit.LimitAmount,
			UsedAmount:      limit.UsedAmount,
			AvailableAmount: limit.AvailableAmount(),
		})
	}

	response.Success(c, http.StatusOK, fmt.Sprintf("Found %d tenor limits for customer", len(limitResponses)), limitResponses)
}

func (h *CustomerHandler) GetAllCustomers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit < 1 || limit > 100 {
		response.Error(c, http.StatusBadRequest, "Invalid limit parameter", "Limit must be between 1 and 100")
		return
	}

	if offset < 0 {
		response.Error(c, http.StatusBadRequest, "Invalid offset parameter", "Offset must be 0 or greater")
		return
	}

	customers, err := h.customerUseCase.GetAllCustomers(c.Request.Context(), limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve customers", err.Error())
		return
	}

	var customerResponses []dto.CustomerResponse
	for _, customer := range customers {
		customerResponses = append(customerResponses, *h.toCustomerResponse(customer))
	}

	response.Success(c, http.StatusOK, fmt.Sprintf("Retrieved %d customers (limit: %d, offset: %d)", len(customerResponses), limit, offset), customerResponses)
}

func (h *CustomerHandler) toCustomerResponse(customer *entity.Customer) *dto.CustomerResponse {
	return &dto.CustomerResponse{
		ID:              customer.ID,
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
