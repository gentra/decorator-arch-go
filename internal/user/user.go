package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service defines the user domain interface
type Service interface {
	Register(ctx context.Context, data RegisterData) (*User, error)
	Login(ctx context.Context, email, password string) (*AuthResult, error)
	GetByID(ctx context.Context, id string) (*User, error)
	UpdateProfile(ctx context.Context, id string, data UpdateProfileData) (*User, error)
	GetPreferences(ctx context.Context, userID string) (*UserPreferences, error)
	UpdatePreferences(ctx context.Context, userID string, prefs UserPreferences) error
}

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RegisterData contains data for user registration
type RegisterData struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
}

// UpdateProfileData contains data for profile updates
type UpdateProfileData struct {
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=2"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=2"`
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
}

// AuthResult contains authentication result data
type AuthResult struct {
	User         *User     `json:"user"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// UserPreferences contains user notification and system preferences
type UserPreferences struct {
	ID                 uuid.UUID       `json:"id"`
	UserID             uuid.UUID       `json:"user_id"`
	EmailNotifications bool            `json:"email_notifications"`
	PushNotifications  bool            `json:"push_notifications"`
	SMSNotifications   bool            `json:"sms_notifications"`
	Theme              string          `json:"theme"` // light, dark, auto
	Language           string          `json:"language"`
	Timezone           string          `json:"timezone"`
	NotificationTypes  map[string]bool `json:"notification_types"` // task_assigned, project_updated, etc.
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

// UserError represents domain-specific user errors
type UserError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func (e UserError) Error() string {
	return e.Message
}

// Common user error codes
var (
	ErrUserNotFound        = UserError{Code: "USER_NOT_FOUND", Message: "User not found"}
	ErrEmailAlreadyExists  = UserError{Code: "EMAIL_EXISTS", Message: "Email already exists"}
	ErrInvalidCredentials  = UserError{Code: "INVALID_CREDENTIALS", Message: "Invalid email or password"}
	ErrInvalidEmail        = UserError{Code: "INVALID_EMAIL", Message: "Invalid email format"}
	ErrWeakPassword        = UserError{Code: "WEAK_PASSWORD", Message: "Password must be at least 8 characters"}
	ErrEmptyFirstName      = UserError{Code: "EMPTY_FIRST_NAME", Message: "First name is required"}
	ErrEmptyLastName       = UserError{Code: "EMPTY_LAST_NAME", Message: "Last name is required"}
	ErrPreferencesNotFound = UserError{Code: "PREFERENCES_NOT_FOUND", Message: "User preferences not found"}
)

// Helper methods for User
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) IsEmailVerified() bool {
	// This would typically check an email verification status
	// For now, we'll assume all users are verified
	return true
}

// Helper methods for UserPreferences
func (p *UserPreferences) IsNotificationEnabled(notificationType string) bool {
	if p.NotificationTypes == nil {
		return false
	}
	enabled, exists := p.NotificationTypes[notificationType]
	return exists && enabled
}

func (p *UserPreferences) EnableNotification(notificationType string) {
	if p.NotificationTypes == nil {
		p.NotificationTypes = make(map[string]bool)
	}
	p.NotificationTypes[notificationType] = true
}

func (p *UserPreferences) DisableNotification(notificationType string) {
	if p.NotificationTypes == nil {
		return
	}
	p.NotificationTypes[notificationType] = false
}

// DefaultUserPreferences returns default preferences for a new user
func DefaultUserPreferences(userID uuid.UUID) *UserPreferences {
	return &UserPreferences{
		ID:                 uuid.New(),
		UserID:             userID,
		EmailNotifications: true,
		PushNotifications:  true,
		SMSNotifications:   false,
		Theme:              "light",
		Language:           "en",
		Timezone:           "UTC",
		NotificationTypes: map[string]bool{
			"task_assigned":   true,
			"task_due_soon":   true,
			"project_updated": true,
			"project_invite":  true,
			"system_updates":  false,
			"marketing":       false,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
