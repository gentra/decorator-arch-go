package console

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gentra/decorator-arch-go/internal/audit"
)

// service implements audit.Service interface using console/stdout logging
type service struct{}

// NewService creates a new console-based audit service
func NewService() audit.Service {
	return &service{}
}

// Log writes the audit entry to console/stdout
func (s *service) Log(ctx context.Context, entry audit.AuditEntry) error {
	entryJSON, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	log.Printf("AUDIT: %s", string(entryJSON))
	return nil
}

// GetAuditLogs retrieves audit logs based on filters (not implemented for console)
func (s *service) GetAuditLogs(ctx context.Context, filters audit.AuditFilters) ([]audit.AuditEntry, error) {
	// Console audit doesn't support retrieval
	return nil, nil
}

// GetAuditLogsByUser retrieves audit logs for a specific user (not implemented for console)
func (s *service) GetAuditLogsByUser(ctx context.Context, userID string, limit int) ([]audit.AuditEntry, error) {
	// Console audit doesn't support retrieval
	return nil, nil
}

// GetAuditLogsByResource retrieves audit logs for a specific resource (not implemented for console)
func (s *service) GetAuditLogsByResource(ctx context.Context, resource, resourceID string, limit int) ([]audit.AuditEntry, error) {
	// Console audit doesn't support retrieval
	return nil, nil
}
