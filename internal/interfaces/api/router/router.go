package router

import (
	"pt-xyz-multifinance/internal/interfaces/api/handler"
	"pt-xyz-multifinance/internal/interfaces/api/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	customerHandler *handler.CustomerHandler,
	transactionHandler *handler.TransactionHandler,
	authHandler *handler.AuthHandler,
) {
	// Global middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.RateLimitMiddleware(100, time.Minute)) // 100 requests per minute
	r.Use(gin.Recovery())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "PT XYZ Multifinance API is running"})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// Customer routes (authentication required)
		customers := v1.Group("/customers")
		// customers.Use(middleware.AuthMiddleware(authUseCase)) // Uncomment when auth is implemented
		{
			customers.POST("", customerHandler.CreateCustomer)
			customers.GET("", customerHandler.GetAllCustomers)
			customers.GET("/:id", customerHandler.GetCustomerByID)
			customers.GET("/:id/limits", customerHandler.GetCustomerLimits)
		}

		// Transaction routes (authentication required)
		transactions := v1.Group("/transactions")
		// transactions.Use(middleware.AuthMiddleware(authUseCase)) // Uncomment when auth is implemented
		{
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("", transactionHandler.GetAllTransactions)
			transactions.GET("/:id", transactionHandler.GetTransactionByID)
			transactions.PUT("/:id/status", transactionHandler.UpdateTransactionStatus)
			transactions.GET("/customer/:customer_id", transactionHandler.GetTransactionsByCustomerID)
		}
	}
}
