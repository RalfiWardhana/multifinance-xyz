package utils_test

import (
	"pt-xyz-multifinance/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateTenor(t *testing.T) {
	// Valid tenors
	validTenors := []int{1, 2, 3, 4}
	for _, tenor := range validTenors {
		err := utils.ValidateTenor(tenor)
		assert.NoError(t, err, "Tenor %d should be valid", tenor)
	}

	// Invalid tenors
	invalidTenors := []int{0, 5, 6, 12, -1}
	for _, tenor := range invalidTenors {
		err := utils.ValidateTenor(tenor)
		assert.Error(t, err, "Tenor %d should be invalid", tenor)
	}
}

func TestValidateCompleteTenors(t *testing.T) {
	// Valid complete tenors
	validTenors := []int{1, 2, 3, 4}
	err := utils.ValidateCompleteTenors(validTenors)
	assert.NoError(t, err)

	// Missing tenors
	incompleteTenors := []int{1, 2, 3}
	err = utils.ValidateCompleteTenors(incompleteTenors)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "incomplete tenors")

	// Duplicate tenors
	duplicateTenors := []int{1, 2, 2, 4}
	err = utils.ValidateCompleteTenors(duplicateTenors)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate tenor")

	// Too many tenors
	tooManyTenors := []int{1, 2, 3, 4, 5}
	err = utils.ValidateCompleteTenors(tooManyTenors)
	assert.Error(t, err)
}

func TestValidateTenorLimits(t *testing.T) {
	// Valid tenor limits
	validLimits := []utils.TenorLimit{
		{Tenor: 1, LimitAmount: 100000},
		{Tenor: 2, LimitAmount: 200000},
		{Tenor: 3, LimitAmount: 300000},
		{Tenor: 4, LimitAmount: 400000},
	}
	err := utils.ValidateTenorLimits(validLimits)
	assert.NoError(t, err)

	// Invalid limit amount
	invalidLimits := []utils.TenorLimit{
		{Tenor: 1, LimitAmount: 0}, // Zero amount
		{Tenor: 2, LimitAmount: 200000},
		{Tenor: 3, LimitAmount: 300000},
		{Tenor: 4, LimitAmount: 400000},
	}
	err = utils.ValidateTenorLimits(invalidLimits)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limit amount must be greater than 0")
}

func TestIsValidNIK(t *testing.T) {
	// Valid NIK
	validNIK := "1234567890123456"
	assert.True(t, utils.IsValidNIK(validNIK))

	// Invalid NIK - wrong length
	invalidNIK1 := "123456789012345" // 15 digits
	assert.False(t, utils.IsValidNIK(invalidNIK1))

	// Invalid NIK - contains letters
	invalidNIK2 := "123456789012345A"
	assert.False(t, utils.IsValidNIK(invalidNIK2))

	// Invalid NIK - too long
	invalidNIK3 := "12345678901234567" // 17 digits
	assert.False(t, utils.IsValidNIK(invalidNIK3))
}

func TestIsValidEmail(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name@domain.co.id",
		"admin+test@company.org",
	}
	for _, email := range validEmails {
		assert.True(t, utils.IsValidEmail(email), "Email %s should be valid", email)
	}

	invalidEmails := []string{
		"invalid-email",
		"@domain.com",
		"user@",
		"user@domain",
		"user name@domain.com",
	}
	for _, email := range invalidEmails {
		assert.False(t, utils.IsValidEmail(email), "Email %s should be invalid", email)
	}
}

func TestValidateAmount(t *testing.T) {

	validAmounts := []float64{0, 100, 1000.50, 999999999999.99}
	for _, amount := range validAmounts {
		err := utils.ValidateAmount(amount)
		assert.NoError(t, err, "Amount %.2f should be valid", amount)
	}

	invalidAmounts := []float64{-1, -100.50, 1000000000000}
	for _, amount := range invalidAmounts {
		err := utils.ValidateAmount(amount)
		assert.Error(t, err, "Amount %.2f should be invalid", amount)
	}
}
