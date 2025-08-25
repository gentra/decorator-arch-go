package notification_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gentra/decorator-arch-go/internal/notification"
)

func TestEmailNotification_IsValid(t *testing.T) {
	tests := []struct {
		name         string
		emailNotif   notification.EmailNotification
		expected     bool
	}{
		{
			name: "Given email notification with to, subject and body, When IsValid is called, Then should return true",
			emailNotif: notification.EmailNotification{
				To:      "test@example.com",
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			expected: true,
		},
		{
			name: "Given email notification with to, subject and HTML body, When IsValid is called, Then should return true",
			emailNotif: notification.EmailNotification{
				To:       "test@example.com",
				Subject:  "Test Subject",
				BodyHTML: "<p>Test HTML Body</p>",
			},
			expected: true,
		},
		{
			name: "Given email notification with empty to, When IsValid is called, Then should return false",
			emailNotif: notification.EmailNotification{
				To:      "",
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			expected: false,
		},
		{
			name: "Given email notification with empty subject, When IsValid is called, Then should return false",
			emailNotif: notification.EmailNotification{
				To:      "test@example.com",
				Subject: "",
				Body:    "Test Body",
			},
			expected: false,
		},
		{
			name: "Given email notification with empty body and HTML body, When IsValid is called, Then should return false",
			emailNotif: notification.EmailNotification{
				To:       "test@example.com",
				Subject:  "Test Subject",
				Body:     "",
				BodyHTML: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.emailNotif.IsValid()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEmailNotification_HasAttachments(t *testing.T) {
	tests := []struct {
		name       string
		emailNotif notification.EmailNotification
		expected   bool
	}{
		{
			name: "Given email notification with attachments, When HasAttachments is called, Then should return true",
			emailNotif: notification.EmailNotification{
				Attachments: []notification.Attachment{
					{Filename: "test.pdf", ContentType: "application/pdf"},
				},
			},
			expected: true,
		},
		{
			name: "Given email notification with empty attachments, When HasAttachments is called, Then should return false",
			emailNotif: notification.EmailNotification{
				Attachments: []notification.Attachment{},
			},
			expected: false,
		},
		{
			name: "Given email notification with nil attachments, When HasAttachments is called, Then should return false",
			emailNotif: notification.EmailNotification{
				Attachments: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.emailNotif.HasAttachments()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEmailNotification_IsScheduled(t *testing.T) {
	tests := []struct {
		name       string
		emailNotif notification.EmailNotification
		expected   bool
	}{
		{
			name: "Given email notification with future scheduled time, When IsScheduled is called, Then should return true",
			emailNotif: notification.EmailNotification{
				ScheduledAt: timePtr(time.Now().Add(time.Hour)),
			},
			expected: true,
		},
		{
			name: "Given email notification with past scheduled time, When IsScheduled is called, Then should return false",
			emailNotif: notification.EmailNotification{
				ScheduledAt: timePtr(time.Now().Add(-time.Hour)),
			},
			expected: false,
		},
		{
			name: "Given email notification with nil scheduled time, When IsScheduled is called, Then should return false",
			emailNotif: notification.EmailNotification{
				ScheduledAt: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.emailNotif.IsScheduled()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPushNotification_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		pushNotif notification.PushNotification
		expected  bool
	}{
		{
			name: "Given push notification with user ID and title, When IsValid is called, Then should return true",
			pushNotif: notification.PushNotification{
				UserID: "user-123",
				Title:  "Test Notification",
			},
			expected: true,
		},
		{
			name: "Given push notification with empty user ID, When IsValid is called, Then should return false",
			pushNotif: notification.PushNotification{
				UserID: "",
				Title:  "Test Notification",
			},
			expected: false,
		},
		{
			name: "Given push notification with empty title, When IsValid is called, Then should return false",
			pushNotif: notification.PushNotification{
				UserID: "user-123",
				Title:  "",
			},
			expected: false,
		},
		{
			name: "Given push notification with both user ID and title empty, When IsValid is called, Then should return false",
			pushNotif: notification.PushNotification{
				UserID: "",
				Title:  "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.pushNotif.IsValid()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPushNotification_HasData(t *testing.T) {
	tests := []struct {
		name      string
		pushNotif notification.PushNotification
		expected  bool
	}{
		{
			name: "Given push notification with data, When HasData is called, Then should return true",
			pushNotif: notification.PushNotification{
				Data: map[string]interface{}{
					"action": "open_app",
				},
			},
			expected: true,
		},
		{
			name: "Given push notification with empty data, When HasData is called, Then should return false",
			pushNotif: notification.PushNotification{
				Data: map[string]interface{}{},
			},
			expected: false,
		},
		{
			name: "Given push notification with nil data, When HasData is called, Then should return false",
			pushNotif: notification.PushNotification{
				Data: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.pushNotif.HasData()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNotificationHistory_IsRead(t *testing.T) {
	tests := []struct {
		name    string
		history notification.NotificationHistory
		expected bool
	}{
		{
			name: "Given notification history with read time, When IsRead is called, Then should return true",
			history: notification.NotificationHistory{
				ReadAt: timePtr(time.Now()),
			},
			expected: true,
		},
		{
			name: "Given notification history with nil read time, When IsRead is called, Then should return false",
			history: notification.NotificationHistory{
				ReadAt: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.history.IsRead()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNotificationHistory_IsSent(t *testing.T) {
	tests := []struct {
		name    string
		history notification.NotificationHistory
		expected bool
	}{
		{
			name: "Given notification history with sent time, When IsSent is called, Then should return true",
			history: notification.NotificationHistory{
				SentAt: timePtr(time.Now()),
			},
			expected: true,
		},
		{
			name: "Given notification history with nil sent time, When IsSent is called, Then should return false",
			history: notification.NotificationHistory{
				SentAt: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.history.IsSent()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNotificationHistory_HasFailed(t *testing.T) {
	tests := []struct {
		name    string
		history notification.NotificationHistory
		expected bool
	}{
		{
			name: "Given notification history with failed status, When HasFailed is called, Then should return true",
			history: notification.NotificationHistory{
				Status: notification.NotificationStatusFailed,
			},
			expected: true,
		},
		{
			name: "Given notification history with sent status, When HasFailed is called, Then should return false",
			history: notification.NotificationHistory{
				Status: notification.NotificationStatusSent,
			},
			expected: false,
		},
		{
			name: "Given notification history with pending status, When HasFailed is called, Then should return false",
			history: notification.NotificationHistory{
				Status: notification.NotificationStatusPending,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.history.HasFailed()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNotificationConfig_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		config   notification.NotificationConfig
		expected bool
	}{
		{
			name: "Given notification config with default from email, When IsValid is called, Then should return true",
			config: notification.NotificationConfig{
				DefaultFromEmail: "noreply@example.com",
			},
			expected: true,
		},
		{
			name: "Given notification config with empty default from email, When IsValid is called, Then should return false",
			config: notification.NotificationConfig{
				DefaultFromEmail: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.config.IsValid()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultNotificationConfig(t *testing.T) {
	t.Run("Given default notification config call, When DefaultNotificationConfig is called, Then should return valid default configuration", func(t *testing.T) {
		// Act
		config := notification.DefaultNotificationConfig()

		// Assert
		assert.Equal(t, "smtp", config.EmailProvider)
		assert.Equal(t, "firebase", config.PushProvider)
		assert.Equal(t, "twilio", config.SMSProvider)
		assert.Equal(t, "noreply@example.com", config.DefaultFromEmail)
		assert.NotNil(t, config.Templates)
		assert.NotNil(t, config.RateLimits)
		
		// Check rate limits
		emailRL, exists := config.RateLimits["email"]
		assert.True(t, exists)
		assert.Equal(t, 60, emailRL.MaxPerMinute)
		assert.Equal(t, 1000, emailRL.MaxPerHour)
		assert.Equal(t, 10000, emailRL.MaxPerDay)
		
		pushRL, exists := config.RateLimits["push"]
		assert.True(t, exists)
		assert.Equal(t, 100, pushRL.MaxPerMinute)
		assert.Equal(t, 5000, pushRL.MaxPerHour)
		assert.Equal(t, 50000, pushRL.MaxPerDay)
		
		smsRL, exists := config.RateLimits["sms"]
		assert.True(t, exists)
		assert.Equal(t, 10, smsRL.MaxPerMinute)
		assert.Equal(t, 100, smsRL.MaxPerHour)
		assert.Equal(t, 500, smsRL.MaxPerDay)
		
		// Check retry config
		assert.Equal(t, 3, config.RetryConfig.MaxRetries)
		assert.Equal(t, time.Second, config.RetryConfig.InitialDelay)
		assert.Equal(t, 2.0, config.RetryConfig.BackoffFactor)
		assert.Equal(t, time.Minute*5, config.RetryConfig.MaxDelay)
	})
}

func TestNotificationTypes(t *testing.T) {
	tests := []struct {
		name         string
		notifType    notification.NotificationType
		expectedStr  string
	}{
		{
			name:         "Given email notification type, When accessing string value, Then should have correct value",
			notifType:    notification.NotificationTypeEmail,
			expectedStr:  "email",
		},
		{
			name:         "Given push notification type, When accessing string value, Then should have correct value",
			notifType:    notification.NotificationTypePush,
			expectedStr:  "push",
		},
		{
			name:         "Given SMS notification type, When accessing string value, Then should have correct value",
			notifType:    notification.NotificationTypeSMS,
			expectedStr:  "sms",
		},
		{
			name:         "Given in-app notification type, When accessing string value, Then should have correct value",
			notifType:    notification.NotificationTypeInApp,
			expectedStr:  "in_app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assert
			assert.Equal(t, tt.expectedStr, string(tt.notifType))
		})
	}
}

func TestNotificationStatuses(t *testing.T) {
	tests := []struct {
		name        string
		status      notification.NotificationStatus
		expectedStr string
	}{
		{
			name:        "Given pending status, When accessing string value, Then should have correct value",
			status:      notification.NotificationStatusPending,
			expectedStr: "pending",
		},
		{
			name:        "Given sent status, When accessing string value, Then should have correct value",
			status:      notification.NotificationStatusSent,
			expectedStr: "sent",
		},
		{
			name:        "Given delivered status, When accessing string value, Then should have correct value",
			status:      notification.NotificationStatusDelivered,
			expectedStr: "delivered",
		},
		{
			name:        "Given failed status, When accessing string value, Then should have correct value",
			status:      notification.NotificationStatusFailed,
			expectedStr: "failed",
		},
		{
			name:        "Given read status, When accessing string value, Then should have correct value",
			status:      notification.NotificationStatusRead,
			expectedStr: "read",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assert
			assert.Equal(t, tt.expectedStr, string(tt.status))
		})
	}
}

func TestPriorities(t *testing.T) {
	tests := []struct {
		name        string
		priority    notification.Priority
		expectedStr string
	}{
		{
			name:        "Given low priority, When accessing string value, Then should have correct value",
			priority:    notification.PriorityLow,
			expectedStr: "low",
		},
		{
			name:        "Given normal priority, When accessing string value, Then should have correct value",
			priority:    notification.PriorityNormal,
			expectedStr: "normal",
		},
		{
			name:        "Given high priority, When accessing string value, Then should have correct value",
			priority:    notification.PriorityHigh,
			expectedStr: "high",
		},
		{
			name:        "Given urgent priority, When accessing string value, Then should have correct value",
			priority:    notification.PriorityUrgent,
			expectedStr: "urgent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assert
			assert.Equal(t, tt.expectedStr, string(tt.priority))
		})
	}
}

func TestAttachment_Structure(t *testing.T) {
	t.Run("Given attachment with all fields, When accessing fields, Then should have correct structure", func(t *testing.T) {
		// Arrange
		content := []byte("test content")
		attachment := notification.Attachment{
			Filename:    "test.pdf",
			ContentType: "application/pdf",
			Content:     content,
			Size:        int64(len(content)),
		}

		// Assert
		assert.Equal(t, "test.pdf", attachment.Filename)
		assert.Equal(t, "application/pdf", attachment.ContentType)
		assert.Equal(t, content, attachment.Content)
		assert.Equal(t, int64(len(content)), attachment.Size)
	})
}

// Helper function for creating time pointers
func timePtr(t time.Time) *time.Time {
	return &t
}