package utils

import (
	"regexp"
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
