package acap

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
)

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a454ef604d7741f45804e708a56f7bf24
type AXParameter struct {
	Ptr              *C.AXParameter
	cStrings         []*cString
	parameterHandles map[string]cgo.Handle
}

// Creates a new AXParameter.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#adf2eefe79f53d60faede33ba71d5928c
func AXParameterNew(appName *string) (*AXParameter, error) {
	var axParam *C.AXParameter
	var gerr *C.GError
	cAppName := newString(appName)
	defer cAppName.Free()

	if axParam = C.ax_parameter_new(cAppName.Ptr, &gerr); axParam == nil {
		return nil, newGError(gerr)
	}

	return &AXParameter{Ptr: axParam, parameterHandles: make(map[string]cgo.Handle)}, nil
}

// Adds a new parameter. Returns failure if the parameter already exists.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a9069e0a2a3c64cacd7d50f1408a9f5fa
func (axp *AXParameter) Add(name string, initialValue string, ptype string) error {
	var gerr *C.GError
	cName := newString(&name)
	cInitialValue := newString(&initialValue)
	cptype := newString(&ptype)
	axp.cStrings = append(axp.cStrings, cName, cInitialValue, cptype)
	if int(C.ax_parameter_add(axp.Ptr, cName.Ptr, cInitialValue.Ptr, cptype.Ptr, &gerr)) == 0 {
		return newGError(gerr)
	}
	return nil
}

// Removes a parameter. Returns FALSE if the parameter doesn't exist.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#aa6f49de80979e1f25ea9d98c449ac8c4
func (axp *AXParameter) Remove(name string) error {
	cName := newString(&name)
	defer cName.Free()
	var gerr *C.GError
	if int(C.ax_parameter_remove(axp.Ptr, cName.Ptr, &gerr)) == 0 {
		return newGError(gerr)
	}
	return nil
}

// Sets the value of a parameter.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a3b767bcb3c99edf38b9fab3f38b7f2d7
func (axp *AXParameter) Set(name string, value string, doSync bool) error {
	cName := newString(&name)
	cValue := newString(&value)
	var gerr *C.GError
	axp.cStrings = append(axp.cStrings, cName, cValue)
	if int(C.ax_parameter_set(axp.Ptr, cName.Ptr, cValue.Ptr, goBooleanToC(doSync), &gerr)) == 0 {
		return newGError(gerr)
	}
	return nil

}

// Retrieves the value of a parameter.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#aa7979d96d425189cfbf5fac6539bbc68
func (axp *AXParameter) Get(name string) (string, error) {
	cName := newString(&name)
	defer cName.Free()
	cValue := newAllocatableCString()
	defer cValue.Free()
	var gerr *C.GError
	if int(C.ax_parameter_get(axp.Ptr, cName.Ptr, cValue.Ptr, &gerr)) == 0 {
		return "", newGError(gerr)
	}
	return cValue.ToGolang(), nil
}

// Lists all parameters for the application.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#ace7b30d51c85d90509de80086b00d655
func (axp *AXParameter) List() ([]string, error) {
	var params []string
	var gerr *C.GError

	gList := C.ax_parameter_list(
		axp.Ptr,
		&gerr,
	)

	if gList == nil {
		return params, errors.New("Glist is nil")
	}
	defer C.g_list_free_full(gList, (C.GDestroyNotify)(C.g_free))

	if err := newGError(gerr); err != nil {
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

// Registers a callback function to be run whenever the given named parameter is changed, eg value updated.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a43ae8096bdcad55ff26bf13eadc3b781
func (axp *AXParameter) RegisterCallback(name string, callback ParameterCallback, userdata any) error {
	var gerr *C.GError
	data := &parameterCallbackData{Callback: callback, Userdata: userdata}
	handle := cgo.NewHandle(data)
	cName := newString(&name)
	defer cName.Free()

	success := C.ax_parameter_register_callback(
		axp.Ptr,
		cName.Ptr,
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
	cName := newString(&name)
	defer cName.Free()
	if handle, exists := axp.parameterHandles[name]; exists {
		handle.Delete()
		delete(axp.parameterHandles, name)
	}
	C.ax_parameter_unregister_callback(
		axp.Ptr,
		cName.Ptr,
	)
	return nil
}

// Frees an AXParameter.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axparameter/html/ax__parameter_8h.html#a78ff4b5a312a1d9aab120436c116a5a2
func (axp *AXParameter) Free() {
	for _, cs := range axp.cStrings {
		cs.Free()
	}
	C.ax_parameter_free(axp.Ptr)
}
