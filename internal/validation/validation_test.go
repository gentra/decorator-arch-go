package validation_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/validation"
	"github.com/gentra/decorator-arch-go/internal/validationrule"
)

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		valErr   validation.ValidationError
		expected string
	}{
		{
			name: "Given validation error with field and message, When Error is called, Then should return formatted error",
			valErr: validation.ValidationError{
				Field:   "email",
				Message: "invalid email format",
			},
			expected: "validation error for field 'email': invalid email format",
		},
		{
			name: "Given validation error with empty field, When Error is called, Then should return formatted error with empty field",
			valErr: validation.ValidationError{
				Field:   "",
				Message: "validation failed",
			},
			expected: "validation error for field '': validation failed",
		},
		{
			name: "Given validation error with empty message, When Error is called, Then should return formatted error with empty message",
			valErr: validation.ValidationError{
				Field:   "password",
				Message: "",
			},
			expected: "validation error for field 'password': ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.valErr.Error()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidationError_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		valErr   validation.ValidationError
		expected bool
	}{
		{
			name: "Given validation error with both field and message empty, When IsEmpty is called, Then should return true",
			valErr: validation.ValidationError{
				Field:   "",
				Message: "",
			},
			expected: true,
		},
		{
			name: "Given validation error with field set, When IsEmpty is called, Then should return false",
			valErr: validation.ValidationError{
				Field:   "email",
				Message: "",
			},
			expected: false,
		},
		{
			name: "Given validation error with message set, When IsEmpty is called, Then should return false",
			valErr: validation.ValidationError{
				Field:   "",
				Message: "validation failed",
			},
			expected: false,
		},
		{
			name: "Given validation error with both field and message set, When IsEmpty is called, Then should return false",
			valErr: validation.ValidationError{
				Field:   "email",
				Message: "invalid email",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.valErr.IsEmpty()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidationError_ChainableMethods(t *testing.T) {
	t.Run("Given validation error, When using chainable methods, Then should set values correctly", func(t *testing.T) {
		// Arrange
		valErr := validation.ValidationError{}

		// Act
		result := valErr.WithField("email").WithValue("invalid@").WithRule("email")

		// Assert
		assert.Equal(t, "email", result.Field)
		assert.Equal(t, "invalid@", result.Value)
		assert.Equal(t, "email", result.Rule)
	})
}

func TestValidationErrors_Error(t *testing.T) {
	tests := []struct {
		name     string
		valErrs  validation.ValidationErrors
		expected string
	}{
		{
			name: "Given validation errors with multiple errors, When Error is called, Then should return joined error messages",
			valErrs: validation.ValidationErrors{
				Errors: []validation.ValidationError{
					{Field: "email", Message: "invalid email"},
					{Field: "password", Message: "too short"},
				},
			},
			expected: "validation error for field 'email': invalid email; validation error for field 'password': too short",
		},
		{
			name: "Given validation errors with single error, When Error is called, Then should return single error message",
			valErrs: validation.ValidationErrors{
				Errors: []validation.ValidationError{
					{Field: "email", Message: "invalid email"},
				},
			},
			expected: "validation error for field 'email': invalid email",
		},
		{
			name: "Given validation errors with no errors, When Error is called, Then should return default message",
			valErrs: validation.ValidationErrors{
				Errors: []validation.ValidationError{},
			},
			expected: "validation errors occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.valErrs.Error()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidationErrors_Add(t *testing.T) {
	t.Run("Given validation errors, When Add is called, Then should append error to list", func(t *testing.T) {
		// Arrange
		valErrs := validation.ValidationErrors{}
		newErr := validation.ValidationError{Field: "email", Message: "invalid email"}

		// Act
		valErrs.Add(newErr)

		// Assert
		assert.Len(t, valErrs.Errors, 1)
		assert.Equal(t, newErr, valErrs.Errors[0])
	})
}

func TestValidationErrors_AddField(t *testing.T) {
	t.Run("Given validation errors, When AddField is called, Then should add field error", func(t *testing.T) {
		// Arrange
		valErrs := validation.ValidationErrors{}

		// Act
		valErrs.AddField("email", "invalid email")

		// Assert
		assert.Len(t, valErrs.Errors, 1)
		assert.Equal(t, "email", valErrs.Errors[0].Field)
		assert.Equal(t, "invalid email", valErrs.Errors[0].Message)
	})
}

func TestValidationErrors_HasErrors(t *testing.T) {
	tests := []struct {
		name     string
		valErrs  validation.ValidationErrors
		expected bool
	}{
		{
			name: "Given validation errors with errors, When HasErrors is called, Then should return true",
			valErrs: validation.ValidationErrors{
				Errors: []validation.ValidationError{
					{Field: "email", Message: "invalid email"},
				},
			},
			expected: true,
		},
		{
			name: "Given validation errors with no errors, When HasErrors is called, Then should return false",
			valErrs: validation.ValidationErrors{
				Errors: []validation.ValidationError{},
			},
			expected: false,
		},
		{
			name: "Given validation errors with nil errors, When HasErrors is called, Then should return false",
			valErrs: validation.ValidationErrors{
				Errors: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.valErrs.HasErrors()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidationErrors_HasFieldError(t *testing.T) {
	tests := []struct {
		name     string
		valErrs  validation.ValidationErrors
		field    string
		expected bool
	}{
		{
			name: "Given validation errors with specific field error, When HasFieldError is called, Then should return true",
			valErrs: validation.ValidationErrors{
				Errors: []validation.ValidationError{
					{Field: "email", Message: "invalid email"},
					{Field: "password", Message: "too short"},
				},
			},
			field:    "email",
			expected: true,
		},
		{
			name: "Given validation errors without specific field error, When HasFieldError is called, Then should return false",
			valErrs: validation.ValidationErrors{
				Errors: []validation.ValidationError{
					{Field: "password", Message: "too short"},
				},
			},
			field:    "email",
			expected: false,
		},
		{
			name: "Given validation errors with no errors, When HasFieldError is called, Then should return false",
			valErrs: validation.ValidationErrors{
				Errors: []validation.ValidationError{},
			},
			field:    "email",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.valErrs.HasFieldError(tt.field)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidationErrors_GetFieldErrors(t *testing.T) {
	tests := []struct {
		name          string
		valErrs       validation.ValidationErrors
		field         string
		expectedCount int
	}{
		{
			name: "Given validation errors with multiple errors for field, When GetFieldErrors is called, Then should return field errors",
			valErrs: validation.ValidationErrors{
				Errors: []validation.ValidationError{
					{Field: "email", Message: "invalid email"},
					{Field: "password", Message: "too short"},
					{Field: "email", Message: "already exists"},
				},
			},
			field:         "email",
			expectedCount: 2,
		},
		{
			name: "Given validation errors with no errors for field, When GetFieldErrors is called, Then should return empty slice",
			valErrs: validation.ValidationErrors{
				Errors: []validation.ValidationError{
					{Field: "password", Message: "too short"},
				},
			},
			field:         "email",
			expectedCount: 0,
		},
		{
			name: "Given validation errors with no errors, When GetFieldErrors is called, Then should return empty slice",
			valErrs: validation.ValidationErrors{
				Errors: []validation.ValidationError{},
			},
			field:         "email",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.valErrs.GetFieldErrors(tt.field)

			// Assert
			assert.Len(t, result, tt.expectedCount)
			if tt.expectedCount > 0 {
				for _, err := range result {
					assert.Equal(t, tt.field, err.Field)
				}
			}
		})
	}
}

func TestValidationResult_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		result   validation.ValidationResult
		expected bool
	}{
		{
			name: "Given validation result valid with no errors, When IsValid is called, Then should return true",
			result: validation.ValidationResult{
				Valid:  true,
				Errors: []validation.ValidationError{},
			},
			expected: true,
		},
		{
			name: "Given validation result valid but with errors, When IsValid is called, Then should return false",
			result: validation.ValidationResult{
				Valid: true,
				Errors: []validation.ValidationError{
					{Field: "email", Message: "invalid email"},
				},
			},
			expected: false,
		},
		{
			name: "Given validation result invalid with no errors, When IsValid is called, Then should return false",
			result: validation.ValidationResult{
				Valid:  false,
				Errors: []validation.ValidationError{},
			},
			expected: false,
		},
		{
			name: "Given validation result invalid with errors, When IsValid is called, Then should return false",
			result: validation.ValidationResult{
				Valid: false,
				Errors: []validation.ValidationError{
					{Field: "email", Message: "invalid email"},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.result.IsValid()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidationResult_AddError(t *testing.T) {
	t.Run("Given validation result, When AddError is called, Then should set valid to false and add error", func(t *testing.T) {
		// Arrange
		result := validation.ValidationResult{Valid: true}
		newErr := validation.ValidationError{Field: "email", Message: "invalid email"}

		// Act
		result.AddError(newErr)

		// Assert
		assert.False(t, result.Valid)
		assert.Len(t, result.Errors, 1)
		assert.Equal(t, newErr, result.Errors[0])
	})
}

func TestValidationConfig_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		config   validation.ValidationConfig
		expected bool
	}{
		{
			name: "Given validation config with default language, When IsValid is called, Then should return true",
			config: validation.ValidationConfig{
				DefaultLanguage: "en",
			},
			expected: true,
		},
		{
			name: "Given validation config with empty default language, When IsValid is called, Then should return false",
			config: validation.ValidationConfig{
				DefaultLanguage: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.config.IsValid()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultValidationConfig(t *testing.T) {
	t.Run("Given default validation config call, When DefaultValidationConfig is called, Then should return valid default configuration", func(t *testing.T) {
		// Act
		config := validation.DefaultValidationConfig()

		// Assert
		assert.False(t, config.StrictMode)
		assert.NotNil(t, config.CustomRules)
		assert.False(t, config.EnableI18n)
		assert.Equal(t, "en", config.DefaultLanguage)
	})
}

func TestValidationConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "Given ErrRequired constant, When accessing value, Then should have correct message",
			constant: validation.ErrRequired,
			expected: "field is required",
		},
		{
			name:     "Given ErrInvalidEmail constant, When accessing value, Then should have correct message",
			constant: validation.ErrInvalidEmail,
			expected: "invalid email format",
		},
		{
			name:     "Given ErrTooShort constant, When accessing value, Then should have correct message",
			constant: validation.ErrTooShort,
			expected: "value is too short",
		},
		{
			name:     "Given ErrTooLong constant, When accessing value, Then should have correct message",
			constant: validation.ErrTooLong,
			expected: "value is too long",
		},
		{
			name:     "Given ErrInvalidUUID constant, When accessing value, Then should have correct message",
			constant: validation.ErrInvalidUUID,
			expected: "invalid UUID format",
		},
		{
			name:     "Given ErrWeakPassword constant, When accessing value, Then should have correct message",
			constant: validation.ErrWeakPassword,
			expected: "password does not meet security requirements",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assert
			assert.Equal(t, tt.expected, tt.constant)
		})
	}
}

type mockValidationRule struct {
	name        string
	description string
}

func (m *mockValidationRule) Validate(ctx context.Context, value interface{}) error {
	return nil
}

func (m *mockValidationRule) Name() string {
	return m.name
}

func (m *mockValidationRule) Description() string {
	return m.description
}

func TestValidationConfig_CustomRules(t *testing.T) {
	t.Run("Given validation config with custom rules, When accessing custom rules, Then should handle custom rules correctly", func(t *testing.T) {
		// Arrange
		mockRule := &mockValidationRule{name: "test-rule", description: "Test rule"}
		config := validation.ValidationConfig{
			CustomRules: map[string]validationrule.Service{
				"test": mockRule,
			},
			DefaultLanguage: "en",
		}

		// Assert
		assert.True(t, config.IsValid())
		assert.NotNil(t, config.CustomRules)
		assert.Len(t, config.CustomRules, 1)
		
		rule, exists := config.CustomRules["test"]
		assert.True(t, exists)
		assert.Equal(t, "test-rule", rule.Name())
		assert.Equal(t, "Test rule", rule.Description())
	})
}