package ratelimit

import (
	"context"
	"time"
)

// Service defines the rate limiting domain interface - the ONLY interface in this domain
type Service interface {
	// Rate limiting operations
	Allow(ctx context.Context, key string) (bool, error)
	Reset(ctx context.Context, key string) error
	GetStatus(ctx context.Context, key string) (*RateLimitStatus, error)
	
	// Configuration management
	SetLimit(ctx context.Context, pattern string, config RateLimitConfig) error
	GetLimit(ctx context.Context, pattern string) (*RateLimitConfig, error)
	RemoveLimit(ctx context.Context, pattern string) error
}

// Domain types and data structures

// RateLimitConfig defines the configuration for a rate limit rule
type RateLimitConfig struct {
	Limit  int           `json:"limit"`  // Number of requests allowed
	Window time.Duration `json:"window"` // Time window for the limit
}

// RateLimitStatus represents the current status of rate limiting for a key
type RateLimitStatus struct {
	Key            string        `json:"key"`
	Limit          int           `json:"limit"`
	Remaining      int           `json:"remaining"`
	ResetTime      time.Time     `json:"reset_time"`
	RetryAfter     time.Duration `json:"retry_after,omitempty"`
	WindowDuration time.Duration `json:"window_duration"`
}

// RateLimitError represents a rate limit exceeded error
type RateLimitError struct {
	Key        string        `json:"key"`
	Limit      int           `json:"limit"`
	Window     time.Duration `json:"window"`
	RetryAfter time.Duration `json:"retry_after"`
}

func (e *RateLimitError) Error() string {
	return "rate limit exceeded for key " + e.Key
}

// RateLimitMetrics contains metrics for rate limiting
type RateLimitMetrics struct {
	TotalRequests   int64         `json:"total_requests"`
	AllowedRequests int64         `json:"allowed_requests"`
	BlockedRequests int64         `json:"blocked_requests"`
	ActiveKeys      int           `json:"active_keys"`
	AverageLatency  time.Duration `json:"average_latency"`
}

// Helper methods for RateLimitConfig
func (c *RateLimitConfig) IsValid() bool {
	return c.Limit > 0 && c.Window > 0
}

func (c *RateLimitConfig) RequestsPerSecond() float64 {
	if c.Window == 0 {
		return 0
	}
	return float64(c.Limit) / c.Window.Seconds()
}

// Helper methods for RateLimitStatus
func (s *RateLimitStatus) IsAllowed() bool {
	return s.Remaining > 0
}

func (s *RateLimitStatus) IsExpired() bool {
	return time.Now().After(s.ResetTime)
}

func (s *RateLimitStatus) TimeUntilReset() time.Duration {
	return time.Until(s.ResetTime)
}

// Default rate limit configurations for common patterns
func GetDefaultRateLimitConfigs() map[string]RateLimitConfig {
	return map[string]RateLimitConfig{
		"user:register":     {Limit: 5, Window: time.Hour},         // 5 registrations per hour per email
		"user:login":        {Limit: 10, Window: 15 * time.Minute}, // 10 login attempts per 15 minutes per email
		"user:read":         {Limit: 100, Window: time.Minute},     // 100 reads per minute per user
		"user:update":       {Limit: 20, Window: time.Hour},        // 20 updates per hour per user
		"user:prefs:read":   {Limit: 50, Window: time.Minute},      // 50 preference reads per minute per user
		"user:prefs:update": {Limit: 10, Window: time.Hour},        // 10 preference updates per hour per user
		"default":           {Limit: 1000, Window: time.Hour},      // Default fallback limit
	}
}
