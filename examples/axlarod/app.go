package main

import (
	"fmt"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	axlarod "github.com/Cacsjep/goxis/pkg/axlaord"
)

func main() {
	app := acapapp.NewAcapApplication()

	l := axlarod.NewLarod()
	if err := l.Initalize(); err != nil {
		app.Syslog.Crit(err.Error())
	}
	app.Syslog.Info("Larod Connected")

	for _, device := range l.Devices {
		app.Syslog.Infof("Larod Device: %s", device.Name)
	}

	if err := l.Disconnect(); err != nil {
		app.Syslog.Error(err.Error())
	}
	app.Syslog.Info("Larod Disconnected")

	ppmap, err := axlarod.NewLarodMapWithEntries([]*axlarod.LarodMapEntries{
		{Key: "image.input.format", Value: "nv12", ValueType: axlarod.LarodMapValueTypeStr},
		{Key: "image.input.size", Value: [2]int64{640, 480}, ValueType: axlarod.LarodMapValueTypeIntArr2},
		{Key: "image.output.format", Value: "rgb-interleaved", ValueType: axlarod.LarodMapValueTypeStr},
		{Key: "image.output.format", Value: "rgb-planar", ValueType: axlarod.LarodMapValueTypeStr},
		{Key: "image.output.size", Value: [2]int64{640, 480}, ValueType: axlarod.LarodMapValueTypeIntArr2},
	})
	if err != nil {
		app.Syslog.Crit(err.Error())
	}

	fmt.Println(ppmap)

}
