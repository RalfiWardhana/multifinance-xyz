package handler

import (
	"net/http"
	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/interfaces/dto"
	"pt-xyz-multifinance/internal/usecase"
	"pt-xyz-multifinance/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Parse birth date
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid birth date format", "Use YYYY-MM-DD format")
		return
	}

	// Create customer entity
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

	// Create limits
	var limits []*entity.CustomerLimit
	for _, limitReq := range req.Limits {
		limits = append(limits, &entity.CustomerLimit{
			TenorMonths: limitReq.TenorMonths,
			LimitAmount: limitReq.LimitAmount,
		})
	}

	// Create customer with limits
	if err := h.customerUseCase.CreateCustomer(c.Request.Context(), customer, limits); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create customer", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Customer created successfully", h.toCustomerResponse(customer))
}

func (h *CustomerHandler) GetCustomerByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid customer ID", err.Error())
		return
	}

	customer, err := h.customerUseCase.GetCustomerByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Customer not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Customer retrieved successfully", h.toCustomerResponse(customer))
}

func (h *CustomerHandler) GetCustomerLimits(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid customer ID", err.Error())
		return
	}

	limits, err := h.customerUseCase.GetCustomerLimits(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Customer limits not found", err.Error())
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
			AvailableAmount: limit.AvailableAmount,
		})
	}

	response.Success(c, http.StatusOK, "Customer limits retrieved successfully", limitResponses)
}

func (h *CustomerHandler) GetAllCustomers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	customers, err := h.customerUseCase.GetAllCustomers(c.Request.Context(), limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve customers", err.Error())
		return
	}

	var customerResponses []dto.CustomerResponse
	for _, customer := range customers {
		customerResponses = append(customerResponses, *h.toCustomerResponse(customer))
	}

	response.Success(c, http.StatusOK, "Customers retrieved successfully", customerResponses)
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
