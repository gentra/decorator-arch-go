package validationrule_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/validationrule"
)

func TestValidationRuleConfig_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		config   validationrule.ValidationRuleConfig
		expected bool
	}{
		{
			name: "Given validation rule config with rule ID and name, When IsValid is called, Then should return true",
			config: validationrule.ValidationRuleConfig{
				RuleID: "email-rule",
				Name:   "Email Validation Rule",
			},
			expected: true,
		},
		{
			name: "Given validation rule config with empty rule ID, When IsValid is called, Then should return false",
			config: validationrule.ValidationRuleConfig{
				RuleID: "",
				Name:   "Email Validation Rule",
			},
			expected: false,
		},
		{
			name: "Given validation rule config with empty name, When IsValid is called, Then should return false",
			config: validationrule.ValidationRuleConfig{
				RuleID: "email-rule",
				Name:   "",
			},
			expected: false,
		},
		{
			name: "Given validation rule config with both rule ID and name empty, When IsValid is called, Then should return false",
			config: validationrule.ValidationRuleConfig{
				RuleID: "",
				Name:   "",
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

func TestValidationRuleConfig_IsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		config   validationrule.ValidationRuleConfig
		expected bool
	}{
		{
			name: "Given validation rule config with enabled true, When IsEnabled is called, Then should return true",
			config: validationrule.ValidationRuleConfig{
				Enabled: true,
			},
			expected: true,
		},
		{
			name: "Given validation rule config with enabled false, When IsEnabled is called, Then should return false",
			config: validationrule.ValidationRuleConfig{
				Enabled: false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.config.IsEnabled()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidationRuleConfig_GetParameter(t *testing.T) {
	tests := []struct {
		name           string
		config         validationrule.ValidationRuleConfig
		key            string
		expectedValue  interface{}
		expectedExists bool
	}{
		{
			name: "Given validation rule config with parameter, When GetParameter is called with existing key, Then should return value and true",
			config: validationrule.ValidationRuleConfig{
				Parameters: map[string]interface{}{
					"min_length": 8,
					"max_length": 255,
					"pattern":    "^[a-zA-Z0-9]+$",
				},
			},
			key:            "min_length",
			expectedValue:  8,
			expectedExists: true,
		},
		{
			name: "Given validation rule config with parameter, When GetParameter is called with non-existing key, Then should return nil and false",
			config: validationrule.ValidationRuleConfig{
				Parameters: map[string]interface{}{
					"min_length": 8,
				},
			},
			key:            "max_length",
			expectedValue:  nil,
			expectedExists: false,
		},
		{
			name: "Given validation rule config with nil parameters, When GetParameter is called, Then should return nil and false",
			config: validationrule.ValidationRuleConfig{
				Parameters: nil,
			},
			key:            "min_length",
			expectedValue:  nil,
			expectedExists: false,
		},
		{
			name: "Given validation rule config with empty parameters, When GetParameter is called, Then should return nil and false",
			config: validationrule.ValidationRuleConfig{
				Parameters: map[string]interface{}{},
			},
			key:            "min_length",
			expectedValue:  nil,
			expectedExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			value, exists := tt.config.GetParameter(tt.key)

			// Assert
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedExists, exists)
		})
	}
}

func TestValidationRuleConfig_SetParameter(t *testing.T) {
	tests := []struct {
		name         string
		config       validationrule.ValidationRuleConfig
		key          string
		value        interface{}
		expectedFunc func(t *testing.T, config validationrule.ValidationRuleConfig)
	}{
		{
			name: "Given validation rule config with existing parameters, When SetParameter is called, Then should set parameter",
			config: validationrule.ValidationRuleConfig{
				Parameters: map[string]interface{}{
					"existing": "value",
				},
			},
			key:   "min_length",
			value: 10,
			expectedFunc: func(t *testing.T, config validationrule.ValidationRuleConfig) {
				assert.NotNil(t, config.Parameters)
				assert.Equal(t, 10, config.Parameters["min_length"])
				assert.Equal(t, "value", config.Parameters["existing"])
			},
		},
		{
			name: "Given validation rule config with nil parameters, When SetParameter is called, Then should create parameters map and set parameter",
			config: validationrule.ValidationRuleConfig{
				Parameters: nil,
			},
			key:   "pattern",
			value: "^[a-zA-Z]+$",
			expectedFunc: func(t *testing.T, config validationrule.ValidationRuleConfig) {
				assert.NotNil(t, config.Parameters)
				assert.Equal(t, "^[a-zA-Z]+$", config.Parameters["pattern"])
			},
		},
		{
			name: "Given validation rule config, When SetParameter is called with overwrite, Then should overwrite existing parameter",
			config: validationrule.ValidationRuleConfig{
				Parameters: map[string]interface{}{
					"max_length": 100,
				},
			},
			key:   "max_length",
			value: 200,
			expectedFunc: func(t *testing.T, config validationrule.ValidationRuleConfig) {
				assert.NotNil(t, config.Parameters)
				assert.Equal(t, 200, config.Parameters["max_length"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			tt.config.SetParameter(tt.key, tt.value)

			// Assert
			tt.expectedFunc(t, tt.config)
		})
	}
}

func TestValidationRuleResult_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		result   validationrule.ValidationRuleResult
		expected bool
	}{
		{
			name: "Given validation rule result with valid true, When IsValid is called, Then should return true",
			result: validationrule.ValidationRuleResult{
				Valid: true,
			},
			expected: true,
		},
		{
			name: "Given validation rule result with valid false, When IsValid is called, Then should return false",
			result: validationrule.ValidationRuleResult{
				Valid: false,
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

func TestValidationRuleResult_HasMetadata(t *testing.T) {
	tests := []struct {
		name     string
		result   validationrule.ValidationRuleResult
		expected bool
	}{
		{
			name: "Given validation rule result with metadata, When HasMetadata is called, Then should return true",
			result: validationrule.ValidationRuleResult{
				Metadata: map[string]interface{}{
					"checked_patterns": 3,
					"match_score":      0.95,
				},
			},
			expected: true,
		},
		{
			name: "Given validation rule result with empty metadata, When HasMetadata is called, Then should return false",
			result: validationrule.ValidationRuleResult{
				Metadata: map[string]interface{}{},
			},
			expected: false,
		},
		{
			name: "Given validation rule result with nil metadata, When HasMetadata is called, Then should return false",
			result: validationrule.ValidationRuleResult{
				Metadata: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.result.HasMetadata()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidationRuleResult_GetMetadata(t *testing.T) {
	tests := []struct {
		name           string
		result         validationrule.ValidationRuleResult
		key            string
		expectedValue  interface{}
		expectedExists bool
	}{
		{
			name: "Given validation rule result with metadata, When GetMetadata is called with existing key, Then should return value and true",
			result: validationrule.ValidationRuleResult{
				Metadata: map[string]interface{}{
					"score":    0.8,
					"attempts": 2,
					"details":  "validation passed",
				},
			},
			key:            "score",
			expectedValue:  0.8,
			expectedExists: true,
		},
		{
			name: "Given validation rule result with metadata, When GetMetadata is called with non-existing key, Then should return nil and false",
			result: validationrule.ValidationRuleResult{
				Metadata: map[string]interface{}{
					"score": 0.8,
				},
			},
			key:            "attempts",
			expectedValue:  nil,
			expectedExists: false,
		},
		{
			name: "Given validation rule result with nil metadata, When GetMetadata is called, Then should return nil and false",
			result: validationrule.ValidationRuleResult{
				Metadata: nil,
			},
			key:            "score",
			expectedValue:  nil,
			expectedExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			value, exists := tt.result.GetMetadata(tt.key)

			// Assert
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedExists, exists)
		})
	}
}

func TestDefaultValidationRuleConfig(t *testing.T) {
	t.Run("Given default validation rule config call, When DefaultValidationRuleConfig is called, Then should return valid default configuration", func(t *testing.T) {
		// Act
		config := validationrule.DefaultValidationRuleConfig()

		// Assert
		assert.True(t, config.Enabled)
		assert.Equal(t, 100, config.Priority)
		assert.NotNil(t, config.Parameters)
		assert.NotNil(t, config.Metadata)
		assert.Empty(t, config.Parameters)
		assert.Empty(t, config.Metadata)
	})
}

func TestValidationRuleError_Error(t *testing.T) {
	tests := []struct {
		name     string
		ruleErr  validationrule.ValidationRuleError
		expected string
	}{
		{
			name: "Given validation rule error with message, When Error is called, Then should return message",
			ruleErr: validationrule.ValidationRuleError{
				Code:    "TEST_ERROR",
				Message: "Test rule error",
			},
			expected: "Test rule error",
		},
		{
			name: "Given validation rule error with empty message, When Error is called, Then should return empty string",
			ruleErr: validationrule.ValidationRuleError{
				Code:    "TEST_ERROR",
				Message: "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.ruleErr.Error()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidationRuleErrors_Constants(t *testing.T) {
	tests := []struct {
		name         string
		err          validationrule.ValidationRuleError
		expectedCode string
	}{
		{
			name:         "Given ErrRuleNotFound, When accessing code, Then should have correct code",
			err:          validationrule.ErrRuleNotFound,
			expectedCode: "RULE_NOT_FOUND",
		},
		{
			name:         "Given ErrRuleDisabled, When accessing code, Then should have correct code",
			err:          validationrule.ErrRuleDisabled,
			expectedCode: "RULE_DISABLED",
		},
		{
			name:         "Given ErrInvalidValue, When accessing code, Then should have correct code",
			err:          validationrule.ErrInvalidValue,
			expectedCode: "INVALID_VALUE",
		},
		{
			name:         "Given ErrRuleExecution, When accessing code, Then should have correct code",
			err:          validationrule.ErrRuleExecution,
			expectedCode: "RULE_EXECUTION",
		},
		{
			name:         "Given ErrInvalidConfig, When accessing code, Then should have correct code",
			err:          validationrule.ErrInvalidConfig,
			expectedCode: "INVALID_CONFIG",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assert
			assert.Equal(t, tt.expectedCode, tt.err.Code)
			assert.NotEmpty(t, tt.err.Message)
		})
	}
}

func TestValidationRuleType_Constants(t *testing.T) {
	tests := []struct {
		name         string
		ruleType     validationrule.ValidationRuleType
		expectedStr  string
	}{
		{
			name:         "Given ValidationRuleTypeFormat constant, When accessing string value, Then should have correct value",
			ruleType:     validationrule.ValidationRuleTypeFormat,
			expectedStr:  "format",
		},
		{
			name:         "Given ValidationRuleTypeLength constant, When accessing string value, Then should have correct value",
			ruleType:     validationrule.ValidationRuleTypeLength,
			expectedStr:  "length",
		},
		{
			name:         "Given ValidationRuleTypeRange constant, When accessing string value, Then should have correct value",
			ruleType:     validationrule.ValidationRuleTypeRange,
			expectedStr:  "range",
		},
		{
			name:         "Given ValidationRuleTypePattern constant, When accessing string value, Then should have correct value",
			ruleType:     validationrule.ValidationRuleTypePattern,
			expectedStr:  "pattern",
		},
		{
			name:         "Given ValidationRuleTypeCustom constant, When accessing string value, Then should have correct value",
			ruleType:     validationrule.ValidationRuleTypeCustom,
			expectedStr:  "custom",
		},
		{
			name:         "Given ValidationRuleTypeRequired constant, When accessing string value, Then should have correct value",
			ruleType:     validationrule.ValidationRuleTypeRequired,
			expectedStr:  "required",
		},
		{
			name:         "Given ValidationRuleTypeConditional constant, When accessing string value, Then should have correct value",
			ruleType:     validationrule.ValidationRuleTypeConditional,
			expectedStr:  "conditional",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assert
			assert.Equal(t, tt.expectedStr, string(tt.ruleType))
		})
	}
}

func TestPriorityConstants(t *testing.T) {
	tests := []struct {
		name         string
		priority     int
		expectedVal  int
	}{
		{
			name:         "Given PriorityHigh constant, When accessing value, Then should have correct value",
			priority:     validationrule.PriorityHigh,
			expectedVal:  10,
		},
		{
			name:         "Given PriorityNormal constant, When accessing value, Then should have correct value",
			priority:     validationrule.PriorityNormal,
			expectedVal:  100,
		},
		{
			name:         "Given PriorityLow constant, When accessing value, Then should have correct value",
			priority:     validationrule.PriorityLow,
			expectedVal:  1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assert
			assert.Equal(t, tt.expectedVal, tt.priority)
		})
	}
}

func TestValidationRuleConfig_CompleteStructure(t *testing.T) {
	t.Run("Given validation rule config with all fields, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		parameters := map[string]interface{}{
			"min_length": 8,
			"max_length": 255,
			"pattern":    "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
		}
		metadata := map[string]string{
			"category": "format",
			"version":  "2.0",
		}

		config := validationrule.ValidationRuleConfig{
			RuleID:      "email-validation",
			Name:        "Email Format Validation",
			Description: "Validates email format using regex pattern",
			Enabled:     true,
			Priority:    validationrule.PriorityHigh,
			Metadata:    metadata,
			Parameters:  parameters,
		}

		// Assert
		assert.Equal(t, "email-validation", config.RuleID)
		assert.Equal(t, "Email Format Validation", config.Name)
		assert.Equal(t, "Validates email format using regex pattern", config.Description)
		assert.True(t, config.Enabled)
		assert.Equal(t, validationrule.PriorityHigh, config.Priority)
		assert.Equal(t, metadata, config.Metadata)
		assert.Equal(t, parameters, config.Parameters)
		
		// Test helper methods
		assert.True(t, config.IsValid())
		assert.True(t, config.IsEnabled())
		
		// Test parameter access
		minLength, exists := config.GetParameter("min_length")
		assert.True(t, exists)
		assert.Equal(t, 8, minLength)
		
		pattern, exists := config.GetParameter("pattern")
		assert.True(t, exists)
		assert.Equal(t, "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", pattern)
		
		nonExistent, exists := config.GetParameter("non_existent")
		assert.False(t, exists)
		assert.Nil(t, nonExistent)
	})
}

func TestValidationRuleResult_CompleteStructure(t *testing.T) {
	t.Run("Given validation rule result with all fields, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		metadata := map[string]interface{}{
			"execution_time": "1.2ms",
			"pattern_match":  true,
			"score":          0.95,
		}

		result := validationrule.ValidationRuleResult{
			RuleID:   "email-validation",
			Valid:    true,
			Message:  "Email format is valid",
			Value:    "test@example.com",
			Metadata: metadata,
		}

		// Assert
		assert.Equal(t, "email-validation", result.RuleID)
		assert.True(t, result.Valid)
		assert.Equal(t, "Email format is valid", result.Message)
		assert.Equal(t, "test@example.com", result.Value)
		assert.Equal(t, metadata, result.Metadata)
		
		// Test helper methods
		assert.True(t, result.IsValid())
		assert.True(t, result.HasMetadata())
		
		// Test metadata access
		execTime, exists := result.GetMetadata("execution_time")
		assert.True(t, exists)
		assert.Equal(t, "1.2ms", execTime)
		
		score, exists := result.GetMetadata("score")
		assert.True(t, exists)
		assert.Equal(t, 0.95, score)
		
		nonExistent, exists := result.GetMetadata("non_existent")
		assert.False(t, exists)
		assert.Nil(t, nonExistent)
	})
}