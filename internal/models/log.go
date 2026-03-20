package models

import "time"

// Level represents log severity
type Level string

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
	LevelFatal Level = "FATAL"
)

// LogEntry is a single log record
type LogEntry struct {
	ID        string            `json:"id"`
	Timestamp time.Time         `json:"timestamp"`
	Level     Level             `json:"level"`
	Service   string            `json:"service"`
	Message   string            `json:"message"`
	Meta      map[string]string `json:"meta,omitempty"` // extra fields like trace_id, user_id, etc.
}

// QueryFilter is used to filter logs in queries
type QueryFilter struct {
	Level   Level  `json:"level"`
	Service string `json:"service"`
	Search  string `json:"search"`
	Limit   int    `json:"limit"`
}

// Stats holds aggregate log statistics
type Stats struct {
	Total    int            `json:"total"`
	ByLevel  map[Level]int  `json:"by_level"`
	ByService map[string]int `json:"by_service"`
}
