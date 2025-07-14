package usecase

import (
	"context"
	"fmt"
	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/domain/repository"
	"pt-xyz-multifinance/internal/interfaces/dto"
	"pt-xyz-multifinance/pkg/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUseCase interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*entity.User, *entity.Customer, error)
	Login(ctx context.Context, username, password string) (string, *entity.User, *entity.Customer, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
	GetUserFromToken(ctx context.Context, tokenString string) (*entity.User, error)
	GetCustomerFromUser(ctx context.Context, userID uint64) (*entity.Customer, error)
}

type authUseCase struct {
	userRepo     repository.UserRepository
	customerRepo repository.CustomerRepository
	limitRepo    repository.LimitRepository
	db           *gorm.DB
	jwtSecret    string
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	customerRepo repository.CustomerRepository,
	limitRepo repository.LimitRepository,
	db *gorm.DB,
) AuthUseCase {
	return &authUseCase{
		userRepo:     userRepo,
		customerRepo: customerRepo,
		limitRepo:    limitRepo,
		db:           db,
		jwtSecret:    "xyz-secret-key-2024",
	}
}

func (uc *authUseCase) Register(ctx context.Context, req *dto.RegisterRequest) (*entity.User, *entity.Customer, error) {
	// Validate password confirmation
	if req.Password != req.ConfirmPassword {
		return nil, nil, fmt.Errorf("password and confirm password do not match")
	}

	// Check if username already exists
	if _, err := uc.userRepo.GetByUsername(ctx, req.Username); err == nil {
		return nil, nil, fmt.Errorf("username already exists")
	}

	// Check if email already exists
	if _, err := uc.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, nil, fmt.Errorf("email already exists")
	}

	// If role is CUSTOMER, validate customer data and NIK uniqueness
	if req.Role == "CUSTOMER" {
		if req.CustomerData == nil {
			return nil, nil, fmt.Errorf("customer data is required for customer role")
		}

		// Check if NIK already exists
		if _, err := uc.customerRepo.GetByNIK(ctx, req.CustomerData.NIK); err == nil {
			return nil, nil, fmt.Errorf("NIK already exists")
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	var user *entity.User
	var customer *entity.Customer

	// Start transaction
	err = uc.db.Transaction(func(tx *gorm.DB) error {
		// Create user
		user = &entity.User{
			Username: req.Username,
			Email:    req.Email,
			Password: string(hashedPassword),
			Role:     entity.UserRole(req.Role),
			IsActive: true,
		}

		if err := uc.userRepo.Create(ctx, user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// If role is CUSTOMER, create customer and limits
		if req.Role == "CUSTOMER" {
			birthDate, err := time.Parse("2006-01-02", req.CustomerData.BirthDate)
			if err != nil {
				return fmt.Errorf("invalid birth date format: %w", err)
			}

			customer = &entity.Customer{
				UserID:          user.ID,
				NIK:             req.CustomerData.NIK,
				FullName:        req.CustomerData.FullName,
				LegalName:       req.CustomerData.LegalName,
				BirthPlace:      req.CustomerData.BirthPlace,
				BirthDate:       birthDate,
				Salary:          req.CustomerData.Salary,
				KTPPhotoPath:    req.CustomerData.KTPPhotoPath,
				SelfiePhotoPath: req.CustomerData.SelfiePhotoPath,
			}

			if err := uc.customerRepo.Create(ctx, customer); err != nil {
				return fmt.Errorf("failed to create customer: %w", err)
			}

			// Create customer limits
			for _, limitReq := range req.CustomerData.Limits {
				limit := &entity.CustomerLimit{
					CustomerID:  customer.ID,
					TenorMonths: limitReq.TenorMonths,
					LimitAmount: limitReq.LimitAmount,
					UsedAmount:  0,
				}

				if err := uc.limitRepo.Create(ctx, limit); err != nil {
					return fmt.Errorf("failed to create customer limit: %w", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	logger.Info("User registered successfully", "userID", user.ID, "role", user.Role)
	return user, customer, nil
}

func (uc *authUseCase) Login(ctx context.Context, username, password string) (string, *entity.User, *entity.Customer, error) {
	// Get user by username
	user, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return "", nil, nil, fmt.Errorf("account is deactivated")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, nil, fmt.Errorf("invalid credentials")
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to generate token: %w", err)
	}

	var customer *entity.Customer
	if user.Role == entity.RoleCustomer {
		customer, err = uc.GetCustomerFromUser(ctx, user.ID)
		if err != nil {
			logger.Error("Failed to get customer data for user", "userID", user.ID, "error", err)
		}
	}

	return tokenString, user, customer, nil
}

func (uc *authUseCase) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (uc *authUseCase) GetUserFromToken(ctx context.Context, tokenString string) (*entity.User, error) {
	claims, err := uc.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	userID := uint64(userIDFloat)
	return uc.userRepo.GetByID(ctx, userID)
}

func (uc *authUseCase) GetCustomerFromUser(ctx context.Context, userID uint64) (*entity.Customer, error) {
	// First, get all customers and find by UserID
	// We need to modify customer repository to support GetByUserID
	customers, err := uc.customerRepo.GetAll(ctx, 1000, 0) // Get many customers to find the one
	if err != nil {
		return nil, err
	}

	for _, customer := range customers {
		if customer.UserID == userID {
			return customer, nil
		}
	}

	return nil, fmt.Errorf("customer not found for user ID: %d", userID)
}
