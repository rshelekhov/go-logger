# go-logger

This Go package provides a flexible and lightweight logging system that supports multiple log levels,
structured log output, and both JSON and text-based formatting. The library allows for easy integration
into applications and supports log output to various io.Writer targets such as files, stdout,
or any custom destination.

The logger processes log messages asynchronously using Go channels. The logging process runs in a separate
goroutine, ensuring that logging does not block the main execution of your application.

## Features
- Multiple Log Levels: Supports DEBUG, INFO, WARNING, ERROR, and FATAL log levels.
- JSON and Plain Text Formatting: Logs can be formatted as structured JSON or plain text.
- Log Output Destination: Logs can be written to files, stdout, or any custom io.Writer.
- Asynchronous Logging: Logs are handled asynchronously using channels for efficient, non-blocking logging.
- Graceful Shutdown: The logger can be gracefully stopped using the Close() method.

## Upcoming Features
- Log Level Customization: Provide support for dynamic adjustment of log levels at runtime.
- Structured File Output: Include both filename and line number in log messages.
- Enhanced Formatting Options: Add customization options for log output formatting, such as timestamps and date formats, colors output, etc.

## Installation
To use this library in your project, run:
```bash
go get github.com/rshelekhov/go-logger
```

## Usage

### Basic Usage Example

The logger can be created with custom configurations, specifying the minimum log level,
the output destination, and whether logs should be in JSON format or plain text.

```go
package main

import (
    "os"
    "github.com/rshelekhov/go-logger"
)

func main() {
    log := logger.New(logger.INFO, os.Stdout, false) // Create a logger for text output
    defer log.Close()

    log.Log(logger.INFO, "INFO message")
    log.Log(logger.ERROR, "ERROR message")
}
```

Examples of different logger configurations can be found in the examples folder:
1. json_to_file: Logs messages in JSON format to a file.
2. json_to_stdout: Outputs logs in JSON format to stdout.
3. text_to_stdout: Outputs plain text logs to stdout.

# Contribution

Feel free to submit issues or pull requests to help improve the library. Some of the planned features
include structured logging with additional metadata and more customizable formatting options.

# License

This project is licensed under the MIT License.