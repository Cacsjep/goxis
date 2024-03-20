package acap

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axstorage/html/index.html

/*
#cgo pkg-config: glib-2.0 gio-2.0 axstorage
#include <glib.h>
#include <axsdk/axstorage.h>
extern void GoStorageSubscriptionCallback(char *storage_id, gpointer user_data, GError *error);
extern void GoStorageSetupCallback(AXStorage *storage, gpointer user_data, GError *error);
extern void GoStorageReleaseCallback(gpointer user_data, GError *error);
*/
import "C"
import (
	"errors"
	"fmt"
	"runtime/cgo"
	"unsafe"
)

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axstorage/html/ax__storage_8h.html#aa14f17bf10b2dcc57d1ca16759da9b8f
type AXStorage struct {
	Ptr      *C.AXStorage
	cStrings []*cString
}

// AXStorageStatusEventId represents the list of events for AXStorage.
type AXStorageStatusEventId int

const (
	AXStorageAvailableEvent   AXStorageStatusEventId = C.AX_STORAGE_AVAILABLE_EVENT
	AXStorageExitingEvent     AXStorageStatusEventId = C.AX_STORAGE_EXITING_EVENT
	AXStorageWritableEvent    AXStorageStatusEventId = C.AX_STORAGE_WRITABLE_EVENT
	AXStorageFullEvent        AXStorageStatusEventId = C.AX_STORAGE_FULL_EVENT
	AXStorageStatusEventIDEnd AXStorageStatusEventId = C.AX_STORAGE_STATUS_EVENT_ID_END
)

// Just a int for the Storage type
type AXStorageType int

const (
	AXStorageTypeLocal    AXStorageType = iota
	AXStorageTypeExternal AXStorageType = iota
	AXStorageTypeUnkown   AXStorageType = iota
)

// Just a string for the storage id
type StorageId string

// subscription callbacks
var subscriptionHandles map[int]cgo.Handle = make(map[int]cgo.Handle)

// Lists all connected storage devices. The returned list and its members must be freed by the caller. Use g_free for the members and g_list_free for the list.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axstorage/html/ax__storage_8h.html#aa32bdb91706c44a86071c1c14889536f
func AxStorageList() ([]StorageId, error) {
	var storage_ids []StorageId
	var gerr *C.GError

	disks_list_ptr := C.ax_storage_list(&gerr)
	if disks_list_ptr == nil {
		return nil, errors.New("Unable to get storage list, ax_storage_list retuns NULL")
	}

	if err := newGError(gerr); err != nil {
		return nil, err
	}

	storagesIdsPrt := uintptr(unsafe.Pointer(disks_list_ptr))
	storagesIdsList := WrapList(storagesIdsPrt)
	storagesIdsList.DataWrapper(wrapString)
	storagesIdsList.Foreach(func(item interface{}) {
		storage_id, ok := item.(string)
		if !ok {
			panic("AxStorageList: item is not of type string")
		}
		storage_ids = append(storage_ids, (StorageId)(storage_id))
	})
	return storage_ids, nil
}

func wrapString(ptr unsafe.Pointer) interface{} {
	return C.GoString((*C.char)(ptr))
}

// GetPath retrieves the location on the storage where the client should save its files.
func (s *AXStorage) GetPath() (string, error) {
	var gerr *C.GError
	cPath := C.ax_storage_get_path(s.Ptr, &gerr)
	if cPath == nil {
		return "", newGError(gerr)
	}
	defer C.g_free(C.gpointer(cPath))
	return C.GoString(cPath), nil
}

// AXStorageGetStatus returns the status of the provided event for a storage.
func AXStorageGetStatus(storage_id StorageId, event AXStorageStatusEventId) (bool, error) {
	var gerr *C.GError
	cStorageId := C.CString(string(storage_id))
	defer C.free(unsafe.Pointer(cStorageId))

	status := C.ax_storage_get_status(cStorageId, C.AXStorageStatusEventId(event), &gerr)
	if gerr != nil {
		return false, newGError(gerr)
	}

	return ctoGoBoolean(status), nil
}

// GetStorageId gets the storage_id from the provided AXStorage.
func (s *AXStorage) GetStorageId() (StorageId, error) {
	var gerr *C.GError
	cID := C.ax_storage_get_storage_id(s.Ptr, &gerr)
	if cID == nil {
		return "", newGError(gerr)
	}
	defer C.g_free(C.gpointer(cID))
	return (StorageId)(C.GoString(cID)), nil
}

// GetType returns the storage type for the given AXStorage.
func (s *AXStorage) GetType() (AXStorageType, error) {
	var gerr *C.GError
	cType := C.ax_storage_get_type(s.Ptr, &gerr)
	if gerr != nil {
		return AXStorageTypeUnkown, newGError(gerr)
	}
	return AXStorageType(cType), nil
}

func NewAxStorageFromC(storage *C.AXStorage) *AXStorage {
	return &AXStorage{Ptr: storage}
}

// StorageSubscriptionCallback is the Go equivalent of AXStorageSubscriptionCallback.
type StorageSubscriptionCallback func(storageID StorageId, userdata any, err error)

// StorageSetupCallback is the Go equivalent of AXStorageSetupCallback.
type StorageSetupCallback func(storage *AXStorage, userdata any, err error)

// StorageReleaseCallback is the Go equivalent of AXStorageReleaseCallback.
type StorageReleaseCallback func(userdata any, err error)

type StorageSubscriptionCallbackData struct {
	Callback StorageSubscriptionCallback
	Userdata any
}

type StorageSetupCallbackData struct {
	Callback StorageSetupCallback
	Userdata any
}

type StorageReleaseCallbackData struct {
	Callback StorageReleaseCallback
	Userdata any
}

type DiskItem struct {
	SubscriptionId int           // Subscription ID for storage events
	Setup          bool          // TRUE: storage was set up async, FALSE otherwise
	Writable       bool          // Storage is writable or not
	Available      bool          // Storage is available or not
	Full           bool          // Storage device is full or not
	Exiting        bool          // Storage is exiting (going to disappear) or not
	Storage        *AXStorage    // AXStorage reference
	StorageType    AXStorageType // Storage type
	StorageId      StorageId     // Storage device name
	StoragePath    string        // Storage path
}

func (d *DiskItem) String() string {
	return fmt.Sprintf("%s, Writeable: %t, Available: %t, Setup: %t, Exiting: %t, Full: %t, StoragePath: %s", d.StorageId, d.Writable, d.Available, d.Setup, d.Exiting, d.Full, d.StoragePath)
}

// NewDiskItem returns a new disk item
func NewDiskItem(storageID StorageId, subscriptionId int) *DiskItem {
	return &DiskItem{
		StorageId:      storageID,
		SubscriptionId: subscriptionId,
	}
}

// Update the Available, Writable, Full, Exiting Event fields
func UpdateDiskItemEvents(d *DiskItem) error {
	var err error
	if d.Available, err = AXStorageGetStatus(d.StorageId, AXStorageAvailableEvent); err != nil {
		return err
	}
	if d.Writable, err = AXStorageGetStatus(d.StorageId, AXStorageWritableEvent); err != nil {
		return err
	}
	if d.Full, err = AXStorageGetStatus(d.StorageId, AXStorageFullEvent); err != nil {
		return err
	}
	if d.Exiting, err = AXStorageGetStatus(d.StorageId, AXStorageExitingEvent); err != nil {
		return err
	}
	return nil
}

//export GoStorageSubscriptionCallback
func GoStorageSubscriptionCallback(storageID *C.char, user_data unsafe.Pointer, gError *C.GError) {
	var err error
	handle := cgo.Handle(user_data)
	callbackData := handle.Value().(*StorageSubscriptionCallbackData)
	if gError != nil {
		err = newGError(gError)
	}
	callbackData.Callback((StorageId)(C.GoString(storageID)), callbackData.Userdata, err)
}

//export GoStorageSetupCallback
func GoStorageSetupCallback(storage *C.AXStorage, user_data unsafe.Pointer, gError *C.GError) {
	var err error
	handle := cgo.Handle(user_data)
	callbackData := handle.Value().(*StorageSetupCallbackData)
	if gError != nil {
		err = newGError(gError)
	}
	callbackData.Callback(&AXStorage{Ptr: storage}, callbackData.Userdata, err)
	handle.Delete()
}

//export GoStorageReleaseCallback
func GoStorageReleaseCallback(user_data unsafe.Pointer, gError *C.GError) {
	var err error
	handle := cgo.Handle(user_data)
	callbackData := handle.Value().(*StorageReleaseCallbackData)
	if gError != nil {
		err = newGError(gError)
	}
	callbackData.Callback(callbackData.Userdata, err)
	handle.Delete()
}

// Subscribe subscribes to storage events for the provided storage ID.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axstorage/html/ax__storage_8h.html#ae64f800e6d88b54cf68cce925bb16312
func AxStorageSubscribe(storageID StorageId, callback StorageSubscriptionCallback, userdata any) (subscriptionId int, err error) {
	cStorageID := C.CString(string(storageID))
	defer C.free(unsafe.Pointer(cStorageID))

	var gerr *C.GError
	data := &StorageSubscriptionCallbackData{
		Callback: callback,
		Userdata: userdata,
	}
	handle := cgo.NewHandle(data)

	subscriptionID := C.ax_storage_subscribe(cStorageID, (C.AXStorageSubscriptionCallback)(C.GoStorageSubscriptionCallback), C.gpointer(handle), &gerr)
	if subscriptionID == 0 {
		return 0, newGError(gerr) // Assume newGError converts a GError to a Go error.
	}

	subscriptionHandles[int(subscriptionID)] = handle
	return int(subscriptionID), nil
}

// Unsubscribe stops subscribing to storage events for the provided subscription ID.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axstorage/html/ax__storage_8h.html#a49676a4baae47e9407f3641fefcd365c
func AxStorageUnsubscribe(subscriptionID int) error {
	var gerr *C.GError
	if handle, exists := subscriptionHandles[subscriptionID]; exists {
		handle.Delete()
		delete(subscriptionHandles, subscriptionID)
	}
	success := C.ax_storage_unsubscribe(C.guint(subscriptionID), &gerr)
	if success == C.FALSE {
		return newGError(gerr)
	}
	return nil
}

// SetupAsync sets up storage for use asynchronously.
// Setup storage for use asynchronously. This method must be called before the storage is to be used in any way (for instance read or write).
// When done using the storage, AxStorageReleaseAsync must be called.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axstorage/html/ax__storage_8h.html#abb7bf1c1cab1961dc4c2fffedaaa73cd
func AxStorageSetupAsync(storageID StorageId, callback StorageSetupCallback, userdata any) error {
	cStorageID := C.CString(string(storageID))
	defer C.free(unsafe.Pointer(cStorageID))

	var gerr *C.GError
	data := &StorageSetupCallbackData{
		Callback: callback,
		Userdata: userdata,
	}
	handle := cgo.NewHandle(data)
	success := C.ax_storage_setup_async(cStorageID, (C.AXStorageSetupCallback)(C.GoStorageSetupCallback), C.gpointer(handle), &gerr)
	if success == C.FALSE {
		return newGError(gerr)
	}
	return nil
}

// ReleaseAsync releases the use of storage asynchronously.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axstorage/html/ax__storage_8h.html#a27909ecc692b78af43ce23bf0369fe95
func (s *AXStorage) AxStorageReleaseAsync(callback StorageReleaseCallback, userdata any) error {
	var gerr *C.GError
	data := &StorageReleaseCallbackData{
		Callback: callback,
		Userdata: userdata,
	}
	handle := cgo.NewHandle(data)
	success := C.ax_storage_release_async(s.Ptr, (C.AXStorageReleaseCallback)(C.GoStorageReleaseCallback), C.gpointer(handle), &gerr)
	if success == C.FALSE {
		return newGError(gerr)
	}
	return nil
}
