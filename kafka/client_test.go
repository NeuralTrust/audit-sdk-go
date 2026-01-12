package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProducer_SetsDefaults(t *testing.T) {
	cfg := &Config{
		Brokers:  []string{"localhost:9092"},
		ClientID: "test-client",
	}

	assert.Equal(t, 0, cfg.TopicNumParts)
	assert.Equal(t, 0, cfg.TopicReplication)
}

func TestBuildKafkaConfig_Basic(t *testing.T) {
	cfg := &Config{
		Brokers:      []string{"broker1:9092", "broker2:9092"},
		ClientID:     "test-client",
		RequiredAcks: 1,
		RetryMax:     3,
	}

	kafkaConfig := buildKafkaConfig(cfg)

	bootstrapServers, _ := kafkaConfig.Get("bootstrap.servers", "")
	assert.Equal(t, "broker1:9092,broker2:9092", bootstrapServers)

	clientID, _ := kafkaConfig.Get("client.id", "")
	assert.Equal(t, "test-client", clientID)

	acks, _ := kafkaConfig.Get("acks", 0)
	assert.Equal(t, 1, acks)

	retries, _ := kafkaConfig.Get("retries", 0)
	assert.Equal(t, 3, retries)

	reconnectBackoff, _ := kafkaConfig.Get("reconnect.backoff.ms", 0)
	assert.Equal(t, 100, reconnectBackoff)

	reconnectBackoffMax, _ := kafkaConfig.Get("reconnect.backoff.max.ms", 0)
	assert.Equal(t, 10000, reconnectBackoffMax)

	keepalive, _ := kafkaConfig.Get("socket.keepalive.enable", false)
	assert.Equal(t, true, keepalive)
}

func TestBuildKafkaConfig_WithSASL(t *testing.T) {
	cfg := &Config{
		Brokers:  []string{"localhost:9092"},
		ClientID: "test",
		SASL: &SASLConfig{
			Enable:    true,
			Mechanism: "PLAIN",
			Username:  "user",
			Password:  "pass",
		},
	}

	kafkaConfig := buildKafkaConfig(cfg)

	protocol, _ := kafkaConfig.Get("security.protocol", "")
	assert.Equal(t, "SASL_SSL", protocol)

	mechanism, _ := kafkaConfig.Get("sasl.mechanisms", "")
	assert.Equal(t, "PLAIN", mechanism)

	username, _ := kafkaConfig.Get("sasl.username", "")
	assert.Equal(t, "user", username)

	password, _ := kafkaConfig.Get("sasl.password", "")
	assert.Equal(t, "pass", password)
}

func TestBuildKafkaConfig_WithTLS(t *testing.T) {
	cfg := &Config{
		Brokers:  []string{"localhost:9092"},
		ClientID: "test",
		TLS: &TLSConfig{
			Enable:             true,
			CAFile:             "/path/to/ca.pem",
			CertFile:           "/path/to/cert.pem",
			KeyFile:            "/path/to/key.pem",
			InsecureSkipVerify: false,
		},
	}

	kafkaConfig := buildKafkaConfig(cfg)

	protocol, _ := kafkaConfig.Get("security.protocol", "")
	assert.Equal(t, "SSL", protocol)

	caLocation, _ := kafkaConfig.Get("ssl.ca.location", "")
	assert.Equal(t, "/path/to/ca.pem", caLocation)

	certLocation, _ := kafkaConfig.Get("ssl.certificate.location", "")
	assert.Equal(t, "/path/to/cert.pem", certLocation)

	keyLocation, _ := kafkaConfig.Get("ssl.key.location", "")
	assert.Equal(t, "/path/to/key.pem", keyLocation)

	verify, _ := kafkaConfig.Get("enable.ssl.certificate.verification", false)
	assert.Equal(t, true, verify)
}

func TestBuildKafkaConfig_WithTLSInsecure(t *testing.T) {
	cfg := &Config{
		Brokers:  []string{"localhost:9092"},
		ClientID: "test",
		TLS: &TLSConfig{
			Enable:             true,
			InsecureSkipVerify: true,
		},
	}

	kafkaConfig := buildKafkaConfig(cfg)

	verify, _ := kafkaConfig.Get("enable.ssl.certificate.verification", true)
	assert.Equal(t, false, verify)
}

func TestBuildKafkaConfig_SASLTakesPriorityOverTLS(t *testing.T) {
	cfg := &Config{
		Brokers:  []string{"localhost:9092"},
		ClientID: "test",
		SASL: &SASLConfig{
			Enable:    true,
			Mechanism: "PLAIN",
			Username:  "user",
			Password:  "pass",
		},
		TLS: &TLSConfig{
			Enable: true,
			CAFile: "/path/to/ca.pem",
		},
	}

	kafkaConfig := buildKafkaConfig(cfg)

	protocol, _ := kafkaConfig.Get("security.protocol", "")
	assert.Equal(t, "SASL_SSL", protocol)

	caLocation, _ := kafkaConfig.Get("ssl.ca.location", "")
	assert.Equal(t, "/path/to/ca.pem", caLocation)
}

func TestBuildKafkaConfig_DisabledSASL(t *testing.T) {
	cfg := &Config{
		Brokers:  []string{"localhost:9092"},
		ClientID: "test",
		SASL: &SASLConfig{
			Enable: false,
		},
	}

	kafkaConfig := buildKafkaConfig(cfg)

	protocol, _ := kafkaConfig.Get("security.protocol", "not_set")
	assert.Equal(t, "not_set", protocol)
}

func TestBuildKafkaConfig_DisabledTLS(t *testing.T) {
	cfg := &Config{
		Brokers:  []string{"localhost:9092"},
		ClientID: "test",
		TLS: &TLSConfig{
			Enable: false,
			CAFile: "/path/to/ca.pem",
		},
	}

	kafkaConfig := buildKafkaConfig(cfg)

	caLocation, _ := kafkaConfig.Get("ssl.ca.location", "not_set")
	assert.Equal(t, "not_set", caLocation)
}


