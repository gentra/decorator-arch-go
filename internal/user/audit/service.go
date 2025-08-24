package audit

import (
	"context"
	"time"

	"github.com/gentra/decorator-arch-go/internal/audit"
	"github.com/gentra/decorator-arch-go/internal/user"
)

// service implements user.Service with audit logging capabilities
type service struct {
	next         user.Service
	auditService audit.Service
}

// NewService creates a new audit-enabled user service
func NewService(next user.Service, auditService audit.Service) user.Service {
	return &service{
		next:         next,
		auditService: auditService,
	}
}

// Register creates a new user with audit logging
func (s *service) Register(ctx context.Context, data user.RegisterData) (*user.User, error) {
	// Call next service
	result, err := s.next.Register(ctx, data)

	// Log audit entry
	s.logAuditEntry(ctx, "user.register", "user", result.ID.String(), map[string]interface{}{
		"email":      data.Email,
		"first_name": data.FirstName,
		"last_name":  data.LastName,
	}, err == nil, err)

	return result, err
}

// Login authenticates a user with audit logging
func (s *service) Login(ctx context.Context, email, password string) (*user.AuthResult, error) {
	// Call next service
	result, err := s.next.Login(ctx, email, password)

	// Log audit entry
	userID := ""
	if result != nil && result.User != nil {
		userID = result.User.ID.String()
	}

	s.logAuditEntry(ctx, "user.login", "user", userID, map[string]interface{}{
		"email": email,
	}, err == nil, err)

	return result, err
}

// GetByID retrieves a user by ID with audit logging
func (s *service) GetByID(ctx context.Context, id string) (*user.User, error) {
	// Call next service
	result, err := s.next.GetByID(ctx, id)

	// Log audit entry
	s.logAuditEntry(ctx, "user.get_by_id", "user", id, map[string]interface{}{
		"requested_user_id": id,
	}, err == nil, err)

	return result, err
}

// UpdateProfile updates user profile with audit logging
func (s *service) UpdateProfile(ctx context.Context, id string, data user.UpdateProfileData) (*user.User, error) {
	// Call next service
	result, err := s.next.UpdateProfile(ctx, id, data)

	// Log audit entry
	changes := make(map[string]interface{})
	if data.FirstName != nil {
		changes["first_name"] = *data.FirstName
	}
	if data.LastName != nil {
		changes["last_name"] = *data.LastName
	}
	if data.Email != nil {
		changes["email"] = *data.Email
	}

	s.logAuditEntry(ctx, "user.update_profile", "user", id, map[string]interface{}{
		"changes": changes,
	}, err == nil, err)

	return result, err
}

// GetPreferences retrieves user preferences with audit logging
func (s *service) GetPreferences(ctx context.Context, userID string) (*user.UserPreferences, error) {
	// Call next service
	result, err := s.next.GetPreferences(ctx, userID)

	// Log audit entry
	s.logAuditEntry(ctx, "user.get_preferences", "user_preferences", userID, map[string]interface{}{
		"requested_user_id": userID,
	}, err == nil, err)

	return result, err
}

// UpdatePreferences updates user preferences with audit logging
func (s *service) UpdatePreferences(ctx context.Context, userID string, prefs user.UserPreferences) error {
	// Call next service
	err := s.next.UpdatePreferences(ctx, userID, prefs)

	// Log audit entry
	s.logAuditEntry(ctx, "user.update_preferences", "user_preferences", userID, map[string]interface{}{
		"theme":    prefs.Theme,
		"language": prefs.Language,
		"timezone": prefs.Timezone,
	}, err == nil, err)

	return err
}

// logAuditEntry logs an audit entry with the provided information
func (s *service) logAuditEntry(ctx context.Context, action, resource, resourceID string, details interface{}, success bool, err error) {
	entry := audit.AuditEntry{
		Timestamp:  time.Now(),
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    details,
		Success:    success,
	}

	if err != nil {
		entry.Error = err.Error()
	}

	// Extract audit context information if available
	if auditCtx := extractAuditContext(ctx); auditCtx != nil {
		entry.UserID = auditCtx.UserID
		entry.IPAddress = auditCtx.IPAddress
		entry.UserAgent = auditCtx.UserAgent
		entry.SessionID = auditCtx.SessionID
	}

	// Log the entry using the audit domain service
	// Don't fail the operation if audit logging fails
	s.auditService.Log(ctx, entry)
}

// AuditContext contains audit-related information from the request context
type AuditContext struct {
	UserID    string
	IPAddress string
	UserAgent string
	SessionID string
}

// Context keys for audit information
type contextKey string

const (
	AuditContextKey contextKey = "audit_context"
)

// extractAuditContext extracts audit information from the context
func extractAuditContext(ctx context.Context) *AuditContext {
	if auditCtx, ok := ctx.Value(AuditContextKey).(AuditContext); ok {
		return &auditCtx
	}
	return nil
}

// WithAuditContext adds audit context information to the request context
func WithAuditContext(ctx context.Context, userID, ipAddress, userAgent, sessionID string) context.Context {
	auditCtx := AuditContext{
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		SessionID: sessionID,
	}

	return context.WithValue(ctx, AuditContextKey, auditCtx)
}
