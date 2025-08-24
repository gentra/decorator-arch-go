package notification

import (
	"context"
	"time"
)

// Service defines the notification domain interface - the ONLY interface in this domain
type Service interface {
	// Email notifications
	SendWelcomeEmail(ctx context.Context, userEmail, userName string) error
	SendPasswordResetEmail(ctx context.Context, userEmail, resetToken string) error
	SendProfileUpdateNotification(ctx context.Context, userID string, changes map[string]interface{}) error
	SendVerificationEmail(ctx context.Context, userEmail, verificationToken string) error
	
	// Push notifications
	SendPushNotification(ctx context.Context, userID string, notification PushNotification) error
	
	// SMS notifications
	SendSMSNotification(ctx context.Context, phoneNumber string, message string) error
	
	// Bulk notifications
	SendBulkEmail(ctx context.Context, emails []EmailNotification) error
	SendBulkPush(ctx context.Context, notifications []PushNotification) error
	
	// Notification management
	GetNotificationHistory(ctx context.Context, userID string, limit int) ([]NotificationHistory, error)
	MarkAsRead(ctx context.Context, notificationID string) error
	GetUnreadCount(ctx context.Context, userID string) (int, error)
}

// Domain types and data structures

// EmailNotification represents an email notification
type EmailNotification struct {
	ID          string                 `json:"id"`
	To          string                 `json:"to"`
	From        string                 `json:"from,omitempty"`
	Subject     string                 `json:"subject"`
	Body        string                 `json:"body"`
	BodyHTML    string                 `json:"body_html,omitempty"`
	Template    string                 `json:"template,omitempty"`
	Variables   map[string]interface{} `json:"variables,omitempty"`
	Priority    Priority               `json:"priority"`
	ScheduledAt *time.Time             `json:"scheduled_at,omitempty"`
	Attachments []Attachment           `json:"attachments,omitempty"`
}

// PushNotification represents a push notification
type PushNotification struct {
	ID       string                 `json:"id"`
	UserID   string                 `json:"user_id"`
	Title    string                 `json:"title"`
	Body     string                 `json:"body"`
	Data     map[string]interface{} `json:"data,omitempty"`
	Badge    int                    `json:"badge,omitempty"`
	Sound    string                 `json:"sound,omitempty"`
	Category string                 `json:"category,omitempty"`
	Priority Priority               `json:"priority"`
}

// SMSNotification represents an SMS notification
type SMSNotification struct {
	ID          string    `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	Message     string    `json:"message"`
	Priority    Priority  `json:"priority"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
}

// NotificationHistory represents a notification in history
type NotificationHistory struct {
	ID           string                 `json:"id"`
	UserID       string                 `json:"user_id"`
	Type         NotificationType       `json:"type"`
	Title        string                 `json:"title"`
	Body         string                 `json:"body"`
	Data         map[string]interface{} `json:"data,omitempty"`
	Status       NotificationStatus     `json:"status"`
	Priority     Priority               `json:"priority"`
	CreatedAt    time.Time              `json:"created_at"`
	SentAt       *time.Time             `json:"sent_at,omitempty"`
	ReadAt       *time.Time             `json:"read_at,omitempty"`
	FailureCount int                    `json:"failure_count"`
	LastError    string                 `json:"last_error,omitempty"`
}

// Attachment represents a file attachment
type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Content     []byte `json:"content"`
	Size        int64  `json:"size"`
}

// NotificationType enum
type NotificationType string

const (
	NotificationTypeEmail NotificationType = "email"
	NotificationTypePush  NotificationType = "push"
	NotificationTypeSMS   NotificationType = "sms"
	NotificationTypeInApp NotificationType = "in_app"
)

// NotificationStatus enum
type NotificationStatus string

const (
	NotificationStatusPending   NotificationStatus = "pending"
	NotificationStatusSent      NotificationStatus = "sent"
	NotificationStatusDelivered NotificationStatus = "delivered"
	NotificationStatusFailed    NotificationStatus = "failed"
	NotificationStatusRead      NotificationStatus = "read"
)

// Priority enum
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityNormal Priority = "normal"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

// NotificationConfig contains configuration for the notification service
type NotificationConfig struct {
	EmailProvider    string                 `json:"email_provider"`    // smtp, sendgrid, ses, etc.
	PushProvider     string                 `json:"push_provider"`     // firebase, apns, etc.
	SMSProvider      string                 `json:"sms_provider"`      // twilio, aws sns, etc.
	DefaultFromEmail string                 `json:"default_from_email"`
	Templates        map[string]string      `json:"templates"`
	RateLimits       map[string]RateLimit   `json:"rate_limits"`
	RetryConfig      RetryConfig            `json:"retry_config"`
}

// RateLimit contains rate limiting configuration for notifications
type RateLimit struct {
	MaxPerMinute int `json:"max_per_minute"`
	MaxPerHour   int `json:"max_per_hour"`
	MaxPerDay    int `json:"max_per_day"`
}

// RetryConfig contains retry configuration
type RetryConfig struct {
	MaxRetries    int           `json:"max_retries"`
	InitialDelay  time.Duration `json:"initial_delay"`
	BackoffFactor float64       `json:"backoff_factor"`
	MaxDelay      time.Duration `json:"max_delay"`
}

// Helper methods for EmailNotification
func (e *EmailNotification) IsValid() bool {
	return e.To != "" && e.Subject != "" && (e.Body != "" || e.BodyHTML != "")
}

func (e *EmailNotification) HasAttachments() bool {
	return len(e.Attachments) > 0
}

func (e *EmailNotification) IsScheduled() bool {
	return e.ScheduledAt != nil && e.ScheduledAt.After(time.Now())
}

// Helper methods for PushNotification
func (p *PushNotification) IsValid() bool {
	return p.UserID != "" && p.Title != ""
}

func (p *PushNotification) HasData() bool {
	return len(p.Data) > 0
}

// Helper methods for NotificationHistory
func (n *NotificationHistory) IsRead() bool {
	return n.ReadAt != nil
}

func (n *NotificationHistory) IsSent() bool {
	return n.SentAt != nil
}

func (n *NotificationHistory) HasFailed() bool {
	return n.Status == NotificationStatusFailed
}

// Helper methods for NotificationConfig
func (c *NotificationConfig) IsValid() bool {
	return c.DefaultFromEmail != ""
}

// Default notification configuration
func DefaultNotificationConfig() NotificationConfig {
	return NotificationConfig{
		EmailProvider:    "smtp",
		PushProvider:     "firebase",
		SMSProvider:      "twilio",
		DefaultFromEmail: "noreply@example.com",
		Templates:        make(map[string]string),
		RateLimits: map[string]RateLimit{
			"email": {MaxPerMinute: 60, MaxPerHour: 1000, MaxPerDay: 10000},
			"push":  {MaxPerMinute: 100, MaxPerHour: 5000, MaxPerDay: 50000},
			"sms":   {MaxPerMinute: 10, MaxPerHour: 100, MaxPerDay: 500},
		},
		RetryConfig: RetryConfig{
			MaxRetries:    3,
			InitialDelay:  time.Second,
			BackoffFactor: 2.0,
			MaxDelay:      time.Minute * 5,
		},
	}
}
