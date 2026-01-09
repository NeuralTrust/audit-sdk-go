package audit

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/NeuralTrust/audit-sdk-go/kafka"
	"github.com/google/uuid"
)

type Client interface {
	Emit(event Event)
	Close() error
}

type client struct {
	config   *Config
	producer Producer
	topics   []string
	closed   bool
	mu       sync.RWMutex
}

func New(cfg *Config) (Client, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	cfg.setDefaults()

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	kafkaCfg := &kafka.Config{
		Brokers:          cfg.Brokers,
		ClientID:         cfg.ClientID,
		RequiredAcks:     cfg.RequiredAcks,
		RetryMax:         cfg.RetryMax,
		TopicAutoCreate:  cfg.TopicAutoCreate,
		TopicNumParts:    cfg.TopicNumParts,
		TopicReplication: cfg.TopicReplication,
	}

	if cfg.TLS != nil {
		kafkaCfg.TLS = &kafka.TLSConfig{
			Enable:             cfg.TLS.Enable,
			CertFile:           cfg.TLS.CertFile,
			KeyFile:            cfg.TLS.KeyFile,
			CAFile:             cfg.TLS.CAFile,
			InsecureSkipVerify: cfg.TLS.InsecureSkipVerify,
		}
	}

	if cfg.SASL != nil {
		kafkaCfg.SASL = &kafka.SASLConfig{
			Enable:    cfg.SASL.Enable,
			Mechanism: cfg.SASL.Mechanism,
			Username:  cfg.SASL.Username,
			Password:  cfg.SASL.Password,
		}
	}

	producer, err := kafka.NewProducer(kafkaCfg)
	if err != nil {
		return nil, err
	}

	topics := []string{cfg.AuditEventsTopic, cfg.AuditLogsIngestTopic}

	if cfg.TopicAutoCreate {
		if err := producer.EnsureTopics(topics); err != nil {
			_ = producer.Close()
			return nil, err
		}
	}

	return &client{
		config:   cfg,
		producer: producer,
		topics:   topics,
	}, nil
}

func (c *client) Emit(event Event) {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return
	}
	c.mu.RUnlock()

	if err := c.validateEvent(&event); err != nil {
		return
	}

	c.enrichEvent(&event)

	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	c.producer.ProduceAsync(c.topics, []byte(event.TeamID), data)
}

func (c *client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	return c.producer.Close()
}

func (c *client) validateEvent(event *Event) error {
	if event.TeamID == "" {
		return ErrEmptyTeamID
	}
	if event.Event.Type == "" {
		return ErrEmptyEventType
	}
	return nil
}

func (c *client) enrichEvent(event *Event) {
	event.Version = Version

	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
}

