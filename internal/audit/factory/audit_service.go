package factory

import (
	"github.com/gentra/decorator-arch-go/internal/audit"
	"github.com/gentra/decorator-arch-go/internal/audit/console"
)

// Config contains all configuration for building the audit service
type Config struct {
	// Output configuration
	OutputTarget string // "console", "file", "database", "external"

	// File output configuration (if OutputTarget = "file")
	LogFilePath string

	// Database configuration (if OutputTarget = "database")
	DatabaseDSN string

	// External service configuration (if OutputTarget = "external")
	ExternalURL    string
	ExternalAPIKey string

	// Feature flags
	Features FeatureFlags
}

// FeatureFlags controls audit service behavior
type FeatureFlags struct {
	EnableConsoleOutput   bool
	EnableFileOutput      bool
	EnableDatabaseOutput  bool
	EnableExternalOutput  bool
	EnableAsyncProcessing bool
	EnableBatching        bool
	EnableCompression     bool
}

// DefaultFeatureFlags returns default feature flag configuration
func DefaultFeatureFlags() FeatureFlags {
	return FeatureFlags{
		EnableConsoleOutput:   true,
		EnableFileOutput:      false,
		EnableDatabaseOutput:  false,
		EnableExternalOutput:  false,
		EnableAsyncProcessing: false,
		EnableBatching:        false,
		EnableCompression:     false,
	}
}

// AuditServiceFactory creates and assembles the complete audit service
type AuditServiceFactory struct {
	config Config
}

// NewFactory creates a new audit service factory with the given configuration
func NewFactory(config Config) *AuditServiceFactory {
	return &AuditServiceFactory{
		config: config,
	}
}

// Build assembles and returns the complete audit service based on configuration
func (f *AuditServiceFactory) Build() (audit.Service, error) {
	// For now, we only have console implementation
	// In the future, we can add strategy pattern here for different outputs

	if f.config.Features.EnableConsoleOutput {
		return console.NewService(), nil
	}

	// Default fallback to console
	return console.NewService(), nil
}

// DefaultConfig returns a sensible default configuration for the audit service
func DefaultConfig() Config {
	return Config{
		OutputTarget: "console",
		LogFilePath:  "/var/log/audit.log",
		Features:     DefaultFeatureFlags(),
	}
}

// ConfigBuilder provides a fluent interface for building audit configuration
type ConfigBuilder struct {
	config Config
}

// NewConfigBuilder creates a new configuration builder with defaults
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: DefaultConfig(),
	}
}

// WithOutputTarget sets the output target
func (b *ConfigBuilder) WithOutputTarget(target string) *ConfigBuilder {
	b.config.OutputTarget = target
	return b
}

// WithLogFilePath sets the log file path
func (b *ConfigBuilder) WithLogFilePath(path string) *ConfigBuilder {
	b.config.LogFilePath = path
	return b
}

// WithDatabaseDSN sets the database connection string
func (b *ConfigBuilder) WithDatabaseDSN(dsn string) *ConfigBuilder {
	b.config.DatabaseDSN = dsn
	return b
}

// WithExternalService sets external service configuration
func (b *ConfigBuilder) WithExternalService(url, apiKey string) *ConfigBuilder {
	b.config.ExternalURL = url
	b.config.ExternalAPIKey = apiKey
	return b
}

// WithFeatures sets the feature flags
func (b *ConfigBuilder) WithFeatures(features FeatureFlags) *ConfigBuilder {
	b.config.Features = features
	return b
}

// EnableAsyncProcessing enables asynchronous audit processing
func (b *ConfigBuilder) EnableAsyncProcessing() *ConfigBuilder {
	b.config.Features.EnableAsyncProcessing = true
	return b
}

// EnableBatching enables audit entry batching
func (b *ConfigBuilder) EnableBatching() *ConfigBuilder {
	b.config.Features.EnableBatching = true
	return b
}

// EnableCompression enables audit entry compression
func (b *ConfigBuilder) EnableCompression() *ConfigBuilder {
	b.config.Features.EnableCompression = true
	return b
}

// Build returns the final configuration
func (b *ConfigBuilder) Build() Config {
	return b.config
}
