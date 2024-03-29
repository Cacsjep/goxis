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
	LOG_PID  = 0x01
	LOG_CONS = 0x02
	LOG_USER = (1 << 3)

	LOG_INFO = 6
	LOG_CRIT = 2
	LOG_WARN = 4
	LOG_ERR  = 3
)

// Syslog struct holds the identifier pointer for syslog entries and a flag to enable
// or disable console logging.
type Syslog struct {
	ident_p        unsafe.Pointer
	consoleLogging bool
}

// NewSyslog initializes a new syslog handler with the specified identifier, option,
// and facility. Console logging is disabled by default.
// - `ident` is a string that identifies the messages in the log.
// - `option` is an integer specifying logging options (e.g., LOG_PID, LOG_CONS).
// - `facility` is an integer specifying the syslog facility.
func NewSyslog(ident string, option int, facility int) *Syslog {
	c_ident := C.CString(ident)
	C.openlog(c_ident, C.int(option), C.int(facility))
	return &Syslog{ident_p: unsafe.Pointer(c_ident), consoleLogging: false}
}

// EnableConsole enables logging to the console for this syslog instance.
func (s *Syslog) EnableConsole() {
	s.consoleLogging = true
}

// DisableConsole disables logging to the console for this syslog instance.
func (s *Syslog) DisableConsole() {
	s.consoleLogging = false
}

// Log sends a log message with the specified priority to the syslog, and optionally
// to the console if console logging is enabled.
// - `priority` is an integer specifying the message's priority (e.g., LOG_INFO, LOG_ERR).
// - `message` is the string message to log.
func (s *Syslog) Log(priority int, message string) {
	if s.consoleLogging {
		fmt.Printf("%s: %s\n", priorityToString(priority), message)
	}
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))
	C.syslog_helper(C.int(priority), cMessage)
}

// Info logs an informational message to the syslog, and optionally to the console.
func (s *Syslog) Info(message string) {
	s.Log(LOG_INFO, message)
}

// Infof logs an informational message, formatted according to a format specifier, to the syslog,
// and optionally to the console.
func (s *Syslog) Infof(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	s.Info(message)
}

// Warn logs a warning message to the syslog, and optionally to the console.
func (s *Syslog) Warn(message string) {
	s.Log(LOG_WARN, message)
}

// Warnf logs a warning message, formatted according to a format specifier, to the syslog,
// and optionally to the console.
func (s *Syslog) Warnf(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	s.Warn(message)
}

// Error logs an error message to the syslog, and optionally to the console.
func (s *Syslog) Error(message string) {
	s.Log(LOG_ERR, message)
}

// Errorf logs an error message, formatted according to a format specifier, to the syslog,
// and optionally to the console.
func (s *Syslog) Errorf(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	s.Error(message)
}

// Crit logs a critical message to the syslog and panics. The message is also logged to the console
// if console logging is enabled.
func (s *Syslog) Crit(message string) {
	s.Log(LOG_CRIT, message)
	panic(message)
}

// Critf logs a critical message, formatted according to a format specifier, to the syslog and panics.
// The message is also logged to the console if console logging is enabled.
func (s *Syslog) Critf(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	s.Crit(message) // Note that s.Crit will log the message and then panic.
}

// Close releases resources associated with the syslog, specifically freeing the identifier.
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
