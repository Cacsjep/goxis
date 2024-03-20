package main

import (
	"time"

	"github.com/Cacsjep/goxis"
	"github.com/Cacsjep/goxis/pkg/acap"
)

// https://github.com/AxisCommunications/acap-native-sdk-examples/blob/main/axstorage/app/axstorage.c#L200
func main() {
	var err error
	var app *goxis.AcapApplication
	var networkshare *acap.DiskItem
	var diskFound bool

	if app, err = goxis.NewAcapApplication(); err != nil {
		panic(err)
	}
	defer app.Close()

	sp := app.NewStorageProvider(false)

	if err := sp.Open(); err != nil {
		app.Syslog.Crit(err.Error())
	}
	defer sp.Close()

	demoFile := "demo.txt"

	go func() int {
		for {
			// Wait for internal subscriptions and storage setups, for more control use storage provider with channel -> app.NewStorageProvider(true)
			time.Sleep(time.Second * 2)

			if networkshare, diskFound = sp.GetDiskItemById("NetworkShare"); !diskFound {
				app.Syslog.Crit("NetworkShare not found")
			}

			// Writes a file
			if w := sp.WriteFile(networkshare, demoFile, []byte("Here is my content....")); w.RwError != goxis.RWErrorNone {
				app.Syslog.Errorf("Unable to create file because %s", w.Error)
				continue
			}
			app.Syslog.Infof("Successfully write file: %s", demoFile)

			// Little sleep so you can look into the storage
			time.Sleep(time.Second * 20)

			// Remove a file
			if r := sp.RemoveFile(networkshare, demoFile); r.RwError != goxis.RWErrorNone {
				app.Syslog.Errorf("Unable to remove file because %s", r.Error)
				continue
			}
			app.Syslog.Infof("Successfully remove file: %s", demoFile)
		}
	}()

	app.Run()
}
