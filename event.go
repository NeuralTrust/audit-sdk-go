package audit

import (
	"encoding/json"
	"time"
)

const (
	Version         = "1.0"
	timestampFormat = "2006-01-02 15:04:05.000000"
)

type ActorType string

const (
	ActorTypeUser    ActorType = "user"
	ActorTypeService ActorType = "service"
	ActorTypeSystem  ActorType = "system"
)

type Event struct {
	Version   string    `json:"version"`
	ID        string    `json:"id"`
	Timestamp Timestamp `json:"timestamp"`
	TeamID    string    `json:"team_id"`
	Event     EventInfo `json:"event"`
	Target    Target    `json:"target"`
	Actor     *Actor    `json:"actor"`
	Context   *Context  `json:"context"`
	Changes   *Changes  `json:"changes,omitempty"`
	Metadata  *Metadata `json:"metadata,omitempty"`
}

type EventInfo struct {
	Type         string `json:"type"`
	Category     string `json:"category"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type Actor struct {
	ID    string    `json:"id"`
	Email string    `json:"email,omitempty"`
	Type  ActorType `json:"type"`
}

type Target struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type Context struct {
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	SessionID string `json:"session_id,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

type Changes struct {
	Previous map[string]interface{} `json:"previous,omitempty"`
	Current  map[string]interface{} `json:"current,omitempty"`
}

type Metadata map[string]interface{}

type Timestamp struct {
	time.Time
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Format(timestampFormat))
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := time.Parse(timestampFormat, s)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}
