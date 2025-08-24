package factory

import (
	"fmt"

	"github.com/gentra/decorator-arch-go/internal/ratelimit"
	"github.com/gentra/decorator-arch-go/internal/ratelimit/memory"
)

// Config contains all configuration for building the rate limit service
type Config struct {
	// Provider configuration
	Provider string // "memory", "redis", "database"

	// Algorithm configuration
	Algorithm string // "sliding_window", "fixed_window", "token_bucket", "leaky_bucket"

	// Memory provider settings
	CleanupInterval string // Duration string like "5m", "1h"

	// Redis provider settings (for future implementation)
	RedisURL       string
	RedisPassword  string
	RedisDB        int
	RedisKeyPrefix string

	// Database provider settings (for future implementation)
	DatabaseDSN string
	TableName   string

	// Default rate limits
	DefaultLimits map[string]ratelimit.RateLimitConfig

	// Global settings
	EnableGlobalLimits bool
	GlobalLimitConfig  ratelimit.RateLimitConfig

	// Feature flags
	Features FeatureFlags
}

// FeatureFlags controls rate limit service behavior
type FeatureFlags struct {
	EnableMemoryProvider    bool
	EnableRedisProvider     bool
	EnableDatabaseProvider  bool
	EnableSlidingWindow     bool
	EnableFixedWindow       bool
	EnableTokenBucket       bool
	EnableLeakyBucket       bool
	EnableDistributedLimits bool
	EnablePerUserLimits     bool
	EnablePerIPLimits       bool
	EnableCustomPatterns    bool
	EnableMetrics           bool
	EnableDynamicLimits     bool
	EnableGracePeriod       bool
}

// DefaultFeatureFlags returns default feature flag configuration
func DefaultFeatureFlags() FeatureFlags {
	return FeatureFlags{
		EnableMemoryProvider:    true,
		EnableRedisProvider:     false,
		EnableDatabaseProvider:  false,
		EnableSlidingWindow:     true,
		EnableFixedWindow:       false,
		EnableTokenBucket:       false,
		EnableLeakyBucket:       false,
		EnableDistributedLimits: false,
		EnablePerUserLimits:     true,
		EnablePerIPLimits:       true,
		EnableCustomPatterns:    true,
		EnableMetrics:           false,
		EnableDynamicLimits:     false,
		EnableGracePeriod:       false,
	}
}

// RateLimitServiceFactory creates and assembles the complete rate limit service
type RateLimitServiceFactory struct {
	config Config
}

// NewFactory creates a new rate limit service factory with the given configuration
func NewFactory(config Config) *RateLimitServiceFactory {
	return &RateLimitServiceFactory{
		config: config,
	}
}

// Build assembles and returns the complete rate limit service based on configuration
func (f *RateLimitServiceFactory) Build() (ratelimit.Service, error) {
	switch f.config.Provider {
	case "memory":
		return f.buildMemoryService()
	case "redis":
		return f.buildRedisService()
	case "database":
		return f.buildDatabaseService()
	default:
		// Default to memory provider
		return f.buildMemoryService()
	}
}

// buildMemoryService creates an in-memory rate limit service
func (f *RateLimitServiceFactory) buildMemoryService() (ratelimit.Service, error) {
	defaultLimits := f.config.DefaultLimits
	if defaultLimits == nil {
		defaultLimits = ratelimit.GetDefaultRateLimitConfigs()
	}

	return memory.NewService(defaultLimits), nil
}

// buildRedisService creates a Redis-based rate limit service (placeholder)
func (f *RateLimitServiceFactory) buildRedisService() (ratelimit.Service, error) {
	// TODO: Implement Redis rate limit service
	return nil, fmt.Errorf("Redis rate limit provider not yet implemented")
}

// buildDatabaseService creates a database-based rate limit service (placeholder)
func (f *RateLimitServiceFactory) buildDatabaseService() (ratelimit.Service, error) {
	// TODO: Implement database rate limit service
	return nil, fmt.Errorf("Database rate limit provider not yet implemented")
}

// DefaultConfig returns a sensible default configuration for the rate limit service
func DefaultConfig() Config {
	return Config{
		Provider:        "memory",
		Algorithm:       "sliding_window",
		CleanupInterval: "5m",
		DefaultLimits:   ratelimit.GetDefaultRateLimitConfigs(),
		Features:        DefaultFeatureFlags(),
	}
}

// ConfigBuilder provides a fluent interface for building rate limit configuration
type ConfigBuilder struct {
	config Config
}

// NewConfigBuilder creates a new configuration builder with defaults
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: DefaultConfig(),
	}
}

// WithProvider sets the rate limit provider
func (b *ConfigBuilder) WithProvider(provider string) *ConfigBuilder {
	b.config.Provider = provider
	return b
}

// WithAlgorithm sets the rate limiting algorithm
func (b *ConfigBuilder) WithAlgorithm(algorithm string) *ConfigBuilder {
	b.config.Algorithm = algorithm
	return b
}

// WithCleanupInterval sets the cleanup interval for memory provider
func (b *ConfigBuilder) WithCleanupInterval(interval string) *ConfigBuilder {
	b.config.CleanupInterval = interval
	return b
}

// WithRedisConfig sets Redis connection configuration
func (b *ConfigBuilder) WithRedisConfig(url, password string, db int, keyPrefix string) *ConfigBuilder {
	b.config.RedisURL = url
	b.config.RedisPassword = password
	b.config.RedisDB = db
	b.config.RedisKeyPrefix = keyPrefix
	return b
}

// WithDatabaseConfig sets database connection configuration
func (b *ConfigBuilder) WithDatabaseConfig(dsn, tableName string) *ConfigBuilder {
	b.config.DatabaseDSN = dsn
	b.config.TableName = tableName
	return b
}

// WithDefaultLimits sets the default rate limit configurations
func (b *ConfigBuilder) WithDefaultLimits(limits map[string]ratelimit.RateLimitConfig) *ConfigBuilder {
	b.config.DefaultLimits = limits
	return b
}

// WithDefaultLimit adds a single default rate limit configuration
func (b *ConfigBuilder) WithDefaultLimit(pattern string, limit ratelimit.RateLimitConfig) *ConfigBuilder {
	if b.config.DefaultLimits == nil {
		b.config.DefaultLimits = make(map[string]ratelimit.RateLimitConfig)
	}
	b.config.DefaultLimits[pattern] = limit
	return b
}

// WithGlobalLimit enables and configures global rate limiting
func (b *ConfigBuilder) WithGlobalLimit(limit ratelimit.RateLimitConfig) *ConfigBuilder {
	b.config.EnableGlobalLimits = true
	b.config.GlobalLimitConfig = limit
	return b
}

// WithFeatures sets the feature flags
func (b *ConfigBuilder) WithFeatures(features FeatureFlags) *ConfigBuilder {
	b.config.Features = features
	return b
}

// EnableDistributedLimits enables distributed rate limiting (requires Redis)
func (b *ConfigBuilder) EnableDistributedLimits() *ConfigBuilder {
	b.config.Provider = "redis"
	b.config.Features.EnableDistributedLimits = true
	b.config.Features.EnableRedisProvider = true
	b.config.Features.EnableMemoryProvider = false
	return b
}

// EnableTokenBucket switches to token bucket algorithm
func (b *ConfigBuilder) EnableTokenBucket() *ConfigBuilder {
	b.config.Algorithm = "token_bucket"
	b.config.Features.EnableTokenBucket = true
	b.config.Features.EnableSlidingWindow = false
	return b
}

// EnableLeakyBucket switches to leaky bucket algorithm
func (b *ConfigBuilder) EnableLeakyBucket() *ConfigBuilder {
	b.config.Algorithm = "leaky_bucket"
	b.config.Features.EnableLeakyBucket = true
	b.config.Features.EnableSlidingWindow = false
	return b
}

// EnableMetrics enables rate limiting metrics collection
func (b *ConfigBuilder) EnableMetrics() *ConfigBuilder {
	b.config.Features.EnableMetrics = true
	return b
}

// EnableDynamicLimits enables dynamic rate limit adjustments
func (b *ConfigBuilder) EnableDynamicLimits() *ConfigBuilder {
	b.config.Features.EnableDynamicLimits = true
	return b
}

// ForDevelopment configures the service for development use
func (b *ConfigBuilder) ForDevelopment() *ConfigBuilder {
	b.config.Provider = "memory"
	b.config.Algorithm = "sliding_window"
	b.config.Features.EnableMetrics = false
	b.config.Features.EnableDistributedLimits = false
	return b
}

// ForProduction configures the service for production use
func (b *ConfigBuilder) ForProduction() *ConfigBuilder {
	b.config.Features.EnableMetrics = true
	b.config.Features.EnableDynamicLimits = true
	b.config.Features.EnableGracePeriod = true
	return b
}

// Build returns the final configuration
func (b *ConfigBuilder) Build() Config {
	return b.config
}
