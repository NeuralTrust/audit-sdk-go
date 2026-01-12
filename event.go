package audit

import "time"

const Version = "1.0"

type ActorType string

const (
	ActorTypeUser    ActorType = "user"
	ActorTypeService ActorType = "service"
	ActorTypeSystem  ActorType = "system"
)

type Event struct {
	Version   string    `json:"version"`
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	TeamID    string    `json:"teamId"`
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
	ErrorMessage string `json:"errorMessage,omitempty"`
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
	IPAddress string `json:"ipAddress,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`
	SessionID string `json:"sessionId,omitempty"`
	RequestID string `json:"requestId,omitempty"`
}

type Changes struct {
	Previous map[string]interface{} `json:"previous,omitempty"`
	Current  map[string]interface{} `json:"current,omitempty"`
}

type Metadata map[string]interface{}
