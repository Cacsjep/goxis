package main

import (
	"os"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	axlarod "github.com/Cacsjep/goxis/pkg/axlaord"
)

var (
	app                     *acapapp.AcapApplication
	l                       *axlarod.Larod
	inference_model_file    *os.File
	inference_model         *axlarod.LarodModel
	inference_model_tensors *axlarod.LarodModelIO
	pp_model                *axlarod.LarodModel
	pp_model_tensors        *axlarod.LarodModelIO
	err                     error
)

func loadInferenceModel() error {
	if inference_model_file, err = os.Open("converted_model.tflite"); err != nil {
		return err
	}

	larod_model_fd := inference_model_file.Fd()
	if inference_model, err = l.LoadModelWithDeviceName(&larod_model_fd, "axis-a8-dlpu-tflite", axlarod.LarodAccessPrivate, "Vdo Example App Model", nil); err != nil {
		return err
	}
	if inference_model_tensors, err = inference_model.CreateModelTensors(); err != nil {
		return err
	}
	return nil
}

func loadPreprocessModel() error {
	var ppmap *axlarod.LarodMap
	if ppmap, err = axlarod.NewLarodMapWithEntries([]*axlarod.LarodMapEntries{
		{Key: "image.input.format", Value: "nv12", ValueType: axlarod.LarodMapValueTypeStr},
		{Key: "image.input.size", Value: [2]int64{640, 480}, ValueType: axlarod.LarodMapValueTypeIntArr2},
		{Key: "image.output.format", Value: "rgb-planar", ValueType: axlarod.LarodMapValueTypeStr},
		{Key: "image.output.size", Value: [2]int64{640, 480}, ValueType: axlarod.LarodMapValueTypeIntArr2},
	}); err != nil {
		return err
	}
	if pp_model, err = l.LoadModelWithDeviceName(nil, "axis-a8-gpu-proc", axlarod.LarodAccessPrivate, "", ppmap); err != nil {
		return err
	}
	if pp_model_tensors, err = pp_model.CreateModelTensors(); err != nil {
		return err
	}
	return nil

}

func main() {
	app = acapapp.NewAcapApplication()
	l = axlarod.NewLarod()

	if err := l.Initalize(); err != nil {
		app.Syslog.Crit(err.Error())
	}

	for _, device := range l.Devices {
		app.Syslog.Infof("Larod Device: %s", device.Name)
	}

	if err := loadPreprocessModel(); err != nil {
		app.Syslog.Crit(err.Error())
	}
	app.Syslog.Infof("Preprocess model loaded: %s", pp_model_tensors.String())

	if err := loadInferenceModel(); err != nil {
		app.Syslog.Crit(err.Error())
	}
	app.Syslog.Infof("Inference model loaded: %s", inference_model_tensors.String())

	if err := l.Disconnect(); err != nil {
		app.Syslog.Error(err.Error())
	}
	app.AddCloseCleanFunc(func() { inference_model_file.Close() })
}
