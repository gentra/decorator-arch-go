package mock

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/gentra/decorator-arch-go/internal/notification"
)

// service implements notification.Service interface with mock operations for testing/development
type service struct {
	config notification.NotificationConfig
}

// NewService creates a new mock notification service
func NewService() notification.Service {
	return &service{
		config: notification.DefaultNotificationConfig(),
	}
}

// SendWelcomeEmail sends a welcome email (mock implementation)
func (s *service) SendWelcomeEmail(ctx context.Context, userEmail, userName string) error {
	log.Printf("MOCK NOTIFICATION: Welcome email sent to %s (%s)", userEmail, userName)
	return nil
}

// SendPasswordResetEmail sends a password reset email (mock implementation)
func (s *service) SendPasswordResetEmail(ctx context.Context, userEmail, resetToken string) error {
	log.Printf("MOCK NOTIFICATION: Password reset email sent to %s with token %s", userEmail, resetToken[:8]+"...")
	return nil
}

// SendProfileUpdateNotification sends a profile update notification (mock implementation)
func (s *service) SendProfileUpdateNotification(ctx context.Context, userID string, changes map[string]interface{}) error {
	log.Printf("MOCK NOTIFICATION: Profile update notification sent to user %s with changes: %+v", userID, changes)
	return nil
}

// SendVerificationEmail sends a verification email (mock implementation)
func (s *service) SendVerificationEmail(ctx context.Context, userEmail, verificationToken string) error {
	log.Printf("MOCK NOTIFICATION: Verification email sent to %s with token %s", userEmail, verificationToken[:8]+"...")
	return nil
}

// SendPushNotification sends a push notification (mock implementation)
func (s *service) SendPushNotification(ctx context.Context, userID string, notification notification.PushNotification) error {
	log.Printf("MOCK NOTIFICATION: Push notification sent to user %s: %s - %s", userID, notification.Title, notification.Body)
	return nil
}

// SendSMSNotification sends an SMS notification (mock implementation)
func (s *service) SendSMSNotification(ctx context.Context, phoneNumber string, message string) error {
	log.Printf("MOCK NOTIFICATION: SMS sent to %s: %s", phoneNumber, message)
	return nil
}

// SendBulkEmail sends bulk emails (mock implementation)
func (s *service) SendBulkEmail(ctx context.Context, emails []notification.EmailNotification) error {
	log.Printf("MOCK NOTIFICATION: Bulk email sent to %d recipients", len(emails))
	for _, email := range emails {
		log.Printf("  - Email to %s: %s", email.To, email.Subject)
	}
	return nil
}

// SendBulkPush sends bulk push notifications (mock implementation)
func (s *service) SendBulkPush(ctx context.Context, notifications []notification.PushNotification) error {
	log.Printf("MOCK NOTIFICATION: Bulk push notifications sent to %d users", len(notifications))
	for _, notif := range notifications {
		log.Printf("  - Push to %s: %s", notif.UserID, notif.Title)
	}
	return nil
}

// GetNotificationHistory returns notification history (mock implementation)
func (s *service) GetNotificationHistory(ctx context.Context, userID string, limit int) ([]notification.NotificationHistory, error) {
	// Return mock notification history
	history := []notification.NotificationHistory{
		{
			ID:        uuid.New().String(),
			UserID:    userID,
			Type:      notification.NotificationTypeEmail,
			Title:     "Welcome!",
			Body:      "Welcome to our platform",
			Status:    notification.NotificationStatusDelivered,
			Priority:  notification.PriorityNormal,
			CreatedAt: time.Now().Add(-time.Hour),
			SentAt:    timePtr(time.Now().Add(-time.Hour + time.Minute)),
			ReadAt:    timePtr(time.Now().Add(-time.Minute * 30)),
		},
		{
			ID:        uuid.New().String(),
			UserID:    userID,
			Type:      notification.NotificationTypePush,
			Title:     "Profile Updated",
			Body:      "Your profile has been successfully updated",
			Status:    notification.NotificationStatusRead,
			Priority:  notification.PriorityNormal,
			CreatedAt: time.Now().Add(-time.Hour * 2),
			SentAt:    timePtr(time.Now().Add(-time.Hour*2 + time.Minute)),
			ReadAt:    timePtr(time.Now().Add(-time.Hour)),
		},
	}

	if limit > 0 && len(history) > limit {
		history = history[:limit]
	}

	return history, nil
}

// MarkAsRead marks a notification as read (mock implementation)
func (s *service) MarkAsRead(ctx context.Context, notificationID string) error {
	log.Printf("MOCK NOTIFICATION: Notification %s marked as read", notificationID)
	return nil
}

// GetUnreadCount returns the count of unread notifications (mock implementation)
func (s *service) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	// Return a mock unread count
	return 3, nil
}

// Helper function to create time pointers
func timePtr(t time.Time) *time.Time {
	return &t
}
