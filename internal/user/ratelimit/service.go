package ratelimit

import (
	"context"
	"fmt"

	"github.com/gentra/decorator-arch-go/internal/ratelimit"
	"github.com/gentra/decorator-arch-go/internal/user"
)

// service implements user.Service with rate limiting capabilities
type service struct {
	next             user.Service
	rateLimitService ratelimit.Service
}

// NewService creates a new rate-limited user service
func NewService(next user.Service, rateLimitService ratelimit.Service) user.Service {
	return &service{
		next:             next,
		rateLimitService: rateLimitService,
	}
}

// Register applies rate limiting for user registration
func (s *service) Register(ctx context.Context, data user.RegisterData) (*user.User, error) {
	key := fmt.Sprintf("user:register:%s", data.Email)

	allowed, err := s.rateLimitService.Allow(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	if !allowed {
		return nil, fmt.Errorf("rate limit exceeded for registration")
	}

	return s.next.Register(ctx, data)
}

// Login applies rate limiting for user login attempts
func (s *service) Login(ctx context.Context, email, password string) (*user.AuthResult, error) {
	key := fmt.Sprintf("user:login:%s", email)

	allowed, err := s.rateLimitService.Allow(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	if !allowed {
		return nil, fmt.Errorf("rate limit exceeded for login")
	}

	return s.next.Login(ctx, email, password)
}

// GetByID applies rate limiting for user data retrieval
func (s *service) GetByID(ctx context.Context, id string) (*user.User, error) {
	key := fmt.Sprintf("user:read:%s", id)

	allowed, err := s.rateLimitService.Allow(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	if !allowed {
		return nil, fmt.Errorf("rate limit exceeded for read")
	}

	return s.next.GetByID(ctx, id)
}

// UpdateProfile applies rate limiting for profile updates
func (s *service) UpdateProfile(ctx context.Context, id string, data user.UpdateProfileData) (*user.User, error) {
	key := fmt.Sprintf("user:update:%s", id)

	allowed, err := s.rateLimitService.Allow(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	if !allowed {
		return nil, fmt.Errorf("rate limit exceeded for update")
	}

	return s.next.UpdateProfile(ctx, id, data)
}

// GetPreferences applies rate limiting for preferences retrieval
func (s *service) GetPreferences(ctx context.Context, userID string) (*user.UserPreferences, error) {
	key := fmt.Sprintf("user:prefs:read:%s", userID)

	allowed, err := s.rateLimitService.Allow(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	if !allowed {
		return nil, fmt.Errorf("rate limit exceeded for preferences read")
	}

	return s.next.GetPreferences(ctx, userID)
}

// UpdatePreferences applies rate limiting for preferences updates
func (s *service) UpdatePreferences(ctx context.Context, userID string, prefs user.UserPreferences) error {
	key := fmt.Sprintf("user:prefs:update:%s", userID)

	allowed, err := s.rateLimitService.Allow(ctx, key)
	if err != nil {
		return fmt.Errorf("rate limiter error: %w", err)
	}

	if !allowed {
		return fmt.Errorf("rate limit exceeded for preferences update")
	}

	return s.next.UpdatePreferences(ctx, userID, prefs)
}
