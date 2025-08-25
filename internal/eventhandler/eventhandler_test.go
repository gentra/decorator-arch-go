package eventhandler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/eventhandler"
)

func TestEventHandlerConfig_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		config   eventhandler.EventHandlerConfig
		expected bool
	}{
		{
			name: "Given event handler config with handler ID and event types, When IsValid is called, Then should return true",
			config: eventhandler.EventHandlerConfig{
				HandlerID:  "user-handler",
				EventTypes: []string{"user.created", "user.updated"},
			},
			expected: true,
		},
		{
			name: "Given event handler config with empty handler ID, When IsValid is called, Then should return false",
			config: eventhandler.EventHandlerConfig{
				HandlerID:  "",
				EventTypes: []string{"user.created", "user.updated"},
			},
			expected: false,
		},
		{
			name: "Given event handler config with empty event types, When IsValid is called, Then should return false",
			config: eventhandler.EventHandlerConfig{
				HandlerID:  "user-handler",
				EventTypes: []string{},
			},
			expected: false,
		},
		{
			name: "Given event handler config with nil event types, When IsValid is called, Then should return false",
			config: eventhandler.EventHandlerConfig{
				HandlerID:  "user-handler",
				EventTypes: nil,
			},
			expected: false,
		},
		{
			name: "Given event handler config with both handler ID empty and nil event types, When IsValid is called, Then should return false",
			config: eventhandler.EventHandlerConfig{
				HandlerID:  "",
				EventTypes: nil,
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

func TestEventHandlerConfig_IsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		config   eventhandler.EventHandlerConfig
		expected bool
	}{
		{
			name: "Given event handler config with enabled true, When IsEnabled is called, Then should return true",
			config: eventhandler.EventHandlerConfig{
				Enabled: true,
			},
			expected: true,
		},
		{
			name: "Given event handler config with enabled false, When IsEnabled is called, Then should return false",
			config: eventhandler.EventHandlerConfig{
				Enabled: false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.config.IsEnabled()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEventHandlerConfig_HandlesEventType(t *testing.T) {
	tests := []struct {
		name      string
		config    eventhandler.EventHandlerConfig
		eventType string
		expected  bool
	}{
		{
			name: "Given event handler config with event type in list, When HandlesEventType is called, Then should return true",
			config: eventhandler.EventHandlerConfig{
				EventTypes: []string{"user.created", "user.updated", "user.deleted"},
			},
			eventType: "user.created",
			expected:  true,
		},
		{
			name: "Given event handler config with event type not in list, When HandlesEventType is called, Then should return false",
			config: eventhandler.EventHandlerConfig{
				EventTypes: []string{"user.created", "user.updated"},
			},
			eventType: "user.deleted",
			expected:  false,
		},
		{
			name: "Given event handler config with empty event types, When HandlesEventType is called, Then should return false",
			config: eventhandler.EventHandlerConfig{
				EventTypes: []string{},
			},
			eventType: "user.created",
			expected:  false,
		},
		{
			name: "Given event handler config with nil event types, When HandlesEventType is called, Then should return false",
			config: eventhandler.EventHandlerConfig{
				EventTypes: nil,
			},
			eventType: "user.created",
			expected:  false,
		},
		{
			name: "Given event handler config with exact match, When HandlesEventType is called, Then should return true",
			config: eventhandler.EventHandlerConfig{
				EventTypes: []string{"auth.user.logged_in"},
			},
			eventType: "auth.user.logged_in",
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.config.HandlesEventType(tt.eventType)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRetryConfig_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		config   eventhandler.RetryConfig
		expected bool
	}{
		{
			name: "Given retry config with valid values, When IsValid is called, Then should return true",
			config: eventhandler.RetryConfig{
				MaxRetries:    3,
				InitialDelay:  "1s",
				BackoffFactor: 2.0,
			},
			expected: true,
		},
		{
			name: "Given retry config with zero max retries, When IsValid is called, Then should return true",
			config: eventhandler.RetryConfig{
				MaxRetries:    0,
				InitialDelay:  "1s",
				BackoffFactor: 2.0,
			},
			expected: true,
		},
		{
			name: "Given retry config with negative max retries, When IsValid is called, Then should return false",
			config: eventhandler.RetryConfig{
				MaxRetries:    -1,
				InitialDelay:  "1s",
				BackoffFactor: 2.0,
			},
			expected: false,
		},
		{
			name: "Given retry config with empty initial delay, When IsValid is called, Then should return false",
			config: eventhandler.RetryConfig{
				MaxRetries:    3,
				InitialDelay:  "",
				BackoffFactor: 2.0,
			},
			expected: false,
		},
		{
			name: "Given retry config with zero backoff factor, When IsValid is called, Then should return false",
			config: eventhandler.RetryConfig{
				MaxRetries:    3,
				InitialDelay:  "1s",
				BackoffFactor: 0,
			},
			expected: false,
		},
		{
			name: "Given retry config with negative backoff factor, When IsValid is called, Then should return false",
			config: eventhandler.RetryConfig{
				MaxRetries:    3,
				InitialDelay:  "1s",
				BackoffFactor: -1.0,
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

func TestDefaultEventHandlerConfig(t *testing.T) {
	t.Run("Given default event handler config call, When DefaultEventHandlerConfig is called, Then should return valid default configuration", func(t *testing.T) {
		// Act
		config := eventhandler.DefaultEventHandlerConfig()

		// Assert
		assert.True(t, config.Enabled)
		assert.Equal(t, 1, config.BatchSize)
		assert.Equal(t, 1, config.Concurrency)
		assert.Equal(t, "30s", config.Timeout)
		
		// Check retry config
		assert.Equal(t, 3, config.RetryConfig.MaxRetries)
		assert.Equal(t, "1s", config.RetryConfig.InitialDelay)
		assert.Equal(t, 2.0, config.RetryConfig.BackoffFactor)
		assert.Equal(t, "5m", config.RetryConfig.MaxDelay)
		
		// The default config should be valid for the parts it defines
		assert.True(t, config.RetryConfig.IsValid())
	})
}

func TestEventHandlerError_Error(t *testing.T) {
	tests := []struct {
		name        string
		handlerErr  eventhandler.EventHandlerError
		expected    string
	}{
		{
			name: "Given event handler error with message, When Error is called, Then should return message",
			handlerErr: eventhandler.EventHandlerError{
				Code:    "TEST_ERROR",
				Message: "Test handler error",
			},
			expected: "Test handler error",
		},
		{
			name: "Given event handler error with empty message, When Error is called, Then should return empty string",
			handlerErr: eventhandler.EventHandlerError{
				Code:    "TEST_ERROR",
				Message: "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.handlerErr.Error()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEventHandlerErrors_Constants(t *testing.T) {
	tests := []struct {
		name         string
		err          eventhandler.EventHandlerError
		expectedCode string
	}{
		{
			name:         "Given ErrHandlerNotFound, When accessing code, Then should have correct code",
			err:          eventhandler.ErrHandlerNotFound,
			expectedCode: "HANDLER_NOT_FOUND",
		},
		{
			name:         "Given ErrHandlingFailed, When accessing code, Then should have correct code",
			err:          eventhandler.ErrHandlingFailed,
			expectedCode: "HANDLING_FAILED",
		},
		{
			name:         "Given ErrInvalidEventType, When accessing code, Then should have correct code",
			err:          eventhandler.ErrInvalidEventType,
			expectedCode: "INVALID_EVENT_TYPE",
		},
		{
			name:         "Given ErrHandlerDisabled, When accessing code, Then should have correct code",
			err:          eventhandler.ErrHandlerDisabled,
			expectedCode: "HANDLER_DISABLED",
		},
		{
			name:         "Given ErrHandlerTimeout, When accessing code, Then should have correct code",
			err:          eventhandler.ErrHandlerTimeout,
			expectedCode: "HANDLER_TIMEOUT",
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

func TestEventHandlerConfig_CompleteStructure(t *testing.T) {
	t.Run("Given event handler config with all fields, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		retryConfig := eventhandler.RetryConfig{
			MaxRetries:    5,
			InitialDelay:  "2s",
			BackoffFactor: 1.5,
			MaxDelay:      "10m",
		}
		metadata := map[string]string{
			"version": "1.0",
			"region":  "us-east-1",
		}

		config := eventhandler.EventHandlerConfig{
			HandlerID:   "notification-handler",
			EventTypes:  []string{"user.created", "user.updated", "user.deleted"},
			Enabled:     true,
			RetryConfig: retryConfig,
			Timeout:     "45s",
			BatchSize:   10,
			Concurrency: 5,
			Metadata:    metadata,
		}

		// Assert
		assert.Equal(t, "notification-handler", config.HandlerID)
		assert.Equal(t, []string{"user.created", "user.updated", "user.deleted"}, config.EventTypes)
		assert.True(t, config.Enabled)
		assert.Equal(t, retryConfig, config.RetryConfig)
		assert.Equal(t, "45s", config.Timeout)
		assert.Equal(t, 10, config.BatchSize)
		assert.Equal(t, 5, config.Concurrency)
		assert.Equal(t, metadata, config.Metadata)
		
		// Test helper methods
		assert.True(t, config.IsValid())
		assert.True(t, config.IsEnabled())
		assert.True(t, config.HandlesEventType("user.created"))
		assert.True(t, config.HandlesEventType("user.updated"))
		assert.True(t, config.HandlesEventType("user.deleted"))
		assert.False(t, config.HandlesEventType("order.created"))
	})
}

func TestRetryConfig_CompleteStructure(t *testing.T) {
	t.Run("Given retry config with all fields, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		config := eventhandler.RetryConfig{
			MaxRetries:    10,
			InitialDelay:  "500ms",
			BackoffFactor: 3.0,
			MaxDelay:      "30m",
		}

		// Assert
		assert.Equal(t, 10, config.MaxRetries)
		assert.Equal(t, "500ms", config.InitialDelay)
		assert.Equal(t, 3.0, config.BackoffFactor)
		assert.Equal(t, "30m", config.MaxDelay)
		assert.True(t, config.IsValid())
	})
}

func TestEventHandlerConfig_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		config   eventhandler.EventHandlerConfig
		testFunc func(t *testing.T, config eventhandler.EventHandlerConfig)
	}{
		{
			name: "Given event handler config with single event type, When HandlesEventType is called with various inputs, Then should handle correctly",
			config: eventhandler.EventHandlerConfig{
				HandlerID:  "single-handler",
				EventTypes: []string{"user.created"},
			},
			testFunc: func(t *testing.T, config eventhandler.EventHandlerConfig) {
				assert.True(t, config.HandlesEventType("user.created"))
				assert.False(t, config.HandlesEventType("user.updated"))
				assert.False(t, config.HandlesEventType(""))
				assert.False(t, config.HandlesEventType("user"))
			},
		},
		{
			name: "Given event handler config with very large batch size, When accessing batch size, Then should handle correctly",
			config: eventhandler.EventHandlerConfig{
				HandlerID:  "batch-handler",
				EventTypes: []string{"batch.event"},
				BatchSize:  10000,
			},
			testFunc: func(t *testing.T, config eventhandler.EventHandlerConfig) {
				assert.Equal(t, 10000, config.BatchSize)
				assert.True(t, config.IsValid())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, tt.config)
		})
	}
}