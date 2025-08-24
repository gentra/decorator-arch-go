package validationrule

import (
	"context"
)

// Service defines the validation rule domain interface - the ONLY interface in this domain
type Service interface {
	Validate(ctx context.Context, value interface{}) error
	Name() string
	Description() string
}

// Domain types and data structures

// ValidationRuleConfig contains configuration for validation rules
type ValidationRuleConfig struct {
	RuleID      string                 `json:"rule_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Enabled     bool                   `json:"enabled"`
	Priority    int                    `json:"priority"`
	Metadata    map[string]string      `json:"metadata,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// ValidationRuleResult contains the result of a validation rule execution
type ValidationRuleResult struct {
	RuleID   string                 `json:"rule_id"`
	Valid    bool                   `json:"valid"`
	Message  string                 `json:"message,omitempty"`
	Value    interface{}            `json:"value,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ValidationRuleError represents domain-specific validation rule errors
type ValidationRuleError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	RuleID  string `json:"rule_id,omitempty"`
	Field   string `json:"field,omitempty"`
}

func (e ValidationRuleError) Error() string {
	return e.Message
}

// Common validation rule error codes
var (
	ErrRuleNotFound  = ValidationRuleError{Code: "RULE_NOT_FOUND", Message: "Validation rule not found"}
	ErrRuleDisabled  = ValidationRuleError{Code: "RULE_DISABLED", Message: "Validation rule is disabled"}
	ErrInvalidValue  = ValidationRuleError{Code: "INVALID_VALUE", Message: "Value is invalid for this rule"}
	ErrRuleExecution = ValidationRuleError{Code: "RULE_EXECUTION", Message: "Error executing validation rule"}
	ErrInvalidConfig = ValidationRuleError{Code: "INVALID_CONFIG", Message: "Invalid rule configuration"}
)

// ValidationRuleType represents different types of validation rules
type ValidationRuleType string

const (
	ValidationRuleTypeFormat      ValidationRuleType = "format"
	ValidationRuleTypeLength      ValidationRuleType = "length"
	ValidationRuleTypeRange       ValidationRuleType = "range"
	ValidationRuleTypePattern     ValidationRuleType = "pattern"
	ValidationRuleTypeCustom      ValidationRuleType = "custom"
	ValidationRuleTypeRequired    ValidationRuleType = "required"
	ValidationRuleTypeConditional ValidationRuleType = "conditional"
)

// Helper methods for ValidationRuleConfig
func (c *ValidationRuleConfig) IsValid() bool {
	return c.RuleID != "" && c.Name != ""
}

func (c *ValidationRuleConfig) IsEnabled() bool {
	return c.Enabled
}

func (c *ValidationRuleConfig) GetParameter(key string) (interface{}, bool) {
	if c.Parameters == nil {
		return nil, false
	}
	value, exists := c.Parameters[key]
	return value, exists
}

func (c *ValidationRuleConfig) SetParameter(key string, value interface{}) {
	if c.Parameters == nil {
		c.Parameters = make(map[string]interface{})
	}
	c.Parameters[key] = value
}

// Helper methods for ValidationRuleResult
func (r *ValidationRuleResult) IsValid() bool {
	return r.Valid
}

func (r *ValidationRuleResult) HasMetadata() bool {
	return len(r.Metadata) > 0
}

func (r *ValidationRuleResult) GetMetadata(key string) (interface{}, bool) {
	if r.Metadata == nil {
		return nil, false
	}
	value, exists := r.Metadata[key]
	return value, exists
}

// Default validation rule configuration
func DefaultValidationRuleConfig() ValidationRuleConfig {
	return ValidationRuleConfig{
		Enabled:    true,
		Priority:   100,
		Parameters: make(map[string]interface{}),
		Metadata:   make(map[string]string),
	}
}

// Common validation rule priorities
const (
	PriorityHigh   = 10
	PriorityNormal = 100
	PriorityLow    = 1000
)
