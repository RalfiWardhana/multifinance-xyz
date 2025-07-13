package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"pt-xyz-multifinance/internal/config"
	"pt-xyz-multifinance/internal/infrastructure/database"
	"pt-xyz-multifinance/internal/infrastructure/repository"
	"pt-xyz-multifinance/internal/interfaces/api/handler"
	"pt-xyz-multifinance/internal/interfaces/api/router"
	"pt-xyz-multifinance/internal/interfaces/dto"
	"pt-xyz-multifinance/internal/usecase"
	"pt-xyz-multifinance/pkg/logger"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupConcurrentTestDB() (*gin.Engine, error) {
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "",
		Database: "pt_xyz_multifinance_test",
	}

	logger.InitLogger("info")
	db, err := database.NewMySQLConnection(*cfg)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	customerRepo := repository.NewCustomerRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	limitRepo := repository.NewLimitRepository(db)

	// Initialize use cases with DB for transactions
	customerUseCase := usecase.NewCustomerUseCase(customerRepo, limitRepo, db)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepo, customerRepo, limitRepo, db)
	authUseCase := usecase.NewAuthUseCase(customerRepo)

	// Initialize handlers
	customerHandler := handler.NewCustomerHandler(customerUseCase)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase)
	authHandler := handler.NewAuthHandler(authUseCase)

	// Setup router
	r := gin.New()
	router.SetupRoutes(r, customerHandler, transactionHandler, authHandler, authUseCase)

	return r, nil
}

// Test concurrent transactions to validate ACID properties
func TestConcurrentTransactionCreation(t *testing.T) {
	router, err := setupConcurrentTestDB()
	if err != nil {
		t.Skip("Database not available for integration tests")
	}

	// Number of concurrent transactions to create
	numConcurrent := 10
	customerID := uint64(1) // Assuming customer exists with limited budget
	tenorMonths := 1
	otrAmount := 50000.0 // Each transaction requests 50k

	var wg sync.WaitGroup
	results := make([]int, numConcurrent)

	// Create concurrent transactions
	for i := 0; i < numConcurrent; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			transactionData := dto.CreateTransactionRequest{
				CustomerID:        customerID,
				TenorMonths:       tenorMonths,
				OTRAmount:         otrAmount,
				AdminFee:          5000,
				InterestAmount:    2500,
				AssetName:         fmt.Sprintf("Test Asset %d", index),
				AssetType:         "WHITE_GOODS",
				TransactionSource: "WEB",
			}

			jsonData, _ := json.Marshal(transactionData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token") // You'll need valid token

			router.ServeHTTP(w, req)
			results[index] = w.Code
		}(i)
	}

	wg.Wait()

	// Analyze results
	successCount := 0
	failureCount := 0

	for _, code := range results {
		if code == 201 {
			successCount++
		} else {
			failureCount++
		}
	}

	// With proper transaction handling, we should have:
	// - Some successful transactions (within limit)
	// - Some failed transactions (limit exceeded)
	// - No corrupted data state

	fmt.Printf("Concurrent test results: %d success, %d failures\n", successCount, failureCount)

	// At least one should succeed and some should fail due to limit constraints
	assert.True(t, successCount > 0, "At least one transaction should succeed")
	assert.True(t, failureCount > 0, "Some transactions should fail due to limit constraints")
	assert.Equal(t, numConcurrent, successCount+failureCount, "All transactions should be processed")
}

// Test transaction rollback scenario
func TestTransactionRollback(t *testing.T) {
	router, err := setupConcurrentTestDB()
	if err != nil {
		t.Skip("Database not available for integration tests")
	}

	// 1. Create a transaction
	transactionData := dto.CreateTransactionRequest{
		CustomerID:        1,
		TenorMonths:       1,
		OTRAmount:         30000,
		AdminFee:          3000,
		InterestAmount:    1500,
		AssetName:         "Test Rollback Asset",
		AssetType:         "WHITE_GOODS",
		TransactionSource: "WEB",
	}

	jsonData, _ := json.Marshal(transactionData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	// Parse response to get transaction ID
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	transactionID := uint64(data["id"].(float64))

	// 2. Reject the transaction (should rollback the used limit)
	rejectData := map[string]string{
		"reason": "Test rollback scenario",
	}

	jsonData, _ = json.Marshal(rejectData)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/api/v1/transactions/%d/reject", transactionID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// 3. Verify that the limit was rolled back by creating another transaction with same amount
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	router.ServeHTTP(w, req)
	// Should succeed again since limit was rolled back
	assert.Equal(t, 201, w.Code)
}

// Benchmark concurrent transaction performance
func BenchmarkConcurrentTransactions(b *testing.B) {
	router, err := setupConcurrentTestDB()
	if err != nil {
		b.Skip("Database not available for benchmark tests")
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			counter++
			transactionData := dto.CreateTransactionRequest{
				CustomerID:        1,
				TenorMonths:       1,
				OTRAmount:         1000, // Small amount to avoid limit issues
				AdminFee:          100,
				InterestAmount:    50,
				AssetName:         fmt.Sprintf("Benchmark Asset %d", counter),
				AssetType:         "WHITE_GOODS",
				TransactionSource: "WEB",
			}

			jsonData, _ := json.Marshal(transactionData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")

			router.ServeHTTP(w, req)
		}
	})
}
