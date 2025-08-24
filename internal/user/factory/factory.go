package factory

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/gentra/decorator-arch-go/internal/audit"
	"github.com/gentra/decorator-arch-go/internal/encryption"
	"github.com/gentra/decorator-arch-go/internal/events"
	"github.com/gentra/decorator-arch-go/internal/notification"
	"github.com/gentra/decorator-arch-go/internal/ratelimit"
	"github.com/gentra/decorator-arch-go/internal/token"
	"github.com/gentra/decorator-arch-go/internal/user"
	userAudit "github.com/gentra/decorator-arch-go/internal/user/audit"
	userEncryption "github.com/gentra/decorator-arch-go/internal/user/encryption"
	userGorm "github.com/gentra/decorator-arch-go/internal/user/gorm"
	userRateLimit "github.com/gentra/decorator-arch-go/internal/user/ratelimit"
	userRedis "github.com/gentra/decorator-arch-go/internal/user/redis"
	"github.com/gentra/decorator-arch-go/internal/user/usecase"
	userValidation "github.com/gentra/decorator-arch-go/internal/user/validation"
	"github.com/gentra/decorator-arch-go/internal/validation"
)

// Config contains all configuration for building the user service
type Config struct {
	// Database configuration
	DB *gorm.DB

	// Redis configuration
	RedisClient *redis.Client
	CacheTTL    time.Duration

	// Domain services - these replace the old interfaces
	AuditService        audit.Service
	EncryptionService   encryption.Service
	RateLimitService    ratelimit.Service
	ValidationService   validation.Service
	NotificationService notification.Service
	TokenService        token.Service
	EventsService       events.Service

	// Feature flags
	Features FeatureFlags
}

// ===== FACTORY STRATEGY LOGIC =====
// The factory assembles decorator chains using domain services

// FeatureFlags controls which layers are enabled
type FeatureFlags struct {
	EnableCache      bool
	EnableAudit      bool
	EnableRateLimit  bool
	EnableEncryption bool
	EnableValidation bool
}

// DefaultFeatureFlags returns default feature flag configuration
func DefaultFeatureFlags() FeatureFlags {
	return FeatureFlags{
		EnableCache:      true,
		EnableAudit:      true,
		EnableRateLimit:  true,
		EnableEncryption: false, // Disabled by default for demo purposes
		EnableValidation: true,
	}
}

// UserServiceFactory creates and assembles the complete user service decorator chain
type UserServiceFactory struct {
	config Config
}

// NewUserServiceFactory creates a new factory with the given configuration
func NewUserServiceFactory(config Config) *UserServiceFactory {
	return &UserServiceFactory{
		config: config,
	}
}

// Build assembles and returns the complete user service with all enabled decorators
func (f *UserServiceFactory) Build() (user.Service, error) {
	// Start with the storage layer (GORM)
	service, err := f.buildStorageLayer()
	if err != nil {
		return nil, fmt.Errorf("failed to build storage layer: %w", err)
	}

	// Add cache layer if enabled
	if f.config.Features.EnableCache {
		service, err = f.addCacheLayer(service)
		if err != nil {
			return nil, fmt.Errorf("failed to add cache layer: %w", err)
		}
	}

	// Add audit layer if enabled
	if f.config.Features.EnableAudit {
		service = f.addAuditLayer(service)
	}

	// Add rate limiting layer if enabled
	if f.config.Features.EnableRateLimit {
		service, err = f.addRateLimitLayer(service)
		if err != nil {
			return nil, fmt.Errorf("failed to add rate limit layer: %w", err)
		}
	}

	// Add encryption layer if enabled
	if f.config.Features.EnableEncryption {
		service, err = f.addEncryptionLayer(service)
		if err != nil {
			return nil, fmt.Errorf("failed to add encryption layer: %w", err)
		}
	}

	// Add validation layer if enabled
	if f.config.Features.EnableValidation {
		service = f.addValidationLayer(service)
	}

	// Add usecase layer (business logic) - always enabled
	service = f.addUseCaseLayer(service)

	return service, nil
}

// BuildMinimal creates a minimal user service with only storage and usecase layers
func (f *UserServiceFactory) BuildMinimal() (user.Service, error) {
	// Start with storage layer
	service, err := f.buildStorageLayer()
	if err != nil {
		return nil, fmt.Errorf("failed to build storage layer: %w", err)
	}

	// Add usecase layer
	service = f.addUseCaseLayer(service)

	return service, nil
}

// BuildForTesting creates a service suitable for testing with configurable layers
func (f *UserServiceFactory) BuildForTesting(enabledLayers []string) (user.Service, error) {
	// Start with storage layer
	service, err := f.buildStorageLayer()
	if err != nil {
		return nil, fmt.Errorf("failed to build storage layer: %w", err)
	}

	// Add layers based on configuration
	for _, layer := range enabledLayers {
		switch layer {
		case "cache":
			service, err = f.addCacheLayer(service)
			if err != nil {
				return nil, fmt.Errorf("failed to add cache layer: %w", err)
			}
		case "audit":
			service = f.addAuditLayer(service)
		case "ratelimit":
			service, err = f.addRateLimitLayer(service)
			if err != nil {
				return nil, fmt.Errorf("failed to add rate limit layer: %w", err)
			}
		case "encryption":
			service, err = f.addEncryptionLayer(service)
			if err != nil {
				return nil, fmt.Errorf("failed to add encryption layer: %w", err)
			}
		case "validation":
			service = f.addValidationLayer(service)
		}
	}

	// Always add usecase layer last
	service = f.addUseCaseLayer(service)

	return service, nil
}

// Layer builders

func (f *UserServiceFactory) buildStorageLayer() (user.Service, error) {
	if f.config.DB == nil {
		return nil, fmt.Errorf("database connection is required")
	}

	return userGorm.NewService(f.config.DB), nil
}

func (f *UserServiceFactory) addCacheLayer(next user.Service) (user.Service, error) {
	if f.config.RedisClient == nil {
		return nil, fmt.Errorf("redis client is required for cache layer")
	}

	cacheTTL := f.config.CacheTTL
	if cacheTTL == 0 {
		cacheTTL = 5 * time.Minute // Default TTL
	}

	return userRedis.NewService(next, f.config.RedisClient, cacheTTL), nil
}

func (f *UserServiceFactory) addAuditLayer(next user.Service) user.Service {
	return userAudit.NewService(next, f.config.AuditService)
}

func (f *UserServiceFactory) addRateLimitLayer(next user.Service) (user.Service, error) {
	return userRateLimit.NewService(next, f.config.RateLimitService), nil
}

func (f *UserServiceFactory) addEncryptionLayer(next user.Service) (user.Service, error) {
	return userEncryption.NewService(next, f.config.EncryptionService), nil
}

func (f *UserServiceFactory) addValidationLayer(next user.Service) user.Service {
	return userValidation.NewService(next, f.config.ValidationService)
}

func (f *UserServiceFactory) addUseCaseLayer(next user.Service) user.Service {
	deps := usecase.Dependencies{
		NotificationService: f.config.NotificationService,
		TokenService:        f.config.TokenService,
		EventPublisher:      f.config.EventsService,
	}
	return usecase.NewService(next, deps)
}

// Helper methods for creating common configurations

// NewDefaultConfig creates a default configuration for the user service factory
func NewDefaultConfig(
	db *gorm.DB,
	redisClient *redis.Client,
	auditSvc audit.Service,
	encryptionSvc encryption.Service,
	rateLimitSvc ratelimit.Service,
	validationSvc validation.Service,
	notificationSvc notification.Service,
	tokenSvc token.Service,
	eventsSvc events.Service,
) Config {
	return Config{
		DB:                  db,
		RedisClient:         redisClient,
		CacheTTL:            5 * time.Minute,
		AuditService:        auditSvc,
		EncryptionService:   encryptionSvc,
		RateLimitService:    rateLimitSvc,
		ValidationService:   validationSvc,
		NotificationService: notificationSvc,
		TokenService:        tokenSvc,
		EventsService:       eventsSvc,
		Features:            DefaultFeatureFlags(),
	}
}

// NewProductionConfig creates a production-ready configuration
func NewProductionConfig(
	db *gorm.DB,
	redisClient *redis.Client,
	auditSvc audit.Service,
	encryptionSvc encryption.Service,
	rateLimitSvc ratelimit.Service,
	validationSvc validation.Service,
	notificationSvc notification.Service,
	tokenSvc token.Service,
	eventsSvc events.Service,
) Config {
	return Config{
		DB:                  db,
		RedisClient:         redisClient,
		CacheTTL:            10 * time.Minute,
		AuditService:        auditSvc,
		EncryptionService:   encryptionSvc,
		RateLimitService:    rateLimitSvc,
		ValidationService:   validationSvc,
		NotificationService: notificationSvc,
		TokenService:        tokenSvc,
		EventsService:       eventsSvc,
		Features: FeatureFlags{
			EnableCache:      true,
			EnableAudit:      true,
			EnableRateLimit:  true,
			EnableEncryption: true,
			EnableValidation: true,
		},
	}
}

// NewTestingConfig creates a configuration suitable for testing
func NewTestingConfig(
	db *gorm.DB,
	auditSvc audit.Service,
	encryptionSvc encryption.Service,
	rateLimitSvc ratelimit.Service,
	validationSvc validation.Service,
	notificationSvc notification.Service,
	tokenSvc token.Service,
	eventsSvc events.Service,
) Config {
	return Config{
		DB:                  db,
		CacheTTL:            time.Minute,
		AuditService:        auditSvc,
		EncryptionService:   encryptionSvc,
		RateLimitService:    rateLimitSvc,
		ValidationService:   validationSvc,
		NotificationService: notificationSvc,
		TokenService:        tokenSvc,
		EventsService:       eventsSvc,
		Features: FeatureFlags{
			EnableCache:      false, // Disable cache for consistent testing
			EnableAudit:      false, // Disable audit to reduce noise
			EnableRateLimit:  false, // Disable rate limiting for testing
			EnableEncryption: false, // Disable encryption for simpler testing
			EnableValidation: true,  // Keep validation for testing business rules
		},
	}
}

// This function is no longer needed since rate limiting configuration
// is now handled by the rate limiting domain service

// ServiceLayerInfo provides information about the layers in the service chain
type ServiceLayerInfo struct {
	Layers []LayerInfo `json:"layers"`
}

type LayerInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// GetServiceInfo returns information about the configured service layers
func (f *UserServiceFactory) GetServiceInfo() ServiceLayerInfo {
	layers := []LayerInfo{
		{
			Name:        "UseCase",
			Description: "Business logic and orchestration layer",
			Enabled:     true, // Always enabled
		},
		{
			Name:        "Validation",
			Description: "Input validation and business rules",
			Enabled:     f.config.Features.EnableValidation,
		},
		{
			Name:        "Encryption",
			Description: "Data encryption for sensitive fields",
			Enabled:     f.config.Features.EnableEncryption,
		},
		{
			Name:        "RateLimit",
			Description: "Rate limiting for API protection",
			Enabled:     f.config.Features.EnableRateLimit,
		},
		{
			Name:        "Audit",
			Description: "Activity logging and audit trail",
			Enabled:     f.config.Features.EnableAudit,
		},
		{
			Name:        "Cache",
			Description: "Redis caching for performance",
			Enabled:     f.config.Features.EnableCache,
		},
		{
			Name:        "Storage",
			Description: "Database persistence layer (GORM)",
			Enabled:     true, // Always enabled
		},
	}

	return ServiceLayerInfo{Layers: layers}
}
