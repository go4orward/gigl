package common

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"
)

type iLogger interface {
	// SetOption() set the datetime format and the flag to show logging position
	//   `time_format` can be an empty string or something like "2006-01-02 15:04:05.000 -0700"
	//   `location`    is true/false to show the logging position like "main.go:251"
	SetOption(time_format string, location bool) iLogger
	// SetTraceFilter() set the filter of 'Logger.Trace()' message pattern, so that we can get selected log messages only.
	//   `regexp_pattern` regular expression pattern for Trace messages.
	SetTraceFilter(regexp_pattern string) iLogger
	// GetLogLevel() returns the LogLevel of the current Logger.
	GetLogLevel() LogLevel
	// IsLogging() returns whether or not the Logger is logging for the given LogLevel.
	IsLogging(loglevel LogLevel) bool
	// ShowStatus() prints brief information about the current Logger on os.Stdout.
	ShowStatus()

	Error(format string, args ...any) // application hits an issue preventing one or more functionalities from properly functioning (most important)
	Warn(format string, args ...any)  // indicating that something unexpected happened in the application, a problem, or a disturbing situation
	Info(format string, args ...any)  // the standard log level to report that something happened, the application entered a certain state, etc
	Debug(format string, args ...any) // for diagnosing issues and troubleshooting
	Trace(format string, args ...any) // most fine grained information (most verbose)
}

var Logger iLogger = &sDummyLogger{} // The Logger
var log_time_format string           // NO datetime info by default, but can be set to "2006-01-02 15:04:05.000 -0700", for example.
var log_show_location bool           // NO locaton  info by default
var log_trace_pattern string         // No trace message filtering by default
var log_trace_regexp *regexp.Regexp  // No trace message filtering by default

type LogLevel int

const (
	LogLevelNone  LogLevel = iota //
	LogLevelError                 // most important
	LogLevelWarn                  //
	LogLevelInfo                  // likely initial (default) log level
	LogLevelDebug                 //
	LogLevelTrace                 // most verbose
)

// SetLogger() sets the current logger to be used.
//
//	`logger` should be the returned value of either NewConsoleLogger() or NewWriterLogger().
func SetLogger(logger iLogger) iLogger {
	if logger != nil {
		Logger = logger
	}
	return Logger
}

// ----------------------------------------------------------------------------
// DummyLogger
// ----------------------------------------------------------------------------

// NewDummyLogger() returns a sDummyLogger, which prints nothing.
func NewDummyLogger() iLogger {
	return &sDummyLogger{}
}

type sDummyLogger struct{} // invisible from outside, so use 'NewDummyLogger()'

func (self *sDummyLogger) SetOption(time_format string, location bool) iLogger { return self }
func (self *sDummyLogger) SetTraceFilter(regexp_pattern string) iLogger        { return self }
func (self *sDummyLogger) GetLogLevel() LogLevel                               { return LogLevelNone }
func (self *sDummyLogger) IsLogging(loglevel LogLevel) bool                    { return false }
func (self *sDummyLogger) ShowStatus()                                         { fmt.Printf("Logger : DummyLogger\n") }

func (sDummyLogger) Error(format string, args ...any) {} // most important
func (sDummyLogger) Warn(format string, args ...any)  {}
func (sDummyLogger) Info(format string, args ...any)  {} // likely to be initial
func (sDummyLogger) Debug(format string, args ...any) {}
func (sDummyLogger) Trace(format string, args ...any) {} // most verbose

// ----------------------------------------------------------------------------
// ConsoleLogger
// ----------------------------------------------------------------------------

// NewConsoleLogger() returns a sConsoleLogger, which prints log messages on os.Stdout.
//
//	`loglevel` can be LogLevel/int or string like "none","error","warn","info","debug","trace".
func NewConsoleLogger(loglevel any) iLogger {
	log_level := LogLevelNone
	switch loglevel.(type) {
	case string:
		log_level = get_log_level_from_string(loglevel.(string))
	case LogLevel:
		log_level = loglevel.(LogLevel)
	}
	return &sConsoleLogger{loglevel: log_level}
}

type sConsoleLogger struct { // invisible from outside, so use 'NewConsoleLogger()'
	loglevel LogLevel
}

func (self *sConsoleLogger) SetOption(time_format string, location bool) iLogger {
	if time_format == "default" {
		time_format = "2006-01-02 15:04:05.000 -0700"
	}
	log_time_format = time_format
	log_show_location = location
	return self
}

func (self *sConsoleLogger) SetTraceFilter(regexp_pattern string) iLogger {
	if regexp_pattern == "" {
		log_trace_regexp = nil
	} else {
		rex, err := regexp.Compile(regexp_pattern)
		if err == nil {
			log_trace_pattern = regexp_pattern
			log_trace_regexp = rex
		}
	}
	return self
}

func (self *sConsoleLogger) GetLogLevel() LogLevel {
	return self.loglevel
}

func (self *sConsoleLogger) IsLogging(loglevel LogLevel) bool {
	return self.loglevel >= loglevel
}

func (self *sConsoleLogger) ShowStatus() {
	fmt.Printf("ConsoleLogger : LogLevel:%d  TraceFilter:%#q\n", Logger.GetLogLevel(), log_trace_pattern)
}

func (self *sConsoleLogger) Error(format string, args ...any) {
	if self.loglevel >= LogLevelError {
		log_print_writer(os.Stdout, false, "[ERROR]", format, args...)
	}
}

func (self *sConsoleLogger) Warn(format string, args ...any) {
	if self.loglevel >= LogLevelWarn {
		log_print_writer(os.Stdout, false, "[WARN] ", format, args...)
	}
}

func (self *sConsoleLogger) Info(format string, args ...any) {
	if self.loglevel >= LogLevelInfo {
		log_print_writer(os.Stdout, false, "[INFO] ", format, args...)
	}
}

func (self *sConsoleLogger) Debug(format string, args ...any) {
	if self.loglevel >= LogLevelDebug {
		log_print_writer(os.Stdout, false, "[DEBUG]", format, args...)
	}
}

func (self *sConsoleLogger) Trace(format string, args ...any) {
	if self.loglevel >= LogLevelTrace {
		log_print_writer(os.Stdout, true, "[TRACE]", format, args...)
	}
}

// ----------------------------------------------------------------------------
// WriterLogger
// ----------------------------------------------------------------------------

// NewWriterLogger() returns a sWriterLogger, which prints log messages to the specified io.Writer.
//
//	`loglevel` can be LogLevel/int or string like "none","error","warn","info","debug","trace".
func NewWriterLogger(loglevel any, writer io.Writer) iLogger {
	log_level := LogLevelNone
	switch loglevel.(type) {
	case string:
		log_level = get_log_level_from_string(loglevel.(string))
	case LogLevel:
		log_level = loglevel.(LogLevel)
	}
	return &sWriterLogger{loglevel: log_level, output: writer}
}

type sWriterLogger struct { // invisible from outside, so use 'NewWriterLogger()'
	loglevel LogLevel
	output   io.Writer
}

func (self *sWriterLogger) SetOption(time_format string, location bool) iLogger {
	if time_format == "default" {
		time_format = "2006-01-02 15:04:05.000 -0700"
	}
	log_time_format = time_format
	log_show_location = location
	return self
}

func (self *sWriterLogger) SetTraceFilter(regexp_pattern string) iLogger {
	if regexp_pattern == "" {
		log_trace_regexp = nil
	} else {
		rex, err := regexp.Compile(regexp_pattern)
		if err == nil {
			log_trace_pattern = regexp_pattern
			log_trace_regexp = rex
		}
	}
	return self
}

func (self *sWriterLogger) GetLogLevel() LogLevel {
	return self.loglevel
}

func (self *sWriterLogger) IsLogging(loglevel LogLevel) bool {
	return self.loglevel >= loglevel
}

func (self *sWriterLogger) ShowStatus() {
	fmt.Printf("WriterLogger  : LogLevel:%d  TraceFilter:%#q\n", Logger.GetLogLevel(), log_trace_pattern)
}

func (self *sWriterLogger) Error(format string, args ...any) {
	if self.loglevel >= LogLevelError {
		log_print_writer(self.output, false, "[ERROR]", format, args...)
	}
}

func (self *sWriterLogger) Warn(format string, args ...any) {
	if self.loglevel >= LogLevelWarn {
		log_print_writer(self.output, false, "[WARN] ", format, args...)
	}
}

func (self *sWriterLogger) Info(format string, args ...any) {
	if self.loglevel >= LogLevelInfo {
		log_print_writer(self.output, false, "[INFO] ", format, args...)
	}
}

func (self *sWriterLogger) Debug(format string, args ...any) {
	if self.loglevel >= LogLevelDebug {
		log_print_writer(self.output, false, "[DEBUG]", format, args...)
	}
}

func (self *sWriterLogger) Trace(format string, args ...any) {
	if self.loglevel >= LogLevelTrace {
		log_print_writer(self.output, true, "[TRACE]", format, args...)
	}
}

// ----------------------------------------------------------------------------
// private functions
// ----------------------------------------------------------------------------

func get_log_level_from_string(loglevel string) LogLevel {
	switch loglevel {
	case "none":
		return LogLevelNone
	case "error":
		return LogLevelError
	case "warn":
		return LogLevelWarn
	case "info":
		return LogLevelInfo
	case "debug":
		return LogLevelDebug
	case "trace":
		return LogLevelTrace
	default:
		return LogLevelNone
	}
}

func log_print_writer(f io.Writer, trace bool, prefix string, format string, args ...any) {
	log_message := fmt.Sprintf(format, args...)
	if !trace || log_trace_regexp == nil || log_trace_regexp.MatchString(log_message) {
		if log_time_format != "" {
			prefix = time.Now().Format(log_time_format) + " " + prefix
		}
		if log_show_location {
			_, file, line, ok := runtime.Caller(3)
			if ok {
				file = filepath.Base(file)
			} else {
				file = "???"
				line = 0
			}
			fmt.Fprintf(f, "%s %-140s %30s:%d\n", prefix, log_message, file, line)
		} else {
			fmt.Fprintf(f, "%s %s\n", prefix, log_message)
		}
	}
}
