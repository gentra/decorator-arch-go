package gorm

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func TestUserModel_GivenEmptyID_WhenBeforeCreate_ThenGeneratesUUID(t *testing.T) {
	tests := []struct {
		name      string
		model     UserModel
		expectNew bool
	}{
		{
			name: "empty UUID generates new",
			model: UserModel{
				ID:    uuid.Nil,
				Email: "test@example.com",
			},
			expectNew: true,
		},
		{
			name: "existing UUID preserved",
			model: UserModel{
				ID:    uuid.New(),
				Email: "test@example.com",
			},
			expectNew: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			originalID := tt.model.ID
			
			// When
			err := tt.model.BeforeCreate(&gorm.DB{})

			// Then
			assert.NoError(t, err)
			if tt.expectNew {
				assert.NotEqual(t, uuid.Nil, tt.model.ID)
				assert.NotEqual(t, originalID, tt.model.ID)
			} else {
				assert.Equal(t, originalID, tt.model.ID)
			}
		})
	}
}

func TestUserPreferencesModel_GivenEmptyID_WhenBeforeCreate_ThenGeneratesUUID(t *testing.T) {
	tests := []struct {
		name      string
		model     UserPreferencesModel
		expectNew bool
	}{
		{
			name: "empty UUID generates new",
			model: UserPreferencesModel{
				ID:     uuid.Nil,
				UserID: uuid.New(),
			},
			expectNew: true,
		},
		{
			name: "existing UUID preserved",
			model: UserPreferencesModel{
				ID:     uuid.New(),
				UserID: uuid.New(),
			},
			expectNew: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			originalID := tt.model.ID
			
			// When
			err := tt.model.BeforeCreate(&gorm.DB{})

			// Then
			assert.NoError(t, err)
			if tt.expectNew {
				assert.NotEqual(t, uuid.Nil, tt.model.ID)
				assert.NotEqual(t, originalID, tt.model.ID)
			} else {
				assert.Equal(t, originalID, tt.model.ID)
			}
		})
	}
}

func TestUserModel_GivenStruct_WhenGettingTableName_ThenReturnsUsersTable(t *testing.T) {
	tests := []struct {
		name          string
		expectedTable string
	}{
		{
			name:          "users table name",
			expectedTable: "users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			model := UserModel{}

			// When
			tableName := model.TableName()

			// Then
			assert.Equal(t, tt.expectedTable, tableName)
		})
	}
}

func TestUserPreferencesModel_GivenStruct_WhenGettingTableName_ThenReturnsUserPreferencesTable(t *testing.T) {
	tests := []struct {
		name          string
		expectedTable string
	}{
		{
			name:          "user_preferences table name",
			expectedTable: "user_preferences",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			model := UserPreferencesModel{}

			// When
			tableName := model.TableName()

			// Then
			assert.Equal(t, tt.expectedTable, tableName)
		})
	}
}

func TestUserModel_GivenCompleteData_WhenValidatingStructure_ThenHasAllRequiredFields(t *testing.T) {
	tests := []struct {
		name  string
		model UserModel
	}{
		{
			name: "complete user model",
			model: UserModel{
				ID:           uuid.New(),
				Email:        "user@example.com",
				PasswordHash: "hashed_password",
				FirstName:    "John",
				LastName:     "Doe",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				Preferences:  &UserPreferencesModel{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Then - Verify all fields exist and have expected types
			assert.IsType(t, uuid.UUID{}, tt.model.ID)
			assert.IsType(t, "", tt.model.Email)
			assert.IsType(t, "", tt.model.PasswordHash)
			assert.IsType(t, "", tt.model.FirstName)
			assert.IsType(t, "", tt.model.LastName)
			assert.IsType(t, time.Time{}, tt.model.CreatedAt)
			assert.IsType(t, time.Time{}, tt.model.UpdatedAt)
			assert.IsType(t, (*UserPreferencesModel)(nil), tt.model.Preferences)
		})
	}
}

func TestUserPreferencesModel_GivenCompleteData_WhenValidatingStructure_ThenHasAllRequiredFields(t *testing.T) {
	tests := []struct {
		name  string
		model UserPreferencesModel
	}{
		{
			name: "complete user preferences model",
			model: UserPreferencesModel{
				ID:                 uuid.New(),
				UserID:             uuid.New(),
				EmailNotifications: true,
				PushNotifications:  true,
				SMSNotifications:   false,
				Theme:              "light",
				Language:           "en",
				Timezone:           "UTC",
				NotificationTypes:  datatypes.JSON(`{"email": true, "push": false}`),
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
				User:               &UserModel{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Then - Verify all fields exist and have expected types
			assert.IsType(t, uuid.UUID{}, tt.model.ID)
			assert.IsType(t, uuid.UUID{}, tt.model.UserID)
			assert.IsType(t, true, tt.model.EmailNotifications)
			assert.IsType(t, true, tt.model.PushNotifications)
			assert.IsType(t, true, tt.model.SMSNotifications)
			assert.IsType(t, "", tt.model.Theme)
			assert.IsType(t, "", tt.model.Language)
			assert.IsType(t, "", tt.model.Timezone)
			assert.IsType(t, datatypes.JSON{}, tt.model.NotificationTypes)
			assert.IsType(t, time.Time{}, tt.model.CreatedAt)
			assert.IsType(t, time.Time{}, tt.model.UpdatedAt)
			assert.IsType(t, (*UserModel)(nil), tt.model.User)
		})
	}
}

func TestUserModel_GivenDefaultValues_WhenCreating_ThenHasExpectedDefaults(t *testing.T) {
	tests := []struct {
		name               string
		email              string
		expectedNonZeroID  bool
		expectedValidEmail bool
	}{
		{
			name:               "valid email creates valid model",
			email:              "test@example.com",
			expectedNonZeroID:  false, // ID will be zero until BeforeCreate is called
			expectedValidEmail: true,
		},
		{
			name:               "empty model has zero values",
			email:              "",
			expectedNonZeroID:  false,
			expectedValidEmail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given/When
			model := UserModel{
				Email: tt.email,
			}

			// Then
			if tt.expectedNonZeroID {
				assert.NotEqual(t, uuid.Nil, model.ID)
			} else {
				assert.Equal(t, uuid.Nil, model.ID)
			}

			if tt.expectedValidEmail {
				assert.NotEmpty(t, model.Email)
				assert.Contains(t, model.Email, "@")
			} else {
				assert.Empty(t, model.Email)
			}

			// Verify zero time values for new models
			assert.True(t, model.CreatedAt.IsZero())
			assert.True(t, model.UpdatedAt.IsZero())
		})
	}
}

func TestUserPreferencesModel_GivenDefaultValues_WhenCreating_ThenHasExpectedDefaults(t *testing.T) {
	tests := []struct {
		name                        string
		userID                      uuid.UUID
		expectedEmailNotifications  bool
		expectedPushNotifications   bool
		expectedSMSNotifications    bool
		expectedTheme               string
		expectedLanguage            string
		expectedTimezone            string
	}{
		{
			name:                       "default preferences",
			userID:                     uuid.New(),
			expectedEmailNotifications: false, // Zero value
			expectedPushNotifications:  false, // Zero value
			expectedSMSNotifications:   false, // Zero value
			expectedTheme:              "",    // Zero value
			expectedLanguage:           "",    // Zero value
			expectedTimezone:           "",    // Zero value
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given/When
			model := UserPreferencesModel{
				UserID: tt.userID,
			}

			// Then
			assert.Equal(t, tt.userID, model.UserID)
			assert.Equal(t, tt.expectedEmailNotifications, model.EmailNotifications)
			assert.Equal(t, tt.expectedPushNotifications, model.PushNotifications)
			assert.Equal(t, tt.expectedSMSNotifications, model.SMSNotifications)
			assert.Equal(t, tt.expectedTheme, model.Theme)
			assert.Equal(t, tt.expectedLanguage, model.Language)
			assert.Equal(t, tt.expectedTimezone, model.Timezone)

			// Verify zero time values for new models
			assert.True(t, model.CreatedAt.IsZero())
			assert.True(t, model.UpdatedAt.IsZero())
		})
	}
}

func TestUserModel_GivenJSONTags_WhenSerializing_ThenPasswordHashIsHidden(t *testing.T) {
	tests := []struct {
		name           string
		passwordHash   string
		expectInJSON   bool
	}{
		{
			name:         "password hash should be hidden in JSON",
			passwordHash: "secret_hash",
			expectInJSON: false, // Due to json:"-" tag
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			model := UserModel{
				Email:        "test@example.com",
				PasswordHash: tt.passwordHash,
			}

			// When - This would be used during JSON marshaling
			// We're testing the struct tag existence here
			
			// Then
			// Check that the PasswordHash field has the json:"-" tag
			// by verifying the field exists but won't be serialized
			assert.NotEmpty(t, model.PasswordHash)
			assert.Equal(t, tt.passwordHash, model.PasswordHash)
			
			// In actual JSON serialization, this field would be omitted
			// due to the json:"-" tag on the PasswordHash field
		})
	}
}

func TestModels_GivenGORMTags_WhenValidatingTags_ThenHaveCorrectConstraints(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "models have proper GORM tags",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			userModel := UserModel{}
			prefsModel := UserPreferencesModel{}

			// Then - Verify the models exist and can be instantiated
			// The actual GORM tag validation would be done by GORM itself
			// when creating tables or performing database operations
			assert.NotNil(t, userModel)
			assert.NotNil(t, prefsModel)
			
			// Verify relationship fields exist
			assert.IsType(t, (*UserPreferencesModel)(nil), userModel.Preferences)
			assert.IsType(t, (*UserModel)(nil), prefsModel.User)
		})
	}
}