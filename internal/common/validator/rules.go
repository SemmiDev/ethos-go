package validator

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// validatePassword validates password with standard rules:
// - Minimum 6 characters
// - At least 1 lowercase letter
// - At least 1 uppercase letter
// - At least 1 digit
// - At least 1 special character
func (v *Validator) validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Minimum length
	if len(password) < 6 {
		return false
	}

	var (
		hasLower   = false
		hasUpper   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasNumber && hasSpecial
}

// validateStrongPassword validates password with stricter rules:
// - Minimum 8 characters
// - At least 2 lowercase letters
// - At least 2 uppercase letters
// - At least 2 digits
// - At least 2 special characters
// - No common patterns or sequences
func (v *Validator) validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Minimum length
	if len(password) < 8 {
		return false
	}

	var (
		lowerCount   = 0
		upperCount   = 0
		numberCount  = 0
		specialCount = 0
	)

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			lowerCount++
		case unicode.IsUpper(char):
			upperCount++
		case unicode.IsNumber(char):
			numberCount++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			specialCount++
		}
	}

	// Check minimum counts
	if lowerCount < 2 || upperCount < 2 || numberCount < 2 || specialCount < 2 {
		return false
	}

	// Check for common weak patterns
	lowerPass := strings.ToLower(password)
	weakPatterns := []string{
		"123456", "654321", "abcdef", "fedcba",
		"qwerty", "asdfgh", "zxcvbn", "password",
		"admin", "user", "guest", "test",
	}

	for _, pattern := range weakPatterns {
		if strings.Contains(lowerPass, pattern) {
			return false
		}
	}

	return true
}

// validateIndonesianPhone validates an Indonesian phone number (e.g., +6281234567890 or 081234567890)
func (v *Validator) validateIndonesianPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// Indonesian phone numbers typically start with +62 or 08, followed by 8-12 digits
	re := regexp.MustCompile(`^(?:\+628|08)\d{8,12}$`)
	return re.MatchString(phone)
}

// validateIndonesianPostalCode validates an Indonesian postal code (5 digits)
func (v *Validator) validateIndonesianPostalCode(fl validator.FieldLevel) bool {
	postalCode := fl.Field().String()
	re := regexp.MustCompile(`^\d{5}$`)
	return re.MatchString(postalCode)
}

// validateIndonesianNIK validates an Indonesian NIK (16 digits)
func (v *Validator) validateIndonesianNIK(fl validator.FieldLevel) bool {
	nik := fl.Field().String()
	re := regexp.MustCompile(`^\d{16}$`)
	return re.MatchString(nik)
}

// validateUsername validates a username (3-20 characters, alphanumeric, underscore, dot allowed)
func (v *Validator) validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	re := regexp.MustCompile(`^[a-zA-Z0-9_.]{3,20}$`)
	return re.MatchString(username)
}

// validateNoHTML ensures the input does not contain HTML tags
func (v *Validator) validateNoHTML(fl validator.FieldLevel) bool {
	input := fl.Field().String()
	re := regexp.MustCompile(`<[^>]+>`)
	return !re.MatchString(input)
}

// validateIndonesianCurrency validates Indonesian currency format (e.g., Rp1.000.000 or Rp 1.000.000)
func (v *Validator) validateIndonesianCurrency(fl validator.FieldLevel) bool {
	currency := fl.Field().String()
	// Allow "Rp" or "Rp " followed by digits with optional thousand separators (.)
	re := regexp.MustCompile(`^Rp\s?\d{1,3}(\.\d{3})*(\,\d{2})?$`)
	return re.MatchString(currency)
}
