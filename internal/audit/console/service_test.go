package console_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gentra/decorator-arch-go/internal/audit"
	"github.com/gentra/decorator-arch-go/internal/audit/console"
)

func TestConsoleAuditService_Log(t *testing.T) {
	tests := []struct {
		name    string
		entry   audit.AuditEntry
		wantErr bool
	}{
		{
			name: "Given valid audit entry with all fields, When Log is called, Then should succeed without error",
			entry: audit.AuditEntry{
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
				Success:   true,
				IPAddress: "192.168.1.1",
				UserAgent: "Mozilla/5.0",
				SessionID: "session-789",
			},
			wantErr: false,
		},
		{
			name: "Given audit entry with minimal required fields, When Log is called, Then should succeed without error",
			entry: audit.AuditEntry{
				Action:    "user.logout",
				Resource:  "auth",
				Timestamp: time.Now(),
				Success:   true,
			},
			wantErr: false,
		},
		{
			name: "Given audit entry with error details, When Log is called, Then should succeed without error",
			entry: audit.AuditEntry{
				ID:        "audit-error-123",
				Timestamp: time.Now(),
				Action:    "user.login",
				Resource:  "auth",
				Success:   false,
				Error:     "Invalid credentials",
			},
			wantErr: false,
		},
		{
			name: "Given audit entry with complex details, When Log is called, Then should succeed without error",
			entry: audit.AuditEntry{
				ID:        "audit-complex-123",
				Timestamp: time.Now(),
				Action:    "data.export",
				Resource:  "reports",
				Details: map[string]interface{}{
					"file_size":    1024,
					"format":       "csv",
					"record_count": 1000,
					"filters": map[string]interface{}{
						"date_from": "2024-01-01",
						"date_to":   "2024-12-31",
					},
				},
				Success: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := console.NewService()
			ctx := context.Background()

			// Act
			err := service.Log(ctx, tt.entry)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConsoleAuditService_GetAuditLogs(t *testing.T) {
	tests := []struct {
		name    string
		filters audit.AuditFilters
		want    []audit.AuditEntry
		wantErr bool
	}{
		{
			name: "Given any audit filters, When GetAuditLogs is called, Then should return nil result and no error",
			filters: audit.AuditFilters{
				UserID:     "user-123",
				Action:     "user.login",
				Resource:   "auth",
				ResourceID: "session-456",
				Success:    &[]bool{true}[0],
				Limit:      100,
				Offset:     0,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Given empty audit filters, When GetAuditLogs is called, Then should return nil result and no error",
			filters: audit.AuditFilters{},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := console.NewService()
			ctx := context.Background()

			// Act
			result, err := service.GetAuditLogs(ctx, tt.filters)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestConsoleAuditService_GetAuditLogsByUser(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		limit   int
		want    []audit.AuditEntry
		wantErr bool
	}{
		{
			name:    "Given valid user ID and limit, When GetAuditLogsByUser is called, Then should return nil result and no error",
			userID:  "user-123",
			limit:   100,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Given empty user ID, When GetAuditLogsByUser is called, Then should return nil result and no error",
			userID:  "",
			limit:   50,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Given zero limit, When GetAuditLogsByUser is called, Then should return nil result and no error",
			userID:  "user-456",
			limit:   0,
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := console.NewService()
			ctx := context.Background()

			// Act
			result, err := service.GetAuditLogsByUser(ctx, tt.userID, tt.limit)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestConsoleAuditService_GetAuditLogsByResource(t *testing.T) {
	tests := []struct {
		name       string
		resource   string
		resourceID string
		limit      int
		want       []audit.AuditEntry
		wantErr    bool
	}{
		{
			name:       "Given valid resource and resource ID, When GetAuditLogsByResource is called, Then should return nil result and no error",
			resource:   "users",
			resourceID: "user-123",
			limit:      100,
			want:       nil,
			wantErr:    false,
		},
		{
			name:       "Given empty resource, When GetAuditLogsByResource is called, Then should return nil result and no error",
			resource:   "",
			resourceID: "user-456",
			limit:      50,
			want:       nil,
			wantErr:    false,
		},
		{
			name:       "Given empty resource ID, When GetAuditLogsByResource is called, Then should return nil result and no error",
			resource:   "auth",
			resourceID: "",
			limit:      25,
			want:       nil,
			wantErr:    false,
		},
		{
			name:       "Given zero limit, When GetAuditLogsByResource is called, Then should return nil result and no error",
			resource:   "reports",
			resourceID: "report-789",
			limit:      0,
			want:       nil,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := console.NewService()
			ctx := context.Background()

			// Act
			result, err := service.GetAuditLogsByResource(ctx, tt.resource, tt.resourceID, tt.limit)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestConsoleAuditService_NewService(t *testing.T) {
	t.Run("Given NewService is called, When service is created, Then should return audit.Service interface", func(t *testing.T) {
		// Act
		service := console.NewService()

		// Assert
		require.NotNil(t, service)
		_, ok := service.(audit.Service)
		assert.True(t, ok, "Service should implement audit.Service interface")
	})
}
