package audit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_SetDefaults(t *testing.T) {
	cfg := &Config{}
	cfg.setDefaults()

	assert.Equal(t, DefaultAuditEventsTopic, cfg.AuditEventsTopic)
	assert.Equal(t, DefaultAuditLogsIngestTopic, cfg.AuditLogsIngestTopic)
	assert.Equal(t, "audit-sdk", cfg.ClientID)
	assert.Equal(t, 100, cfg.BatchSize)
	assert.Equal(t, 1*time.Second, cfg.BatchTimeout)
	assert.Equal(t, 3, cfg.RetryMax)
	assert.Equal(t, 100*time.Millisecond, cfg.RetryBackoff)
	assert.Equal(t, 1, cfg.RequiredAcks)
}

func TestConfig_SetDefaults_FromEnv(t *testing.T) {
	t.Setenv(EnvAuditEventsTopic, "custom_events")
	t.Setenv(EnvAuditLogsIngestTopic, "custom_logs")

	cfg := &Config{}
	cfg.setDefaults()

	assert.Equal(t, "custom_events", cfg.AuditEventsTopic)
	assert.Equal(t, "custom_logs", cfg.AuditLogsIngestTopic)
}

func TestConfig_SetDefaults_ConfigOverridesEnv(t *testing.T) {
	t.Setenv(EnvAuditEventsTopic, "env_events")

	cfg := &Config{
		AuditEventsTopic: "config_events",
	}
	cfg.setDefaults()

	assert.Equal(t, "config_events", cfg.AuditEventsTopic)
}

func TestConfig_Validate_NoBrokers(t *testing.T) {
	cfg := &Config{}
	err := cfg.validate()

	assert.Equal(t, ErrNoBrokers, err)
}

func TestConfig_Validate_Success(t *testing.T) {
	cfg := &Config{
		Brokers: []string{"localhost:9092"},
	}
	err := cfg.validate()

	assert.NoError(t, err)
}

func TestResolveValue(t *testing.T) {
	tests := []struct {
		name         string
		configValue  string
		envKey       string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "config value takes priority",
			configValue:  "from_config",
			envKey:       "TEST_KEY",
			envValue:     "from_env",
			defaultValue: "default",
			expected:     "from_config",
		},
		{
			name:         "env value when no config",
			configValue:  "",
			envKey:       "TEST_KEY_2",
			envValue:     "from_env",
			defaultValue: "default",
			expected:     "from_env",
		},
		{
			name:         "default when no config or env",
			configValue:  "",
			envKey:       "TEST_KEY_3",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv(tt.envKey, tt.envValue)
			}

			result := resolveValue(tt.configValue, tt.envKey, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}
