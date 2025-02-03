// Package axlarod provides a Go wrapper for interacting with the Larod API, facilitating the management and execution of deep learning models on Axis devices.
//
// The package includes types and functions for:
// - Establishing and terminating connections with the Larod service.
// - Querying, listing, and managing devices capable of running inference models.
// - Managing the lifecycle of models, including creation, execution, and resource cleanup.
// - Manipulating and transferring data between Go structures and the underlying hardware through Larod's memory management systems.
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

// NewLarod creates a new Larod instance with uninitialized connection and device list.
func NewLarod() *Larod {
	return &Larod{
		Devices: make([]*LarodDevice, 0),
	}
}

// Initialize establishes a connection to the Larod service and retrieves a list of available devices.
// It handles errors that occur during the connection or device listing process.
func (l *Larod) Initalize() error {
	if err := l.Connect(); err != nil {
		return fmt.Errorf("failed to connect to Larod: %w", err)
	}
	if _, err := l.ListDevices(); err != nil {
		return fmt.Errorf("failed to list devices: %w", err)
	}
	return nil
}

// Connect attempts to establish a connection with the Larod service.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/larod/html/larod_8h.html#ab7337d6d6663c0e45326022503927851
func (l *Larod) Connect() error {
	var conn *C.larodConnection
	var cError *C.larodError
	if C.larodConnect(&conn, &cError) == C.bool(false) {
		return newLarodError(cError)
	}
	l.conn = &LarodConnection{ptr: conn}
	return nil
}

// Disconnect terminates the connection with the Larod service.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/larod/html/larod_8h.html#ab8f97b4b4d15798384ca25f32ca77bba
func (l *Larod) Disconnect() error {
	var cError *C.larodError
	if l.conn == nil {
		return nil
	}
	if l.conn.ptr == nil {
		return nil
	}
	if C.larodDisconnect(&l.conn.ptr, &cError) == C.bool(false) {
		return newLarodError(cError)
	}
	return nil
}

// Connection returns the current connection to the Larod service.
func (l *Larod) Connection() *LarodConnection {
	return l.conn
}

// DestroyInputTensors cleans up resources associated with input tensors of a LarodModel.
func (l *Larod) DestroyInputTensors(model *LarodModel) error {
	var cError *C.larodError
	if C.larodDestroyTensors(l.conn.ptr, &model.inputTensorsPtr, C.size_t(model.InputsCount), &cError) == C.bool(false) {
		return newLarodError(cError)
	}
	return nil
}

// DestroyOutputTensors cleans up resources associated with output tensors of a LarodModel.
func (l *Larod) DestroyOutputTensors(model *LarodModel) error {
	var cError *C.larodError
	if C.larodDestroyTensors(l.conn.ptr, &model.outputTensorsPtr, C.size_t(model.OutputsCount), &cError) == C.bool(false) {
		return newLarodError(cError)
	}
	return nil
}

// DestroyModel cleans up resources associated with a LarodModel.
func (l *Larod) DestroyModel(model *LarodModel) error {
	for _, m := range model.maps {
		m.Destroy()
	}

	model.Destroy()

	if model.Job != nil {
		model.Job.Destroy()
	}

	if err := l.DestroyInputTensors(model); err != nil {
		return err
	}

	for _, t := range model.Inputs {
		if err := t.MemMapFile.UnmapMemory(); err != nil {
			return err
		}
		if err := t.MemMapFile.File.Close(); err != nil {
			return err
		}
	}

	if err := l.DestroyOutputTensors(model); err != nil {
		return err
	}
	for _, t := range model.Outputs {
		if err := t.MemMapFile.UnmapMemory(); err != nil {
			return err
		}
		if err := t.MemMapFile.File.Close(); err != nil {
			return err
		}
	}

	return nil
}
