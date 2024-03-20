package goxis

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Cacsjep/goxis/pkg/acap"
)

// StorageProvider represents a handler for managing storage devices, enabling operations
// such as file writing, removal, and subscriptions to storage events.
type StorageProvider struct {
	app              *AcapApplication // Reference to the main application.
	DiskItems        []*acap.DiskItem // List of disk items representing storage devices.
	subscribtions    []int            // Subscription list for unsubscribe
	DiskItemsEvents  chan *acap.DiskItem
	UseChannelEvents bool
}

// NewStorageProvider initializes and returns a new StorageProvider associated with a given AcapApplication.
// Whem useChannelEvents is true the DiskItemsEvents channel got events from subscriptions callbacks
// in form of *acap.DiskItem.
// Its an unbufferd channel with cap 10
func (a *AcapApplication) NewStorageProvider(useChannelEvents bool) *StorageProvider {
	return &StorageProvider{
		app:              a,
		DiskItemsEvents:  make(chan *acap.DiskItem, 10),
		UseChannelEvents: useChannelEvents,
	}
}

// RwError represents various errors that can occur during read/write operations on storage devices.
type RwError int

const (
	RWErrorNone RwError = iota
	RWErrorFull
	RWErrorNotAvalible
	RWErrorNotSetuped
	RWErrorNotWriteable
	RWErrorNotUpdateable
	RWErrorOs
)

// RwResult encapsulates the result of a read/write operation, including any errors that occurred.
type RwResult struct {
	RwError RwError // Specific write error encountered, if any.
	Error   error   // General error encountered during the operation.
	Data    []byte  // Data read from storage, applicable for read operations.
}

// checkRwPossibility evaluates if read/write operations can be performed on the provided DiskItem.
// It returns an RwResult indicating any potential issues that would prevent operations.
func checkRwPossibility(di *acap.DiskItem) *RwResult {

	// Force update to get sure we are have correct states
	if err := acap.UpdateDiskItemEvents(di); err != nil {
		return &RwResult{RwError: RWErrorNotUpdateable, Error: err}
	}

	if !di.Available {
		return &RwResult{RwError: RWErrorNotAvalible, Error: errors.New("Storage/Disk is not avalible")}
	}

	if !di.Writable {
		return &RwResult{RwError: RWErrorNotWriteable, Error: errors.New("Storage/Disk is not writeable")}
	}

	if di.Full {
		return &RwResult{RwError: RWErrorFull, Error: errors.New("Storage/Disk is full")}
	}

	if !di.Setup {
		return &RwResult{RwError: RWErrorNotSetuped, Error: errors.New("Storage/Disk is not setuped")}
	}
	return &RwResult{RwError: RWErrorNone}
}

// WriteFile writes the given content to a file at the specified path on the disk item.
// It returns an RwResult indicating the outcome of the write operation.
func (sp *StorageProvider) WriteFile(di *acap.DiskItem, filePath string, content []byte) *RwResult {
	var rwPossible *RwResult
	if rwPossible = checkRwPossibility(di); rwPossible.RwError == RWErrorNone {
		if err := os.WriteFile(filepath.Join(di.StoragePath, filePath), content, 0644); err != nil {
			return &RwResult{RwError: RWErrorOs, Error: err}
		}
		return &RwResult{RwError: RWErrorNone}
	}
	return rwPossible
}

// RemoveFile deletes the specified file from the disk item.
// It returns an RwResult indicating the outcome of the remove operation.
func (sp *StorageProvider) RemoveFile(di *acap.DiskItem, filePath string) *RwResult {
	var rwPossible *RwResult
	if rwPossible = checkRwPossibility(di); rwPossible.RwError == RWErrorNone {
		if err := os.Remove(filepath.Join(di.StoragePath, filePath)); err != nil {
			return &RwResult{RwError: RWErrorOs, Error: err}
		}
		return &RwResult{RwError: RWErrorNone}
	}
	return rwPossible
}

// ReadFile reads the content of a specified file from the disk item.
// It returns an RwResult containing the read data and any errors that occurred.
func (sp *StorageProvider) ReadFile(di *acap.DiskItem, filePath string) *RwResult {
	var rwPossible *RwResult
	var err error
	var dat []byte
	if rwPossible = checkRwPossibility(di); rwPossible.RwError == RWErrorNone {
		if dat, err = os.ReadFile(filepath.Join(di.StoragePath, filePath)); err != nil {
			return &RwResult{RwError: RWErrorOs, Error: err}
		}
		return &RwResult{RwError: RWErrorNone, Data: dat}
	}
	return rwPossible
}

// Open searches for storage devices,  creates callbacks for event subscriptions,
// and manages them. This method attempts to establish
// communication with all available storages and subscribes to their respective events for monitoring
// changes in their state or attributes.
func (sp *StorageProvider) Open() error {
	var err error
	var storageIds []acap.StorageId

	if storageIds, err = acap.AxStorageList(); err != nil {
		return fmt.Errorf("Unable get storages list: %s", err.Error())
	}

	if len(storageIds) == 0 {
		return errors.New("No storage found")
	}

	for _, storageId := range storageIds {
		subscriptionId, err := acap.AxStorageSubscribe(storageId, storageSubscribeCallback, sp)
		if err != nil {
			sp.app.Syslog.Warnf("Unable to create storage subscription callback: %s for storage: %s", err.Error(), storageId)
		} else {
			sp.app.Syslog.Infof("Successfully create storage subscription for storage: %s, subsciption-id: %d", storageId, subscriptionId)
			sp.DiskItems = append(sp.DiskItems, acap.NewDiskItem(storageId, subscriptionId))
			sp.subscribtions = append(sp.subscribtions, subscriptionId)
		}
	}
	return nil
}

// Unsubscribe Stop subscribing to storage events.
func (sp *StorageProvider) Unsubscribe(d *acap.DiskItem) error {
	return acap.AxStorageUnsubscribe(d.SubscriptionId)
}

// UnsubscribeAll Stop subscribing to all storages events.
func (sp *StorageProvider) UnsubscribeAll() {
	for _, d := range sp.DiskItems {
		if err := sp.Unsubscribe(d); err != nil {
			sp.app.Syslog.Warnf("Failed to unsubscribe event of %s. Error: %s", d.StorageId, err.Error())
		}
	}
}

// Release async release a disk/storage
func (sp *StorageProvider) Release(diskItem *acap.DiskItem) error {
	if diskItem.Setup {
		return diskItem.Storage.AxStorageReleaseAsync(releaseCallback, &storageUserData{storageProvider: sp, diskItem: diskItem})
	}
	return nil
}

// UnsubscribeAll Stop subscribing to all storages events.
func (sp *StorageProvider) ReleaseAll() {
	for _, d := range sp.DiskItems {
		if err := d.Storage.AxStorageReleaseAsync(releaseCallback, &storageUserData{storageProvider: sp, diskItem: d}); err != nil {
			sp.app.Syslog.Warnf("Failed to unsubscribe event of %s. Error: %s", d.StorageId, err.Error())
		}
	}
}

// Close unsubscribes and release all storages/disks
func (sp *StorageProvider) Close() {
	sp.UnsubscribeAll()
	sp.ReleaseAll()
}

// GetDiskItem searches for a DiskItem by its storageId among the managed storage devices.
// It returns the found DiskItem and a boolean indicating whether the search was successful.
func (sp *StorageProvider) GetDiskItem(storageId acap.StorageId) (*acap.DiskItem, bool) {
	for _, d := range sp.DiskItems {
		if d.StorageId == storageId {
			return d, true
		}
	}
	return nil, false
}

// Setup prepares a given DiskItem for use by performing necessary initializations such as
// setting up its directory structure and ensuring it's ready for read/write operations.
// This method performs asynchronous setup and is intended to be called when the disk is
// determined to be in a state suitable for setup (e.g., writable and not full).
func (sp *StorageProvider) Setup(diskItem *acap.DiskItem) error {
	// Writable implies that the disk is available
	if diskItem.Writable && !diskItem.Full && !diskItem.Exiting && !diskItem.Setup {
		return acap.AxStorageSetupAsync(diskItem.StorageId, setupCallback, &storageUserData{storageProvider: sp, diskItem: diskItem})
	}

	if !diskItem.Writable {
		return fmt.Errorf("Storage %s is not writeable/available", diskItem.StorageId)
	}

	if diskItem.Full {
		return fmt.Errorf("Storage %s is full", diskItem.StorageId)
	}

	if diskItem.Exiting {
		return fmt.Errorf("Storage %s is exiting", diskItem.StorageId)
	}
	return nil
}

// setupCallback is a callback function for handling setup completion.
// It updates the disk item's status and logs any errors encountered during setup.
func setupCallback(storage *acap.AXStorage, userdata any, setupErr error) {
	var err error
	sup := userdata.(*storageUserData)

	if setupErr != nil {
		sup.storageProvider.app.Syslog.Warnf("Failed to setup disk: %s. Error: %s", sup.diskItem.StorageId, setupErr.Error())
		return
	}

	if storage.Ptr == nil {
		sup.storageProvider.app.Syslog.Warnf("Failed to setup disk: %s. Error: Storage ptr is NULL", sup.diskItem.StorageId)
		return
	}

	if sup.diskItem.StorageId, err = storage.GetStorageId(); err != nil {
		sup.storageProvider.app.Syslog.Warnf("Failed to get storage_id %s. Error: %", sup.diskItem.StorageId, err.Error())
		return
	}

	if sup.diskItem.StoragePath, err = storage.GetPath(); err != nil {
		sup.storageProvider.app.Syslog.Warnf("Failed to get storage %s path. Error: %s", sup.diskItem.StorageId, err.Error())
		return
	}

	if sup.diskItem.StorageType, err = storage.GetType(); err != nil {
		sup.storageProvider.app.Syslog.Warnf("Failed to get storage %s type. Error: %s", sup.diskItem.StorageId, err.Error())
		return
	}
	sup.diskItem.Storage = storage
	sup.diskItem.Setup = true
	if sup.storageProvider.UseChannelEvents {
		sup.storageProvider.DiskItemsEvents <- sup.diskItem
	}
}

// storageUserData is a helper struct used to pass additional data to callbacks.
type storageUserData struct {
	storageProvider *StorageProvider
	diskItem        *acap.DiskItem
}

// ReleaseOnExiting releases a DiskItem if it is exiting and has been previously set up.
// This is a critical operation to ensure resources are properly released before the disk
// becomes unavailable. It should be invoked as part of the storage management lifecycle,
// especially when handling storage removal or disconnection events.
func (sp *StorageProvider) ReleaseOnExiting(diskItem *acap.DiskItem) {
	if diskItem.Exiting && diskItem.Setup {
		if err := diskItem.Storage.AxStorageReleaseAsync(releaseCallback, &storageUserData{storageProvider: sp, diskItem: diskItem}); err != nil {
			sp.app.Syslog.Warn(err.Error())
		}

	}
}

// releaseCallback is a callback function for handling the release of storage resources.
// It logs the outcome of the release operation, indicating success or detailing any errors.
func releaseCallback(userdata any, err error) {
	sup := userdata.(*storageUserData)
	if err != nil {
		sup.storageProvider.app.Syslog.Warnf("Failed to release %s. Error %s.", sup.diskItem.StorageId, err.Error())
	} else {
		sup.diskItem.Setup = false
		sup.storageProvider.app.Syslog.Infof("Release of %s was successful", sup.diskItem.StorageId)
	}
}

// storageSubscribeCallback is a callback function for handling storage event subscriptions.
// It updates the disk item's status based on the events and manages the lifecycle of the disk item,
// including setup and release as necessary.
func storageSubscribeCallback(storageID acap.StorageId, userdata any, subscribe_err error) {
	var diskItem *acap.DiskItem
	var diskExists bool
	var err error

	sp := userdata.(*StorageProvider)

	if subscribe_err != nil {
		sp.app.Syslog.Error(subscribe_err.Error())
		return
	}

	// Check if disk item exist otherwise we create a new one with populated fields,
	// when it exists we update the event fields with UpdateDiskItemEvents.
	diskItem, diskExists = sp.GetDiskItem(storageID)
	if diskExists {
		if err = acap.UpdateDiskItemEvents(diskItem); err != nil {
			sp.app.Syslog.Warnf("Unable to update disk-item: %s", storageID)
		}
	} else {
		sp.app.Syslog.Warnf("Disk not found in storage provider: %s", storageID)
		return
	}

	sp.ReleaseOnExiting(diskItem)
	if err := sp.Setup(diskItem); err != nil {
		sp.app.Syslog.Warnf("Unable to setup storage %s because: %s", diskItem.StorageId, err.Error())
	}

	if sp.UseChannelEvents {
		sp.DiskItemsEvents <- diskItem
	}
}
