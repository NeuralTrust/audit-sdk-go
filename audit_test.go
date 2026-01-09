package audit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockProducer struct {
	producedMessages []producedMessage
	ensuredTopics    []string
	closed           bool
}

type producedMessage struct {
	topics []string
	key    []byte
	value  []byte
}

func (m *mockProducer) ProduceAsync(topics []string, key, value []byte) {
	m.producedMessages = append(m.producedMessages, producedMessage{
		topics: topics,
		key:    key,
		value:  value,
	})
}

func (m *mockProducer) EnsureTopics(topics []string) error {
	m.ensuredTopics = topics
	return nil
}

func (m *mockProducer) Close() error {
	m.closed = true
	return nil
}

func TestClient_Emit_Success(t *testing.T) {
	mock := &mockProducer{}
	c := &client{
		config: &Config{
			AuditEventsTopic:     "events",
			AuditLogsIngestTopic: "logs",
		},
		producer: mock,
		topics:   []string{"events", "logs"},
	}

	event := Event{
		TeamID: "team-123",
		Event: EventInfo{
			Type:     "test.event",
			Category: "test",
		},
		Actor: Actor{
			ID:   "user-1",
			Type: ActorTypeUser,
		},
	}

	c.Emit(event)

	require.Len(t, mock.producedMessages, 1)
	assert.Equal(t, []string{"events", "logs"}, mock.producedMessages[0].topics)
	assert.Equal(t, []byte("team-123"), mock.producedMessages[0].key)
	assert.Contains(t, string(mock.producedMessages[0].value), `"teamId":"team-123"`)
	assert.Contains(t, string(mock.producedMessages[0].value), `"type":"test.event"`)
}

func TestClient_Emit_EnrichesEvent(t *testing.T) {
	mock := &mockProducer{}
	c := &client{
		config:   &Config{},
		producer: mock,
		topics:   []string{"topic1"},
	}

	event := Event{
		TeamID: "team-123",
		Event: EventInfo{
			Type: "test.event",
		},
	}

	c.Emit(event)

	require.Len(t, mock.producedMessages, 1)
	value := string(mock.producedMessages[0].value)

	assert.Contains(t, value, `"version":"1.0"`)
	assert.Contains(t, value, `"id":"`)
	assert.Contains(t, value, `"timestamp":"`)
}

func TestClient_Emit_EmptyTeamID_DoesNotProduce(t *testing.T) {
	mock := &mockProducer{}
	c := &client{
		config:   &Config{},
		producer: mock,
		topics:   []string{"topic1"},
	}

	event := Event{
		TeamID: "",
		Event: EventInfo{
			Type: "test.event",
		},
	}

	c.Emit(event)

	assert.Len(t, mock.producedMessages, 0)
}

func TestClient_Emit_EmptyEventType_DoesNotProduce(t *testing.T) {
	mock := &mockProducer{}
	c := &client{
		config:   &Config{},
		producer: mock,
		topics:   []string{"topic1"},
	}

	event := Event{
		TeamID: "team-123",
		Event: EventInfo{
			Type: "",
		},
	}

	c.Emit(event)

	assert.Len(t, mock.producedMessages, 0)
}

func TestClient_Emit_WhenClosed_DoesNotProduce(t *testing.T) {
	mock := &mockProducer{}
	c := &client{
		config:   &Config{},
		producer: mock,
		topics:   []string{"topic1"},
		closed:   true,
	}

	event := Event{
		TeamID: "team-123",
		Event: EventInfo{
			Type: "test.event",
		},
	}

	c.Emit(event)

	assert.Len(t, mock.producedMessages, 0)
}

func TestClient_Close(t *testing.T) {
	mock := &mockProducer{}
	c := &client{
		config:   &Config{},
		producer: mock,
		topics:   []string{"topic1"},
	}

	err := c.Close()

	assert.NoError(t, err)
	assert.True(t, mock.closed)
	assert.True(t, c.closed)
}

func TestClient_Close_AlreadyClosed(t *testing.T) {
	mock := &mockProducer{}
	c := &client{
		config:   &Config{},
		producer: mock,
		topics:   []string{"topic1"},
		closed:   true,
	}

	err := c.Close()

	assert.NoError(t, err)
	assert.False(t, mock.closed)
}

func TestValidateEvent(t *testing.T) {
	c := &client{}

	tests := []struct {
		name    string
		event   Event
		wantErr error
	}{
		{
			name: "valid event",
			event: Event{
				TeamID: "team-123",
				Event:  EventInfo{Type: "test"},
			},
			wantErr: nil,
		},
		{
			name: "empty team id",
			event: Event{
				TeamID: "",
				Event:  EventInfo{Type: "test"},
			},
			wantErr: ErrEmptyTeamID,
		},
		{
			name: "empty event type",
			event: Event{
				TeamID: "team-123",
				Event:  EventInfo{Type: ""},
			},
			wantErr: ErrEmptyEventType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.validateEvent(&tt.event)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestEnrichEvent(t *testing.T) {
	c := &client{}

	event := Event{
		TeamID: "team-123",
	}

	c.enrichEvent(&event)

	assert.Equal(t, Version, event.Version)
	assert.NotEmpty(t, event.ID)
	assert.False(t, event.Timestamp.IsZero())
	assert.WithinDuration(t, time.Now().UTC(), event.Timestamp, time.Second)
}

func TestEnrichEvent_PreservesExistingValues(t *testing.T) {
	c := &client{}

	existingTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	event := Event{
		ID:        "existing-id",
		Timestamp: existingTime,
		TeamID:    "team-123",
	}

	c.enrichEvent(&event)

	assert.Equal(t, Version, event.Version)
	assert.Equal(t, "existing-id", event.ID)
	assert.Equal(t, existingTime, event.Timestamp)
}

