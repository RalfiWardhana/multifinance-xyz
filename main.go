package main

import (
	"log"
	"pt-xyz-multifinance/internal/config"
	"pt-xyz-multifinance/internal/infrastructure/database"
	"pt-xyz-multifinance/internal/infrastructure/repository"
	"pt-xyz-multifinance/internal/interfaces/api/handler"
	"pt-xyz-multifinance/internal/interfaces/api/router"
	"pt-xyz-multifinance/internal/usecase"
	"pt-xyz-multifinance/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize logger
	logger.InitLogger(cfg.LogLevel)

	// Initialize database
	db, err := database.NewMySQLConnection(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.CloseConnection(db)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	limitRepo := repository.NewLimitRepository(db)

	// Initialize use cases (pass DB instance for transaction handling)
	customerUseCase := usecase.NewCustomerUseCase(customerRepo, limitRepo, db)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepo, customerRepo, limitRepo, db)
	authUseCase := usecase.NewAuthUseCase(userRepo, customerRepo, limitRepo, db)

	// Initialize handlers
	customerHandler := handler.NewCustomerHandler(customerUseCase)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase, customerUseCase)
	authHandler := handler.NewAuthHandler(authUseCase)

	// Initialize Gin router
	r := gin.New()
	router.SetupRoutes(r, customerHandler, transactionHandler, authHandler, authUseCase)

	// Start server
	logger.Info("Starting PT XYZ Multifinance API server on port " + cfg.Server.Port)
	logger.Info("Available endpoints:")
	logger.Info("  POST /api/v1/auth/register - Register new user")
	logger.Info("  POST /api/v1/auth/login - Login")
	logger.Info("  GET  /api/v1/admin/customers - Get all customers (Admin only)")
	logger.Info("  GET  /api/v1/customers/me - Get my profile (Customer)")
	logger.Info("  POST /api/v1/transactions - Create transaction")
	logger.Info("  GET  /health - Health check")

	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
