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
type StorageSubscriptionCallback func(storageID StorageId, userdata any, err error, cleanup func())

// StorageSetupCallback is the Go equivalent of AXStorageSetupCallback.
type StorageSetupCallback func(storage *AXStorage, userdata any, err error, cleanup func())

// StorageReleaseCallback is the Go equivalent of AXStorageReleaseCallback.
type StorageReleaseCallback func(userdata any, err error, cleanup func())

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
	SubscriptionId uint          // Subscription ID for storage events
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

func NewDiskItemFromSubscription(storageID StorageId) (*DiskItem, error) {
	available, err := AXStorageGetStatus(storageID, AXStorageAvailableEvent)
	if err != nil {
		return nil, err
	}
	writable, err := AXStorageGetStatus(storageID, AXStorageWritableEvent)
	if err != nil {
		return nil, err
	}
	full, err := AXStorageGetStatus(storageID, AXStorageFullEvent)
	if err != nil {
		return nil, err
	}
	exiting, err := AXStorageGetStatus(storageID, AXStorageExitingEvent)
	if err != nil {
		return nil, err
	}

	return &DiskItem{
		Available: available,
		Writable:  writable,
		Full:      full,
		Exiting:   exiting,
		StorageId: storageID,
	}, nil
}

//export GoStorageSubscriptionCallback
func GoStorageSubscriptionCallback(storageID *C.char, user_data unsafe.Pointer, gError *C.GError) {
	handle := cgo.Handle(user_data)
	callbackData := handle.Value().(*StorageSubscriptionCallbackData)
	callbackData.Callback((StorageId)(C.GoString(storageID)), callbackData.Userdata, newGError(gError).AsError(), handle.Delete)
}

//export GoStorageSetupCallback
func GoStorageSetupCallback(storage *C.AXStorage, user_data unsafe.Pointer, gError *C.GError) {
	handle := cgo.Handle(user_data)
	callbackData := handle.Value().(*StorageSetupCallbackData)
	callbackData.Callback(&AXStorage{Ptr: storage}, callbackData.Userdata, newGError(gError).AsError(), handle.Delete)
}

//export GoStorageReleaseCallback
func GoStorageReleaseCallback(user_data unsafe.Pointer, gError *C.GError) {
	handle := cgo.Handle(user_data)
	callbackData := handle.Value().(StorageReleaseCallbackData)
	callbackData.Callback(callbackData.Userdata, newGError(gError).AsError(), handle.Delete)
}

// Subscribe subscribes to storage events for the provided storage ID.
func AxStorageSubscribe(storageID StorageId, callback StorageSubscriptionCallback, userdata any) (uint, error) {
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
	return uint(subscriptionID), nil
}

// Unsubscribe stops subscribing to storage events for the provided subscription ID.
func AxStorageUnsubscribe(subscriptionID int) error {
	var gerr *C.GError
	success := C.ax_storage_unsubscribe(C.guint(subscriptionID), &gerr)
	if success == C.FALSE {
		return newGError(gerr)
	}
	return nil
}

// SetupAsync sets up storage for use asynchronously.
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
