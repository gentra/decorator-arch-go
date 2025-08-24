package factory

import (
	"fmt"
	"time"

	"github.com/gentra/decorator-arch-go/internal/eventhandler"
)

// Config contains all configuration for building the event handler service
type Config struct {
	// Handler configuration
	HandlerType string // "sync", "async", "batch", "stream"

	// Processing configuration
	Concurrency    int
	BufferSize     int
	BatchSize      int
	BatchTimeout   time.Duration
	ProcessTimeout time.Duration

	// Retry configuration
	MaxRetries    int
	InitialDelay  time.Duration
	BackoffFactor float64
	MaxDelay      time.Duration

	// Event filtering
	EventTypes     []string
	EventPatterns  []string
	IgnorePatterns []string

	// Performance settings
	EnableMetrics   bool
	EnableTracing   bool
	MetricsInterval time.Duration

	// Storage configuration (for event sourcing)
	StorageProvider string // "memory", "file", "database"
	StoragePath     string

	// Feature flags
	Features FeatureFlags
}

// FeatureFlags controls event handler service behavior
type FeatureFlags struct {
	EnableSyncProcessing      bool
	EnableAsyncProcessing     bool
	EnableBatchProcessing     bool
	EnableStreamProcessing    bool
	EnableRetryLogic          bool
	EnableDeadLetterQueue     bool
	EnableEventFiltering      bool
	EnableEventTransformation bool
	EnableEventValidation     bool
	EnableEventSourcing       bool
	EnableMetrics             bool
	EnableTracing             bool
	EnableCircuitBreaker      bool
	EnableRateLimiting        bool
	EnableOrdering            bool
	EnableDeduplication       bool
}

// DefaultFeatureFlags returns default feature flag configuration
func DefaultFeatureFlags() FeatureFlags {
	return FeatureFlags{
		EnableSyncProcessing:      true,
		EnableAsyncProcessing:     false,
		EnableBatchProcessing:     false,
		EnableStreamProcessing:    false,
		EnableRetryLogic:          true,
		EnableDeadLetterQueue:     false,
		EnableEventFiltering:      true,
		EnableEventTransformation: false,
		EnableEventValidation:     true,
		EnableEventSourcing:       false,
		EnableMetrics:             false,
		EnableTracing:             false,
		EnableCircuitBreaker:      false,
		EnableRateLimiting:        false,
		EnableOrdering:            false,
		EnableDeduplication:       false,
	}
}

// EventHandlerServiceFactory creates and assembles the complete event handler service
type EventHandlerServiceFactory struct {
	config Config
}

// NewFactory creates a new event handler service factory with the given configuration
func NewFactory(config Config) *EventHandlerServiceFactory {
	return &EventHandlerServiceFactory{
		config: config,
	}
}

// Build assembles and returns the complete event handler service based on configuration
func (f *EventHandlerServiceFactory) Build() (eventhandler.Service, error) {
	switch f.config.HandlerType {
	case "sync":
		return f.buildSyncHandler()
	case "async":
		return f.buildAsyncHandler()
	case "batch":
		return f.buildBatchHandler()
	case "stream":
		return f.buildStreamHandler()
	default:
		// Default to sync handler
		return f.buildSyncHandler()
	}
}

// buildSyncHandler creates a synchronous event handler (placeholder)
func (f *EventHandlerServiceFactory) buildSyncHandler() (eventhandler.Service, error) {
	// TODO: Implement synchronous event handler
	return nil, fmt.Errorf("synchronous event handler not yet implemented")
}

// buildAsyncHandler creates an asynchronous event handler (placeholder)
func (f *EventHandlerServiceFactory) buildAsyncHandler() (eventhandler.Service, error) {
	// TODO: Implement asynchronous event handler
	return nil, fmt.Errorf("asynchronous event handler not yet implemented")
}

// buildBatchHandler creates a batch event handler (placeholder)
func (f *EventHandlerServiceFactory) buildBatchHandler() (eventhandler.Service, error) {
	// TODO: Implement batch event handler
	return nil, fmt.Errorf("batch event handler not yet implemented")
}

// buildStreamHandler creates a stream event handler (placeholder)
func (f *EventHandlerServiceFactory) buildStreamHandler() (eventhandler.Service, error) {
	// TODO: Implement stream event handler
	return nil, fmt.Errorf("stream event handler not yet implemented")
}

// DefaultConfig returns a sensible default configuration for the event handler service
func DefaultConfig() Config {
	return Config{
		HandlerType:     "sync",
		Concurrency:     1,
		BufferSize:      100,
		BatchSize:       10,
		BatchTimeout:    5 * time.Second,
		ProcessTimeout:  30 * time.Second,
		MaxRetries:      3,
		InitialDelay:    1 * time.Second,
		BackoffFactor:   2.0,
		MaxDelay:        5 * time.Minute,
		EventTypes:      []string{},
		EventPatterns:   []string{},
		IgnorePatterns:  []string{},
		EnableMetrics:   false,
		EnableTracing:   false,
		MetricsInterval: 1 * time.Minute,
		StorageProvider: "memory",
		Features:        DefaultFeatureFlags(),
	}
}

// ConfigBuilder provides a fluent interface for building event handler configuration
type ConfigBuilder struct {
	config Config
}

// NewConfigBuilder creates a new configuration builder with defaults
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: DefaultConfig(),
	}
}

// WithHandlerType sets the event handler type
func (b *ConfigBuilder) WithHandlerType(handlerType string) *ConfigBuilder {
	b.config.HandlerType = handlerType
	return b
}

// WithConcurrency sets the number of concurrent handlers
func (b *ConfigBuilder) WithConcurrency(concurrency int) *ConfigBuilder {
	b.config.Concurrency = concurrency
	return b
}

// WithBufferSize sets the event buffer size
func (b *ConfigBuilder) WithBufferSize(size int) *ConfigBuilder {
	b.config.BufferSize = size
	return b
}

// WithBatchConfig sets batch processing configuration
func (b *ConfigBuilder) WithBatchConfig(batchSize int, batchTimeout time.Duration) *ConfigBuilder {
	b.config.BatchSize = batchSize
	b.config.BatchTimeout = batchTimeout
	return b
}

// WithProcessTimeout sets the processing timeout
func (b *ConfigBuilder) WithProcessTimeout(timeout time.Duration) *ConfigBuilder {
	b.config.ProcessTimeout = timeout
	return b
}

// WithRetryConfig sets retry configuration
func (b *ConfigBuilder) WithRetryConfig(maxRetries int, initialDelay time.Duration, backoffFactor float64, maxDelay time.Duration) *ConfigBuilder {
	b.config.MaxRetries = maxRetries
	b.config.InitialDelay = initialDelay
	b.config.BackoffFactor = backoffFactor
	b.config.MaxDelay = maxDelay
	return b
}

// WithEventTypes sets the event types to handle
func (b *ConfigBuilder) WithEventTypes(eventTypes []string) *ConfigBuilder {
	b.config.EventTypes = eventTypes
	return b
}

// WithEventPatterns sets event patterns to match
func (b *ConfigBuilder) WithEventPatterns(patterns []string) *ConfigBuilder {
	b.config.EventPatterns = patterns
	return b
}

// WithIgnorePatterns sets event patterns to ignore
func (b *ConfigBuilder) WithIgnorePatterns(patterns []string) *ConfigBuilder {
	b.config.IgnorePatterns = patterns
	return b
}

// WithMetrics enables metrics collection
func (b *ConfigBuilder) WithMetrics(enable bool, interval time.Duration) *ConfigBuilder {
	b.config.EnableMetrics = enable
	b.config.MetricsInterval = interval
	b.config.Features.EnableMetrics = enable
	return b
}

// WithTracing enables distributed tracing
func (b *ConfigBuilder) WithTracing(enable bool) *ConfigBuilder {
	b.config.EnableTracing = enable
	b.config.Features.EnableTracing = enable
	return b
}

// WithStorageProvider sets the storage provider for event sourcing
func (b *ConfigBuilder) WithStorageProvider(provider string, path string) *ConfigBuilder {
	b.config.StorageProvider = provider
	b.config.StoragePath = path
	return b
}

// WithFeatures sets the feature flags
func (b *ConfigBuilder) WithFeatures(features FeatureFlags) *ConfigBuilder {
	b.config.Features = features
	return b
}

// EnableAsyncProcessing switches to asynchronous processing
func (b *ConfigBuilder) EnableAsyncProcessing() *ConfigBuilder {
	b.config.HandlerType = "async"
	b.config.Features.EnableAsyncProcessing = true
	b.config.Features.EnableSyncProcessing = false
	return b
}

// EnableBatchProcessing switches to batch processing
func (b *ConfigBuilder) EnableBatchProcessing() *ConfigBuilder {
	b.config.HandlerType = "batch"
	b.config.Features.EnableBatchProcessing = true
	b.config.Features.EnableSyncProcessing = false
	return b
}

// EnableStreamProcessing switches to stream processing
func (b *ConfigBuilder) EnableStreamProcessing() *ConfigBuilder {
	b.config.HandlerType = "stream"
	b.config.Features.EnableStreamProcessing = true
	b.config.Features.EnableSyncProcessing = false
	return b
}

// EnableRetryLogic enables retry logic with dead letter queue
func (b *ConfigBuilder) EnableRetryLogic() *ConfigBuilder {
	b.config.Features.EnableRetryLogic = true
	b.config.Features.EnableDeadLetterQueue = true
	return b
}

// EnableEventFiltering enables event filtering capabilities
func (b *ConfigBuilder) EnableEventFiltering() *ConfigBuilder {
	b.config.Features.EnableEventFiltering = true
	return b
}

// EnableEventSourcing enables event sourcing capabilities
func (b *ConfigBuilder) EnableEventSourcing() *ConfigBuilder {
	b.config.Features.EnableEventSourcing = true
	return b
}

// EnableCircuitBreaker enables circuit breaker pattern
func (b *ConfigBuilder) EnableCircuitBreaker() *ConfigBuilder {
	b.config.Features.EnableCircuitBreaker = true
	return b
}

// EnableOrdering enables event ordering guarantees
func (b *ConfigBuilder) EnableOrdering() *ConfigBuilder {
	b.config.Features.EnableOrdering = true
	return b
}

// EnableDeduplication enables event deduplication
func (b *ConfigBuilder) EnableDeduplication() *ConfigBuilder {
	b.config.Features.EnableDeduplication = true
	return b
}

// ForDevelopment configures the service for development use
func (b *ConfigBuilder) ForDevelopment() *ConfigBuilder {
	b.config.HandlerType = "sync"
	b.config.Concurrency = 1
	b.config.BufferSize = 10
	b.config.EnableMetrics = false
	b.config.EnableTracing = false
	b.config.Features.EnableMetrics = false
	b.config.Features.EnableTracing = false
	return b
}

// ForProduction configures the service for production use
func (b *ConfigBuilder) ForProduction() *ConfigBuilder {
	b.config.HandlerType = "async"
	b.config.Concurrency = 10
	b.config.BufferSize = 1000
	b.config.EnableMetrics = true
	b.config.EnableTracing = true
	b.config.MetricsInterval = 1 * time.Minute
	b.config.Features.EnableMetrics = true
	b.config.Features.EnableTracing = true
	b.config.Features.EnableRetryLogic = true
	b.config.Features.EnableCircuitBreaker = true
	return b
}

// Build returns the final configuration
func (b *ConfigBuilder) Build() Config {
	return b.config
}
