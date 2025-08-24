package audit_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/audit"
)

func TestAuditEntry_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		entry   audit.AuditEntry
		isValid bool
	}{
		{
			name: "Given audit entry with all required fields, When IsValid is called, Then should return true",
			entry: audit.AuditEntry{
				Action:    "user.login",
				Resource:  "auth",
				Timestamp: time.Now(),
			},
			isValid: true,
		},
		{
			name: "Given audit entry with missing action, When IsValid is called, Then should return false",
			entry: audit.AuditEntry{
				Resource:  "auth",
				Timestamp: time.Now(),
			},
			isValid: false,
		},
		{
			name: "Given audit entry with missing resource, When IsValid is called, Then should return false",
			entry: audit.AuditEntry{
				Action:    "user.login",
				Timestamp: time.Now(),
			},
			isValid: false,
		},
		{
			name: "Given audit entry with zero timestamp, When IsValid is called, Then should return false",
			entry: audit.AuditEntry{
				Action:   "user.login",
				Resource: "auth",
			},
			isValid: false,
		},
		{
			name: "Given audit entry with empty action, When IsValid is called, Then should return false",
			entry: audit.AuditEntry{
				Action:    "",
				Resource:  "auth",
				Timestamp: time.Now(),
			},
			isValid: false,
		},
		{
			name: "Given audit entry with empty resource, When IsValid is called, Then should return false",
			entry: audit.AuditEntry{
				Action:    "user.login",
				Resource:  "",
				Timestamp: time.Now(),
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.entry.IsValid()

			// Assert
			assert.Equal(t, tt.isValid, result)
		})
	}
}

func TestAuditEntry_SetSuccess(t *testing.T) {
	t.Run("Given audit entry with error, When SetSuccess is called, Then should set success to true and clear error", func(t *testing.T) {
		// Arrange
		entry := audit.AuditEntry{
			Success: false,
			Error:   "Previous error occurred",
		}

		// Act
		entry.SetSuccess()

		// Assert
		assert.True(t, entry.Success)
		assert.Empty(t, entry.Error)
	})
}

func TestAuditEntry_SetError(t *testing.T) {
	tests := []struct {
		name        string
		entry       audit.AuditEntry
		err         error
		expectedErr string
	}{
		{
			name: "Given audit entry and error, When SetError is called, Then should set success to false and set error message",
			entry: audit.AuditEntry{
				Success: true,
				Error:   "",
			},
			err:         errors.New("Authentication failed"),
			expectedErr: "Authentication failed",
		},
		{
			name: "Given audit entry and nil error, When SetError is called, Then should set success to false and keep previous error",
			entry: audit.AuditEntry{
				Success: true,
				Error:   "Previous error",
			},
			err:         nil,
			expectedErr: "Previous error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			tt.entry.SetError(tt.err)

			// Assert
			assert.False(t, tt.entry.Success)
			assert.Equal(t, tt.expectedErr, tt.entry.Error)
		})
	}
}

func TestAuditContext_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		context audit.AuditContext
		isValid bool
	}{
		{
			name: "Given audit context with user ID, When IsValid is called, Then should return true",
			context: audit.AuditContext{
				CurrentUserID: "user-123",
			},
			isValid: true,
		},
		{
			name: "Given audit context with IP address, When IsValid is called, Then should return true",
			context: audit.AuditContext{
				IPAddress: "192.168.1.1",
			},
			isValid: true,
		},
		{
			name: "Given audit context with both user ID and IP address, When IsValid is called, Then should return true",
			context: audit.AuditContext{
				CurrentUserID: "user-123",
				IPAddress:     "192.168.1.1",
			},
			isValid: true,
		},
		{
			name: "Given audit context with only user agent, When IsValid is called, Then should return false",
			context: audit.AuditContext{
				UserAgent: "Mozilla/5.0",
			},
			isValid: false,
		},
		{
			name: "Given audit context with only session ID, When IsValid is called, Then should return false",
			context: audit.AuditContext{
				SessionID: "session-123",
			},
			isValid: false,
		},
		{
			name:    "Given empty audit context, When IsValid is called, Then should return false",
			context: audit.AuditContext{},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.context.IsValid()

			// Assert
			assert.Equal(t, tt.isValid, result)
		})
	}
}

func TestWithAuditContext(t *testing.T) {
	t.Run("Given context and audit information, When WithAuditContext is called, Then should return context with audit information", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := "user-123"
		ipAddress := "192.168.1.1"
		userAgent := "Mozilla/5.0"
		sessionID := "session-456"

		// Act
		resultCtx := audit.WithAuditContext(ctx, userID, ipAddress, userAgent, sessionID)

		// Assert
		auditCtx := audit.ExtractAuditContext(resultCtx)
		assert.Equal(t, userID, auditCtx.CurrentUserID)
		assert.Equal(t, ipAddress, auditCtx.IPAddress)
		assert.Equal(t, userAgent, auditCtx.UserAgent)
		assert.Equal(t, sessionID, auditCtx.SessionID)
	})
}

func TestExtractAuditContext(t *testing.T) {
	tests := []struct {
		name           string
		context        context.Context
		expectedResult audit.AuditContext
	}{
		{
			name: "Given context with audit information, When ExtractAuditContext is called, Then should return audit context",
			context: audit.WithAuditContext(
				context.Background(),
				"user-123",
				"192.168.1.1",
				"Mozilla/5.0",
				"session-456",
			),
			expectedResult: audit.AuditContext{
				CurrentUserID: "user-123",
				IPAddress:     "192.168.1.1",
				UserAgent:     "Mozilla/5.0",
				SessionID:     "session-456",
			},
		},
		{
			name:           "Given context without audit information, When ExtractAuditContext is called, Then should return empty audit context",
			context:        context.Background(),
			expectedResult: audit.AuditContext{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := audit.ExtractAuditContext(tt.context)

			// Assert
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestAuditEntry_CompleteExample(t *testing.T) {
	t.Run("Given complete audit entry workflow, When all methods are called, Then should work correctly", func(t *testing.T) {
		// Arrange
		entry := audit.AuditEntry{
			ID:         "audit-123",
			Timestamp:  time.Now(),
			UserID:     "user-456",
			Action:     "user.login",
			Resource:   "auth",
			ResourceID: "session-789",
			Details: map[string]interface{}{
				"ip_address": "192.168.1.1",
				"user_agent": "Mozilla/5.0",
			},
			IPAddress: "192.168.1.1",
			UserAgent: "Mozilla/5.0",
			SessionID: "session-789",
		}

		// Act & Assert - Test validation
		assert.True(t, entry.IsValid(), "Entry should be valid with all required fields")

		// Act & Assert - Test success setting
		entry.SetSuccess()
		assert.True(t, entry.Success)
		assert.Empty(t, entry.Error)

		// Act & Assert - Test error setting
		testError := errors.New("Test error occurred")
		entry.SetError(testError)
		assert.False(t, entry.Success)
		assert.Equal(t, testError.Error(), entry.Error)

		// Act & Assert - Test success setting again
		entry.SetSuccess()
		assert.True(t, entry.Success)
		assert.Empty(t, entry.Error)
	})
}

func TestAuditContext_CompleteExample(t *testing.T) {
	t.Run("Given complete audit context workflow, When all methods are called, Then should work correctly", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := "user-123"
		ipAddress := "192.168.1.1"
		userAgent := "Mozilla/5.0"
		sessionID := "session-456"

		// Act
		auditCtx := audit.WithAuditContext(ctx, userID, ipAddress, userAgent, sessionID)
		extractedCtx := audit.ExtractAuditContext(auditCtx)

		// Assert
		assert.True(t, extractedCtx.IsValid(), "Extracted context should be valid")
		assert.Equal(t, userID, extractedCtx.CurrentUserID)
		assert.Equal(t, ipAddress, extractedCtx.IPAddress)
		assert.Equal(t, userAgent, extractedCtx.UserAgent)
		assert.Equal(t, sessionID, extractedCtx.SessionID)
	})
}
