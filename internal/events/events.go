package events

import (
	"context"
	"time"

	"github.com/gentra/decorator-arch-go/internal/eventhandler"
)

// Service defines the events domain interface - the ONLY interface in this domain
type Service interface {
	// Event publishing
	Publish(ctx context.Context, event Event) error
	PublishBatch(ctx context.Context, events []Event) error

	// Event subscription and consumption
	Subscribe(ctx context.Context, topics []string, handler eventhandler.Service) error
	Unsubscribe(ctx context.Context, subscriptionID string) error

	// Event querying and replay
	GetEvents(ctx context.Context, filters EventFilters) ([]Event, error)
	GetEventsByAggregate(ctx context.Context, aggregateID string, limit int) ([]Event, error)
	ReplayEvents(ctx context.Context, aggregateID string, fromVersion int, handler eventhandler.Service) error
}

// Domain types and data structures

// Event represents a domain event
type Event struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	AggregateID   string                 `json:"aggregate_id"`
	AggregateType string                 `json:"aggregate_type"`
	Version       int                    `json:"version"`
	Data          map[string]interface{} `json:"data"`
	Metadata      EventMetadata          `json:"metadata"`
	Timestamp     time.Time              `json:"timestamp"`
}

// EventMetadata contains metadata about an event
type EventMetadata struct {
	UserID        string            `json:"user_id,omitempty"`
	CorrelationID string            `json:"correlation_id,omitempty"`
	CausationID   string            `json:"causation_id,omitempty"`
	Source        string            `json:"source,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
	IPAddress     string            `json:"ip_address,omitempty"`
	UserAgent     string            `json:"user_agent,omitempty"`
}

// EventFilters for querying events
type EventFilters struct {
	EventTypes     []string   `json:"event_types,omitempty"`
	AggregateID    string     `json:"aggregate_id,omitempty"`
	AggregateTypes []string   `json:"aggregate_types,omitempty"`
	StartTime      *time.Time `json:"start_time,omitempty"`
	EndTime        *time.Time `json:"end_time,omitempty"`
	UserID         string     `json:"user_id,omitempty"`
	CorrelationID  string     `json:"correlation_id,omitempty"`
	Limit          int        `json:"limit,omitempty"`
	Offset         int        `json:"offset,omitempty"`
}

// EventSubscription represents an event subscription
type EventSubscription struct {
	ID        string               `json:"id"`
	Topics    []string             `json:"topics"`
	Handler   eventhandler.Service `json:"-"`
	CreatedAt time.Time            `json:"created_at"`
	Active    bool                 `json:"active"`
}

// EventConfig contains configuration for the event service
type EventConfig struct {
	Provider      string            `json:"provider"`      // inmemory, redis, kafka, etc.
	BufferSize    int               `json:"buffer_size"`   // Buffer size for async processing
	RetryConfig   RetryConfig       `json:"retry_config"`  // Retry configuration
	Serialization string            `json:"serialization"` // json, protobuf, etc.
	Compression   bool              `json:"compression"`   // Enable compression
	Persistence   bool              `json:"persistence"`   // Enable event persistence
	Topics        map[string]string `json:"topics"`        // Topic configuration
}

// RetryConfig contains retry configuration for failed events
type RetryConfig struct {
	MaxRetries    int           `json:"max_retries"`
	InitialDelay  time.Duration `json:"initial_delay"`
	BackoffFactor float64       `json:"backoff_factor"`
	MaxDelay      time.Duration `json:"max_delay"`
}

// EventError represents domain-specific event errors
type EventError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Event   *Event `json:"event,omitempty"`
}

func (e EventError) Error() string {
	return e.Message
}

// Common event error codes
var (
	ErrEventNotFound      = EventError{Code: "EVENT_NOT_FOUND", Message: "Event not found"}
	ErrInvalidEvent       = EventError{Code: "INVALID_EVENT", Message: "Invalid event data"}
	ErrHandlerNotFound    = EventError{Code: "HANDLER_NOT_FOUND", Message: "Event handler not found"}
	ErrPublishFailed      = EventError{Code: "PUBLISH_FAILED", Message: "Failed to publish event"}
	ErrSubscriptionFailed = EventError{Code: "SUBSCRIPTION_FAILED", Message: "Failed to create subscription"}
	ErrVersionConflict    = EventError{Code: "VERSION_CONFLICT", Message: "Event version conflict"}
)

// Helper methods for Event
func (e *Event) IsValid() bool {
	return e.ID != "" && e.Type != "" && e.AggregateID != "" && !e.Timestamp.IsZero()
}

func (e *Event) WithMetadata(metadata EventMetadata) *Event {
	e.Metadata = metadata
	return e
}

func (e *Event) WithUserContext(userID, correlationID string) *Event {
	e.Metadata.UserID = userID
	e.Metadata.CorrelationID = correlationID
	return e
}

// Helper methods for EventFilters
func (f *EventFilters) IsValid() bool {
	return len(f.EventTypes) > 0 || f.AggregateID != "" || len(f.AggregateTypes) > 0
}

func (f *EventFilters) WithTimeRange(start, end time.Time) *EventFilters {
	f.StartTime = &start
	f.EndTime = &end
	return f
}

func (f *EventFilters) WithPagination(limit, offset int) *EventFilters {
	f.Limit = limit
	f.Offset = offset
	return f
}

// Helper methods for EventConfig
func (c *EventConfig) IsValid() bool {
	return c.Provider != "" && c.BufferSize > 0
}

// Default event configuration
func DefaultEventConfig() EventConfig {
	return EventConfig{
		Provider:      "inmemory",
		BufferSize:    1000,
		Serialization: "json",
		Compression:   false,
		Persistence:   false,
		RetryConfig: RetryConfig{
			MaxRetries:    3,
			InitialDelay:  time.Second,
			BackoffFactor: 2.0,
			MaxDelay:      time.Minute * 5,
		},
		Topics: map[string]string{
			"user.events":     "user-domain-events",
			"auth.events":     "auth-domain-events",
			"payment.events":  "payment-domain-events",
			"document.events": "document-domain-events",
		},
	}
}

// Common event types for different domains
const (
	// User domain events
	EventTypeUserRegistered   = "user.registered"
	EventTypeUserUpdated      = "user.updated"
	EventTypeUserDeleted      = "user.deleted"
	EventTypeUserPrefsUpdated = "user.preferences.updated"

	// Auth domain events
	EventTypeUserLoggedIn    = "auth.user.logged_in"
	EventTypeUserLoggedOut   = "auth.user.logged_out"
	EventTypePasswordChanged = "auth.password.changed"
	EventTypeTokenRefreshed  = "auth.token.refreshed"

	// System events
	EventTypeSystemStarted = "system.started"
	EventTypeSystemStopped = "system.stopped"
	EventTypeErrorOccurred = "system.error.occurred"
)
