package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gentra/decorator-arch-go/internal/auth"
	authmock "github.com/gentra/decorator-arch-go/internal/auth/mock"
	"github.com/gentra/decorator-arch-go/internal/auth/usecase"
	"github.com/gentra/decorator-arch-go/internal/user"
)

func TestJWTAuthStrategy_Authenticate_Simple(t *testing.T) {
	t.Run("Given valid JWT credentials, When Authenticate is called with jwt strategy, Then should authenticate successfully", func(t *testing.T) {
		// Arrange
		mockUserService := new(authmock.MockUserService)
		secret := []byte("test-secret-key-for-testing")
		tokenManager := usecase.NewJWTTokenManager(secret, time.Hour, 24*time.Hour)

		// Generate a real JWT token
		testToken, _, _ := tokenManager.GenerateAuthToken("550e8400-e29b-41d4-a716-446655440000", "test@example.com")

		// Mock user retrieval
		user := &user.User{
			ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockUserService.On("GetByID", mock.Anything, "550e8400-e29b-41d4-a716-446655440000").Return(user, nil)

		jwtAuth := usecase.NewJWTAuthStrategy(mockUserService, tokenManager)

		// Act
		result, err := jwtAuth.Authenticate(context.Background(), "jwt", auth.JWTCredentials{
			Token: testToken,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.User.ID)
		assert.Equal(t, "test@example.com", result.User.Email)
		assert.Equal(t, testToken, result.Token)
		assert.Equal(t, "jwt", result.Strategy)

		mockUserService.AssertExpectations(t)
	})

	t.Run("Given unsupported strategy, When Authenticate is called, Then should return unsupported strategy error", func(t *testing.T) {
		// Arrange
		mockUserService := new(authmock.MockUserService)
		secret := []byte("test-secret-key-for-testing")
		tokenManager := usecase.NewJWTTokenManager(secret, time.Hour, 24*time.Hour)
		jwtAuth := usecase.NewJWTAuthStrategy(mockUserService, tokenManager)

		// Act
		result, err := jwtAuth.Authenticate(context.Background(), "basic", auth.JWTCredentials{
			Token: "valid-jwt-token",
		})

		// Assert
		assert.Error(t, err)
		assert.Equal(t, auth.ErrUnsupportedStrategy, err)
		assert.Nil(t, result)
	})
}

func TestJWTAuthStrategy_GetSupportedStrategies_Simple(t *testing.T) {
	t.Run("Given JWTAuthStrategy, When GetSupportedStrategies is called, Then should return only jwt strategy", func(t *testing.T) {
		// Arrange
		mockUserService := new(authmock.MockUserService)
		secret := []byte("test-secret-key-for-testing")
		tokenManager := usecase.NewJWTTokenManager(secret, time.Hour, 24*time.Hour)
		jwtAuth := usecase.NewJWTAuthStrategy(mockUserService, tokenManager)

		// Act
		strategies := jwtAuth.GetSupportedStrategies()

		// Assert
		assert.Equal(t, []string{"jwt"}, strategies)
	})
}
