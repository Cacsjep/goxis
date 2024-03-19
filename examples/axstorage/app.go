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
	app.Run()
}
