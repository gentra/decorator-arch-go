package factory

import (
	"fmt"

	"github.com/gentra/decorator-arch-go/internal/validation"
	"github.com/gentra/decorator-arch-go/internal/validation/standard"
	"github.com/gentra/decorator-arch-go/internal/validationrule"
)

// Config contains all configuration for building the validation service
type Config struct {
	// Provider configuration
	Provider string // "standard", "custom", "external"

	// Validation engine settings
	Engine string // "go-playground", "ozzo", "custom"

	// Validation behavior
	StrictMode      bool
	EnableI18n      bool
	DefaultLanguage string

	// Custom rules configuration
	CustomRules   map[string]validationrule.Service
	CustomRuleDir string

	// External provider settings (for future implementation)
	ExternalURL    string
	ExternalAPIKey string

	// Performance settings
	CacheRules    bool
	CacheTTL      string
	ParallelMode  bool
	MaxGoroutines int

	// Feature flags
	Features FeatureFlags
}

// FeatureFlags controls validation service behavior
type FeatureFlags struct {
	EnableStandardProvider     bool
	EnableCustomProvider       bool
	EnableExternalProvider     bool
	EnableGoPlaygroundEngine   bool
	EnableOzzoEngine           bool
	EnableCustomRules          bool
	EnableI18nSupport          bool
	EnableRuleCaching          bool
	EnableParallelValidation   bool
	EnableFieldLevelErrors     bool
	EnableStructLevelErrors    bool
	EnableCrossFieldValidation bool
	EnableConditionalRules     bool
	EnableAsyncValidation      bool
	EnableMetrics              bool
}

// DefaultFeatureFlags returns default feature flag configuration
func DefaultFeatureFlags() FeatureFlags {
	return FeatureFlags{
		EnableStandardProvider:     true,
		EnableCustomProvider:       false,
		EnableExternalProvider:     false,
		EnableGoPlaygroundEngine:   true,
		EnableOzzoEngine:           false,
		EnableCustomRules:          true,
		EnableI18nSupport:          false,
		EnableRuleCaching:          true,
		EnableParallelValidation:   false,
		EnableFieldLevelErrors:     true,
		EnableStructLevelErrors:    true,
		EnableCrossFieldValidation: true,
		EnableConditionalRules:     false,
		EnableAsyncValidation:      false,
		EnableMetrics:              false,
	}
}

// ValidationServiceFactory creates and assembles the complete validation service
type ValidationServiceFactory struct {
	config Config
}

// NewFactory creates a new validation service factory with the given configuration
func NewFactory(config Config) *ValidationServiceFactory {
	return &ValidationServiceFactory{
		config: config,
	}
}

// Build assembles and returns the complete validation service based on configuration
func (f *ValidationServiceFactory) Build() (validation.Service, error) {
	switch f.config.Provider {
	case "standard":
		return f.buildStandardService()
	case "custom":
		return f.buildCustomService()
	case "external":
		return f.buildExternalService()
	default:
		// Default to standard provider
		return f.buildStandardService()
	}
}

// buildStandardService creates a standard validation service
func (f *ValidationServiceFactory) buildStandardService() (validation.Service, error) {
	switch f.config.Engine {
	case "go-playground":
		return standard.NewService(), nil
	case "ozzo":
		return f.buildOzzoService()
	default:
		// Default to go-playground engine
		return standard.NewService(), nil
	}
}

// buildOzzoService creates an Ozzo validation service (placeholder)
func (f *ValidationServiceFactory) buildOzzoService() (validation.Service, error) {
	// TODO: Implement Ozzo validation service
	return nil, fmt.Errorf("Ozzo validation engine not yet implemented")
}

// buildCustomService creates a custom validation service (placeholder)
func (f *ValidationServiceFactory) buildCustomService() (validation.Service, error) {
	// TODO: Implement custom validation service
	return nil, fmt.Errorf("custom validation provider not yet implemented")
}

// buildExternalService creates an external validation service (placeholder)
func (f *ValidationServiceFactory) buildExternalService() (validation.Service, error) {
	// TODO: Implement external validation service
	return nil, fmt.Errorf("external validation provider not yet implemented")
}

// DefaultConfig returns a sensible default configuration for the validation service
func DefaultConfig() Config {
	return Config{
		Provider:        "standard",
		Engine:          "go-playground",
		StrictMode:      false,
		EnableI18n:      false,
		DefaultLanguage: "en",
		CustomRules:     make(map[string]validationrule.Service),
		CacheRules:      true,
		CacheTTL:        "1h",
		ParallelMode:    false,
		MaxGoroutines:   10,
		Features:        DefaultFeatureFlags(),
	}
}

// ConfigBuilder provides a fluent interface for building validation configuration
type ConfigBuilder struct {
	config Config
}

// NewConfigBuilder creates a new configuration builder with defaults
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: DefaultConfig(),
	}
}

// WithProvider sets the validation provider
func (b *ConfigBuilder) WithProvider(provider string) *ConfigBuilder {
	b.config.Provider = provider
	return b
}

// WithEngine sets the validation engine
func (b *ConfigBuilder) WithEngine(engine string) *ConfigBuilder {
	b.config.Engine = engine
	return b
}

// WithStrictMode enables or disables strict validation mode
func (b *ConfigBuilder) WithStrictMode(strict bool) *ConfigBuilder {
	b.config.StrictMode = strict
	return b
}

// WithI18n enables internationalization support
func (b *ConfigBuilder) WithI18n(enable bool, defaultLanguage string) *ConfigBuilder {
	b.config.EnableI18n = enable
	b.config.DefaultLanguage = defaultLanguage
	b.config.Features.EnableI18nSupport = enable
	return b
}

// WithCustomRule adds a custom validation rule
func (b *ConfigBuilder) WithCustomRule(name string, rule validationrule.Service) *ConfigBuilder {
	if b.config.CustomRules == nil {
		b.config.CustomRules = make(map[string]validationrule.Service)
	}
	b.config.CustomRules[name] = rule
	b.config.Features.EnableCustomRules = true
	return b
}

// WithCustomRuleDir sets the directory for loading custom rules
func (b *ConfigBuilder) WithCustomRuleDir(dir string) *ConfigBuilder {
	b.config.CustomRuleDir = dir
	b.config.Features.EnableCustomRules = true
	return b
}

// WithExternalProvider sets external validation provider configuration
func (b *ConfigBuilder) WithExternalProvider(url, apiKey string) *ConfigBuilder {
	b.config.ExternalURL = url
	b.config.ExternalAPIKey = apiKey
	b.config.Provider = "external"
	b.config.Features.EnableExternalProvider = true
	return b
}

// WithCaching enables rule caching with TTL
func (b *ConfigBuilder) WithCaching(enable bool, ttl string) *ConfigBuilder {
	b.config.CacheRules = enable
	b.config.CacheTTL = ttl
	b.config.Features.EnableRuleCaching = enable
	return b
}

// WithParallelValidation enables parallel validation
func (b *ConfigBuilder) WithParallelValidation(enable bool, maxGoroutines int) *ConfigBuilder {
	b.config.ParallelMode = enable
	b.config.MaxGoroutines = maxGoroutines
	b.config.Features.EnableParallelValidation = enable
	return b
}

// WithFeatures sets the feature flags
func (b *ConfigBuilder) WithFeatures(features FeatureFlags) *ConfigBuilder {
	b.config.Features = features
	return b
}

// EnableOzzoEngine switches to Ozzo validation engine
func (b *ConfigBuilder) EnableOzzoEngine() *ConfigBuilder {
	b.config.Engine = "ozzo"
	b.config.Features.EnableOzzoEngine = true
	b.config.Features.EnableGoPlaygroundEngine = false
	return b
}

// EnableCrossFieldValidation enables cross-field validation rules
func (b *ConfigBuilder) EnableCrossFieldValidation() *ConfigBuilder {
	b.config.Features.EnableCrossFieldValidation = true
	return b
}

// EnableConditionalRules enables conditional validation rules
func (b *ConfigBuilder) EnableConditionalRules() *ConfigBuilder {
	b.config.Features.EnableConditionalRules = true
	return b
}

// EnableAsyncValidation enables asynchronous validation
func (b *ConfigBuilder) EnableAsyncValidation() *ConfigBuilder {
	b.config.Features.EnableAsyncValidation = true
	b.config.ParallelMode = true
	b.config.Features.EnableParallelValidation = true
	return b
}

// EnableMetrics enables validation metrics collection
func (b *ConfigBuilder) EnableMetrics() *ConfigBuilder {
	b.config.Features.EnableMetrics = true
	return b
}

// ForDevelopment configures the service for development use
func (b *ConfigBuilder) ForDevelopment() *ConfigBuilder {
	b.config.Provider = "standard"
	b.config.Engine = "go-playground"
	b.config.StrictMode = false
	b.config.CacheRules = false
	b.config.ParallelMode = false
	b.config.Features.EnableMetrics = false
	return b
}

// ForProduction configures the service for production use
func (b *ConfigBuilder) ForProduction() *ConfigBuilder {
	b.config.StrictMode = true
	b.config.CacheRules = true
	b.config.CacheTTL = "1h"
	b.config.ParallelMode = true
	b.config.MaxGoroutines = 10
	b.config.Features.EnableRuleCaching = true
	b.config.Features.EnableParallelValidation = true
	b.config.Features.EnableMetrics = true
	return b
}

// Build returns the final configuration
func (b *ConfigBuilder) Build() Config {
	return b.config
}
