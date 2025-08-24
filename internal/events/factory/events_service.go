package factory

import (
	"fmt"

	"github.com/gentra/decorator-arch-go/internal/events"
	"github.com/gentra/decorator-arch-go/internal/events/memory"
)

// Config contains all configuration for building the events service
type Config struct {
	// Provider configuration
	Provider string // "memory", "redis", "kafka", "nats", "rabbitmq"

	// Memory provider settings
	BufferSize int

	// Redis provider settings (for future implementation)
	RedisURL      string
	RedisPassword string
	RedisDB       int

	// Kafka provider settings (for future implementation)
	KafkaBrokers []string
	KafkaTopic   string

	// NATS provider settings (for future implementation)
	NATSServers []string
	NATSSubject string

	// Event processing configuration
	EventConfig events.EventConfig

	// Feature flags
	Features FeatureFlags
}

// FeatureFlags controls events service behavior
type FeatureFlags struct {
	EnableMemoryProvider   bool
	EnableRedisProvider    bool
	EnableKafkaProvider    bool
	EnableNATSProvider     bool
	EnableRabbitMQProvider bool
	EnablePersistence      bool
	EnableCompression      bool
	EnableRetryLogic       bool
	EnableDeadLetterQueue  bool
	EnableEventValidation  bool
	EnableMetrics          bool
	EnableTracing          bool
}

// DefaultFeatureFlags returns default feature flag configuration
func DefaultFeatureFlags() FeatureFlags {
	return FeatureFlags{
		EnableMemoryProvider:   true,
		EnableRedisProvider:    false,
		EnableKafkaProvider:    false,
		EnableNATSProvider:     false,
		EnableRabbitMQProvider: false,
		EnablePersistence:      false,
		EnableCompression:      false,
		EnableRetryLogic:       true,
		EnableDeadLetterQueue:  false,
		EnableEventValidation:  true,
		EnableMetrics:          false,
		EnableTracing:          false,
	}
}

// EventsServiceFactory creates and assembles the complete events service
type EventsServiceFactory struct {
	config Config
}

// NewFactory creates a new events service factory with the given configuration
func NewFactory(config Config) *EventsServiceFactory {
	return &EventsServiceFactory{
		config: config,
	}
}

// Build assembles and returns the complete events service based on configuration
func (f *EventsServiceFactory) Build() (events.Service, error) {
	switch f.config.Provider {
	case "memory":
		return f.buildMemoryService()
	case "redis":
		return f.buildRedisService()
	case "kafka":
		return f.buildKafkaService()
	case "nats":
		return f.buildNATSService()
	case "rabbitmq":
		return f.buildRabbitMQService()
	default:
		// Default to memory provider
		return f.buildMemoryService()
	}
}

// buildMemoryService creates an in-memory events service
func (f *EventsServiceFactory) buildMemoryService() (events.Service, error) {
	eventConfig := f.config.EventConfig

	// Apply feature flags to event config
	if f.config.Features.EnablePersistence {
		eventConfig.Persistence = true
	}
	if f.config.Features.EnableCompression {
		eventConfig.Compression = true
	}

	// Set buffer size if specified
	if f.config.BufferSize > 0 {
		eventConfig.BufferSize = f.config.BufferSize
	}

	return memory.NewService(eventConfig), nil
}

// buildRedisService creates a Redis-based events service (placeholder)
func (f *EventsServiceFactory) buildRedisService() (events.Service, error) {
	// TODO: Implement Redis events service
	return nil, fmt.Errorf("Redis events provider not yet implemented")
}

// buildKafkaService creates a Kafka-based events service (placeholder)
func (f *EventsServiceFactory) buildKafkaService() (events.Service, error) {
	// TODO: Implement Kafka events service
	return nil, fmt.Errorf("Kafka events provider not yet implemented")
}

// buildNATSService creates a NATS-based events service (placeholder)
func (f *EventsServiceFactory) buildNATSService() (events.Service, error) {
	// TODO: Implement NATS events service
	return nil, fmt.Errorf("NATS events provider not yet implemented")
}

// buildRabbitMQService creates a RabbitMQ-based events service (placeholder)
func (f *EventsServiceFactory) buildRabbitMQService() (events.Service, error) {
	// TODO: Implement RabbitMQ events service
	return nil, fmt.Errorf("RabbitMQ events provider not yet implemented")
}

// DefaultConfig returns a sensible default configuration for the events service
func DefaultConfig() Config {
	return Config{
		Provider:    "memory",
		BufferSize:  1000,
		EventConfig: events.DefaultEventConfig(),
		Features:    DefaultFeatureFlags(),
	}
}

// ConfigBuilder provides a fluent interface for building events configuration
type ConfigBuilder struct {
	config Config
}

// NewConfigBuilder creates a new configuration builder with defaults
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: DefaultConfig(),
	}
}

// WithProvider sets the events provider
func (b *ConfigBuilder) WithProvider(provider string) *ConfigBuilder {
	b.config.Provider = provider
	return b
}

// WithBufferSize sets the buffer size for in-memory processing
func (b *ConfigBuilder) WithBufferSize(size int) *ConfigBuilder {
	b.config.BufferSize = size
	return b
}

// WithRedisConfig sets Redis connection configuration
func (b *ConfigBuilder) WithRedisConfig(url, password string, db int) *ConfigBuilder {
	b.config.RedisURL = url
	b.config.RedisPassword = password
	b.config.RedisDB = db
	return b
}

// WithKafkaConfig sets Kafka connection configuration
func (b *ConfigBuilder) WithKafkaConfig(brokers []string, topic string) *ConfigBuilder {
	b.config.KafkaBrokers = brokers
	b.config.KafkaTopic = topic
	return b
}

// WithNATSConfig sets NATS connection configuration
func (b *ConfigBuilder) WithNATSConfig(servers []string, subject string) *ConfigBuilder {
	b.config.NATSServers = servers
	b.config.NATSSubject = subject
	return b
}

// WithEventConfig sets the event configuration
func (b *ConfigBuilder) WithEventConfig(eventConfig events.EventConfig) *ConfigBuilder {
	b.config.EventConfig = eventConfig
	return b
}

// WithFeatures sets the feature flags
func (b *ConfigBuilder) WithFeatures(features FeatureFlags) *ConfigBuilder {
	b.config.Features = features
	return b
}

// EnablePersistence enables event persistence
func (b *ConfigBuilder) EnablePersistence() *ConfigBuilder {
	b.config.Features.EnablePersistence = true
	return b
}

// EnableCompression enables event compression
func (b *ConfigBuilder) EnableCompression() *ConfigBuilder {
	b.config.Features.EnableCompression = true
	return b
}

// EnableRetryLogic enables retry logic for failed events
func (b *ConfigBuilder) EnableRetryLogic() *ConfigBuilder {
	b.config.Features.EnableRetryLogic = true
	return b
}

// EnableMetrics enables metrics collection
func (b *ConfigBuilder) EnableMetrics() *ConfigBuilder {
	b.config.Features.EnableMetrics = true
	return b
}

// EnableTracing enables distributed tracing
func (b *ConfigBuilder) EnableTracing() *ConfigBuilder {
	b.config.Features.EnableTracing = true
	return b
}

// ForDevelopment configures the service for development use
func (b *ConfigBuilder) ForDevelopment() *ConfigBuilder {
	b.config.Provider = "memory"
	b.config.BufferSize = 100
	b.config.Features.EnablePersistence = false
	b.config.Features.EnableMetrics = false
	b.config.Features.EnableTracing = false
	return b
}

// ForProduction configures the service for production use
func (b *ConfigBuilder) ForProduction() *ConfigBuilder {
	b.config.BufferSize = 10000
	b.config.Features.EnablePersistence = true
	b.config.Features.EnableCompression = true
	b.config.Features.EnableRetryLogic = true
	b.config.Features.EnableMetrics = true
	b.config.Features.EnableTracing = true
	return b
}

// Build returns the final configuration
func (b *ConfigBuilder) Build() Config {
	return b.config
}
