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

	// NoProgress refers that the node is not a progress bar
	NoProgress = 0
	// DefaultProgress refers that the node has progress bar without known progress size
	DefaultProgress = 100
)

// LogScopeStarted sends start message to with time stamp to node specified by scopes
type LogScopeStarted struct {
	// scopes
	scopes []string
	// total is setted for progress bar
	total int64
	// time is the time of start.
	time   time.Time
}

// NewLogScopeStarted will create a LogScopeStarted with LogScopeStarted. scopes is path of log.
// If the total is not 0, then it starts a log with progress bar
func NewLogScopeStarted(total int64,scopes ...string) *LogScopeStarted {
	return &LogScopeStarted{
		scopes: scopes,
		total: total,
		time:   time.Now(),
	}
}

// GetScopes will return scopes path of entry
func (entry *LogScopeStarted) GetScopes() []string {
	return entry.scopes
}

// GetProgressSize returns size of progress
func (entry *LogScopeStarted) GetProgressSize() int64 {
	return entry.total
}

// LogScopeFinished sends finished message and finish status(succeed or failed) to node specified with path 
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

// LogEntryMessage is a struct sends new message with certain level to node specified by scopes
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

// LogProcessMessage sends progress message to node specified by scopes
type LogProcessMessage struct {
	Progress int64
	Addprogress int64
	Percentage int
	Addpercentage int
	scopes []string
}

// NewLogProcessMessage creates a log process
func NewLogProcessMessage(scopes... string) *LogProcessMessage {
	return &LogProcessMessage{
		scopes: scopes,
	}
}

// GetScopes returns scopes of entry
func (entry *LogProcessMessage) GetScopes() []string {
	return entry.scopes
}