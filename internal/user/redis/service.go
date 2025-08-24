package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/gentra/decorator-arch-go/internal/user"
)

// service implements the user.Service interface with Redis caching
type service struct {
	next   user.Service
	client *redis.Client
	ttl    time.Duration
}

// NewService creates a new Redis-backed user service
func NewService(next user.Service, client *redis.Client, ttl time.Duration) user.Service {
	return &service{
		next:   next,
		client: client,
		ttl:    ttl,
	}
}

// Register creates a new user (cache invalidation pattern)
func (s *service) Register(ctx context.Context, data user.RegisterData) (*user.User, error) {
	// Call next service to register user
	result, err := s.next.Register(ctx, data)
	if err != nil {
		return nil, err
	}

	// Cache the newly created user
	if err := s.cacheUser(ctx, result); err != nil {
		// Log error but don't fail the registration
		// In production, you'd use a proper logger
		fmt.Printf("Failed to cache user after registration: %v\n", err)
	}

	// Invalidate email cache if it exists
	emailCacheKey := s.getEmailCacheKey(data.Email)
	s.client.Del(ctx, emailCacheKey)

	return result, nil
}

// Login authenticates a user (cache aside pattern)
func (s *service) Login(ctx context.Context, email, password string) (*user.AuthResult, error) {
	// For login, we don't cache credentials for security reasons
	// We go directly to the next service
	result, err := s.next.Login(ctx, email, password)
	if err != nil {
		return nil, err
	}

	// Cache the user data after successful login
	if result.User != nil {
		if err := s.cacheUser(ctx, result.User); err != nil {
			fmt.Printf("Failed to cache user after login: %v\n", err)
		}
	}

	return result, nil
}

// GetByID retrieves a user by ID (cache aside pattern)
func (s *service) GetByID(ctx context.Context, id string) (*user.User, error) {
	// Try to get from cache first
	cacheKey := s.getUserCacheKey(id)
	cached, err := s.client.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit - deserialize and return
		var cachedUser user.User
		if err := json.Unmarshal([]byte(cached), &cachedUser); err == nil {
			return &cachedUser, nil
		}
		// If deserialization fails, continue to fetch from next service
		fmt.Printf("Failed to deserialize cached user: %v\n", err)
	} else if err != redis.Nil {
		// Log cache error but continue to next service
		fmt.Printf("Cache error for user %s: %v\n", id, err)
	}

	// Cache miss or error - get from next service
	result, err := s.next.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if err := s.cacheUser(ctx, result); err != nil {
		fmt.Printf("Failed to cache user %s: %v\n", id, err)
	}

	return result, nil
}

// UpdateProfile updates user profile (cache invalidation pattern)
func (s *service) UpdateProfile(ctx context.Context, id string, data user.UpdateProfileData) (*user.User, error) {
	// Call next service to update profile
	result, err := s.next.UpdateProfile(ctx, id, data)
	if err != nil {
		return nil, err
	}

	// Invalidate cache for this user
	cacheKey := s.getUserCacheKey(id)
	if err := s.client.Del(ctx, cacheKey).Err(); err != nil {
		fmt.Printf("Failed to invalidate cache for user %s: %v\n", id, err)
	}

	// If email was updated, invalidate old email cache
	if data.Email != nil {
		// We can't know the old email without another query, so we just cache the new data
		if err := s.cacheUser(ctx, result); err != nil {
			fmt.Printf("Failed to cache updated user %s: %v\n", id, err)
		}
	}

	return result, nil
}

// GetPreferences retrieves user preferences (cache aside pattern)
func (s *service) GetPreferences(ctx context.Context, userID string) (*user.UserPreferences, error) {
	// Try to get from cache first
	cacheKey := s.getPreferencesCacheKey(userID)
	cached, err := s.client.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit - deserialize and return
		var cachedPrefs user.UserPreferences
		if err := json.Unmarshal([]byte(cached), &cachedPrefs); err == nil {
			return &cachedPrefs, nil
		}
		// If deserialization fails, continue to fetch from next service
		fmt.Printf("Failed to deserialize cached preferences: %v\n", err)
	} else if err != redis.Nil {
		// Log cache error but continue to next service
		fmt.Printf("Cache error for preferences %s: %v\n", userID, err)
	}

	// Cache miss or error - get from next service
	result, err := s.next.GetPreferences(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if err := s.cachePreferences(ctx, userID, result); err != nil {
		fmt.Printf("Failed to cache preferences %s: %v\n", userID, err)
	}

	return result, nil
}

// UpdatePreferences updates user preferences (cache invalidation pattern)
func (s *service) UpdatePreferences(ctx context.Context, userID string, prefs user.UserPreferences) error {
	// Call next service to update preferences
	err := s.next.UpdatePreferences(ctx, userID, prefs)
	if err != nil {
		return err
	}

	// Invalidate cache for these preferences
	cacheKey := s.getPreferencesCacheKey(userID)
	if err := s.client.Del(ctx, cacheKey).Err(); err != nil {
		fmt.Printf("Failed to invalidate preferences cache for user %s: %v\n", userID, err)
	}

	// Cache the updated preferences
	if err := s.cachePreferences(ctx, userID, &prefs); err != nil {
		fmt.Printf("Failed to cache updated preferences %s: %v\n", userID, err)
	}

	return nil
}

// Helper methods for caching operations

func (s *service) cacheUser(ctx context.Context, u *user.User) error {
	// Serialize user to JSON
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	// Store in cache with TTL
	cacheKey := s.getUserCacheKey(u.ID.String())
	return s.client.Set(ctx, cacheKey, data, s.ttl).Err()
}

func (s *service) cachePreferences(ctx context.Context, userID string, prefs *user.UserPreferences) error {
	// Serialize preferences to JSON
	data, err := json.Marshal(prefs)
	if err != nil {
		return err
	}

	// Store in cache with TTL
	cacheKey := s.getPreferencesCacheKey(userID)
	return s.client.Set(ctx, cacheKey, data, s.ttl).Err()
}

func (s *service) getUserCacheKey(userID string) string {
	return fmt.Sprintf("user:%s", userID)
}

func (s *service) getPreferencesCacheKey(userID string) string {
	return fmt.Sprintf("user_preferences:%s", userID)
}

func (s *service) getEmailCacheKey(email string) string {
	return fmt.Sprintf("user_email:%s", email)
}
