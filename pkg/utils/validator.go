package utils

import (
	"fmt"
	"pt-xyz-multifinance/pkg/constants"
	"regexp"
	"sort"
	"unicode"
)

func IsValidNIK(nik string) bool {
	if len(nik) != 16 {
		return false
	}

	nikRegex := regexp.MustCompile(`^\d{16}$`)
	return nikRegex.MatchString(nik)
}

func IsValidPhoneNumber(phone string) bool {
	phoneRegex := regexp.MustCompile(`^(\+62|62|0)8[1-9][0-9]{6,9}$`)
	return phoneRegex.MatchString(phone)
}

func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateTenor validates single tenor months (only 1, 2, 3, 4 allowed)
func ValidateTenor(tenor int) error {
	if !constants.IsValidTenor(tenor) {
		return fmt.Errorf("invalid tenor %d months, only 1, 2, 3, 4 months allowed", tenor)
	}
	return nil
}

// ValidateCompleteTenors validates that customer has ALL required tenors (1,2,3,4)
func ValidateCompleteTenors(tenors []int) error {
	// Check if we have exactly 4 tenors
	if len(tenors) != 4 {
		missing := getMissingTenors(tenors)
		if len(missing) > 0 {
			return fmt.Errorf("incomplete tenors provided. Missing tenors: %v. All tenors (1,2,3,4) are required", missing)
		}
		return fmt.Errorf("exactly 4 tenors (1,2,3,4) are required, got %d", len(tenors))
	}

	// Validate each tenor individually
	for _, tenor := range tenors {
		if err := ValidateTenor(tenor); err != nil {
			return err
		}
	}

	// Check for duplicates
	tenorMap := make(map[int]bool)
	for _, tenor := range tenors {
		if tenorMap[tenor] {
			return fmt.Errorf("duplicate tenor %d found. Each tenor (1,2,3,4) must appear exactly once", tenor)
		}
		tenorMap[tenor] = true
	}

	// Check if all required tenors are present
	requiredTenors := constants.GetValidTenors()
	for _, required := range requiredTenors {
		if !tenorMap[required] {
			missing := getMissingTenors(tenors)
			return fmt.Errorf("missing required tenor %d. All tenors (1,2,3,4) are mandatory. Missing: %v", required, missing)
		}
	}

	return nil
}

// getMissingTenors returns which tenors are missing from the provided list
func getMissingTenors(providedTenors []int) []int {
	requiredTenors := constants.GetValidTenors()
	providedMap := make(map[int]bool)

	for _, tenor := range providedTenors {
		providedMap[tenor] = true
	}

	var missing []int
	for _, required := range requiredTenors {
		if !providedMap[required] {
			missing = append(missing, required)
		}
	}

	sort.Ints(missing)
	return missing
}

// ValidateTenorLimits validates tenor and limit pairs for customer creation
func ValidateTenorLimits(tenorLimits []TenorLimit) error {
	if len(tenorLimits) == 0 {
		return fmt.Errorf("at least one tenor limit is required")
	}

	// Extract tenors for validation
	tenors := make([]int, len(tenorLimits))
	for i, tl := range tenorLimits {
		tenors[i] = tl.Tenor

		// Validate limit amount
		if tl.LimitAmount <= 0 {
			return fmt.Errorf("limit amount must be greater than 0 for tenor %d months", tl.Tenor)
		}
	}

	// Validate complete tenors
	return ValidateCompleteTenors(tenors)
}

// TenorLimit represents a tenor with its limit amount
type TenorLimit struct {
	Tenor       int     `json:"tenor_months"`
	LimitAmount float64 `json:"limit_amount"`
}

// ValidateAmount validates monetary amount
func ValidateAmount(amount float64) error {
	if amount < 0 {
		return fmt.Errorf("amount cannot be negative")
	}
	if amount > 999999999999.99 {
		return fmt.Errorf("amount too large")
	}
	return nil
}

// ValidateAssetType validates asset type
func ValidateAssetType(assetType string) error {
	validTypes := []string{"WHITE_GOODS", "MOTOR", "MOBIL"}
	for _, validType := range validTypes {
		if assetType == validType {
			return nil
		}
	}
	return fmt.Errorf("invalid asset type: %s, allowed: WHITE_GOODS, MOTOR, MOBIL", assetType)
}

// ValidateTransactionSource validates transaction source
func ValidateTransactionSource(source string) error {
	validSources := []string{"ECOMMERCE", "WEB", "DEALER"}
	for _, validSource := range validSources {
		if source == validSource {
			return nil
		}
	}
	return fmt.Errorf("invalid transaction source: %s, allowed: ECOMMERCE, WEB, DEALER", source)
}
