package main

import (
	"fmt"

	"github.com/Cacsjep/goxis"
	"github.com/Cacsjep/goxis/pkg/acap"
)

var (
	err          error
	subscription int
	app          *goxis.AcapApplication
)

func main() {
	if app, err = goxis.NewAcapApplication(); err != nil {
		panic(err)
	}
	defer app.Close()

	storages, err := acap.AxStorageList()
	fmt.Println(storages, err)
	for _, sid := range storages {
		ssss, err := acap.AXStorageGetStatus(sid, acap.AXStorageAvailableEvent)
		fmt.Println(sid, ssss, err)
	}

	subid, err := acap.AxStorageSubscribe(storages[0], func(storageID acap.StorageId, userdata any, err error, cleanup func()) {
		if err != nil {
			app.Syslog.Error(err.Error())
			panic(err)
		}
		diskItem, err := acap.NewDiskItemFromSubscription(storageID)
		fmt.Println(diskItem)
		cleanup()
	}, "my user data")

	fmt.Println(subid, err)

	err = acap.AxStorageSetupAsync(storages[0], func(storage *acap.AXStorage, userdata any, err error, cleanup func()) {
		fmt.Println(err, "err")
		if err != nil {
			app.Syslog.Error(err.Error())
			panic(err)
		}
		fmt.Println("AxStorageSetupAsync")
		fmt.Println(storage.GetStorageId())
		fmt.Println(storage.GetPath())
		fmt.Println(storage.GetType())
		cleanup()
	}, "my user data")
	fmt.Println(err)
	app.Run()
}
