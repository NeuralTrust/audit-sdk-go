package audit

import (
	"os"
	"time"
)

const (
	DefaultAuditEventsTopic     = "audit_events"
	DefaultAuditLogsIngestTopic = "audit_logs_ingest"

	EnvAuditEventsTopic     = "AUDIT_EVENTS_TOPIC"
	EnvAuditLogsIngestTopic = "AUDIT_LOGS_INGEST_TOPIC"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type Config struct {
	Brokers              []string
	AuditEventsTopic     string
	AuditLogsIngestTopic string
	ClientID             string
	BatchSize            int
	BatchTimeout         time.Duration
	RetryMax             int
	RetryBackoff         time.Duration
	RequiredAcks         int
	TopicAutoCreate      bool
	TopicNumParts        int
	TopicReplication     int
	TLS                  *TLSConfig
	SASL                 *SASLConfig
	LogLevel             LogLevel
}

type TLSConfig struct {
	Enable             bool
	CertFile           string
	KeyFile            string
	CAFile             string
	InsecureSkipVerify bool
}

type SASLConfig struct {
	Enable    bool
	Mechanism string
	Username  string
	Password  string
}

func (c *Config) setDefaults() {
	c.AuditEventsTopic = resolveValue(c.AuditEventsTopic, EnvAuditEventsTopic, DefaultAuditEventsTopic)
	c.AuditLogsIngestTopic = resolveValue(c.AuditLogsIngestTopic, EnvAuditLogsIngestTopic, DefaultAuditLogsIngestTopic)

	if c.ClientID == "" {
		c.ClientID = "audit-sdk"
	}
	if c.BatchSize == 0 {
		c.BatchSize = 100
	}
	if c.BatchTimeout == 0 {
		c.BatchTimeout = 1 * time.Second
	}
	if c.RetryMax == 0 {
		c.RetryMax = 3
	}
	if c.RetryBackoff == 0 {
		c.RetryBackoff = 100 * time.Millisecond
	}
	if c.RequiredAcks == 0 {
		c.RequiredAcks = 1
	}
	if c.LogLevel == "" {
		c.LogLevel = LogLevelInfo
	}
}

func resolveValue(configValue, envKey, defaultValue string) string {
	if configValue != "" {
		return configValue
	}
	if envValue := os.Getenv(envKey); envValue != "" {
		return envValue
	}
	return defaultValue
}

func (c *Config) validate() error {
	if len(c.Brokers) == 0 {
		return ErrNoBrokers
	}
	return nil
}


