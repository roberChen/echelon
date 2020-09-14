package echelon

// genericLogEntry is a log entry contains whether started log, finished log or running log
type genericLogEntry struct {
	LogStarted  *LogScopeStarted
	LogFinished *LogScopeFinished
	LogEntry    *LogEntryMessage
}

// LogRendered interface defines a log which can start/finish and render
type LogRendered interface {
	// RenderScopeStarted will start a render job
	RenderScopeStarted(entry *LogScopeStarted)
	// RenderScopeFinished will finish render job
	RenderScopeFinished(entry *LogScopeFinished)
	// RenderMessage will render messages from entry
	RenderMessage(entry *LogEntryMessage)
}

// Logger is a log object with a log level, scopes and entries chan.entries Channel 
// will render all entries it receive after calling (*Logger).streamEntries function
type Logger struct {
	level          LogLevel
	scopes         []string
	// entriesChannel will render all entries it receive after calling (*Logger).streamEntries function
	entriesChannel chan *genericLogEntry
}

// NewLogger creates a log object with new generated entries channel. And use renderer as renderer of logger
func NewLogger(level LogLevel, renderer LogRendered) *Logger {
	logger := &Logger{
		level:          level,
		entriesChannel: make(chan *genericLogEntry),
	}
	go logger.streamEntries(renderer)
	return logger
}

// Scoped creates a sub log with name scope
func (logger *Logger) Scoped(scope string) *Logger {
	result := &Logger{
		level:          logger.level,
		scopes:         append(logger.scopes, scope),
		entriesChannel: logger.entriesChannel,
	}
	result.entriesChannel <- &genericLogEntry{
		LogStarted: NewLogScopeStarted(result.scopes...),
	}
	return result
}

// streamEntries will continiously render all entry receinved from logger entries channel.
func (logger *Logger) streamEntries(renderer LogRendered) {
	for {
		// receive entry from logger.entriesChannel
		entry := <-logger.entriesChannel
		if entry.LogStarted != nil {
			renderer.RenderScopeStarted(entry.LogStarted)
		}
		if entry.LogFinished != nil {
			renderer.RenderScopeFinished(entry.LogFinished)
		}
		if entry.LogEntry != nil {
			renderer.RenderMessage(entry.LogEntry)
		}
	}
}

// Tracef will print trace info
func (logger *Logger) Tracef(format string, args ...interface{}) {
	logger.Logf(TraceLevel, format, args...)
}

// Debugf will print debug info
func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.Logf(DebugLevel, format, args...)
}

// Infof will print info
func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.Logf(InfoLevel, format, args...)
}

// Warnf will print warnning info
func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.Logf(WarnLevel, format, args...)
}

// Errorf will print error
func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.Logf(ErrorLevel, format, args...)
}

// Logf sends a log entry message with LogLevel level to logger entries chan
func (logger *Logger) Logf(level LogLevel, format string, args ...interface{}) {
	if logger.IsLogLevelEnabled(level) {
		logger.entriesChannel <- &genericLogEntry{
			LogEntry: NewLogEntryMessage(logger.scopes, level, format, args...),
		}
	}
}

// Finish will finsh a log with success status (true for succeed, false for failed),
// it will sends a NewLogScopeFinished to entries chan of logger
func (logger *Logger) Finish(success bool) {
	logger.entriesChannel <- &genericLogEntry{
		LogFinished: NewLogScopeFinished(success, logger.scopes...),
	}
}

// IsLogLevelEnabled returns wheter a log will print to Writer
func (logger *Logger) IsLogLevelEnabled(level LogLevel) bool {
	return level <= logger.level
}
