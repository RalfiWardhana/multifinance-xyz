package usecase

import (
	"context"
	"fmt"
	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/domain/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthUseCase interface {
	Login(ctx context.Context, nik string) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error) // Changed return type
	GetCustomerFromToken(ctx context.Context, tokenString string) (*entity.Customer, error)
}

type authUseCase struct {
	customerRepo repository.CustomerRepository
	jwtSecret    string
}

func NewAuthUseCase(customerRepo repository.CustomerRepository) AuthUseCase {
	return &authUseCase{
		customerRepo: customerRepo,
		jwtSecret:    "xyz-secret-key", // Should come from config
	}
}

func (uc *authUseCase) Login(ctx context.Context, nik string) (string, error) {
	customer, err := uc.customerRepo.GetByNIK(ctx, nik)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"customer_id": customer.ID,
		"nik":         customer.NIK,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
		"iat":         time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

// Fixed ValidateToken function
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
		return claims, nil // Return claims directly, not pointer
	}

	return nil, fmt.Errorf("invalid token")
}

// Fixed GetCustomerFromToken function
func (uc *authUseCase) GetCustomerFromToken(ctx context.Context, tokenString string) (*entity.Customer, error) {
	claims, err := uc.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Type assertion with safety check
	customerIDFloat, ok := claims["customer_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid customer_id in token")
	}

	customerID := uint64(customerIDFloat)
	return uc.customerRepo.GetByID(ctx, customerID)
}
