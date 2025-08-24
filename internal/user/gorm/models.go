package gorm

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// UserModel represents the GORM model for users table
type UserModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	FirstName    string    `gorm:"not null" json:"first_name"`
	LastName     string    `gorm:"not null" json:"last_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	Preferences *UserPreferencesModel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"preferences,omitempty"`
}

// UserPreferencesModel represents the GORM model for user_preferences table
type UserPreferencesModel struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID             uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	EmailNotifications bool           `gorm:"default:true" json:"email_notifications"`
	PushNotifications  bool           `gorm:"default:true" json:"push_notifications"`
	SMSNotifications   bool           `gorm:"default:false" json:"sms_notifications"`
	Theme              string         `gorm:"default:light" json:"theme"`
	Language           string         `gorm:"default:en" json:"language"`
	Timezone           string         `gorm:"default:UTC" json:"timezone"`
	NotificationTypes  datatypes.JSON `json:"notification_types"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`

	// Relationships
	User *UserModel `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BeforeCreate will set a UUID rather than numeric ID for UserModel
func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// BeforeCreate will set a UUID rather than numeric ID for UserPreferencesModel
func (p *UserPreferencesModel) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// TableName overrides the table name used by UserModel to `users`
func (UserModel) TableName() string {
	return "users"
}

// TableName overrides the table name used by UserPreferencesModel to `user_preferences`
func (UserPreferencesModel) TableName() string {
	return "user_preferences"
}
