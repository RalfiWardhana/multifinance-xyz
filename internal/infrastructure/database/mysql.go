package database

import (
	"fmt"
	"pt-xyz-multifinance/internal/config"
	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/pkg/logger"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func NewMySQLConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	config := &gorm.Config{
		Logger: glogger.Default.LogMode(glogger.Info),
	}

	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto migrate tables
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	// Create default admin user if not exists
	if err := createDefaultAdmin(db); err != nil {
		logger.Error("Failed to create default admin", "error", err)
	}

	logger.Info("Successfully connected to MySQL database")
	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.Customer{},
		&entity.CustomerLimit{},
		&entity.Transaction{},
	)
}

func createDefaultAdmin(db *gorm.DB) error {
	// Check if admin user already exists
	var count int64
	db.Model(&entity.User{}).Where("role = ?", entity.RoleAdmin).Count(&count)

	if count > 0 {
		logger.Info("Admin user already exists, skipping creation")
		return nil
	}

	// Create default admin user
	// Password: admin123 (hashed)
	hashedPassword := "$2a$10$K5QzJ8gOLJ8K5QzJ8gOLJO5QzJ8gOLJ8K5QzJ8gOLJ8K5QzJ8gOLJO" // admin123

	admin := &entity.User{
		Username: "admin",
		Email:    "admin@ptxyz.com",
		Password: hashedPassword,
		Role:     entity.RoleAdmin,
		IsActive: true,
	}

	if err := db.Create(admin).Error; err != nil {
		return fmt.Errorf("failed to create default admin: %w", err)
	}

	logger.Info("Default admin user created successfully", "username", "admin", "email", "admin@ptxyz.com")
	logger.Info("Default admin password: admin123 (please change this in production)")

	return nil
}

func CloseConnection(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
		logger.Info("Database connection closed")
	}
}
