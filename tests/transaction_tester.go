package testing

import (
	"context"
	"fmt"
	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/usecase"
	"sync"
	"time"
)

type TransactionTester struct {
	transactionUseCase usecase.TransactionUseCase
}

func NewTransactionTester(transactionUseCase usecase.TransactionUseCase) *TransactionTester {
	return &TransactionTester{
		transactionUseCase: transactionUseCase,
	}
}

// TestConcurrentTransactionCreation tests multiple concurrent transactions
func (tt *TransactionTester) TestConcurrentTransactionCreation(customerID uint64, numTransactions int, amount float64) (*ConcurrentTestResult, error) {
	var wg sync.WaitGroup
	results := make([]TransactionResult, numTransactions)
	start := time.Now()

	for i := 0; i < numTransactions; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			transaction := &entity.Transaction{
				CustomerID:        customerID,
				TenorMonths:       1,
				OTRAmount:         amount,
				AdminFee:          amount * 0.1,
				InterestAmount:    amount * 0.05,
				AssetName:         fmt.Sprintf("Concurrent Test Asset %d", index),
				AssetType:         entity.AssetWhiteGoods,
				TransactionSource: entity.SourceWeb,
			}

			err := tt.transactionUseCase.CreateTransaction(context.Background(), transaction)
			results[index] = TransactionResult{
				Index:         index,
				TransactionID: transaction.ID,
				Success:       err == nil,
				Error:         err,
				Timestamp:     time.Now(),
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	return &ConcurrentTestResult{
		Results:      results,
		Duration:     duration,
		TotalCount:   numTransactions,
		SuccessCount: tt.countSuccessful(results),
		FailureCount: tt.countFailed(results),
	}, nil
}

func (tt *TransactionTester) countSuccessful(results []TransactionResult) int {
	count := 0
	for _, result := range results {
		if result.Success {
			count++
		}
	}
	return count
}

func (tt *TransactionTester) countFailed(results []TransactionResult) int {
	count := 0
	for _, result := range results {
		if !result.Success {
			count++
		}
	}
	return count
}

type TransactionResult struct {
	Index         int
	TransactionID uint64
	Success       bool
	Error         error
	Timestamp     time.Time
}

type ConcurrentTestResult struct {
	Results      []TransactionResult
	Duration     time.Duration
	TotalCount   int
	SuccessCount int
	FailureCount int
}

func (ctr *ConcurrentTestResult) PrintSummary() {
	fmt.Printf("=== Concurrent Transaction Test Results ===\n")
	fmt.Printf("Total Transactions: %d\n", ctr.TotalCount)
	fmt.Printf("Successful: %d\n", ctr.SuccessCount)
	fmt.Printf("Failed: %d\n", ctr.FailureCount)
	fmt.Printf("Success Rate: %.2f%%\n", float64(ctr.SuccessCount)/float64(ctr.TotalCount)*100)
	fmt.Printf("Duration: %v\n", ctr.Duration)
	fmt.Printf("Transactions/sec: %.2f\n", float64(ctr.TotalCount)/ctr.Duration.Seconds())

	if ctr.FailureCount > 0 {
		fmt.Printf("\nFailure Details:\n")
		for _, result := range ctr.Results {
			if !result.Success {
				fmt.Printf("Transaction %d failed: %v\n", result.Index, result.Error)
			}
		}
	}
}
