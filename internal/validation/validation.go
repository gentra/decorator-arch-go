package validation

import (
	"context"
	"fmt"
	"strings"

	"github.com/gentra/decorator-arch-go/internal/validationrule"
)

// Service defines the validation domain interface - the ONLY interface in this domain
type Service interface {
	// General validation operations
	ValidateStruct(ctx context.Context, data interface{}) error
	ValidateField(ctx context.Context, field string, value interface{}, rules string) error

	// User domain specific validations
	ValidateUserRegistration(ctx context.Context, data interface{}) error
	ValidateUserUpdate(ctx context.Context, data interface{}) error
	ValidateUserPreferences(ctx context.Context, data interface{}) error
	ValidateUserID(ctx context.Context, id string) error
	ValidateEmail(ctx context.Context, email string) error
	ValidatePassword(ctx context.Context, password string) error

	// Configuration
	AddCustomRule(name string, rule validationrule.Service) error
	RemoveCustomRule(name string) error
}

// Domain types and data structures

// ValidationError represents a validation error with field-specific details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
	Rule    string `json:"rule,omitempty"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (e ValidationErrors) Error() string {
	if len(e.Errors) == 0 {
		return "validation errors occurred"
	}

	messages := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		messages[i] = err.Error()
	}

	return strings.Join(messages, "; ")
}

// ValidationResult contains the result of a validation operation
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

// ValidationConfig contains configuration for the validation service
type ValidationConfig struct {
	StrictMode      bool                              `json:"strict_mode"`      // Fail on first error vs collect all errors
	CustomRules     map[string]validationrule.Service `json:"custom_rules"`     // Custom validation rules
	EnableI18n      bool                              `json:"enable_i18n"`      // Enable internationalization
	DefaultLanguage string                            `json:"default_language"` // Default language for error messages
}

// Helper methods for ValidationError
func (e *ValidationError) IsEmpty() bool {
	return e.Field == "" && e.Message == ""
}

func (e *ValidationError) WithField(field string) *ValidationError {
	e.Field = field
	return e
}

func (e *ValidationError) WithValue(value string) *ValidationError {
	e.Value = value
	return e
}

func (e *ValidationError) WithRule(rule string) *ValidationError {
	e.Rule = rule
	return e
}

// Helper methods for ValidationErrors
func (e *ValidationErrors) Add(err ValidationError) {
	e.Errors = append(e.Errors, err)
}

func (e *ValidationErrors) AddField(field, message string) {
	e.Add(ValidationError{Field: field, Message: message})
}

func (e *ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

func (e *ValidationErrors) HasFieldError(field string) bool {
	for _, err := range e.Errors {
		if err.Field == field {
			return true
		}
	}
	return false
}

func (e *ValidationErrors) GetFieldErrors(field string) []ValidationError {
	var fieldErrors []ValidationError
	for _, err := range e.Errors {
		if err.Field == field {
			fieldErrors = append(fieldErrors, err)
		}
	}
	return fieldErrors
}

// Helper methods for ValidationResult
func (r *ValidationResult) IsValid() bool {
	return r.Valid && len(r.Errors) == 0
}

func (r *ValidationResult) AddError(err ValidationError) {
	r.Valid = false
	r.Errors = append(r.Errors, err)
}

// Helper methods for ValidationConfig
func (c *ValidationConfig) IsValid() bool {
	return c.DefaultLanguage != ""
}

// DefaultValidationConfig returns default validation configuration
func DefaultValidationConfig() ValidationConfig {
	return ValidationConfig{
		StrictMode:      false,
		CustomRules:     make(map[string]validationrule.Service),
		EnableI18n:      false,
		DefaultLanguage: "en",
	}
}

// Common validation error messages
var (
	ErrRequired     = "field is required"
	ErrInvalidEmail = "invalid email format"
	ErrTooShort     = "value is too short"
	ErrTooLong      = "value is too long"
	ErrInvalidUUID  = "invalid UUID format"
	ErrWeakPassword = "password does not meet security requirements"
)
