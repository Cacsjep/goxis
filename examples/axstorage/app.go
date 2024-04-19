package main

import (
	"time"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axstorage"
)

// This example demonstrates how to use the storage provider to interact with axstroage.
//
// Orginal C Example: https://github.com/AxisCommunications/acap-native-sdk-examples/tree/main/axstorage
func main() {

	// Initialize a new ACAP application instance.
	// AcapApplication initializes the ACAP application with there name, eventloop, and syslog etc..
	app := acapapp.NewAcapApplication()

	// Storage provider setup
	app.NewStorageProvider(false)
	if err := app.StorageProvider.Open(); err != nil {
		app.Syslog.Crit(err.Error())
	}

	var networkshare *axstorage.DiskItem
	var diskFound bool
	demoFile := "demo.txt"

	go func() int {
		for {
			// Wait for internal subscriptions and storage setups, for more control use storage provider with channel -> app.NewStorageProvider(true)
			time.Sleep(time.Second * 2)

			if networkshare, diskFound = app.StorageProvider.GetDiskItemById("NetworkShare"); !diskFound {
				app.Syslog.Crit("NetworkShare not found")
				continue
			}

			// Writes a file
			if w := app.StorageProvider.WriteFile(networkshare, demoFile, []byte("Here is my content....")); w.RwError != acapapp.RWErrorNone {
				app.Syslog.Errorf("Unable to create file because %s", w.Error)
				continue
			}
			app.Syslog.Infof("Successfully write file: %s", demoFile)

			// Little sleep so you can look into the storage
			time.Sleep(time.Second * 20)

			// Remove a file
			if r := app.StorageProvider.RemoveFile(networkshare, demoFile); r.RwError != acapapp.RWErrorNone {
				app.Syslog.Errorf("Unable to remove file because %s", r.Error)
				continue
			}
			app.Syslog.Infof("Successfully remove file: %s", demoFile)
		}
	}()

	// Run gmain loop with signal handler attached.
	// This will block the main thread until the application is stopped.
	// The application can be stopped by sending a signal to the process (e.g. SIGINT).
	// AxStorage needs a running event loop to handle the callbacks corretly
	app.Run()
}
