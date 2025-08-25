package events_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/events"
)

func TestEvent_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		event    events.Event
		expected bool
	}{
		{
			name: "Given event with ID, type, aggregate ID and timestamp, When IsValid is called, Then should return true",
			event: events.Event{
				ID:          "event-123",
				Type:        "user.created",
				AggregateID: "user-456",
				Timestamp:   time.Now(),
			},
			expected: true,
		},
		{
			name: "Given event with empty ID, When IsValid is called, Then should return false",
			event: events.Event{
				ID:          "",
				Type:        "user.created",
				AggregateID: "user-456",
				Timestamp:   time.Now(),
			},
			expected: false,
		},
		{
			name: "Given event with empty type, When IsValid is called, Then should return false",
			event: events.Event{
				ID:          "event-123",
				Type:        "",
				AggregateID: "user-456",
				Timestamp:   time.Now(),
			},
			expected: false,
		},
		{
			name: "Given event with empty aggregate ID, When IsValid is called, Then should return false",
			event: events.Event{
				ID:          "event-123",
				Type:        "user.created",
				AggregateID: "",
				Timestamp:   time.Now(),
			},
			expected: false,
		},
		{
			name: "Given event with zero timestamp, When IsValid is called, Then should return false",
			event: events.Event{
				ID:          "event-123",
				Type:        "user.created",
				AggregateID: "user-456",
				Timestamp:   time.Time{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.event.IsValid()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEvent_WithMetadata(t *testing.T) {
	t.Run("Given event and metadata, When WithMetadata is called, Then should set metadata and return event", func(t *testing.T) {
		// Arrange
		event := events.Event{
			ID:          "event-123",
			Type:        "user.created",
			AggregateID: "user-456",
		}
		metadata := events.EventMetadata{
			UserID:        "user-789",
			CorrelationID: "corr-123",
			Source:        "api",
		}

		// Act
		result := event.WithMetadata(metadata)

		// Assert
		assert.Equal(t, metadata, result.Metadata)
		assert.Equal(t, &event, result) // Should return pointer to same event
	})
}

func TestEvent_WithUserContext(t *testing.T) {
	t.Run("Given event, user ID and correlation ID, When WithUserContext is called, Then should set user context and return event", func(t *testing.T) {
		// Arrange
		event := events.Event{
			ID:          "event-123",
			Type:        "user.created",
			AggregateID: "user-456",
		}
		userID := "user-789"
		correlationID := "corr-123"

		// Act
		result := event.WithUserContext(userID, correlationID)

		// Assert
		assert.Equal(t, userID, result.Metadata.UserID)
		assert.Equal(t, correlationID, result.Metadata.CorrelationID)
		assert.Equal(t, &event, result) // Should return pointer to same event
	})
}

func TestEventFilters_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		filters  events.EventFilters
		expected bool
	}{
		{
			name: "Given event filters with event types, When IsValid is called, Then should return true",
			filters: events.EventFilters{
				EventTypes: []string{"user.created", "user.updated"},
			},
			expected: true,
		},
		{
			name: "Given event filters with aggregate ID, When IsValid is called, Then should return true",
			filters: events.EventFilters{
				AggregateID: "user-123",
			},
			expected: true,
		},
		{
			name: "Given event filters with aggregate types, When IsValid is called, Then should return true",
			filters: events.EventFilters{
				AggregateTypes: []string{"user", "order"},
			},
			expected: true,
		},
		{
			name: "Given event filters with no criteria, When IsValid is called, Then should return false",
			filters: events.EventFilters{},
			expected: false,
		},
		{
			name: "Given event filters with empty event types, aggregate ID and aggregate types, When IsValid is called, Then should return false",
			filters: events.EventFilters{
				EventTypes:     []string{},
				AggregateID:    "",
				AggregateTypes: []string{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.filters.IsValid()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEventFilters_WithTimeRange(t *testing.T) {
	t.Run("Given event filters and time range, When WithTimeRange is called, Then should set time range and return filters", func(t *testing.T) {
		// Arrange
		filters := events.EventFilters{
			EventTypes: []string{"user.created"},
		}
		startTime := time.Now().Add(-time.Hour)
		endTime := time.Now()

		// Act
		result := filters.WithTimeRange(startTime, endTime)

		// Assert
		assert.NotNil(t, result.StartTime)
		assert.NotNil(t, result.EndTime)
		assert.Equal(t, startTime, *result.StartTime)
		assert.Equal(t, endTime, *result.EndTime)
		assert.Equal(t, &filters, result) // Should return pointer to same filters
	})
}

func TestEventFilters_WithPagination(t *testing.T) {
	t.Run("Given event filters and pagination params, When WithPagination is called, Then should set pagination and return filters", func(t *testing.T) {
		// Arrange
		filters := events.EventFilters{
			EventTypes: []string{"user.created"},
		}
		limit := 50
		offset := 100

		// Act
		result := filters.WithPagination(limit, offset)

		// Assert
		assert.Equal(t, limit, result.Limit)
		assert.Equal(t, offset, result.Offset)
		assert.Equal(t, &filters, result) // Should return pointer to same filters
	})
}

func TestEventConfig_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		config   events.EventConfig
		expected bool
	}{
		{
			name: "Given event config with provider and buffer size, When IsValid is called, Then should return true",
			config: events.EventConfig{
				Provider:   "inmemory",
				BufferSize: 1000,
			},
			expected: true,
		},
		{
			name: "Given event config with empty provider, When IsValid is called, Then should return false",
			config: events.EventConfig{
				Provider:   "",
				BufferSize: 1000,
			},
			expected: false,
		},
		{
			name: "Given event config with zero buffer size, When IsValid is called, Then should return false",
			config: events.EventConfig{
				Provider:   "inmemory",
				BufferSize: 0,
			},
			expected: false,
		},
		{
			name: "Given event config with negative buffer size, When IsValid is called, Then should return false",
			config: events.EventConfig{
				Provider:   "inmemory",
				BufferSize: -1,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.config.IsValid()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultEventConfig(t *testing.T) {
	t.Run("Given default event config call, When DefaultEventConfig is called, Then should return valid default configuration", func(t *testing.T) {
		// Act
		config := events.DefaultEventConfig()

		// Assert
		assert.Equal(t, "inmemory", config.Provider)
		assert.Equal(t, 1000, config.BufferSize)
		assert.Equal(t, "json", config.Serialization)
		assert.False(t, config.Compression)
		assert.False(t, config.Persistence)
		
		// Check retry config
		assert.Equal(t, 3, config.RetryConfig.MaxRetries)
		assert.Equal(t, time.Second, config.RetryConfig.InitialDelay)
		assert.Equal(t, 2.0, config.RetryConfig.BackoffFactor)
		assert.Equal(t, time.Minute*5, config.RetryConfig.MaxDelay)
		
		// Check topics
		assert.NotNil(t, config.Topics)
		expectedTopics := map[string]string{
			"user.events":     "user-domain-events",
			"auth.events":     "auth-domain-events",
			"payment.events":  "payment-domain-events",
			"document.events": "document-domain-events",
		}
		
		for key, expectedValue := range expectedTopics {
			value, exists := config.Topics[key]
			assert.True(t, exists, "Topic %s should exist", key)
			assert.Equal(t, expectedValue, value, "Topic %s should have correct value", key)
		}
		
		// Validate the config
		assert.True(t, config.IsValid())
	})
}

func TestEventError_Error(t *testing.T) {
	tests := []struct {
		name     string
		eventErr events.EventError
		expected string
	}{
		{
			name: "Given event error with message, When Error is called, Then should return message",
			eventErr: events.EventError{
				Code:    "TEST_ERROR",
				Message: "Test event error",
			},
			expected: "Test event error",
		},
		{
			name: "Given event error with empty message, When Error is called, Then should return empty string",
			eventErr: events.EventError{
				Code:    "TEST_ERROR",
				Message: "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.eventErr.Error()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEventErrors_Constants(t *testing.T) {
	tests := []struct {
		name         string
		err          events.EventError
		expectedCode string
	}{
		{
			name:         "Given ErrEventNotFound, When accessing code, Then should have correct code",
			err:          events.ErrEventNotFound,
			expectedCode: "EVENT_NOT_FOUND",
		},
		{
			name:         "Given ErrInvalidEvent, When accessing code, Then should have correct code",
			err:          events.ErrInvalidEvent,
			expectedCode: "INVALID_EVENT",
		},
		{
			name:         "Given ErrHandlerNotFound, When accessing code, Then should have correct code",
			err:          events.ErrHandlerNotFound,
			expectedCode: "HANDLER_NOT_FOUND",
		},
		{
			name:         "Given ErrPublishFailed, When accessing code, Then should have correct code",
			err:          events.ErrPublishFailed,
			expectedCode: "PUBLISH_FAILED",
		},
		{
			name:         "Given ErrSubscriptionFailed, When accessing code, Then should have correct code",
			err:          events.ErrSubscriptionFailed,
			expectedCode: "SUBSCRIPTION_FAILED",
		},
		{
			name:         "Given ErrVersionConflict, When accessing code, Then should have correct code",
			err:          events.ErrVersionConflict,
			expectedCode: "VERSION_CONFLICT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assert
			assert.Equal(t, tt.expectedCode, tt.err.Code)
			assert.NotEmpty(t, tt.err.Message)
		})
	}
}

func TestEventConstants(t *testing.T) {
	tests := []struct {
		name         string
		constant     string
		expectedStr  string
	}{
		{
			name:         "Given EventTypeUserRegistered constant, When accessing string value, Then should have correct value",
			constant:     events.EventTypeUserRegistered,
			expectedStr:  "user.registered",
		},
		{
			name:         "Given EventTypeUserUpdated constant, When accessing string value, Then should have correct value",
			constant:     events.EventTypeUserUpdated,
			expectedStr:  "user.updated",
		},
		{
			name:         "Given EventTypeUserDeleted constant, When accessing string value, Then should have correct value",
			constant:     events.EventTypeUserDeleted,
			expectedStr:  "user.deleted",
		},
		{
			name:         "Given EventTypeUserPrefsUpdated constant, When accessing string value, Then should have correct value",
			constant:     events.EventTypeUserPrefsUpdated,
			expectedStr:  "user.preferences.updated",
		},
		{
			name:         "Given EventTypeUserLoggedIn constant, When accessing string value, Then should have correct value",
			constant:     events.EventTypeUserLoggedIn,
			expectedStr:  "auth.user.logged_in",
		},
		{
			name:         "Given EventTypeUserLoggedOut constant, When accessing string value, Then should have correct value",
			constant:     events.EventTypeUserLoggedOut,
			expectedStr:  "auth.user.logged_out",
		},
		{
			name:         "Given EventTypePasswordChanged constant, When accessing string value, Then should have correct value",
			constant:     events.EventTypePasswordChanged,
			expectedStr:  "auth.password.changed",
		},
		{
			name:         "Given EventTypeTokenRefreshed constant, When accessing string value, Then should have correct value",
			constant:     events.EventTypeTokenRefreshed,
			expectedStr:  "auth.token.refreshed",
		},
		{
			name:         "Given EventTypeSystemStarted constant, When accessing string value, Then should have correct value",
			constant:     events.EventTypeSystemStarted,
			expectedStr:  "system.started",
		},
		{
			name:         "Given EventTypeSystemStopped constant, When accessing string value, Then should have correct value",
			constant:     events.EventTypeSystemStopped,
			expectedStr:  "system.stopped",
		},
		{
			name:         "Given EventTypeErrorOccurred constant, When accessing string value, Then should have correct value",
			constant:     events.EventTypeErrorOccurred,
			expectedStr:  "system.error.occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assert
			assert.Equal(t, tt.expectedStr, tt.constant)
		})
	}
}

func TestEventSubscription_Structure(t *testing.T) {
	t.Run("Given event subscription with all fields, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		createdAt := time.Now()
		subscription := events.EventSubscription{
			ID:        "sub-123",
			Topics:    []string{"user.events", "auth.events"},
			CreatedAt: createdAt,
			Active:    true,
		}

		// Assert
		assert.Equal(t, "sub-123", subscription.ID)
		assert.Equal(t, []string{"user.events", "auth.events"}, subscription.Topics)
		assert.Equal(t, createdAt, subscription.CreatedAt)
		assert.True(t, subscription.Active)
	})
}

func TestEventMetadata_Structure(t *testing.T) {
	t.Run("Given event metadata with all fields, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		headers := map[string]string{
			"Content-Type": "application/json",
			"X-Source":     "api",
		}
		metadata := events.EventMetadata{
			UserID:        "user-123",
			CorrelationID: "corr-456",
			CausationID:   "cause-789",
			Source:        "api-gateway",
			Headers:       headers,
			IPAddress:     "192.168.1.1",
			UserAgent:     "Mozilla/5.0",
		}

		// Assert
		assert.Equal(t, "user-123", metadata.UserID)
		assert.Equal(t, "corr-456", metadata.CorrelationID)
		assert.Equal(t, "cause-789", metadata.CausationID)
		assert.Equal(t, "api-gateway", metadata.Source)
		assert.Equal(t, headers, metadata.Headers)
		assert.Equal(t, "192.168.1.1", metadata.IPAddress)
		assert.Equal(t, "Mozilla/5.0", metadata.UserAgent)
	})
}

func TestEvent_CompleteStructure(t *testing.T) {
	t.Run("Given event with all fields, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		timestamp := time.Now()
		data := map[string]interface{}{
			"user_id": "123",
			"email":   "test@example.com",
		}
		metadata := events.EventMetadata{
			UserID:        "user-456",
			CorrelationID: "corr-789",
		}

		event := events.Event{
			ID:            "event-123",
			Type:          "user.created",
			AggregateID:   "user-456",
			AggregateType: "user",
			Version:       1,
			Data:          data,
			Metadata:      metadata,
			Timestamp:     timestamp,
		}

		// Assert
		assert.Equal(t, "event-123", event.ID)
		assert.Equal(t, "user.created", event.Type)
		assert.Equal(t, "user-456", event.AggregateID)
		assert.Equal(t, "user", event.AggregateType)
		assert.Equal(t, 1, event.Version)
		assert.Equal(t, data, event.Data)
		assert.Equal(t, metadata, event.Metadata)
		assert.Equal(t, timestamp, event.Timestamp)
		assert.True(t, event.IsValid())
	})
}