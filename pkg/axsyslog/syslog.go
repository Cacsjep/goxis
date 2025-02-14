// Package syslog provides a wrapper for system log functionality, allowing for logging
// messages with various priorities and optional console output.
package axsyslog

/*
#include <stdlib.h>
#include <syslog.h>

void syslog_helper(int priority, const char *message) {
    syslog(priority, "%s", message);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

const (
	// Options for openlog.
	LOG_PID  = 0x01
	LOG_CONS = 0x02
	// Facilities.
	LOG_USER = (1 << 3)

	// Syslog priorities.
	LOG_EMERG  = 0
	LOG_ALERT  = 1
	LOG_CRIT   = 2
	LOG_ERR    = 3
	LOG_WARN   = 4
	LOG_NOTICE = 5
	LOG_INFO   = 6
	LOG_DEBUG  = 7
)

// Syslog holds the identifier pointer for syslog entries, a flag to enable
// or disable console logging, and the minimum priority for messages to be logged.
type Syslog struct {
	ident_p        unsafe.Pointer
	consoleLogging bool
	minPriority    int
}

// NewSyslog initializes a new syslog handler.
//   - `ident` is a string that identifies the messages in the log.
//   - `option` is an integer specifying logging options (e.g., LOG_PID, LOG_CONS).
//   - `facility` is an integer specifying the syslog facility (e.g., LOG_USER).
//   - `minPriority` is the minimum priority (severity) that will be logged.
//     Messages with a numerical value higher than minPriority (less critical)
//     will be ignored.
func NewSyslog(ident string, option int, facility int, minPriority int) *Syslog {
	c_ident := C.CString(ident)
	C.openlog(c_ident, C.int(option), C.int(facility))
	return &Syslog{ident_p: unsafe.Pointer(c_ident), consoleLogging: false, minPriority: minPriority}
}

// EnableConsole enables logging to the console for this syslog instance.
func (s *Syslog) EnableConsole() {
	s.consoleLogging = true
}

// DisableConsole disables logging to the console for this syslog instance.
func (s *Syslog) DisableConsole() {
	s.consoleLogging = false
}

// Log sends a log message with the specified priority to syslog, and optionally
// to the console if enabled. If the message's priority is lower (less critical)
// than the configured minimum, the message is skipped.
// - `priority` is an integer specifying the message's priority (e.g., LOG_INFO, LOG_ERR).
// - `message` is the string message to log.
func (s *Syslog) Log(priority int, message string) {
	// Filter out messages that don't meet the minimum priority.
	// (Remember: lower numerical values are more critical.)
	if priority > s.minPriority {
		return
	}

	if s.consoleLogging {
		fmt.Printf("%s: %s\n", priorityToString(priority), message)
	}
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))
	C.syslog_helper(C.int(priority), cMessage)
}

// Info logs an informational message.
func (s *Syslog) Info(message string) {
	s.Log(LOG_INFO, message)
}

// Infof logs a formatted informational message.
func (s *Syslog) Infof(format string, a ...interface{}) {
	s.Info(fmt.Sprintf(format, a...))
}

// Warn logs a warning message.
func (s *Syslog) Warn(message string) {
	s.Log(LOG_WARN, message)
}

// Warnf logs a formatted warning message.
func (s *Syslog) Warnf(format string, a ...interface{}) {
	s.Warn(fmt.Sprintf(format, a...))
}

// Error logs an error message.
func (s *Syslog) Error(message string) {
	s.Log(LOG_ERR, message)
}

// Errorf logs a formatted error message.
func (s *Syslog) Errorf(format string, a ...interface{}) {
	s.Error(fmt.Sprintf(format, a...))
}

// Debug logs a debug message.
func (s *Syslog) Debug(message string) {
	s.Log(LOG_DEBUG, message)
}

// Debugf logs a formatted debug message.
func (s *Syslog) Debugf(format string, a ...interface{}) {
	s.Debug(fmt.Sprintf(format, a...))
}

// Crit logs a critical message and panics.
func (s *Syslog) Crit(message string) {
	s.Log(LOG_CRIT, message)
	panic(message)
}

// Critf logs a formatted critical message and panics.
func (s *Syslog) Critf(format string, a ...interface{}) {
	s.Crit(fmt.Sprintf(format, a...))
}

// Close releases resources associated with the syslog.
func (s *Syslog) Close() {
	C.free(s.ident_p)
	C.closelog()
}

// priorityToString maps a priority integer to its corresponding textual representation.
// Returns the textual representation of the given priority.
func priorityToString(priority int) string {
	switch priority {
	case LOG_INFO:
		return "INFO"
	case LOG_CRIT:
		return "CRITICAL"
	case LOG_WARN:
		return "WARNING"
	case LOG_ERR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
