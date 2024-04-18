package main

import (
	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axlarod"
	"github.com/Cacsjep/goxis/pkg/axvdo"
)

var (
	err error                    // err commonly holds errors encountered during the runtime.
	lea *larodExampleApplication // lea is an instance of the application handling video processing and model inference.
)

// larodExampleApplication struct defines the structure for this example.
// It includes configuration for application, models, video stream, and other operational parameters.
type larodExampleApplication struct {
	app               *acapapp.AcapApplication       // app represents the acap application
	PPModel           *axlarod.LarodModel            // PPModel is the preprocessing model.
	DetectionModel    *axlarod.LarodModel            // DetectionModel is the model used for detecting objects in video frames.
	streamWidth       int                            // streamWidth specifies the width of the video stream.
	streamHeight      int                            // streamHeight specifies the height of the video stream.
	fps               int                            // fps represents the frame rate of the video stream.
	sconfig           *axvdo.VideoSteamConfiguration // sconfig holds the configuration for the video stream.
	pp_result         *axlarod.JobResult             // pp_result holds the result of the preprocessing model job.
	infer_result      *axlarod.JobResult             // infer_result holds the result of the detection model job.
	prediction_result *PredictionResult              // prediction_result stores the output of the inference process.
}

// Initialize prepares and initializes all necessary components for the application.
// It sets up models, video streaming and processing configurations.
// Returns a configured instance of larodExampleApplication or an error if initialization fails.
func Initalize() (*larodExampleApplication, error) {
	lea := &larodExampleApplication{streamWidth: 480, streamHeight: 270, fps: 5}
	lea.app = acapapp.NewAcapApplication()

	if err := lea.app.InitalizeLarod(); err != nil {
		return nil, err
	}

	if err = lea.InitalizePPModel(axlarod.PreProccessOutputFormatRgbInterleaved); err != nil {
		return nil, err
	}

	for _, d := range lea.app.Larod.Devices {
		lea.app.Syslog.Infof("Device: %s", d.Name)
	}

	if err = lea.InitalizeDetectionModel("converted_model.tflite", "axis-a8-dlpu-tflite"); err != nil {
		return nil, err
	}

	if err = lea.InitalizeAndStartVdo(); err != nil {
		return nil, err
	}
	return lea, nil
}

// main serves as the entry point of the application. It initializes the application and handles the main loop.
// In case of failure during initialization, the application will panic.
func main() {
	if lea, err = Initalize(); err != nil {
		panic(err)
	}

	defer lea.app.Close()

	for {
		select {
		case frame := <-lea.app.FrameProvider.FrameStreamChannel:
			if frame.Error != nil {
				lea.app.Syslog.Errorf("Unexpected Vdo Error: %s", frame.Error.Error())
				continue
			}

			// Execute the prepossessing model job
			if lea.pp_result, err = lea.PreProcess(frame); err != nil {
				lea.app.Syslog.Errorf("Failed to execute PPModel: %s", err.Error())
				return
			}

			// Execute the detection model job
			if lea.infer_result, err = lea.Inference(); err != nil {
				lea.app.Syslog.Errorf("Failed to execute Detection Model: %s", err.Error())
				return
			}

			// Retrieve the prediction result
			if lea.prediction_result, err = lea.InferenceOutputRead(lea.infer_result.OutputData.(*InferenceResult)); err != nil {
				lea.app.Syslog.Errorf("Failed to convert prediction result: %s", err.Error())
				return
			}

			lea.app.Syslog.Infof("Frame: %d, PP time: %.fms, Infer time: %.fms, Persons: %.1f%%, Car: %.1f%%",
				frame.SequenceNbr,
				lea.pp_result.ExecutionTime,
				lea.infer_result.ExecutionTime,
				lea.prediction_result.Persons,
				lea.prediction_result.Car,
			)
		}
	}
}
