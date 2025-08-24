package mock_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/audit"
	auditmock "github.com/gentra/decorator-arch-go/internal/audit/mock"
)

func TestMockAuditService_InterfaceCompliance(t *testing.T) {
	t.Run("Given MockAuditService, When created, Then should implement audit.Service interface", func(t *testing.T) {
		// Arrange
		mockService := &auditmock.MockAuditService{}

		// Act & Assert
		var _ audit.Service = mockService
		assert.Implements(t, (*audit.Service)(nil), mockService)
	})
}

func TestMockAuditService_MethodCalls(t *testing.T) {
	t.Run("Given MockAuditService with expectations, When methods are called, Then should return expected values", func(t *testing.T) {
		// Arrange
		mockService := &auditmock.MockAuditService{}
		ctx := context.Background()
		entry := audit.AuditEntry{
			Action:    "test.action",
			Resource:  "test.resource",
			Timestamp: time.Now(),
		}
		filters := audit.AuditFilters{}

		// Set up expectations
		mockService.On("Log", ctx, entry).Return(nil)
		mockService.On("GetAuditLogs", ctx, filters).Return([]audit.AuditEntry{}, nil)
		mockService.On("GetAuditLogsByUser", ctx, "user-123", 10).Return([]audit.AuditEntry{}, nil)
		mockService.On("GetAuditLogsByResource", ctx, "resource", "id-123", 10).Return([]audit.AuditEntry{}, nil)

		// Act
		err := mockService.Log(ctx, entry)
		logs, logsErr := mockService.GetAuditLogs(ctx, filters)
		userLogs, userLogsErr := mockService.GetAuditLogsByUser(ctx, "user-123", 10)
		resourceLogs, resourceLogsErr := mockService.GetAuditLogsByResource(ctx, "resource", "id-123", 10)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, logsErr)
		assert.NoError(t, userLogsErr)
		assert.NoError(t, resourceLogsErr)
		assert.Empty(t, logs)
		assert.Empty(t, userLogs)
		assert.Empty(t, resourceLogs)

		// Verify all expectations were met
		mockService.AssertExpectations(t)
	})
}
