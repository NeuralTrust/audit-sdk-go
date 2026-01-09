package kafka

import (
	"context"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Config struct {
	Brokers          []string
	ClientID         string
	RequiredAcks     int
	RetryMax         int
	TopicAutoCreate  bool
	TopicNumParts    int
	TopicReplication int
	TLS              *TLSConfig
	SASL             *SASLConfig
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

type Producer struct {
	kafkaProducer *kafka.Producer
	adminClient   *kafka.AdminClient
	config        *Config
}

func NewProducer(cfg *Config) (*Producer, error) {
	if cfg.TopicNumParts == 0 {
		cfg.TopicNumParts = 3
	}
	if cfg.TopicReplication == 0 {
		cfg.TopicReplication = 1
	}

	kafkaConfig := buildKafkaConfig(cfg)

	p, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		return nil, err
	}

	admin, err := kafka.NewAdminClientFromProducer(p)
	if err != nil {
		p.Close()
		return nil, err
	}

	return &Producer{
		kafkaProducer: p,
		adminClient:   admin,
		config:        cfg,
	}, nil
}

func buildKafkaConfig(cfg *Config) *kafka.ConfigMap {
	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(cfg.Brokers, ","),
		"client.id":         cfg.ClientID,
		"acks":              cfg.RequiredAcks,
		"retries":           cfg.RetryMax,

		"reconnect.backoff.ms":     100,
		"reconnect.backoff.max.ms": 10000,
		"socket.keepalive.enable":  true,
		"metadata.max.age.ms":      300000,
	}

	if cfg.SASL != nil && cfg.SASL.Enable {
		(*kafkaConfig)["security.protocol"] = "SASL_SSL"
		(*kafkaConfig)["sasl.mechanisms"] = cfg.SASL.Mechanism
		(*kafkaConfig)["sasl.username"] = cfg.SASL.Username
		(*kafkaConfig)["sasl.password"] = cfg.SASL.Password
	}

	if cfg.TLS != nil && cfg.TLS.Enable {
		if cfg.SASL == nil || !cfg.SASL.Enable {
			(*kafkaConfig)["security.protocol"] = "SSL"
		}
		if cfg.TLS.CAFile != "" {
			(*kafkaConfig)["ssl.ca.location"] = cfg.TLS.CAFile
		}
		if cfg.TLS.CertFile != "" {
			(*kafkaConfig)["ssl.certificate.location"] = cfg.TLS.CertFile
		}
		if cfg.TLS.KeyFile != "" {
			(*kafkaConfig)["ssl.key.location"] = cfg.TLS.KeyFile
		}
		(*kafkaConfig)["enable.ssl.certificate.verification"] = !cfg.TLS.InsecureSkipVerify
	}

	return kafkaConfig
}

func (p *Producer) EnsureTopics(topics []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var topicsToCreate []kafka.TopicSpecification

	for _, topic := range topics {
		metadata, err := p.adminClient.GetMetadata(&topic, false, 5000)
		if err != nil {
			topicsToCreate = append(topicsToCreate, kafka.TopicSpecification{
				Topic:             topic,
				NumPartitions:     p.config.TopicNumParts,
				ReplicationFactor: p.config.TopicReplication,
			})
			continue
		}

		if _, exists := metadata.Topics[topic]; !exists {
			topicsToCreate = append(topicsToCreate, kafka.TopicSpecification{
				Topic:             topic,
				NumPartitions:     p.config.TopicNumParts,
				ReplicationFactor: p.config.TopicReplication,
			})
		}
	}

	if len(topicsToCreate) == 0 {
		return nil
	}

	results, err := p.adminClient.CreateTopics(ctx, topicsToCreate)
	if err != nil {
		return err
	}

	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			return result.Error
		}
	}

	return nil
}

func (p *Producer) ProduceAsync(topics []string, key, value []byte) {
	for _, topic := range topics {
		t := topic
		msg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &t, Partition: kafka.PartitionAny},
			Key:            key,
			Value:          value,
		}
		_ = p.kafkaProducer.Produce(msg, nil)
	}
}

func (p *Producer) Close() error {
	p.kafkaProducer.Flush(5000)
	p.adminClient.Close()
	p.kafkaProducer.Close()
	return nil
}

