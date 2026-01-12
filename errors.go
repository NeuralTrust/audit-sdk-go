package audit

import "errors"

var (
	ErrNoBrokers      = errors.New("audit: no kafka brokers configured")
	ErrEmptyTeamID    = errors.New("audit: teamId is required")
	ErrEmptyEventType = errors.New("audit: event type is required")
)


