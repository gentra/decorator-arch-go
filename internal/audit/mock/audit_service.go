package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/gentra/decorator-arch-go/internal/audit"
)

// MockAuditService is a mock implementation of audit.Service
type MockAuditService struct {
	mock.Mock
}

// Log mocks the Log method
func (m *MockAuditService) Log(ctx context.Context, entry audit.AuditEntry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

// GetAuditLogs mocks the GetAuditLogs method
func (m *MockAuditService) GetAuditLogs(ctx context.Context, filters audit.AuditFilters) ([]audit.AuditEntry, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]audit.AuditEntry), args.Error(1)
}

// GetAuditLogsByUser mocks the GetAuditLogsByUser method
func (m *MockAuditService) GetAuditLogsByUser(ctx context.Context, userID string, limit int) ([]audit.AuditEntry, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]audit.AuditEntry), args.Error(1)
}

// GetAuditLogsByResource mocks the GetAuditLogsByResource method
func (m *MockAuditService) GetAuditLogsByResource(ctx context.Context, resource, resourceID string, limit int) ([]audit.AuditEntry, error) {
	args := m.Called(ctx, resource, resourceID, limit)
	return args.Get(0).([]audit.AuditEntry), args.Error(1)
}
