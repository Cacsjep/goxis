package main

import (
	"fmt"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	axlarod "github.com/Cacsjep/goxis/pkg/axlaord"
	"github.com/Cacsjep/goxis/pkg/axvdo"
)

var (
	err error
	lea *larodExampleApplication
)

type larodExampleApplication struct {
	app              *acapapp.AcapApplication
	PPModel          *axlarod.LarodModel
	DetectionModel   *axlarod.LarodModel
	streamWidth      int
	streamHeight     int
	fps              int
	sconfig          *axvdo.VideoSteamConfiguration
	fp               *acapapp.FrameProvider
	pp_result        *axlarod.JobResult
	infer_result     *axlarod.JobResult
	detection_result *DetectionResult
}

// Initalize all we need
func Initalize() (*larodExampleApplication, error) {
	lea := &larodExampleApplication{streamWidth: 480, streamHeight: 270}
	lea.app = acapapp.NewAcapApplication()

	if err := lea.app.InitalizeLarod(); err != nil {
		return nil, err
	}

	if err = lea.InitalizePPModel(); err != nil {
		return nil, err
	}

	for _, d := range lea.app.Larod.Devices {
		fmt.Println(d.Name)
	}

	if err = lea.InitalizeDetectionModel("converted_model.tflite", "idon"); err != nil {
		return nil, err
	}

	if err = lea.InitalizeAndStartVdo(); err != nil {
		return nil, err
	}
	return lea, nil
}

func main() {
	if lea, err = Initalize(); err != nil {
		panic(err)
	}

	defer lea.app.Close()

	for {
		select {
		case frame := <-lea.fp.FrameStreamChannel:
			if frame.Error != nil {
				lea.app.Syslog.Errorf("Unexpected Vdo Error: %s", frame.Error.Error())
				continue
			}

			// Execute the prepossessing model job
			if lea.pp_result, err = lea.PreProcess(frame); err != nil {
				lea.app.Syslog.Errorf("Failed to execute PPModel: %s", err.Error())
				continue
			}

			if lea.infer_result, err = lea.Inference(lea.pp_result.OutputData); err != nil {
				lea.app.Syslog.Errorf("Failed to execute Detection Model: %s", err.Error())
				continue
			}

			lea.detection_result, err = lea.DetectionResultConverter(lea.infer_result.OutputData)

			lea.app.Syslog.Infof("Frame: %s, PP exec time: %d, Inference exec time: %d, Persons: %f, Car: %f",
				frame.String(),
				lea.pp_result.ExecutionTime,
				lea.infer_result.ExecutionTime,
				lea.detection_result.Persons,
				lea.detection_result.Car,
			)

			return
		}
	}
}
