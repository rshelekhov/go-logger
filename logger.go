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

	//
	fatalEncountered bool
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
// It continuously listens for log messages and processes them based on their level.
// If the log level is FATAL, it processes all remaining messages in the log channel
// before terminating the program. The function also listens for a termination signal
// via the `done` channel to stop the logger. This method runs as a separate goroutine,
// allowing non-blocking, asynchronous logging.
func (l *Logger) run() {
	for {
		select {
		case logMessage, ok := <-l.logChan:
			if !ok { // logChan has been closed
				return
			}
			// Process the received log message
			l.processLogMessage(logMessage)

			// Terminate the program if a FATAL log level is encountered
			if l.fatalEncountered {
				l.processRemainingLogMessages()
				os.Exit(1)
			}
		case <-l.done:
			// Stop the logger when termination signal is received

			// Process remaining log messages in the log channel
			l.processRemainingLogMessages()
			return
		}
	}
}

// processLogMessage checks the log message level and processes it accordingly.
// If the log level is below the logger's set level, the message is ignored.
// If the message meets the required log level, it is either formatted as JSON or plain text
// based on the logger's configuration, then written to the output.
// In case of a FATAL log level, the program is terminated after the message is logged.
func (l *Logger) processLogMessage(logMessage LogMessage) {
	if logMessage.Level < l.level {
		// Ignore messages that are below the current log level
		return
	}

	var logRecord string
	var err error

	// Format the log message as JSON or plain text depending on the logger's configuration
	if l.json {
		logRecord, err = l.formatJSON(logMessage)
		if err != nil {
			l.logError(fmt.Errorf("error formatting log message as JSON: %w", err))
			return
		}
	} else {
		logRecord = l.formatText(logMessage)
	}

	// Write the formatted log message to the specified output writer
	if _, err = fmt.Fprint(l.writer, logRecord); err != nil {
		l.logError(fmt.Errorf("error writing log message: %w", err))
	}
}

// processRemainingLogMessages drains and processes all messages in the logChan until it is empty.
// This function is used to ensure that no log messages are left unprocessed before the logger
// terminates, either due to a FATAL log level or a graceful shutdown via the done channel.
func (l *Logger) processRemainingLogMessages() {
	for logMessage := range l.logChan {
		l.processLogMessage(logMessage)
	}
}

// formatJSON formats the given log message as a JSON string.
// It marshals the LogMessage struct into JSON format, returning the string representation.
// If an error occurs during marshaling, it returns an error.
func (l *Logger) formatJSON(logMessage LogMessage) (string, error) {
	jsonData, err := json.Marshal(logMessage)
	if err != nil {
		return "", fmt.Errorf("error marshaling log message: %w", err)
	}
	return string(jsonData) + "\n", nil
}

// formatText formats the log message as a plain text string.
// The format includes the timestamp, log level, and the log message content.
// The timestamp is formatted using the RFC3339 standard. The log level is converted to a string.
func (l *Logger) formatText(logMessage LogMessage) string {
	return fmt.Sprintf("%s [%s] %s\n", // TODO: add opts for date format
		logMessage.Time.Format(time.RFC3339),
		GetLevelString(logMessage.Level),
		logMessage.Msg,
	)
}

// Log sends a log message to the log channel with the specified log level and message content.
// It creates a new LogMessage instance with the current timestamp.
func (l *Logger) Log(level int, message string) {
	l.logChan <- LogMessage{
		Level: level,
		Msg:   message,
		Time:  time.Now(),
	}
}

// Debug logs a message at the DEBUG level.
//
// This method allows you to log detailed information useful for debugging purposes.
// The DEBUG level is the lowest level, and these logs are typically used during
// development or troubleshooting to gain insight into the flow and state of the application.
//
// Parameters:
//
//	message (string): The log message to be recorded.
func (l *Logger) Debug(message string) {
	l.Log(DEBUG, message)
}

// Info logs a message at the INFO level.
//
// This method is used to log general informational messages that highlight the
// progress of the application at a high level. INFO-level logs typically represent
// successful operations and normal workflow updates.
//
// Parameters:
//
//	message (string): The log message to be recorded.
func (l *Logger) Info(message string) {
	l.Log(INFO, message)
}

// Warning logs a message at the WARNING level.
//
// This method allows logging of potentially harmful situations that don't cause
// an immediate error but indicate that something unexpected has occurred.
// It can be used to highlight conditions that could lead to future issues.
//
// Parameters:
//
//	message (string): The log message to be recorded.
func (l *Logger) Warning(message string) {
	l.Log(WARNING, message)
}

// Error logs a message at the ERROR level.
//
// This method is used to log error events that might still allow the application
// to continue running but indicate a significant issue. ERROR-level logs are
// crucial for tracking exceptions and failures within the application.
//
// Parameters:
//
//	message (string): The log message to be recorded.
func (l *Logger) Error(message string) {
	l.Log(ERROR, message)
}

// Fatal logs a message at the FATAL level and terminates the program.
//
// This method logs critical issues that result in the termination of the application.
// After logging the message, the program will be forcibly exited with a status of 1.
//
// Parameters:
//
//	message (string): The log message to be recorded.
func (l *Logger) Fatal(message string) {
	// Set the fatalEncountered flag to true, indicating a fatal error has occurred.
	l.fatalEncountered = true

	// Log the message with the FATAL level.
	l.Log(FATAL, message)

	// Close the logger channels to stop further logging and release resources.
	l.Close()
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

// Close signals the logger to terminate by closing the channels.
// This allows the run method to exit gracefully and stop processing log messages.
func (l *Logger) Close() {
	close(l.logChan)
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
