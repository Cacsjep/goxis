package axlarod

/*
#cgo pkg-config: liblarod
#include "larod.h"
*/
import "C"
import (
	"fmt"
	"os"
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

func (l *Larod) Connection() *LarodConnection {
	return l.conn
}

func (l *Larod) DestroyModel(model *LarodModel) error {
	var cError *C.larodError

	for _, m := range model.maps {
		C.larodDestroyMap(&m.ptr)
	}

	C.larodDestroyModel(&model.ptr)

	if model.Job != nil {
		C.larodDestroyJobRequest(&model.Job.ptr)
	}

	if C.larodDestroyTensors(l.conn.ptr, &model.inputTensorPtr, C.uint(model.LarodModelIO.InputsCount), &cError) == C.bool(false) {
		return newLarodError(cError)
	}
	for _, t := range model.LarodModelIO.Inputs {

		if err := t.TmpFile.UnmapMemory(); err != nil {
			return err
		}

		if err := t.TmpFile.File.Close(); err != nil {
			return err
		}
		if err := os.Remove(t.TmpFile.File.Name()); err != nil {
			return err
		}
	}

	if C.larodDestroyTensors(l.conn.ptr, &model.outputTensorPtr, C.uint(model.LarodModelIO.OutputsCount), &cError) == C.bool(false) {
		return newLarodError(cError)
	}
	for _, t := range model.LarodModelIO.Outputs {

		if err := t.TmpFile.UnmapMemory(); err != nil {
			return err
		}

		if err := t.TmpFile.File.Close(); err != nil {
			return err
		}
		if err := os.Remove(t.TmpFile.File.Name()); err != nil {
			return err
		}
	}
	return nil
}
