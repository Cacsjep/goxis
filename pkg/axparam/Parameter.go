package axparam

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
	"unsafe"

	"github.com/Cacsjep/goxis/pkg/clib"
)

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a454ef604d7741f45804e708a56f7bf24
type AXParameter struct {
	Ptr              *C.AXParameter
	cStrings         []*clib.String
	parameterHandles map[string]cgo.Handle
}

// Creates a new AXParameter.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#adf2eefe79f53d60faede33ba71d5928c
func AXParameterNew(appName string) (*AXParameter, error) {
	cAppName := clib.NewString(&appName)
	defer cAppName.Free()
	cError := clib.NewError()
	axParam := C.ax_parameter_new(
		(*C.char)(cAppName.Ptr),
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)
	if err := cError.IsError(); err != nil {
		return nil, err
	}
	return &AXParameter{Ptr: axParam, parameterHandles: make(map[string]cgo.Handle)}, nil
}

// Adds a new parameter. Returns failure if the parameter already exists.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a9069e0a2a3c64cacd7d50f1408a9f5fa
func (axparameter *AXParameter) Add(name string, initialValue string, ptype string) error {
	cName := clib.NewString(&name)
	cInitialValue := clib.NewString(&initialValue)
	cptype := clib.NewString(&ptype)
	axparameter.cStrings = append(axparameter.cStrings, cName, cInitialValue, cptype)

	cError := clib.NewError()
	success := C.ax_parameter_add(
		axparameter.Ptr,
		(*C.char)(cName.Ptr),
		(*C.char)(cInitialValue.Ptr),
		(*C.char)(cptype.Ptr),
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)
	if err := cError.IsErrorOrNotSuccess(int(success), "Unable to add parameter"); err != nil {
		return err
	}
	return nil
}

// Removes a parameter. Returns FALSE if the parameter doesn't exist.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#aa6f49de80979e1f25ea9d98c449ac8c4
func (axparameter *AXParameter) Remove(name string) error {
	cName := clib.NewString(&name)
	defer cName.Free()
	cError := clib.NewError()
	success := C.ax_parameter_remove(
		axparameter.Ptr,
		(*C.char)(cName.Ptr),
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)
	if err := cError.IsErrorOrNotSuccess(int(success), "Unable to remove parameter"); err != nil {
		return err
	}
	return nil
}

// Sets the value of a parameter.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a3b767bcb3c99edf38b9fab3f38b7f2d7
func (axparameter *AXParameter) Set(name string, value string, doSync bool) error {
	cName := clib.NewString(&name)
	cValue := clib.NewString(&value)
	cError := clib.NewError()

	axparameter.cStrings = append(axparameter.cStrings, cName, cValue)
	success := C.ax_parameter_set(
		axparameter.Ptr,
		(*C.char)(cName.Ptr),
		(*C.char)(cValue.Ptr),
		(C.gboolean)(clib.GoBooleanToC(doSync)),
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)

	if err := cError.IsErrorOrNotSuccess(int(success), "Unable to set parameter"); err != nil {
		return err
	}
	return nil

}

// Retrieves the value of a parameter.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#aa7979d96d425189cfbf5fac6539bbc68
func (axparameter *AXParameter) Get(name string) (string, error) {
	cName := clib.NewString(&name)
	defer cName.Free()
	cValue := clib.NewAllocatableCString()
	defer cValue.Free()
	cError := clib.NewError()

	success := C.ax_parameter_get(
		axparameter.Ptr,
		(*C.char)(cName.Ptr),
		(**C.char)(unsafe.Pointer(cValue.Ptr)),
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)

	if err := cError.IsErrorOrNotSuccess(int(success), "Unable to get parameter"); err != nil {
		return "", err
	}
	return cValue.ToGolang(), nil
}

// Lists all parameters for the application.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#ace7b30d51c85d90509de80086b00d655
func (axparameter *AXParameter) List() ([]string, error) {
	var params []string
	cError := clib.NewError()
	gList := C.ax_parameter_list(
		axparameter.Ptr,
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)

	if gList == nil {
		return params, errors.New("Glist is nil")
	}
	defer C.g_list_free_full(gList, (C.GDestroyNotify)(C.g_free))

	if err := cError.IsError(); err != nil {
		return params, err
	}

	for l := gList; l != nil; l = l.next {
		param := C.g_list_nth_data(gList, C.guint(0))
		if param != nil {
			params = append(params, C.GoString((*C.char)(param)))
		}
	}

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

// Registers a callback function to be run whenever a parameter value is updated.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a43ae8096bdcad55ff26bf13eadc3b781
func (axparameter *AXParameter) RegisterCallback(name string, callback ParameterCallback, userdata any) error {
	cError := clib.NewError()
	data := &parameterCallbackData{Callback: callback, Userdata: userdata}
	handle := cgo.NewHandle(data)
	cName := clib.NewString(&name)
	defer cName.Free()

	success := C.ax_parameter_register_callback(
		axparameter.Ptr,
		(*C.char)(cName.Ptr),
		(*C.AXParameterCallback)(unsafe.Pointer(C.GoParameterCallback)),
		(C.gpointer)(unsafe.Pointer(handle)),
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)

	if err := cError.IsErrorOrNotSuccess(int(success), "Failed to register callback"); err != nil {
		handle.Delete()
		return err
	}

	axparameter.parameterHandles[name] = handle
	return nil
}

// Unregisters the parameter callback function.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a4e3b5330dfbe45fad22753c83a1ee065
func (axparameter *AXParameter) UnregisterCallback(name string) error {
	cName := clib.NewString(&name)
	defer cName.Free()

	if handle, exists := axparameter.parameterHandles[name]; exists {
		handle.Delete()
		delete(axparameter.parameterHandles, name)
	}
	C.ax_parameter_unregister_callback(
		axparameter.Ptr,
		(*C.char)(cName.Ptr),
	)
	return nil
}

// Frees an AXParameter.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a78ff4b5a312a1d9aab120436c116a5a2
func (axparameter *AXParameter) Free() {
	for _, cs := range axparameter.cStrings {
		cs.Free()
	}
	C.ax_parameter_free(axparameter.Ptr)
}
