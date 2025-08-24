package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/auth"
	authmock "github.com/gentra/decorator-arch-go/internal/auth/mock"
	"github.com/gentra/decorator-arch-go/internal/auth/usecase"
)

func TestOAuthAuthStrategy_Authenticate_Simple(t *testing.T) {
	t.Run("Given valid OAuth credentials with configured provider, When Authenticate is called with oauth strategy, Then should delegate to provider", func(t *testing.T) {
		// Arrange
		mockUserService := new(authmock.MockUserService)
		secret := []byte("test-secret-key-for-testing")
		tokenManager := usecase.NewJWTTokenManager(secret, time.Hour, 24*time.Hour)

		// Create a mock OAuth provider
		mockProvider := new(authmock.MockAuthStrategy)
		expectedResult := &auth.AuthResult{
			User: &auth.User{
				ID:    "user-123",
				Email: "test@example.com",
			},
			Token:     "oauth-generated-token",
			ExpiresAt: time.Now().Add(time.Hour),
			Strategy:  "oauth",
		}
		mockProvider.On("Authenticate", context.Background(), "oauth", auth.OAuthCredentials{
			Provider:    "google",
			AccessToken: "oauth-access-token",
		}).Return(expectedResult, nil)

		oauthProviders := map[string]auth.Service{
			"google": mockProvider,
		}
		oauthAuth := usecase.NewOAuthAuthStrategy(mockUserService, tokenManager, oauthProviders)

		// Act
		result, err := oauthAuth.Authenticate(context.Background(), "oauth", auth.OAuthCredentials{
			Provider:    "google",
			AccessToken: "oauth-access-token",
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "user-123", result.User.ID)
		assert.Equal(t, "test@example.com", result.User.Email)
		assert.Equal(t, "oauth-generated-token", result.Token)
		assert.Equal(t, "oauth", result.Strategy)

		mockProvider.AssertExpectations(t)
	})

	t.Run("Given unsupported OAuth provider, When Authenticate is called, Then should return provider not found error", func(t *testing.T) {
		// Arrange
		mockUserService := new(authmock.MockUserService)
		secret := []byte("test-secret-key-for-testing")
		tokenManager := usecase.NewJWTTokenManager(secret, time.Hour, 24*time.Hour)

		oauthAuth := usecase.NewOAuthAuthStrategy(mockUserService, tokenManager, make(map[string]auth.Service))

		// Act
		result, err := oauthAuth.Authenticate(context.Background(), "oauth", auth.OAuthCredentials{
			Provider:    "unsupported-provider",
			AccessToken: "oauth-access-token",
		})

		// Assert
		assert.Error(t, err)
		assert.Equal(t, auth.ErrOAuthProviderNotFound, err)
		assert.Nil(t, result)
	})

	t.Run("Given unsupported strategy, When Authenticate is called, Then should return unsupported strategy error", func(t *testing.T) {
		// Arrange
		mockUserService := new(authmock.MockUserService)
		secret := []byte("test-secret-key-for-testing")
		tokenManager := usecase.NewJWTTokenManager(secret, time.Hour, 24*time.Hour)
		oauthAuth := usecase.NewOAuthAuthStrategy(mockUserService, tokenManager, make(map[string]auth.Service))

		// Act
		result, err := oauthAuth.Authenticate(context.Background(), "basic", auth.OAuthCredentials{
			Provider:    "google",
			AccessToken: "oauth-access-token",
		})

		// Assert
		assert.Error(t, err)
		assert.Equal(t, auth.ErrUnsupportedStrategy, err)
		assert.Nil(t, result)
	})
}

func TestOAuthAuthStrategy_GetSupportedStrategies_Simple(t *testing.T) {
	t.Run("Given OAuthAuthStrategy, When GetSupportedStrategies is called, Then should return only oauth strategy", func(t *testing.T) {
		// Arrange
		mockUserService := new(authmock.MockUserService)
		secret := []byte("test-secret-key-for-testing")
		tokenManager := usecase.NewJWTTokenManager(secret, time.Hour, 24*time.Hour)
		oauthAuth := usecase.NewOAuthAuthStrategy(mockUserService, tokenManager, make(map[string]auth.Service))

		// Act
		strategies := oauthAuth.GetSupportedStrategies()

		// Assert
		assert.Equal(t, []string{"oauth"}, strategies)
	})
}
