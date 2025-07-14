package router

import (
	"pt-xyz-multifinance/internal/interfaces/api/handler"
	"pt-xyz-multifinance/internal/interfaces/api/middleware"
	"pt-xyz-multifinance/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	customerHandler *handler.CustomerHandler,
	transactionHandler *handler.TransactionHandler,
	authHandler *handler.AuthHandler,
	authUseCase usecase.AuthUseCase,
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
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Admin-only routes
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(authUseCase))
		admin.Use(middleware.AdminOnly())
		{
			// Admin can access all customers
			admin.GET("/customers", customerHandler.GetAllCustomers)
			admin.GET("/customers/:id", customerHandler.GetCustomerByID)
			admin.GET("/customers/:id/limits", customerHandler.GetCustomerLimits)

			// Admin can access all transactions
			admin.GET("/transactions", transactionHandler.GetAllTransactions)
			admin.GET("/transactions/:id", transactionHandler.GetTransactionByID)
			admin.PUT("/transactions/:id/status", transactionHandler.UpdateTransactionStatus)
			admin.GET("/transactions/customer/:customer_id", transactionHandler.GetTransactionsByCustomerID)

			// Admin can create customers directly (without user registration)
			admin.POST("/customers", customerHandler.CreateCustomer)
		}

		// Customer routes (authentication required + ownership check)
		customers := v1.Group("/customers")
		customers.Use(middleware.AuthMiddleware(authUseCase))
		customers.Use(middleware.CustomerOwnershipMiddleware())
		{
			// Customers can only access their own data
			customers.GET("/:id", customerHandler.GetCustomerByID)
			customers.GET("/:id/limits", customerHandler.GetCustomerLimits)

			// Customer can view their own profile (using their customer ID from token)
			customers.GET("/me", customerHandler.GetMyProfile)
		}

		// Transaction routes (authentication required + ownership check)
		transactions := v1.Group("/transactions")
		transactions.Use(middleware.AuthMiddleware(authUseCase))
		{
			// Both admin and customer can create transactions
			// But customers can only create for themselves (will be validated in handler)
			transactions.POST("", transactionHandler.CreateTransaction)

			// Protected routes with ownership middleware
			protected := transactions.Group("")
			protected.Use(middleware.CustomerOwnershipMiddleware())
			{
				protected.GET("/:id", transactionHandler.GetTransactionByID)
				protected.GET("/customer/:customer_id", transactionHandler.GetTransactionsByCustomerID)
			}

			// Admin-only transaction management
			adminTrans := transactions.Group("/admin")
			adminTrans.Use(middleware.AdminOnly())
			{
				adminTrans.GET("", transactionHandler.GetAllTransactions)
				adminTrans.PUT("/:id/status", transactionHandler.UpdateTransactionStatus)
			}
		}

		// User management routes (admin only)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware(authUseCase))
		users.Use(middleware.AdminOnly())
		{
			// TODO: Add user management handlers
			// users.GET("", userHandler.GetAllUsers)
			// users.GET("/:id", userHandler.GetUserByID)
			// users.PUT("/:id/status", userHandler.UpdateUserStatus)
		}

		// Deprecated: Keep old register endpoint for backward compatibility
		v1.POST("/register", authHandler.Register)
	}
}
