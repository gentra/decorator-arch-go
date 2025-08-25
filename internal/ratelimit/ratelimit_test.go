package ratelimit_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/ratelimit"
)

func TestRateLimitConfig_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		config   ratelimit.RateLimitConfig
		expected bool
	}{
		{
			name: "Given rate limit config with positive limit and window, When IsValid is called, Then should return true",
			config: ratelimit.RateLimitConfig{
				Limit:  100,
				Window: time.Minute,
			},
			expected: true,
		},
		{
			name: "Given rate limit config with zero limit, When IsValid is called, Then should return false",
			config: ratelimit.RateLimitConfig{
				Limit:  0,
				Window: time.Minute,
			},
			expected: false,
		},
		{
			name: "Given rate limit config with negative limit, When IsValid is called, Then should return false",
			config: ratelimit.RateLimitConfig{
				Limit:  -1,
				Window: time.Minute,
			},
			expected: false,
		},
		{
			name: "Given rate limit config with zero window, When IsValid is called, Then should return false",
			config: ratelimit.RateLimitConfig{
				Limit:  100,
				Window: 0,
			},
			expected: false,
		},
		{
			name: "Given rate limit config with negative window, When IsValid is called, Then should return false",
			config: ratelimit.RateLimitConfig{
				Limit:  100,
				Window: -time.Minute,
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

func TestRateLimitConfig_RequestsPerSecond(t *testing.T) {
	tests := []struct {
		name     string
		config   ratelimit.RateLimitConfig
		expected float64
	}{
		{
			name: "Given rate limit config with 60 requests per minute, When RequestsPerSecond is called, Then should return 1.0",
			config: ratelimit.RateLimitConfig{
				Limit:  60,
				Window: time.Minute,
			},
			expected: 1.0,
		},
		{
			name: "Given rate limit config with 100 requests per hour, When RequestsPerSecond is called, Then should return correct value",
			config: ratelimit.RateLimitConfig{
				Limit:  100,
				Window: time.Hour,
			},
			expected: 100.0 / 3600.0, // 100 / (60 * 60)
		},
		{
			name: "Given rate limit config with 10 requests per second, When RequestsPerSecond is called, Then should return 10.0",
			config: ratelimit.RateLimitConfig{
				Limit:  10,
				Window: time.Second,
			},
			expected: 10.0,
		},
		{
			name: "Given rate limit config with zero window, When RequestsPerSecond is called, Then should return 0",
			config: ratelimit.RateLimitConfig{
				Limit:  100,
				Window: 0,
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.config.RequestsPerSecond()

			// Assert
			assert.InDelta(t, tt.expected, result, 0.0001) // Allow small floating point differences
		})
	}
}

func TestRateLimitStatus_IsAllowed(t *testing.T) {
	tests := []struct {
		name     string
		status   ratelimit.RateLimitStatus
		expected bool
	}{
		{
			name: "Given rate limit status with remaining > 0, When IsAllowed is called, Then should return true",
			status: ratelimit.RateLimitStatus{
				Remaining: 5,
			},
			expected: true,
		},
		{
			name: "Given rate limit status with remaining = 0, When IsAllowed is called, Then should return false",
			status: ratelimit.RateLimitStatus{
				Remaining: 0,
			},
			expected: false,
		},
		{
			name: "Given rate limit status with negative remaining, When IsAllowed is called, Then should return false",
			status: ratelimit.RateLimitStatus{
				Remaining: -1,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.status.IsAllowed()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRateLimitStatus_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		status   ratelimit.RateLimitStatus
		expected bool
	}{
		{
			name: "Given rate limit status with future reset time, When IsExpired is called, Then should return false",
			status: ratelimit.RateLimitStatus{
				ResetTime: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "Given rate limit status with past reset time, When IsExpired is called, Then should return true",
			status: ratelimit.RateLimitStatus{
				ResetTime: time.Now().Add(-time.Hour),
			},
			expected: true,
		},
		{
			name: "Given rate limit status with current reset time, When IsExpired is called, Then should return true",
			status: ratelimit.RateLimitStatus{
				ResetTime: time.Now(),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.status.IsExpired()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRateLimitStatus_TimeUntilReset(t *testing.T) {
	tests := []struct {
		name     string
		status   ratelimit.RateLimitStatus
		expected bool // true if positive, false if negative/zero
	}{
		{
			name: "Given rate limit status with future reset time, When TimeUntilReset is called, Then should return positive duration",
			status: ratelimit.RateLimitStatus{
				ResetTime: time.Now().Add(time.Hour),
			},
			expected: true,
		},
		{
			name: "Given rate limit status with past reset time, When TimeUntilReset is called, Then should return negative duration",
			status: ratelimit.RateLimitStatus{
				ResetTime: time.Now().Add(-time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.status.TimeUntilReset()

			// Assert
			if tt.expected {
				assert.True(t, result > 0)
			} else {
				assert.True(t, result <= 0)
			}
		})
	}
}

func TestRateLimitError_Error(t *testing.T) {
	tests := []struct {
		name     string
		ratErr   ratelimit.RateLimitError
		expected string
	}{
		{
			name: "Given rate limit error with key, When Error is called, Then should return formatted error message",
			ratErr: ratelimit.RateLimitError{
				Key:        "user:123",
				Limit:      100,
				Window:     time.Minute,
				RetryAfter: time.Second * 30,
			},
			expected: "rate limit exceeded for key user:123",
		},
		{
			name: "Given rate limit error with empty key, When Error is called, Then should return formatted error message",
			ratErr: ratelimit.RateLimitError{
				Key:        "",
				Limit:      100,
				Window:     time.Minute,
				RetryAfter: time.Second * 30,
			},
			expected: "rate limit exceeded for key ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.ratErr.Error()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetDefaultRateLimitConfigs(t *testing.T) {
	t.Run("Given default rate limit configs call, When GetDefaultRateLimitConfigs is called, Then should return valid default configurations", func(t *testing.T) {
		// Act
		configs := ratelimit.GetDefaultRateLimitConfigs()

		// Assert
		assert.NotNil(t, configs)
		
		expectedConfigs := map[string]struct {
			limit  int
			window time.Duration
		}{
			"user:register":     {5, time.Hour},
			"user:login":        {10, 15 * time.Minute},
			"user:read":         {100, time.Minute},
			"user:update":       {20, time.Hour},
			"user:prefs:read":   {50, time.Minute},
			"user:prefs:update": {10, time.Hour},
			"default":           {1000, time.Hour},
		}
		
		for pattern, expected := range expectedConfigs {
			config, exists := configs[pattern]
			assert.True(t, exists, "Pattern %s should exist", pattern)
			assert.Equal(t, expected.limit, config.Limit, "Pattern %s should have correct limit", pattern)
			assert.Equal(t, expected.window, config.Window, "Pattern %s should have correct window", pattern)
			assert.True(t, config.IsValid(), "Pattern %s should have valid config", pattern)
		}
	})
}

func TestRateLimitStatus_CompleteStructure(t *testing.T) {
	t.Run("Given rate limit status with all fields, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		resetTime := time.Now().Add(time.Hour)
		retryAfter := time.Minute * 30
		windowDuration := time.Hour
		
		status := ratelimit.RateLimitStatus{
			Key:            "user:123:login",
			Limit:          100,
			Remaining:      25,
			ResetTime:      resetTime,
			RetryAfter:     retryAfter,
			WindowDuration: windowDuration,
		}

		// Assert
		assert.Equal(t, "user:123:login", status.Key)
		assert.Equal(t, 100, status.Limit)
		assert.Equal(t, 25, status.Remaining)
		assert.Equal(t, resetTime, status.ResetTime)
		assert.Equal(t, retryAfter, status.RetryAfter)
		assert.Equal(t, windowDuration, status.WindowDuration)
		assert.True(t, status.IsAllowed())
		assert.False(t, status.IsExpired())
		assert.True(t, status.TimeUntilReset() > 0)
	})
}

func TestRateLimitMetrics_Structure(t *testing.T) {
	t.Run("Given rate limit metrics with all fields, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		metrics := ratelimit.RateLimitMetrics{
			TotalRequests:   1000,
			AllowedRequests: 850,
			BlockedRequests: 150,
			ActiveKeys:      25,
			AverageLatency:  time.Millisecond * 5,
		}

		// Assert
		assert.Equal(t, int64(1000), metrics.TotalRequests)
		assert.Equal(t, int64(850), metrics.AllowedRequests)
		assert.Equal(t, int64(150), metrics.BlockedRequests)
		assert.Equal(t, 25, metrics.ActiveKeys)
		assert.Equal(t, time.Millisecond*5, metrics.AverageLatency)
	})
}

func TestRateLimitConfig_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		config   ratelimit.RateLimitConfig
		testFunc func(t *testing.T, config ratelimit.RateLimitConfig)
	}{
		{
			name: "Given rate limit config with very small window, When RequestsPerSecond is called, Then should handle correctly",
			config: ratelimit.RateLimitConfig{
				Limit:  1,
				Window: time.Millisecond,
			},
			testFunc: func(t *testing.T, config ratelimit.RateLimitConfig) {
				rps := config.RequestsPerSecond()
				assert.Equal(t, 1000.0, rps) // 1 request per millisecond = 1000 per second
			},
		},
		{
			name: "Given rate limit config with very large limit, When RequestsPerSecond is called, Then should handle correctly",
			config: ratelimit.RateLimitConfig{
				Limit:  1000000,
				Window: time.Hour,
			},
			testFunc: func(t *testing.T, config ratelimit.RateLimitConfig) {
				rps := config.RequestsPerSecond()
				expected := 1000000.0 / 3600.0
				assert.InDelta(t, expected, rps, 0.1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, tt.config)
		})
	}
}