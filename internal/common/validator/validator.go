package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
	locale   string
}

func New(locale string) *Validator {
	validate := validator.New()

	// Register a custom tag name function to use "json" tags for field names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	v := &Validator{
		validate: validate,
		locale:   locale,
	}

	v.registerCustomValidators()

	return v
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) ToKV() map[string]ValidationError {
	m := make(map[string]ValidationError, len(ve))
	for _, err := range ve {
		m[err.Field] = err
	}
	return m
}

func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, "; ")
}

func (v *Validator) Validate(i any) error {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	var validationErrors ValidationErrors
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			validationError := ValidationError{
				Field:   fieldErr.Field(),
				Tag:     fieldErr.Tag(),
				Value:   fmt.Sprintf("%v", fieldErr.Value()),
				Message: v.getErrorMessage(fieldErr),
			}
			validationErrors = append(validationErrors, validationError)
		}
	}

	return validationErrors
}

func (v *Validator) ValidateAndGetErrors(i any) ValidationErrors {
	err := v.Validate(i)
	if err == nil {
		return nil
	}

	if validationErrs, ok := err.(ValidationErrors); ok {
		return validationErrs
	}

	return nil
}

// getErrorMessage generates localized error messages based on validation tag and locale
func (v *Validator) getErrorMessage(fe validator.FieldError) string {
	field := fe.Field()
	param := fe.Param()
	tag := fe.Tag()

	switch strings.ToLower(v.locale) {
	case "id":
		return v.indonesianErrorMessage(field, tag, param)
	case "en":
		return v.englishErrorMessage(field, tag, param)
	default:
		return v.indonesianErrorMessage(field, tag, param)
	}
}

// registerCustomValidators registers all custom validation tags
func (v *Validator) registerCustomValidators() {
	// Password validation with comprehensive rules
	v.validate.RegisterValidation("password", v.validatePassword)
	// Strong password with stricter rules
	v.validate.RegisterValidation("strong_password", v.validateStrongPassword)
	// Indonesian phone number
	v.validate.RegisterValidation("phone_id", v.validateIndonesianPhone)
	// Indonesian postal code
	v.validate.RegisterValidation("postal_code_id", v.validateIndonesianPostalCode)
	// Indonesian NIK
	v.validate.RegisterValidation("nik", v.validateIndonesianNIK)
	// Username
	v.validate.RegisterValidation("username", v.validateUsername)
	// No HTML tags
	v.validate.RegisterValidation("no_html", v.validateNoHTML)
	// Indonesian currency format
	v.validate.RegisterValidation("currency_id", v.validateIndonesianCurrency)
}

// IsValidationErrors checks if the error is of type ValidationErrors
func IsValidationErrors(err error) bool {
	_, ok := err.(ValidationErrors)
	return ok
}

// GetValidationErrors safely converts error to ValidationErrors
func GetValidationErrors(err error) (ValidationErrors, bool) {
	if validationErrs, ok := err.(ValidationErrors); ok {
		return validationErrs, true
	}
	return nil, false
}
