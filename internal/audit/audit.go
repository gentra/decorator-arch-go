package audit

import (
	"context"
	"time"
)

// Service defines the audit domain interface - the ONLY interface in this domain
type Service interface {
	Log(ctx context.Context, entry AuditEntry) error
	GetAuditLogs(ctx context.Context, filters AuditFilters) ([]AuditEntry, error)
	GetAuditLogsByUser(ctx context.Context, userID string, limit int) ([]AuditEntry, error)
	GetAuditLogsByResource(ctx context.Context, resource, resourceID string, limit int) ([]AuditEntry, error)
}

// Domain types and data structures

// AuditEntry represents an audit log entry
type AuditEntry struct {
	ID         string      `json:"id"`
	Timestamp  time.Time   `json:"timestamp"`
	UserID     string      `json:"user_id,omitempty"`
	Action     string      `json:"action"`
	Resource   string      `json:"resource"`
	ResourceID string      `json:"resource_id,omitempty"`
	Details    interface{} `json:"details,omitempty"`
	Success    bool        `json:"success"`
	Error      string      `json:"error,omitempty"`
	IPAddress  string      `json:"ip_address,omitempty"`
	UserAgent  string      `json:"user_agent,omitempty"`
	SessionID  string      `json:"session_id,omitempty"`
}

// AuditFilters for querying audit logs
type AuditFilters struct {
	UserID     string     `json:"user_id,omitempty"`
	Action     string     `json:"action,omitempty"`
	Resource   string     `json:"resource,omitempty"`
	ResourceID string     `json:"resource_id,omitempty"`
	Success    *bool      `json:"success,omitempty"`
	StartTime  *time.Time `json:"start_time,omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	Limit      int        `json:"limit,omitempty"`
	Offset     int        `json:"offset,omitempty"`
}

// AuditContext contains audit-related information from the request context
type AuditContext struct {
	CurrentUserID string
	IPAddress     string
	UserAgent     string
	SessionID     string
}

// Context keys for audit information
type contextKey string

const (
	AuditContextKey contextKey = "audit_context"
)

// Helper methods for AuditEntry
func (e *AuditEntry) IsValid() bool {
	return e.Action != "" && e.Resource != "" && !e.Timestamp.IsZero()
}

func (e *AuditEntry) SetSuccess() {
	e.Success = true
	e.Error = ""
}

func (e *AuditEntry) SetError(err error) {
	e.Success = false
	if err != nil {
		e.Error = err.Error()
	}
}

// Helper methods for AuditContext
func (ctx AuditContext) IsValid() bool {
	return ctx.CurrentUserID != "" || ctx.IPAddress != ""
}

// Helper functions for context management

// WithAuditContext adds audit context information to the request context
func WithAuditContext(ctx context.Context, userID, ipAddress, userAgent, sessionID string) context.Context {
	auditCtx := AuditContext{
		CurrentUserID: userID,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		SessionID:     sessionID,
	}

	return context.WithValue(ctx, AuditContextKey, auditCtx)
}

// ExtractAuditContext extracts audit information from the context
func ExtractAuditContext(ctx context.Context) AuditContext {
	if auditCtx, ok := ctx.Value(AuditContextKey).(AuditContext); ok {
		return auditCtx
	}

	// Return empty context if not found
	return AuditContext{}
}
