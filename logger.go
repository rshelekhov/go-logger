package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Log levels
const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// LogMessage structure represents a log entry with level, message, and timestamp
type LogMessage struct {
	Level int       `json:"level"`
	Msg   string    `json:"msg"`
	Time  time.Time `json:"time"`
}

// Logger struct encapsulates the logging functionality
type Logger struct {
	// Minimum log level for logging messages
	level int

	// Channel for sending log messages
	logChan chan LogMessage

	// Channel to signal when to stop the logger
	done chan struct{}

	// Interface for writing log outputs
	writer io.Writer

	// Flag to determine if logs should be in JSON format
	// true => output JSON; false => output text
	json bool
}

// New function creates a new Logger instance with the specified parameters
func New(level int, writer io.Writer, json bool) *Logger {
	logger := &Logger{
		level:   level,
		logChan: make(chan LogMessage, 100),
		done:    make(chan struct{}),
		writer:  writer,
		json:    json,
	}

	// Start the logging process in a new goroutine
	go logger.run()

	// Return the new Logger instance
	return logger
}

// run method processes incoming log messages
func (l *Logger) run() {
	for {
		select {
		case logMessage := <-l.logChan:
			// Check if the message level is above or equal to the logger's level.
			// If level is higher than current log level, write log to writer, else do nothing
			if logMessage.Level >= l.level {
				var logRecord string

				if l.json {
					// TODO: Move to separate function
					// TODO: add converts log level to string

					// If JSON format is specified, marshal the log message to JSON
					jsonData, err := json.Marshal(logMessage)
					if err != nil {
						// Log the error if marshaling fails
						l.logError(fmt.Errorf("error marshaling log message: %w", err))

						// Skip to the next message
						continue
					}
					logRecord = string(jsonData) + "\n"
				} else {
					// TODO: Move to separate function

					// If text format is specified, format the log message as plain text
					logRecord = fmt.Sprintf("%s [%s] %s\n", // TODO: add opts for date format
						logMessage.Time.Format(time.RFC3339),
						GetLevelString(logMessage.Level),
						logMessage.Msg,
					)
				}

				// Write the log record to the specified writer
				if _, err := fmt.Fprint(l.writer, logRecord); err != nil {
					l.logError(fmt.Errorf("error writing log message: %w", err))
				}

				// If the log level is FATAL, terminate the program
				if logMessage.Level == FATAL {
					os.Exit(1)
				}
			}

		// Check if termination signal received
		case <-l.done:
			// Stop the logger
			return
		}
	}
}

// TODO: refactor this for using different methods for each level

// Log method sends a log message to the log channel
func (l *Logger) Log(level int, message string) {
	l.logChan <- LogMessage{
		Level: level,
		Msg:   message,
		Time:  time.Now(),
	}
}

// logError method handles errors encountered during logging in the run method
func (l *Logger) logError(err error) {
	errorLog := LogMessage{
		Level: ERROR,
		Msg:   err.Error(),
		Time:  time.Now(),
	}

	// Write the error log to the writer, if it fails log to stdout
	if _, writeErr := fmt.Fprintf(l.writer, "Logger error: %s\n", errorLog.Msg); writeErr != nil {
		fmt.Printf("Error writing to log writer: %v\n", writeErr)
	}
}

// Close method signals to terminate the logger
func (l *Logger) Close() {
	close(l.done)
}

// GetLevelString converts log level integer to string representation
func GetLevelString(level int) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}
