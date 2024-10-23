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

// New creates a new Logger instance with the specified logging level, output writer, and format.
//
// Parameters:
//
//	level (int): The minimum log level for messages to be logged. Messages with a level
//	             lower than this will be ignored. Valid levels are DEBUG, INFO, WARNING,
//	             ERROR, and FATAL.
//	writer (io.Writer): An interface for writing log outputs. This can be any writer,
//	                    such as os.Stdout, a file, etc.
//	json (bool): A boolean flag indicating the desired output format. If true, logs
//	             will be formatted as JSON; if false, logs will be in plain text format.
//
// Returns:
//
//	*Logger: A pointer to the newly created Logger instance
func New(level int, writer io.Writer, json bool) *Logger {
	logger := &Logger{
		level: level,

		// Buffered channel with a capacity of 100 messages, allowing non-blocking log writes.
		// This improves performance by enabling the logger to continue processing logs
		// asynchronously without waiting for the log consumer to be ready.
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

// run processes incoming log messages from the log channel.
// It checks if the log message level is above or equal to the logger's set level.
// If it is, the log message is formatted and written to the specified output.
// The method runs in a separate goroutine, allowing for asynchronous logging.
// It also handles FATAL log levels by terminating the program.
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

// Log sends a log message to the log channel with the specified log level and message content.
// It creates a new LogMessage instance with the current timestamp.
func (l *Logger) Log(level int, message string) {
	l.logChan <- LogMessage{
		Level: level,
		Msg:   message,
		Time:  time.Now(),
	}
}

// logError handles errors encountered during the logging process.
// It constructs an error LogMessage and writes it to the specified writer.
// If writing fails, it logs the error to stdout.
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

// Close signals the logger to terminate by closing the done channel.
// This allows the run method to exit gracefully and stop processing log messages.
func (l *Logger) Close() {
	close(l.done)
}

// GetLevelString converts a log level integer to its string representation.
// This function returns the corresponding string for each log level,
// or "UNKNOWN" if the level is not recognized.
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
