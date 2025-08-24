package memory

import (
	"context"
	"sync"
	"time"

	"github.com/gentra/decorator-arch-go/internal/ratelimit"
)

// service implements ratelimit.Service interface using in-memory storage
type service struct {
	limits   map[string]ratelimit.RateLimitConfig
	counters map[string]*rateLimitCounter
	mu       sync.RWMutex
}

// rateLimitCounter tracks requests for a specific key
type rateLimitCounter struct {
	requests []time.Time
	mu       sync.Mutex
}

// NewService creates a new in-memory rate limiter
func NewService(defaultLimits map[string]ratelimit.RateLimitConfig) ratelimit.Service {
	if defaultLimits == nil {
		defaultLimits = ratelimit.GetDefaultRateLimitConfigs()
	}

	return &service{
		limits:   defaultLimits,
		counters: make(map[string]*rateLimitCounter),
	}
}

// Allow checks if a request is allowed for the given key
func (s *service) Allow(ctx context.Context, key string) (bool, error) {
	s.mu.RLock()
	config, exists := s.limits[getPatternFromKey(key)]
	s.mu.RUnlock()

	if !exists {
		// If no specific limit is configured, allow the request
		return true, nil
	}

	s.mu.Lock()
	counter, exists := s.counters[key]
	if !exists {
		counter = &rateLimitCounter{
			requests: make([]time.Time, 0),
		}
		s.counters[key] = counter
	}
	s.mu.Unlock()

	counter.mu.Lock()
	defer counter.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-config.Window)

	// Remove old requests outside the window
	validRequests := make([]time.Time, 0, len(counter.requests))
	for _, reqTime := range counter.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	counter.requests = validRequests

	// Check if we're under the limit
	if len(counter.requests) < config.Limit {
		counter.requests = append(counter.requests, now)
		return true, nil
	}

	return false, nil
}

// Reset clears the rate limit counter for a key
func (s *service) Reset(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.counters, key)
	return nil
}

// GetStatus returns the current rate limit status for a key
func (s *service) GetStatus(ctx context.Context, key string) (*ratelimit.RateLimitStatus, error) {
	s.mu.RLock()
	config, exists := s.limits[getPatternFromKey(key)]
	counter := s.counters[key]
	s.mu.RUnlock()

	if !exists {
		// No limit configured
		return &ratelimit.RateLimitStatus{
			Key:            key,
			Limit:          -1, // No limit
			Remaining:      -1,
			ResetTime:      time.Now().Add(time.Hour),
			WindowDuration: time.Hour,
		}, nil
	}

	now := time.Now()
	remaining := config.Limit
	resetTime := now.Add(config.Window)

	if counter != nil {
		counter.mu.Lock()
		// Count valid requests
		cutoff := now.Add(-config.Window)
		validRequests := 0
		for _, reqTime := range counter.requests {
			if reqTime.After(cutoff) {
				validRequests++
			}
		}
		remaining = config.Limit - validRequests
		if remaining < 0 {
			remaining = 0
		}
		counter.mu.Unlock()
	}

	status := &ratelimit.RateLimitStatus{
		Key:            key,
		Limit:          config.Limit,
		Remaining:      remaining,
		ResetTime:      resetTime,
		WindowDuration: config.Window,
	}

	if remaining == 0 {
		status.RetryAfter = config.Window
	}

	return status, nil
}

// SetLimit sets a rate limit configuration for a pattern
func (s *service) SetLimit(ctx context.Context, pattern string, config ratelimit.RateLimitConfig) error {
	if !config.IsValid() {
		return &ratelimit.RateLimitError{
			Key:    pattern,
			Limit:  config.Limit,
			Window: config.Window,
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.limits[pattern] = config
	return nil
}

// GetLimit returns the rate limit configuration for a pattern
func (s *service) GetLimit(ctx context.Context, pattern string) (*ratelimit.RateLimitConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if config, exists := s.limits[pattern]; exists {
		return &config, nil
	}

	return nil, nil // No limit configured
}

// RemoveLimit removes a rate limit configuration for a pattern
func (s *service) RemoveLimit(ctx context.Context, pattern string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.limits, pattern)
	return nil
}

// Helper function to extract the pattern from a rate limit key
func getPatternFromKey(key string) string {
	// Extract the pattern part from keys like "user:register:email@example.com"
	// This is a simple implementation - in production you might use regex
	if len(key) > 13 && key[:13] == "user:register" {
		return "user:register"
	}
	if len(key) > 10 && key[:10] == "user:login" {
		return "user:login"
	}
	if len(key) > 9 && key[:9] == "user:read" {
		return "user:read"
	}
	if len(key) > 11 && key[:11] == "user:update" {
		return "user:update"
	}
	if len(key) > 15 && key[:15] == "user:prefs:read" {
		return "user:prefs:read"
	}
	if len(key) > 17 && key[:17] == "user:prefs:update" {
		return "user:prefs:update"
	}

	return "default"
}
