package handler

import (
	"net/http"
	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/interfaces/dto"
	"pt-xyz-multifinance/internal/usecase"
	"pt-xyz-multifinance/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionUseCase usecase.TransactionUseCase
}

func NewTransactionHandler(transactionUseCase usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		transactionUseCase: transactionUseCase,
	}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Create transaction entity
	transaction := &entity.Transaction{
		CustomerID:        req.CustomerID,
		TenorMonths:       req.TenorMonths,
		OTRAmount:         req.OTRAmount,
		AdminFee:          req.AdminFee,
		InterestAmount:    req.InterestAmount,
		AssetName:         req.AssetName,
		AssetType:         req.AssetType,
		TransactionSource: req.TransactionSource,
		Status:            entity.StatusPending,
	}

	if err := h.transactionUseCase.CreateTransaction(c.Request.Context(), transaction); err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to create transaction", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Transaction created successfully", h.toTransactionResponse(transaction))
}

func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid transaction ID", err.Error())
		return
	}

	transaction, err := h.transactionUseCase.GetTransactionByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Transaction not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Transaction retrieved successfully", h.toTransactionResponse(transaction))
}

func (h *TransactionHandler) GetTransactionsByCustomerID(c *gin.Context) {
	idParam := c.Param("customer_id")
	customerID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid customer ID", err.Error())
		return
	}

	transactions, err := h.transactionUseCase.GetTransactionsByCustomerID(c.Request.Context(), customerID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Transactions not found", err.Error())
		return
	}

	var transactionResponses []dto.TransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, *h.toTransactionResponse(transaction))
	}

	response.Success(c, http.StatusOK, "Transactions retrieved successfully", transactionResponses)
}

func (h *TransactionHandler) UpdateTransactionStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid transaction ID", err.Error())
		return
	}

	var req dto.UpdateTransactionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if err := h.transactionUseCase.UpdateTransactionStatus(c.Request.Context(), id, req.Status); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update transaction status", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Transaction status updated successfully", nil)
}

func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	transactions, err := h.transactionUseCase.GetAllTransactions(c.Request.Context(), limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve transactions", err.Error())
		return
	}

	var transactionResponses []dto.TransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, *h.toTransactionResponse(transaction))
	}

	response.Success(c, http.StatusOK, "Transactions retrieved successfully", transactionResponses)
}

func (h *TransactionHandler) toTransactionResponse(transaction *entity.Transaction) *dto.TransactionResponse {
	return &dto.TransactionResponse{
		ID:                transaction.ID,
		ContractNumber:    transaction.ContractNumber,
		CustomerID:        transaction.CustomerID,
		TenorMonths:       transaction.TenorMonths,
		OTRAmount:         transaction.OTRAmount,
		AdminFee:          transaction.AdminFee,
		InstallmentAmount: transaction.InstallmentAmount,
		InterestAmount:    transaction.InterestAmount,
		AssetName:         transaction.AssetName,
		AssetType:         transaction.AssetType,
		Status:            transaction.Status,
		TransactionSource: transaction.TransactionSource,
		Customer: dto.CustomerResponse{
			ID:        transaction.Customer.ID,
			NIK:       transaction.Customer.NIK,
			FullName:  transaction.Customer.FullName,
			LegalName: transaction.Customer.LegalName,
		},
		CreatedAt: transaction.CreatedAt,
		UpdatedAt: transaction.UpdatedAt,
	}
}
