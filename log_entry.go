package echelon

import (
	"fmt"
	"time"
)

// LogLevel is level of log
type LogLevel uint32

const (
	// ErrorLevel is level to print error
	ErrorLevel LogLevel = iota
	// WarnLevel will print warnnings
	WarnLevel
	// InfoLevel will print info
	InfoLevel
	// DebugLevel will output debug info
	DebugLevel
	// TraceLevel will output trace info
	TraceLevel
)

// LogScopeStarted is a log with scope path and it's time
type LogScopeStarted struct {
	// scopes
	scopes []string
	// time is the time of start.
	time   time.Time
}

// NewLogScopeStarted will create a LogScopeStarted with LogScopeStarted.scopes is scpoes
func NewLogScopeStarted(scopes ...string) *LogScopeStarted {
	return &LogScopeStarted{
		scopes: scopes,
		time:   time.Now(),
	}
}

// GetScopes will return scopes path of entry
func (entry *LogScopeStarted) GetScopes() []string {
	return entry.scopes
}

// LogScopeFinished is a finished structure with path and finish status(succeed or failed)
type LogScopeFinished struct {
	scopes  []string
	success bool
}

// NewLogScopeFinished will create LogScopeFinished
func NewLogScopeFinished(success bool, scopes ...string) *LogScopeFinished {
	return &LogScopeFinished{
		scopes:  scopes,
		success: success,
	}
}

// Success returns wheter to LogScopeFinished has finished successfully
func (entry *LogScopeFinished) Success() bool {
	return entry.success
}

// GetScopes will returns scopes of LogScopeFinished
func (entry *LogScopeFinished) GetScopes() []string {
	return entry.scopes
}

// LogEntryMessage is a log with certain level, path
type LogEntryMessage struct {
	// Level is level o Log
	Level     LogLevel
	// formatt and arguments defines the message of log
	format    string
	arguments []interface{}
	// scopes are the names of scopes, which points out the path of log.
	scopes    []string
}

// NewLogEntryMessage creates a new log entry with path 'scopes', log level 'level', and message format, a...
func NewLogEntryMessage(scopes []string, level LogLevel, format string, a ...interface{}) *LogEntryMessage {
	return &LogEntryMessage{
		Level:     level,
		format:    format,
		arguments: a,
		scopes:    scopes,
	}
}

// GetMessage returns messages of log
func (entry *LogEntryMessage) GetMessage() string {
	return fmt.Sprintf(entry.format, entry.arguments...)
}

// GetScopes returns the path of scope
func (entry *LogEntryMessage) GetScopes() []string {
	return entry.scopes
}
