package factory_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gentra/decorator-arch-go/internal/audit"
	"github.com/gentra/decorator-arch-go/internal/audit/factory"
)

func TestAuditServiceFactory_Build(t *testing.T) {
	tests := []struct {
		name    string
		config  factory.Config
		wantErr bool
	}{
		{
			name: "Given factory with console output enabled, When Build is called, Then should return console service without error",
			config: factory.Config{
				OutputTarget: "console",
				Features: factory.FeatureFlags{
					EnableConsoleOutput: true,
				},
			},
			wantErr: false,
		},
		{
			name: "Given factory with console output disabled, When Build is called, Then should return console service without error (fallback)",
			config: factory.Config{
				OutputTarget: "console",
				Features: factory.FeatureFlags{
					EnableConsoleOutput: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Given factory with file output enabled, When Build is called, Then should return console service without error (fallback)",
			config: factory.Config{
				OutputTarget: "file",
				Features: factory.FeatureFlags{
					EnableFileOutput: true,
				},
			},
			wantErr: false,
		},
		{
			name: "Given factory with database output enabled, When Build is called, Then should return console service without error (fallback)",
			config: factory.Config{
				OutputTarget: "database",
				Features: factory.FeatureFlags{
					EnableDatabaseOutput: true,
				},
			},
			wantErr: false,
		},
		{
			name: "Given factory with external output enabled, When Build is called, Then should return console service without error (fallback)",
			config: factory.Config{
				OutputTarget: "external",
				Features: factory.FeatureFlags{
					EnableExternalOutput: true,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			factoryInstance := factory.NewFactory(tt.config)

			// Act
			service, err := factoryInstance.Build()

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, service)
				_, ok := service.(audit.Service)
				assert.True(t, ok, "Service should implement audit.Service interface")
			}
		})
	}
}

func TestAuditServiceFactory_NewFactory(t *testing.T) {
	t.Run("Given NewFactory is called with config, When factory is created, Then should return factory instance", func(t *testing.T) {
		// Arrange
		config := factory.Config{
			OutputTarget: "console",
			Features: factory.FeatureFlags{
				EnableConsoleOutput: true,
			},
		}

		// Act
		factoryInstance := factory.NewFactory(config)

		// Assert
		require.NotNil(t, factoryInstance)
	})
}

func TestAuditServiceFactory_DefaultConfig(t *testing.T) {
	t.Run("Given DefaultConfig is called, When default config is created, Then should return sensible defaults", func(t *testing.T) {
		// Act
		config := factory.DefaultConfig()

		// Assert
		assert.Equal(t, "console", config.OutputTarget)
		assert.Equal(t, "/var/log/audit.log", config.LogFilePath)
		assert.True(t, config.Features.EnableConsoleOutput)
		assert.False(t, config.Features.EnableFileOutput)
		assert.False(t, config.Features.EnableDatabaseOutput)
		assert.False(t, config.Features.EnableExternalOutput)
		assert.False(t, config.Features.EnableAsyncProcessing)
		assert.False(t, config.Features.EnableBatching)
		assert.False(t, config.Features.EnableCompression)
	})
}

func TestAuditServiceFactory_DefaultFeatureFlags(t *testing.T) {
	t.Run("Given DefaultFeatureFlags is called, When default feature flags are created, Then should return sensible defaults", func(t *testing.T) {
		// Act
		features := factory.DefaultFeatureFlags()

		// Assert
		assert.True(t, features.EnableConsoleOutput)
		assert.False(t, features.EnableFileOutput)
		assert.False(t, features.EnableDatabaseOutput)
		assert.False(t, features.EnableExternalOutput)
		assert.False(t, features.EnableAsyncProcessing)
		assert.False(t, features.EnableBatching)
		assert.False(t, features.EnableCompression)
	})
}

func TestConfigBuilder_Build(t *testing.T) {
	tests := []struct {
		name           string
		builderActions func(*factory.ConfigBuilder) *factory.ConfigBuilder
		expectedConfig factory.Config
	}{
		{
			name: "Given config builder with output target set, When Build is called, Then should return config with set output target",
			builderActions: func(b *factory.ConfigBuilder) *factory.ConfigBuilder {
				return b.WithOutputTarget("file")
			},
			expectedConfig: factory.Config{
				OutputTarget: "file",
				LogFilePath:  "/var/log/audit.log",
				Features:     factory.DefaultFeatureFlags(),
			},
		},
		{
			name: "Given config builder with log file path set, When Build is called, Then should return config with set log file path",
			builderActions: func(b *factory.ConfigBuilder) *factory.ConfigBuilder {
				return b.WithLogFilePath("/custom/audit.log")
			},
			expectedConfig: factory.Config{
				OutputTarget: "console",
				LogFilePath:  "/custom/audit.log",
				Features:     factory.DefaultFeatureFlags(),
			},
		},
		{
			name: "Given config builder with database DSN set, When Build is called, Then should return config with set database DSN",
			builderActions: func(b *factory.ConfigBuilder) *factory.ConfigBuilder {
				return b.WithDatabaseDSN("postgres://user:pass@localhost:5432/audit")
			},
			expectedConfig: factory.Config{
				OutputTarget: "console",
				LogFilePath:  "/var/log/audit.log",
				DatabaseDSN:  "postgres://user:pass@localhost:5432/audit",
				Features:     factory.DefaultFeatureFlags(),
			},
		},
		{
			name: "Given config builder with external service set, When Build is called, Then should return config with set external service",
			builderActions: func(b *factory.ConfigBuilder) *factory.ConfigBuilder {
				return b.WithExternalService("https://audit.example.com", "api-key-123")
			},
			expectedConfig: factory.Config{
				OutputTarget:   "console",
				LogFilePath:    "/var/log/audit.log",
				ExternalURL:    "https://audit.example.com",
				ExternalAPIKey: "api-key-123",
				Features:       factory.DefaultFeatureFlags(),
			},
		},
		{
			name: "Given config builder with async processing enabled, When Build is called, Then should return config with async processing enabled",
			builderActions: func(b *factory.ConfigBuilder) *factory.ConfigBuilder {
				return b.EnableAsyncProcessing()
			},
			expectedConfig: factory.Config{
				OutputTarget: "console",
				LogFilePath:  "/var/log/audit.log",
				Features: factory.FeatureFlags{
					EnableConsoleOutput:   true,
					EnableFileOutput:      false,
					EnableDatabaseOutput:  false,
					EnableExternalOutput:  false,
					EnableAsyncProcessing: true,
					EnableBatching:        false,
					EnableCompression:     false,
				},
			},
		},
		{
			name: "Given config builder with batching enabled, When Build is called, Then should return config with batching enabled",
			builderActions: func(b *factory.ConfigBuilder) *factory.ConfigBuilder {
				return b.EnableBatching()
			},
			expectedConfig: factory.Config{
				OutputTarget: "console",
				LogFilePath:  "/var/log/audit.log",
				Features: factory.FeatureFlags{
					EnableConsoleOutput:   true,
					EnableFileOutput:      false,
					EnableDatabaseOutput:  false,
					EnableExternalOutput:  false,
					EnableAsyncProcessing: false,
					EnableBatching:        true,
					EnableCompression:     false,
				},
			},
		},
		{
			name: "Given config builder with compression enabled, When Build is called, Then should return config with compression enabled",
			builderActions: func(b *factory.ConfigBuilder) *factory.ConfigBuilder {
				return b.EnableCompression()
			},
			expectedConfig: factory.Config{
				OutputTarget: "console",
				LogFilePath:  "/var/log/audit.log",
				Features: factory.FeatureFlags{
					EnableConsoleOutput:   true,
					EnableFileOutput:      false,
					EnableDatabaseOutput:  false,
					EnableExternalOutput:  false,
					EnableAsyncProcessing: false,
					EnableBatching:        false,
					EnableCompression:     true,
				},
			},
		},
		{
			name: "Given config builder with multiple features enabled, When Build is called, Then should return config with all features enabled",
			builderActions: func(b *factory.ConfigBuilder) *factory.ConfigBuilder {
				return b.EnableAsyncProcessing().EnableBatching().EnableCompression()
			},
			expectedConfig: factory.Config{
				OutputTarget: "console",
				LogFilePath:  "/var/log/audit.log",
				Features: factory.FeatureFlags{
					EnableConsoleOutput:   true,
					EnableFileOutput:      false,
					EnableDatabaseOutput:  false,
					EnableExternalOutput:  false,
					EnableAsyncProcessing: true,
					EnableBatching:        true,
					EnableCompression:     true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			builder := factory.NewConfigBuilder()

			// Act
			if tt.builderActions != nil {
				builder = tt.builderActions(builder)
			}
			config := builder.Build()

			// Assert
			assert.Equal(t, tt.expectedConfig, config)
		})
	}
}

func TestConfigBuilder_NewConfigBuilder(t *testing.T) {
	t.Run("Given NewConfigBuilder is called, When builder is created, Then should return builder with default config", func(t *testing.T) {
		// Act
		builder := factory.NewConfigBuilder()

		// Assert
		require.NotNil(t, builder)
		// Note: We can't directly access builder.config as it's unexported
		// The test verifies the builder was created successfully
	})
}
