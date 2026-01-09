# Audit SDK Go

A Go SDK for publishing audit events to Kafka. Designed for use with TrustGate and TrustGate-EE.

## Version

**v0.1.0**

## Requirements

- Go 1.21+
- librdkafka (required by confluent-kafka-go)

### Installing librdkafka

**macOS:**
```bash
brew install librdkafka
```

**Ubuntu/Debian:**
```bash
apt-get install librdkafka-dev
```

**Alpine:**
```bash
apk add librdkafka-dev
```

## Installation

```bash
go get github.com/NeuralTrust/audit-sdk-go
```

## Quick Start

```go
package main

import (
    "log"

    audit "github.com/NeuralTrust/audit-sdk-go"
)

func main() {
    client, err := audit.New(&audit.Config{
        Brokers:         []string{"localhost:9092"},
        TopicAutoCreate: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    client.Emit(audit.Event{
        TeamID: "team-123",
        Event: audit.EventInfo{
            Type:        "gateway.created",
            Category:    "gateway",
            Description: "A new gateway was created",
            Status:      "success",
        },
        Actor: audit.Actor{
            ID:    "user-456",
            Email: "user@example.com",
            Type:  audit.ActorTypeUser,
        },
        Target: audit.Target{
            Type: "gateway",
            ID:   "gw-789",
            Name: "my-gateway",
        },
    })
}
```

## Configuration

### Config Options

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Brokers` | `[]string` | **required** | Kafka broker addresses |
| `AuditEventsTopic` | `string` | `audit_events` | Topic for audit events |
| `AuditLogsIngestTopic` | `string` | `audit_logs_ingest` | Topic for audit logs ingestion |
| `ClientID` | `string` | `audit-sdk` | Kafka client identifier |
| `TopicAutoCreate` | `bool` | `false` | Auto-create topics if they don't exist |
| `TopicNumParts` | `int` | `3` | Number of partitions for auto-created topics |
| `TopicReplication` | `int` | `1` | Replication factor for auto-created topics |
| `RetryMax` | `int` | `3` | Maximum retries for failed messages |
| `RequiredAcks` | `int` | `1` | Required acks (0=none, 1=leader, -1=all) |
| `TLS` | `*TLSConfig` | `nil` | TLS configuration |
| `SASL` | `*SASLConfig` | `nil` | SASL authentication configuration |

### Environment Variables

Topics can be configured via environment variables:

| Variable | Default |
|----------|---------|
| `AUDIT_EVENTS_TOPIC` | `audit_events` |
| `AUDIT_LOGS_INGEST_TOPIC` | `audit_logs_ingest` |

Priority: Config > Environment Variable > Default

### TLS Configuration

```go
client, err := audit.New(&audit.Config{
    Brokers: []string{"kafka:9093"},
    TLS: &audit.TLSConfig{
        Enable:   true,
        CAFile:   "/path/to/ca.pem",
        CertFile: "/path/to/cert.pem",
        KeyFile:  "/path/to/key.pem",
    },
})
```

### SASL Authentication

```go
client, err := audit.New(&audit.Config{
    Brokers: []string{"kafka:9092"},
    SASL: &audit.SASLConfig{
        Enable:    true,
        Mechanism: "PLAIN",
        Username:  "user",
        Password:  "password",
    },
})
```

## Event Structure

```go
type Event struct {
    Version   string    `json:"version"`   // Auto-set to "1.0"
    ID        string    `json:"id"`        // Auto-generated UUID if empty
    Timestamp time.Time `json:"timestamp"` // Auto-set to current time if zero
    TeamID    string    `json:"teamId"`    // Required
    Event     EventInfo `json:"event"`     // Required (Type field is required)
    Actor     Actor     `json:"actor"`
    Target    Target    `json:"target"`
    Context   Context   `json:"context"`
    Changes   *Changes  `json:"changes,omitempty"`
    Metadata  Metadata  `json:"metadata,omitempty"`
}
```

### Actor Types

```go
audit.ActorTypeUser    // "user"
audit.ActorTypeService // "service"
audit.ActorTypeSystem  // "system"
```

## API

### `audit.New(cfg *Config) (Client, error)`

Creates a new audit client. Returns an error if configuration is invalid or Kafka connection fails.

### `client.Emit(event Event)`

Publishes an audit event to Kafka. This method is **asynchronous** (fire-and-forget) and will not block.

Events are published to both configured topics (`AuditEventsTopic` and `AuditLogsIngestTopic`).

Required fields:
- `TeamID`
- `Event.Type`

### `client.Close() error`

Closes the Kafka producer and flushes pending messages. Should be called before application shutdown.

## Development

### Run Tests

```bash
make test
```

### Generate Coverage Report

```bash
make test-coverage
```

### Build

```bash
make build
```

### Lint

```bash
make lint
```

## License

MIT License - NeuralTrust

