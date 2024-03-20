package main

import (
	"time"

	"github.com/Cacsjep/goxis"
)

// https://github.com/AxisCommunications/acap-native-sdk-examples/blob/main/axstorage/app/axstorage.c#L200
func main() {
	var err error
	var app *goxis.AcapApplication

	if app, err = goxis.NewAcapApplication(); err != nil {
		panic(err)
	}
	defer app.Close()

	sp := app.NewStorageProvider(true)

	if err := sp.Open(); err != nil {
		app.Syslog.Crit(err.Error())
	}
	defer sp.Close()

	demoFile := "demo.txt"

	go func() int {
		for {
			select {
			case diskItem := <-sp.DiskItemsEvents:
				app.Syslog.Info(diskItem.String())

				// When disk is setuped we try to write and remove the demo file
				if diskItem.Setup {

					// Writes a file
					if w := sp.WriteFile(diskItem, demoFile, []byte("Here is my content....")); w.RwError != goxis.RWErrorNone {
						app.Syslog.Errorf("Unable to create file because %s", w.Error)
						continue
					}

					app.Syslog.Infof("Successfully write file: %s", demoFile)
					time.Sleep(time.Second * 20)

					// Remove a file
					if r := sp.RemoveFile(diskItem, demoFile); r.RwError != goxis.RWErrorNone {
						app.Syslog.Errorf("Unable to remove file because %s", r.Error)
						continue
					}
					app.Syslog.Infof("Successfully remove file: %s", demoFile)
				}
			}
		}
	}()

	app.Run()
}
