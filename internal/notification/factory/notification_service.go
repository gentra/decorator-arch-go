package factory

import (
	"fmt"

	"github.com/gentra/decorator-arch-go/internal/notification"
	"github.com/gentra/decorator-arch-go/internal/notification/mock"
)

// Config contains all configuration for building the notification service
type Config struct {
	// Provider configuration
	EmailProvider string // "mock", "smtp", "sendgrid", "ses", "mailgun"
	PushProvider  string // "mock", "firebase", "apns", "onesignal"
	SMSProvider   string // "mock", "twilio", "sns", "nexmo"

	// SMTP configuration (if EmailProvider = "smtp")
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string

	// SendGrid configuration (if EmailProvider = "sendgrid")
	SendGridAPIKey string

	// AWS SES configuration (if EmailProvider = "ses")
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string

	// Firebase configuration (if PushProvider = "firebase")
	FirebaseProjectID     string
	FirebaseCredentialKey []byte

	// Twilio configuration (if SMSProvider = "twilio")
	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioFromNumber string

	// General notification settings
	DefaultFromEmail  string
	DefaultFromName   string
	MaxRetries        int
	RetryDelaySeconds int

	// Template configuration
	TemplateDir string
	Templates   map[string]string

	// Feature flags
	Features FeatureFlags
}

// FeatureFlags controls notification service behavior
type FeatureFlags struct {
	EnableEmailNotifications bool
	EnablePushNotifications  bool
	EnableSMSNotifications   bool
	EnableMockProvider       bool
	EnableTemplateEngine     bool
	EnableRetryLogic         bool
	EnableBulkOperations     bool
	EnableRateLimiting       bool
	EnableNotificationQueue  bool
	EnableDeliveryTracking   bool
	EnableAnalytics          bool
}

// DefaultFeatureFlags returns default feature flag configuration
func DefaultFeatureFlags() FeatureFlags {
	return FeatureFlags{
		EnableEmailNotifications: true,
		EnablePushNotifications:  true,
		EnableSMSNotifications:   true,
		EnableMockProvider:       true,
		EnableTemplateEngine:     true,
		EnableRetryLogic:         true,
		EnableBulkOperations:     true,
		EnableRateLimiting:       true,
		EnableNotificationQueue:  false,
		EnableDeliveryTracking:   false,
		EnableAnalytics:          false,
	}
}

// NotificationServiceFactory creates and assembles the complete notification service
type NotificationServiceFactory struct {
	config Config
}

// NewFactory creates a new notification service factory with the given configuration
func NewFactory(config Config) *NotificationServiceFactory {
	return &NotificationServiceFactory{
		config: config,
	}
}

// Build assembles and returns the complete notification service based on configuration
func (f *NotificationServiceFactory) Build() (notification.Service, error) {
	// For now, we only have mock implementation
	// In the future, we can add strategy pattern here for different providers

	if f.config.Features.EnableMockProvider {
		return f.buildMockService()
	}

	// Check for specific provider implementations
	switch f.config.EmailProvider {
	case "mock":
		return f.buildMockService()
	case "smtp":
		return f.buildSMTPService()
	case "sendgrid":
		return f.buildSendGridService()
	case "ses":
		return f.buildSESService()
	default:
		// Default to mock service
		return f.buildMockService()
	}
}

// buildMockService creates a mock notification service for testing/development
func (f *NotificationServiceFactory) buildMockService() (notification.Service, error) {
	return mock.NewService(), nil
}

// buildSMTPService creates an SMTP-based notification service (placeholder)
func (f *NotificationServiceFactory) buildSMTPService() (notification.Service, error) {
	// TODO: Implement SMTP notification service
	return nil, fmt.Errorf("SMTP notification provider not yet implemented")
}

// buildSendGridService creates a SendGrid-based notification service (placeholder)
func (f *NotificationServiceFactory) buildSendGridService() (notification.Service, error) {
	// TODO: Implement SendGrid notification service
	return nil, fmt.Errorf("SendGrid notification provider not yet implemented")
}

// buildSESService creates an AWS SES-based notification service (placeholder)
func (f *NotificationServiceFactory) buildSESService() (notification.Service, error) {
	// TODO: Implement AWS SES notification service
	return nil, fmt.Errorf("AWS SES notification provider not yet implemented")
}

// DefaultConfig returns a sensible default configuration for the notification service
func DefaultConfig() Config {
	return Config{
		EmailProvider:     "mock",
		PushProvider:      "mock",
		SMSProvider:       "mock",
		DefaultFromEmail:  "noreply@example.com",
		DefaultFromName:   "Application",
		MaxRetries:        3,
		RetryDelaySeconds: 5,
		Templates:         make(map[string]string),
		Features:          DefaultFeatureFlags(),
	}
}

// ConfigBuilder provides a fluent interface for building notification configuration
type ConfigBuilder struct {
	config Config
}

// NewConfigBuilder creates a new configuration builder with defaults
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: DefaultConfig(),
	}
}

// WithEmailProvider sets the email provider
func (b *ConfigBuilder) WithEmailProvider(provider string) *ConfigBuilder {
	b.config.EmailProvider = provider
	return b
}

// WithPushProvider sets the push notification provider
func (b *ConfigBuilder) WithPushProvider(provider string) *ConfigBuilder {
	b.config.PushProvider = provider
	return b
}

// WithSMSProvider sets the SMS provider
func (b *ConfigBuilder) WithSMSProvider(provider string) *ConfigBuilder {
	b.config.SMSProvider = provider
	return b
}

// WithSMTPConfig sets SMTP configuration
func (b *ConfigBuilder) WithSMTPConfig(host string, port int, username, password string) *ConfigBuilder {
	b.config.SMTPHost = host
	b.config.SMTPPort = port
	b.config.SMTPUsername = username
	b.config.SMTPPassword = password
	return b
}

// WithSendGridAPIKey sets SendGrid API key
func (b *ConfigBuilder) WithSendGridAPIKey(apiKey string) *ConfigBuilder {
	b.config.SendGridAPIKey = apiKey
	return b
}

// WithAWSConfig sets AWS SES configuration
func (b *ConfigBuilder) WithAWSConfig(region, accessKeyID, secretAccessKey string) *ConfigBuilder {
	b.config.AWSRegion = region
	b.config.AWSAccessKeyID = accessKeyID
	b.config.AWSSecretAccessKey = secretAccessKey
	return b
}

// WithFirebaseConfig sets Firebase configuration
func (b *ConfigBuilder) WithFirebaseConfig(projectID string, credentialKey []byte) *ConfigBuilder {
	b.config.FirebaseProjectID = projectID
	b.config.FirebaseCredentialKey = credentialKey
	return b
}

// WithTwilioConfig sets Twilio configuration
func (b *ConfigBuilder) WithTwilioConfig(accountSID, authToken, fromNumber string) *ConfigBuilder {
	b.config.TwilioAccountSID = accountSID
	b.config.TwilioAuthToken = authToken
	b.config.TwilioFromNumber = fromNumber
	return b
}

// WithDefaultSender sets default sender information
func (b *ConfigBuilder) WithDefaultSender(email, name string) *ConfigBuilder {
	b.config.DefaultFromEmail = email
	b.config.DefaultFromName = name
	return b
}

// WithRetryConfig sets retry configuration
func (b *ConfigBuilder) WithRetryConfig(maxRetries, retryDelaySeconds int) *ConfigBuilder {
	b.config.MaxRetries = maxRetries
	b.config.RetryDelaySeconds = retryDelaySeconds
	return b
}

// WithTemplateDir sets the template directory
func (b *ConfigBuilder) WithTemplateDir(dir string) *ConfigBuilder {
	b.config.TemplateDir = dir
	return b
}

// WithTemplate adds a template mapping
func (b *ConfigBuilder) WithTemplate(name, path string) *ConfigBuilder {
	if b.config.Templates == nil {
		b.config.Templates = make(map[string]string)
	}
	b.config.Templates[name] = path
	return b
}

// WithFeatures sets the feature flags
func (b *ConfigBuilder) WithFeatures(features FeatureFlags) *ConfigBuilder {
	b.config.Features = features
	return b
}

// EnableAllChannels enables all notification channels
func (b *ConfigBuilder) EnableAllChannels() *ConfigBuilder {
	b.config.Features.EnableEmailNotifications = true
	b.config.Features.EnablePushNotifications = true
	b.config.Features.EnableSMSNotifications = true
	return b
}

// EnableTemplateEngine enables template processing
func (b *ConfigBuilder) EnableTemplateEngine() *ConfigBuilder {
	b.config.Features.EnableTemplateEngine = true
	return b
}

// EnableProductionFeatures enables production-ready features
func (b *ConfigBuilder) EnableProductionFeatures() *ConfigBuilder {
	b.config.Features.EnableRetryLogic = true
	b.config.Features.EnableRateLimiting = true
	b.config.Features.EnableNotificationQueue = true
	b.config.Features.EnableDeliveryTracking = true
	b.config.Features.EnableAnalytics = true
	return b
}

// ForDevelopment configures the service for development use
func (b *ConfigBuilder) ForDevelopment() *ConfigBuilder {
	b.config.EmailProvider = "mock"
	b.config.PushProvider = "mock"
	b.config.SMSProvider = "mock"
	b.config.Features.EnableMockProvider = true
	b.config.Features.EnableNotificationQueue = false
	b.config.Features.EnableAnalytics = false
	return b
}

// Build returns the final configuration
func (b *ConfigBuilder) Build() Config {
	return b.config
}
