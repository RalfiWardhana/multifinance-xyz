// File: main.go (Fixed version)
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
	customerRepo := repository.NewCustomerRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	limitRepo := repository.NewLimitRepository(db)

	// Initialize use cases
	customerUseCase := usecase.NewCustomerUseCase(customerRepo, limitRepo)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepo, customerRepo, limitRepo)
	authUseCase := usecase.NewAuthUseCase(customerRepo)

	// Initialize handlers
	customerHandler := handler.NewCustomerHandler(customerUseCase)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase)
	authHandler := handler.NewAuthHandler(authUseCase)

	// Initialize Gin router
	r := gin.New()
	router.SetupRoutes(r, customerHandler, transactionHandler, authHandler, authUseCase) // Added authUseCase

	// Start server
	logger.Info("Starting server on port " + cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
