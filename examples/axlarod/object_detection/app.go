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

// ! Note this example only works on Artpec-8
// This example demonstrates how detect objects in video frames using the larod package.
// and overlay bounding boxes with axoverlay pkg.
//
// Orginal C Example: https://github.com/AxisCommunications/acap-native-sdk-examples/tree/main/object-detection
func main() {
	if lea, err = Initalize(); err != nil {
		panic(err)
	}

	// For correct singal handling and overlay drawing, the g main loop is required to run in the background.
	lea.app.RunInBackground()

	// Defer the cleanup of the application to ensure all resources are released when the application exits in the below for loop.
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
			if lea.prediction_result, err = lea.InferenceOutputRead(lea.infer_result.OutputData.(*CocoResult)); err != nil {
				lea.app.Syslog.Errorf("Failed to convert prediction result: %s", err.Error())
				return
			}

			// Draw overlay
			if err = lea.overlayProvider.Redraw(); err != nil {
				lea.app.Syslog.Errorf("Failed to redraw overlay: %s", err.Error())
			}

			lea.app.Syslog.Infof("Frame: %d, PP time: %.fms, Infer. exec time: %.fms, Detections: %d",
				frame.SequenceNbr,
				lea.pp_result.ExecutionTime,
				lea.infer_result.ExecutionTime,
				len(lea.detections),
			)

		}
	}
}

// larodExampleApplication struct defines the structure for this example.
// It includes configuration for application, models, video stream, and other operational parameters.
type larodExampleApplication struct {
	app               *acapapp.AcapApplication       // app represents the acap application
	PPModel           *axlarod.LarodModel            // PPModel is the preprocessing model.
	DetectionModel    *axlarod.LarodModel            // DetectionModel is the model used for detecting objects in video frames.
	streamWidth       int                            // streamWidth specifies the width of the video stream.
	streamHeight      int                            // streamHeight specifies the height of the video stream.
	cocoInputWidth    int                            // cocoInputWidth specifies the width of the input tensor for the detection model.
	cocoInputHeight   int                            // cocoInptHeight specifies the height of the input tensor for the detection model.
	fps               int                            // fps represents the frame rate of the video stream.
	sconfig           *axvdo.VideoSteamConfiguration // sconfig holds the configuration for the video stream.
	pp_result         *axlarod.JobResult             // pp_result holds the result of the preprocessing model job.
	infer_result      *axlarod.JobResult             // infer_result holds the result of the detection model job.
	prediction_result *PredictionResult              // prediction_result stores the output of the inference process.
	threshold         float32                        // threshold is the minimum score required for an object to be considered detected.
	overlayProvider   *acapapp.OverlayProvider       // overlayProvider is used to draw overlay on the video stream.
	detections        []Detection                    // detections stores the detected objects.
}

// Initialize prepares and initializes all necessary components for the application.
// It sets up models, video streaming and processing configurations.
// Returns a configured instance of larodExampleApplication or an error if initialization fails.
func Initalize() (*larodExampleApplication, error) {

	lea := &larodExampleApplication{fps: 12, threshold: 0.4, cocoInputWidth: 300, cocoInputHeight: 300, detections: []Detection{}}

	// Initialize a new ACAP application instance.
	// AcapApplication initializes the ACAP application with there name, eventloop, and syslog etc..
	lea.app = acapapp.NewAcapApplication()

	// Determine the stream resolution
	if err := lea.SetupStreamResolution(); err != nil {
		return nil, err
	}

	// Initialize/Connecting Larod
	if err = lea.app.InitalizeLarod(); err != nil {
		return nil, err
	}

	// Initialize the preprocessing model
	if err = lea.InitalizePPModel(axlarod.PreProccessOutputFormatRgbInterleaved); err != nil {
		return nil, err
	}

	// Print the available devices
	for _, d := range lea.app.Larod.Devices {
		lea.app.Syslog.Infof("Device: %s", d.Name)
	}

	// Initialize the detection model
	if err = lea.InitalizeDetectionModel("ssd_mobilenet_v2_coco_quant_postprocess.tflite", "axis-a8-dlpu-tflite"); err != nil {
		return nil, err
	}

	// Initialize and start the video stream
	if err = lea.InitalizeAndStartVdo(); err != nil {
		return nil, err
	}

	// Initialize the overlay provider
	if err = lea.InitOverlay(); err != nil {
		return nil, err
	}
	return lea, nil
}
