package standard

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/gentra/decorator-arch-go/internal/validation"
	"github.com/gentra/decorator-arch-go/internal/validationrule"
)

// service implements validation.Service interface using go-playground/validator
type service struct {
	validator   *validator.Validate
	customRules map[string]validationrule.Service
}

// NewService creates a new standard validation service
func NewService() validation.Service {
	v := validator.New()

	// Register custom validation functions
	v.RegisterValidation("strong_password", validateStrongPassword)
	v.RegisterValidation("clean_name", validateCleanName)
	v.RegisterValidation("theme", validateTheme)
	v.RegisterValidation("language", validateLanguage)

	return &service{
		validator:   v,
		customRules: make(map[string]validationrule.Service),
	}
}

// ValidateStruct validates a struct using struct tags
func (s *service) ValidateStruct(ctx context.Context, data interface{}) error {
	if err := s.validator.Struct(data); err != nil {
		// Convert validator errors to our validation errors
		var validationErrors validation.ValidationErrors
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors.Add(validation.ValidationError{
				Field:   err.Field(),
				Message: getErrorMessage(err),
				Value:   fmt.Sprintf("%v", err.Value()),
				Rule:    err.Tag(),
			})
		}
		return validationErrors
	}
	return nil
}

// ValidateField validates a single field
func (s *service) ValidateField(ctx context.Context, field string, value interface{}, rules string) error {
	if err := s.validator.Var(value, rules); err != nil {
		return validation.ValidationError{
			Field:   field,
			Message: getErrorMessage(err.(validator.ValidationErrors)[0]),
			Value:   fmt.Sprintf("%v", value),
		}
	}
	return nil
}

// ValidateUserRegistration validates user registration data
func (s *service) ValidateUserRegistration(ctx context.Context, data interface{}) error {
	// Use reflection or type assertion to extract fields
	// For simplicity, assuming we receive a map or struct with known fields
	// In a real implementation, you'd use reflection or type-specific validation

	return s.ValidateStruct(ctx, data)
}

// ValidateUserUpdate validates user update data
func (s *service) ValidateUserUpdate(ctx context.Context, data interface{}) error {
	return s.ValidateStruct(ctx, data)
}

// ValidateUserPreferences validates user preferences
func (s *service) ValidateUserPreferences(ctx context.Context, data interface{}) error {
	return s.ValidateStruct(ctx, data)
}

// ValidateUserID validates a user ID format
func (s *service) ValidateUserID(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return validation.ValidationError{
			Field:   "user_id",
			Message: "must be a valid UUID",
			Value:   id,
			Rule:    "uuid",
		}
	}
	return nil
}

// ValidateEmail validates email format and business rules
func (s *service) ValidateEmail(ctx context.Context, email string) error {
	// Basic format validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return validation.ValidationError{
			Field:   "email",
			Message: "must be a valid email address",
			Value:   email,
			Rule:    "email",
		}
	}

	// Length check
	if len(email) > 254 {
		return validation.ValidationError{
			Field:   "email",
			Message: "must be no more than 254 characters",
			Value:   email,
			Rule:    "max",
		}
	}

	// Check for common suspicious patterns
	if strings.Contains(email, "..") {
		return validation.ValidationError{
			Field:   "email",
			Message: "cannot contain consecutive dots",
			Value:   email,
			Rule:    "format",
		}
	}

	return nil
}

// ValidatePassword validates password strength
func (s *service) ValidatePassword(ctx context.Context, password string) error {
	var errorMessages []string

	// Length check
	if len(password) < 8 {
		errorMessages = append(errorMessages, "must be at least 8 characters long")
	}

	if len(password) > 128 {
		errorMessages = append(errorMessages, "must be no more than 128 characters long")
	}

	// Character type checks
	var hasLower, hasUpper, hasDigit, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasLower {
		errorMessages = append(errorMessages, "must contain at least one lowercase letter")
	}
	if !hasUpper {
		errorMessages = append(errorMessages, "must contain at least one uppercase letter")
	}
	if !hasDigit {
		errorMessages = append(errorMessages, "must contain at least one digit")
	}
	if !hasSpecial {
		errorMessages = append(errorMessages, "must contain at least one special character")
	}

	// Check for common weak passwords
	weakPasswords := []string{
		"password", "123456", "qwerty", "abc123",
		"password123", "admin", "letmein", "welcome",
	}

	for _, weak := range weakPasswords {
		if strings.ToLower(password) == weak {
			errorMessages = append(errorMessages, "password is too common")
			break
		}
	}

	if len(errorMessages) > 0 {
		return validation.ValidationError{
			Field:   "password",
			Message: strings.Join(errorMessages, "; "),
			Rule:    "strong_password",
		}
	}

	return nil
}

// AddCustomRule adds a custom validation rule
func (s *service) AddCustomRule(name string, rule validationrule.Service) error {
	s.customRules[name] = rule
	return nil
}

// RemoveCustomRule removes a custom validation rule
func (s *service) RemoveCustomRule(name string) error {
	delete(s.customRules, name)
	return nil
}

// Custom validation functions for the validator package

func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var hasLower, hasUpper, hasDigit, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasDigit && hasSpecial
}

func validateCleanName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s'-]+$`)
	return nameRegex.MatchString(name)
}

func validateTheme(fl validator.FieldLevel) bool {
	theme := fl.Field().String()
	validThemes := []string{"light", "dark", "auto"}
	return contains(validThemes, theme)
}

func validateLanguage(fl validator.FieldLevel) bool {
	language := fl.Field().String()
	return len(language) == 2 && isAlpha(language)
}

// Helper functions

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters long", err.Param())
	case "max":
		return fmt.Sprintf("must be no more than %s characters long", err.Param())
	case "strong_password":
		return "password does not meet security requirements"
	case "clean_name":
		return "can only contain letters, spaces, hyphens, and apostrophes"
	case "theme":
		return "must be one of: light, dark, auto"
	case "language":
		return "must be a 2-letter language code"
	default:
		return fmt.Sprintf("validation failed for rule: %s", err.Tag())
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func isAlpha(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	return true
}
