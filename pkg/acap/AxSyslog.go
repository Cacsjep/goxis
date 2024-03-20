package acap

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

type Syslog struct {
	ident_p unsafe.Pointer
}

func NewSyslog(ident string, option int, facility int) *Syslog {
	c_ident := C.CString(ident)
	C.openlog(c_ident, C.int(option), C.int(facility))
	return &Syslog{ident_p: unsafe.Pointer(c_ident)}
}

func (s *Syslog) Close() {
	C.free(s.ident_p)
	C.closelog()
}

func (s *Syslog) Log(priority int, message string) {
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))
	C.syslog_helper(C.int(priority), cMessage)
}

func (s *Syslog) Info(message string) {
	s.Log(LOG_INFO, message)
}

func (s *Syslog) Infof(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	s.Info(message)
}

func (s *Syslog) Warn(message string) {
	s.Log(LOG_WARN, message)
}

func (s *Syslog) Warnf(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	s.Warn(message)
}

func (s *Syslog) Error(message string) {
	s.Log(LOG_ERR, message)
}

func (s *Syslog) Errorf(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	s.Error(message)
}

// Crit also do a golang panic
func (s *Syslog) Crit(message string) {
	s.Log(LOG_CRIT, message)
	panic(message)
}

func (s *Syslog) Critf(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	s.Crit(message) // Note that s.Crit will log the message and then panic.
}
