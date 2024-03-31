/*
Package axparameter provides a Go interface to the AXParameter library, enabling the management of application-specific parameters on Axis devices.
This package allows applications to add, remove, set, and get parameters, as well as to register callbacks for parameter changes.
It is particularly useful for applications that need to manage configuration settings dynamically or respond to changes in their environment.

Usage:
1. Initialize an AXParameter instance for your application.
2. Add parameters with initial values and types, optionally using control words for special behavior.
3. Register callbacks for parameters you need to monitor for changes, providing a function that will react to updates.
4. Use the Set and Get functions to update parameters dynamically and retrieve their current values as needed.
5. Properly free the AXParameter instance when your application shuts down or no longer needs to manage parameters.

This package is essential for ACAP developers who need to manage application parameters, providing a flexible and dynamic way to handle configuration settings. It abstracts away the complexities of the underlying C library, offering a simple and idiomatic Go API.

Requirements:
- An Axis device capable of running ACAP applications.
- The ACAP SDK and appropriate development tools for compilation and deployment.

Example:
See the package examples for detailed usage patterns, including how to add parameters, register callbacks for updates, and dynamically change parameter values based on application logic.
*/
package axparameter

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/index.html

/*
#cgo pkg-config: glib-2.0 gio-2.0 axparameter
#include <glib.h>
#include <axsdk/axparameter.h>
extern void GoParameterCallback(gchar *name, gchar *value, gpointer user_data);
*/
import "C"
import (
	"errors"
	"fmt"
	"runtime/cgo"
	"strconv"
	"unsafe"

	"github.com/Cacsjep/goxis/pkg/glib"
)

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a454ef604d7741f45804e708a56f7bf24
type AXParameter struct {
	Ptr              *C.AXParameter
	parameterHandles map[string]cgo.Handle
}

// Creates a new AXParameter.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#adf2eefe79f53d60faede33ba71d5928c
func AXParameterNew(appName string) (*AXParameter, error) {
	var axParam *C.AXParameter
	var gerr *C.GError

	cName := C.CString(appName)
	defer C.free(unsafe.Pointer(cName))

	if axParam = C.ax_parameter_new(cName, &gerr); axParam == nil {
		return nil, newGError(gerr)
	}

	return &AXParameter{Ptr: axParam, parameterHandles: make(map[string]cgo.Handle)}, nil
}

// Adds a new parameter. Returns failure if the parameter already exists.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a9069e0a2a3c64cacd7d50f1408a9f5fa
func (axp *AXParameter) Add(name string, initialValue string, ptype string) error {
	var gerr *C.GError

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cInitialValue := C.CString(initialValue)
	defer C.free(unsafe.Pointer(cInitialValue))

	cPtype := C.CString(ptype)
	defer C.free(unsafe.Pointer(cPtype))

	if int(C.ax_parameter_add(axp.Ptr, cName, cInitialValue, cPtype, &gerr)) == 0 {
		return newGError(gerr)
	}
	return nil
}

// Removes a parameter. Returns FALSE if the parameter doesn't exist.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#aa6f49de80979e1f25ea9d98c449ac8c4
func (axp *AXParameter) Remove(name string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var gerr *C.GError
	if int(C.ax_parameter_remove(axp.Ptr, cName, &gerr)) == 0 {
		return newGError(gerr)
	}
	return nil
}

// Sets the value of a parameter.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a3b767bcb3c99edf38b9fab3f38b7f2d7
func (axp *AXParameter) Set(name string, value string, doSync bool) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	var gerr *C.GError
	if int(C.ax_parameter_set(axp.Ptr, cName, cValue, C.gboolean(map[bool]int{true: 1, false: 0}[doSync]), &gerr)) == 0 {
		return newGError(gerr)
	}
	return nil

}

// Retrieves the value of a parameter.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#aa7979d96d425189cfbf5fac6539bbc68
func (axp *AXParameter) Get(name string) (string, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var cValue *C.char
	defer func() {
		if cValue != nil {
			C.free(unsafe.Pointer(cValue))
		}
	}()

	var gerr *C.GError
	if int(C.ax_parameter_get(axp.Ptr, cName, &cValue, &gerr)) == 0 {
		return "", newGError(gerr)
	}
	return C.GoString(cValue), nil
}

func (axp *AXParameter) GetAsFloat(name string) (float64, error) {
	var err error
	var str_val string

	if str_val, err = axp.Get(name); err != nil {
		return 0, err
	}
	return strconv.ParseFloat(str_val, 64)
}

func (axp *AXParameter) GetAsInt(name string) (int, error) {
	var err error
	var str_val string

	if str_val, err = axp.Get(name); err != nil {
		return 0, err
	}
	return strconv.Atoi(str_val)
}

// Lists all parameters for the application.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#ace7b30d51c85d90509de80086b00d655
func (axp *AXParameter) List() ([]string, error) {
	var params []string
	var gerr *C.GError

	gList_ptr := C.ax_parameter_list(axp.Ptr, &gerr)
	if gList_ptr == nil {
		return params, errors.New("Unable to list parameters, ax_parameter_list returned nil")
	}

	if err := newGError(gerr); err != nil {
		return params, err
	}
	paramsListPtr := uintptr(unsafe.Pointer(gList_ptr))
	paramsList := glib.WrapList(paramsListPtr)
	paramsList.DataWrapper(wrapString)
	paramsList.Data()
	paramsList.Foreach(func(item interface{}) {
		param, ok := item.(string)
		if !ok {
			panic("param: item is not of type string")
		}
		params = append(params, param)
	})
	paramsList.Free()
	return params, nil
}

// The typedef for a callback function registered by ax_parameter_register_callback()
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a82ba0fc58e5e8749b1313825b6a7a670
type ParameterCallback func(name string, value string, userdata any)

type parameterCallbackData struct {
	Callback ParameterCallback
	Userdata any
}

//export GoParameterCallback
func GoParameterCallback(name *C.gchar, value *C.gchar, user_data unsafe.Pointer) {
	h := cgo.Handle(user_data)
	data := h.Value().(*parameterCallbackData)
	if data == nil {
		fmt.Println("Error: in value conv (GoParameterCallback)")
		return
	}
	data.Callback(C.GoString(name), C.GoString(value), data.Userdata)
}

// Registers a callback function to be run whenever the given named parameter is changed, eg value updated.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a43ae8096bdcad55ff26bf13eadc3b781
func (axp *AXParameter) RegisterCallback(name string, callback ParameterCallback, userdata any) error {
	var gerr *C.GError
	data := &parameterCallbackData{Callback: callback, Userdata: userdata}
	handle := cgo.NewHandle(data)

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	success := C.ax_parameter_register_callback(
		axp.Ptr,
		cName,
		(*C.AXParameterCallback)(unsafe.Pointer(C.GoParameterCallback)),
		(C.gpointer)(unsafe.Pointer(handle)),
		&gerr,
	)

	if int(success) == 0 {
		return newGError(gerr)
	}
	axp.parameterHandles[name] = handle
	return nil
}

// Unregisters the parameter callback function.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a4e3b5330dfbe45fad22753c83a1ee065
func (axp *AXParameter) UnregisterCallback(name string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	if handle, exists := axp.parameterHandles[name]; exists {
		handle.Delete()
		delete(axp.parameterHandles, name)
	}
	C.ax_parameter_unregister_callback(
		axp.Ptr,
		cName,
	)
	return nil
}

// Frees an AXParameter.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a78ff4b5a312a1d9aab120436c116a5a2
func (axp *AXParameter) Free() {
	C.ax_parameter_free(axp.Ptr)
}

func wrapString(ptr unsafe.Pointer) interface{} {
	return C.GoString((*C.char)(ptr))
}
