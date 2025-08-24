package eventhandler

import (
	"context"
)

// Service defines the event handler domain interface - the ONLY interface in this domain
type Service interface {
	Handle(ctx context.Context, event interface{}) error
	GetHandledEventTypes() []string
}

// Domain types and data structures

// EventHandlerConfig contains configuration for event handlers
type EventHandlerConfig struct {
	HandlerID   string            `json:"handler_id"`
	EventTypes  []string          `json:"event_types"`
	Enabled     bool              `json:"enabled"`
	RetryConfig RetryConfig       `json:"retry_config"`
	Timeout     string            `json:"timeout"`
	BatchSize   int               `json:"batch_size"`
	Concurrency int               `json:"concurrency"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// RetryConfig contains retry configuration for failed event handling
type RetryConfig struct {
	MaxRetries    int     `json:"max_retries"`
	InitialDelay  string  `json:"initial_delay"`
	BackoffFactor float64 `json:"backoff_factor"`
	MaxDelay      string  `json:"max_delay"`
}

// EventHandlerError represents domain-specific event handler errors
type EventHandlerError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Handler string `json:"handler,omitempty"`
	Event   string `json:"event,omitempty"`
}

func (e EventHandlerError) Error() string {
	return e.Message
}

// Common event handler error codes
var (
	ErrHandlerNotFound  = EventHandlerError{Code: "HANDLER_NOT_FOUND", Message: "Event handler not found"}
	ErrHandlingFailed   = EventHandlerError{Code: "HANDLING_FAILED", Message: "Event handling failed"}
	ErrInvalidEventType = EventHandlerError{Code: "INVALID_EVENT_TYPE", Message: "Invalid event type for handler"}
	ErrHandlerDisabled  = EventHandlerError{Code: "HANDLER_DISABLED", Message: "Event handler is disabled"}
	ErrHandlerTimeout   = EventHandlerError{Code: "HANDLER_TIMEOUT", Message: "Event handler timed out"}
)

// Helper methods for EventHandlerConfig
func (c *EventHandlerConfig) IsValid() bool {
	return c.HandlerID != "" && len(c.EventTypes) > 0
}

func (c *EventHandlerConfig) IsEnabled() bool {
	return c.Enabled
}

func (c *EventHandlerConfig) HandlesEventType(eventType string) bool {
	for _, et := range c.EventTypes {
		if et == eventType {
			return true
		}
	}
	return false
}

// Helper methods for RetryConfig
func (r *RetryConfig) IsValid() bool {
	return r.MaxRetries >= 0 && r.InitialDelay != "" && r.BackoffFactor > 0
}

// Default event handler configuration
func DefaultEventHandlerConfig() EventHandlerConfig {
	return EventHandlerConfig{
		Enabled:     true,
		BatchSize:   1,
		Concurrency: 1,
		Timeout:     "30s",
		RetryConfig: RetryConfig{
			MaxRetries:    3,
			InitialDelay:  "1s",
			BackoffFactor: 2.0,
			MaxDelay:      "5m",
		},
	}
}
