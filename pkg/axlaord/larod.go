package axlarod

/*
#cgo pkg-config: liblarod
#include "larod.h"
*/
import "C"
import (
	"fmt"
)

type Larod struct {
	conn    *LarodConnection
	Devices []*LarodDevice
}

type LarodConnection struct {
	ptr *C.larodConnection
}

type LarodDevice struct {
	ptr  *C.larodDevice
	Name string
}

// Connect and load devices
func (l *Larod) Initalize() error {
	if err := l.Connect(); err != nil {
		return fmt.Errorf("failed to connect to Larod: %w", err)
	}
	if _, err := l.ListDevices(); err != nil {
		return fmt.Errorf("failed to list devices: %w", err)
	}
	return nil
}

func NewLarod() *Larod {
	return &Larod{}
}

func (l *Larod) Connect() error {
	var conn *C.larodConnection
	var cError *C.larodError
	if C.larodConnect(&conn, &cError) == C.bool(false) {
		return newLarodError(cError)
	}
	l.conn = &LarodConnection{ptr: conn}
	return nil
}

func (l *Larod) Disconnect() error {
	var cError *C.larodError
	if C.larodDisconnect(&l.conn.ptr, &cError) == C.bool(false) {
		return newLarodError(cError)
	}
	return nil
}
