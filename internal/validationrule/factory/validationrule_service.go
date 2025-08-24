package factory

import (
	"fmt"

	"github.com/gentra/decorator-arch-go/internal/validationrule"
)

// Config contains all configuration for building the validation rule service
type Config struct {
	// Rule configuration
	RuleType    string // "format", "length", "range", "pattern", "custom", "conditional"
	RuleName    string
	Description string

	// Rule parameters
	Parameters map[string]interface{}

	// Format rule parameters
	FormatType string // "email", "url", "phone", "uuid", "date", "time"

	// Length rule parameters
	MinLength int
	MaxLength int

	// Range rule parameters
	MinValue interface{}
	MaxValue interface{}

	// Pattern rule parameters
	Regex      string
	RegexFlags string

	// Custom rule parameters
	CustomLogic   func(interface{}) error
	CustomMessage string

	// Conditional rule parameters
	Condition       func(interface{}) bool
	ConditionalRule validationrule.Service

	// Performance settings
	CacheResults bool
	CacheTTL     string
	CompileRegex bool

	// Feature flags
	Features FeatureFlags
}

// FeatureFlags controls validation rule service behavior
type FeatureFlags struct {
	EnableFormatRules      bool
	EnableLengthRules      bool
	EnableRangeRules       bool
	EnablePatternRules     bool
	EnableCustomRules      bool
	EnableConditionalRules bool
	EnableCompositeRules   bool
	EnableAsyncRules       bool
	EnableCachedResults    bool
	EnableRegexCompilation bool
	EnableI18nMessages     bool
	EnableMetrics          bool
	EnableParameterBinding bool
	EnableRuleChaining     bool
}

// DefaultFeatureFlags returns default feature flag configuration
func DefaultFeatureFlags() FeatureFlags {
	return FeatureFlags{
		EnableFormatRules:      true,
		EnableLengthRules:      true,
		EnableRangeRules:       true,
		EnablePatternRules:     true,
		EnableCustomRules:      true,
		EnableConditionalRules: false,
		EnableCompositeRules:   false,
		EnableAsyncRules:       false,
		EnableCachedResults:    true,
		EnableRegexCompilation: true,
		EnableI18nMessages:     false,
		EnableMetrics:          false,
		EnableParameterBinding: true,
		EnableRuleChaining:     false,
	}
}

// ValidationRuleServiceFactory creates and assembles the complete validation rule service
type ValidationRuleServiceFactory struct {
	config Config
}

// NewFactory creates a new validation rule service factory with the given configuration
func NewFactory(config Config) *ValidationRuleServiceFactory {
	return &ValidationRuleServiceFactory{
		config: config,
	}
}

// Build assembles and returns the complete validation rule service based on configuration
func (f *ValidationRuleServiceFactory) Build() (validationrule.Service, error) {
	switch f.config.RuleType {
	case "format":
		return f.buildFormatRule()
	case "length":
		return f.buildLengthRule()
	case "range":
		return f.buildRangeRule()
	case "pattern":
		return f.buildPatternRule()
	case "custom":
		return f.buildCustomRule()
	case "conditional":
		return f.buildConditionalRule()
	default:
		return nil, fmt.Errorf("unknown rule type: %s", f.config.RuleType)
	}
}

// buildFormatRule creates a format validation rule (placeholder)
func (f *ValidationRuleServiceFactory) buildFormatRule() (validationrule.Service, error) {
	// TODO: Implement format validation rule
	return nil, fmt.Errorf("format validation rule not yet implemented")
}

// buildLengthRule creates a length validation rule (placeholder)
func (f *ValidationRuleServiceFactory) buildLengthRule() (validationrule.Service, error) {
	// TODO: Implement length validation rule
	return nil, fmt.Errorf("length validation rule not yet implemented")
}

// buildRangeRule creates a range validation rule (placeholder)
func (f *ValidationRuleServiceFactory) buildRangeRule() (validationrule.Service, error) {
	// TODO: Implement range validation rule
	return nil, fmt.Errorf("range validation rule not yet implemented")
}

// buildPatternRule creates a pattern validation rule (placeholder)
func (f *ValidationRuleServiceFactory) buildPatternRule() (validationrule.Service, error) {
	// TODO: Implement pattern validation rule
	return nil, fmt.Errorf("pattern validation rule not yet implemented")
}

// buildCustomRule creates a custom validation rule (placeholder)
func (f *ValidationRuleServiceFactory) buildCustomRule() (validationrule.Service, error) {
	// TODO: Implement custom validation rule
	return nil, fmt.Errorf("custom validation rule not yet implemented")
}

// buildConditionalRule creates a conditional validation rule (placeholder)
func (f *ValidationRuleServiceFactory) buildConditionalRule() (validationrule.Service, error) {
	// TODO: Implement conditional validation rule
	return nil, fmt.Errorf("conditional validation rule not yet implemented")
}

// DefaultConfig returns a sensible default configuration for the validation rule service
func DefaultConfig() Config {
	return Config{
		RuleType:     "format",
		RuleName:     "default",
		Description:  "Default validation rule",
		Parameters:   make(map[string]interface{}),
		FormatType:   "email",
		MinLength:    0,
		MaxLength:    255,
		CacheResults: true,
		CacheTTL:     "1h",
		CompileRegex: true,
		Features:     DefaultFeatureFlags(),
	}
}

// ConfigBuilder provides a fluent interface for building validation rule configuration
type ConfigBuilder struct {
	config Config
}

// NewConfigBuilder creates a new configuration builder with defaults
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: DefaultConfig(),
	}
}

// WithRuleType sets the validation rule type
func (b *ConfigBuilder) WithRuleType(ruleType string) *ConfigBuilder {
	b.config.RuleType = ruleType
	return b
}

// WithRuleName sets the rule name and description
func (b *ConfigBuilder) WithRuleName(name, description string) *ConfigBuilder {
	b.config.RuleName = name
	b.config.Description = description
	return b
}

// WithParameters sets rule parameters
func (b *ConfigBuilder) WithParameters(params map[string]interface{}) *ConfigBuilder {
	b.config.Parameters = params
	return b
}

// WithParameter adds a single parameter
func (b *ConfigBuilder) WithParameter(key string, value interface{}) *ConfigBuilder {
	if b.config.Parameters == nil {
		b.config.Parameters = make(map[string]interface{})
	}
	b.config.Parameters[key] = value
	return b
}

// WithFormatType sets the format type for format rules
func (b *ConfigBuilder) WithFormatType(formatType string) *ConfigBuilder {
	b.config.FormatType = formatType
	b.config.RuleType = "format"
	return b
}

// WithLengthRange sets the length range for length rules
func (b *ConfigBuilder) WithLengthRange(min, max int) *ConfigBuilder {
	b.config.MinLength = min
	b.config.MaxLength = max
	b.config.RuleType = "length"
	return b
}

// WithValueRange sets the value range for range rules
func (b *ConfigBuilder) WithValueRange(min, max interface{}) *ConfigBuilder {
	b.config.MinValue = min
	b.config.MaxValue = max
	b.config.RuleType = "range"
	return b
}

// WithRegexPattern sets the regex pattern for pattern rules
func (b *ConfigBuilder) WithRegexPattern(pattern, flags string) *ConfigBuilder {
	b.config.Regex = pattern
	b.config.RegexFlags = flags
	b.config.RuleType = "pattern"
	return b
}

// WithCustomLogic sets custom validation logic
func (b *ConfigBuilder) WithCustomLogic(logic func(interface{}) error, message string) *ConfigBuilder {
	b.config.CustomLogic = logic
	b.config.CustomMessage = message
	b.config.RuleType = "custom"
	return b
}

// WithConditionalRule sets conditional validation
func (b *ConfigBuilder) WithConditionalRule(condition func(interface{}) bool, rule validationrule.Service) *ConfigBuilder {
	b.config.Condition = condition
	b.config.ConditionalRule = rule
	b.config.RuleType = "conditional"
	return b
}

// WithCaching enables result caching
func (b *ConfigBuilder) WithCaching(enable bool, ttl string) *ConfigBuilder {
	b.config.CacheResults = enable
	b.config.CacheTTL = ttl
	b.config.Features.EnableCachedResults = enable
	return b
}

// WithFeatures sets the feature flags
func (b *ConfigBuilder) WithFeatures(features FeatureFlags) *ConfigBuilder {
	b.config.Features = features
	return b
}

// AsEmailRule creates an email format validation rule
func (b *ConfigBuilder) AsEmailRule() *ConfigBuilder {
	b.config.RuleType = "format"
	b.config.FormatType = "email"
	b.config.RuleName = "email"
	b.config.Description = "Email format validation"
	return b
}

// AsURLRule creates a URL format validation rule
func (b *ConfigBuilder) AsURLRule() *ConfigBuilder {
	b.config.RuleType = "format"
	b.config.FormatType = "url"
	b.config.RuleName = "url"
	b.config.Description = "URL format validation"
	return b
}

// AsPhoneRule creates a phone format validation rule
func (b *ConfigBuilder) AsPhoneRule() *ConfigBuilder {
	b.config.RuleType = "format"
	b.config.FormatType = "phone"
	b.config.RuleName = "phone"
	b.config.Description = "Phone number format validation"
	return b
}

// AsUUIDRule creates a UUID format validation rule
func (b *ConfigBuilder) AsUUIDRule() *ConfigBuilder {
	b.config.RuleType = "format"
	b.config.FormatType = "uuid"
	b.config.RuleName = "uuid"
	b.config.Description = "UUID format validation"
	return b
}

// AsRequiredRule creates a required field validation rule
func (b *ConfigBuilder) AsRequiredRule() *ConfigBuilder {
	b.config.RuleType = "custom"
	b.config.RuleName = "required"
	b.config.Description = "Required field validation"
	b.config.CustomMessage = "field is required"
	return b
}

// AsPasswordRule creates a password strength validation rule
func (b *ConfigBuilder) AsPasswordRule() *ConfigBuilder {
	b.config.RuleType = "custom"
	b.config.RuleName = "password"
	b.config.Description = "Password strength validation"
	b.config.MinLength = 8
	b.config.MaxLength = 128
	return b
}

// EnableConditionalLogic enables conditional rule evaluation
func (b *ConfigBuilder) EnableConditionalLogic() *ConfigBuilder {
	b.config.Features.EnableConditionalRules = true
	return b
}

// EnableCompositeRules enables composite rule combining
func (b *ConfigBuilder) EnableCompositeRules() *ConfigBuilder {
	b.config.Features.EnableCompositeRules = true
	return b
}

// EnableAsyncValidation enables asynchronous rule evaluation
func (b *ConfigBuilder) EnableAsyncValidation() *ConfigBuilder {
	b.config.Features.EnableAsyncRules = true
	return b
}

// EnableMetrics enables rule performance metrics
func (b *ConfigBuilder) EnableMetrics() *ConfigBuilder {
	b.config.Features.EnableMetrics = true
	return b
}

// EnableRuleChaining enables rule chaining capabilities
func (b *ConfigBuilder) EnableRuleChaining() *ConfigBuilder {
	b.config.Features.EnableRuleChaining = true
	return b
}

// ForDevelopment configures the service for development use
func (b *ConfigBuilder) ForDevelopment() *ConfigBuilder {
	b.config.CacheResults = false
	b.config.CompileRegex = false
	b.config.Features.EnableCachedResults = false
	b.config.Features.EnableMetrics = false
	return b
}

// ForProduction configures the service for production use
func (b *ConfigBuilder) ForProduction() *ConfigBuilder {
	b.config.CacheResults = true
	b.config.CacheTTL = "1h"
	b.config.CompileRegex = true
	b.config.Features.EnableCachedResults = true
	b.config.Features.EnableRegexCompilation = true
	b.config.Features.EnableMetrics = true
	return b
}

// Build returns the final configuration
func (b *ConfigBuilder) Build() Config {
	return b.config
}
