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
// This example demonstrates how detect objects in video frames using the larod package and yolov5.
// This example use frameprovider with post processing mode!
func main() {
	if lea, err = Initalize(); err != nil {
		panic(err)
	}

	// For correct singal handling and overlay drawing, the g main loop is required to run in the background.
	lea.app.RunInBackground()

	for {
		select {
		case frame := <-lea.app.FrameProvider.FrameStreamChannel:
			// Frameprovide works in post processing mode.
			// the output of the frame provider and larod preprocessing model is
			// automatically set in memmap file via lea.app.FrameProvider.PostProcessModel.Outputs[0].MemMapFile in d_model.go
			if frame.Error != nil {
				lea.app.Syslog.Errorf("Unexpected Vdo Error: %s", frame.Error.Error())
				continue
			}

			// Execute the detection model job
			if lea.infer_result, err = lea.Inference(); err != nil {
				lea.app.Syslog.Errorf("Failed to execute Detection Model: %s", err.Error())
				continue
			}

			// Draw overlay
			if err = lea.overlayProvider.Redraw(); err != nil {
				lea.app.Syslog.Errorf("Failed to redraw overlay: %s", err.Error())
			}

			lea.app.Syslog.Infof("Frame: %d, Infer. exec time: %.fms, Detections: %d",
				frame.SequenceNbr,
				lea.infer_result.ExecutionTime,
				len(lea.detections),
			)

		}
	}
}

// larodExampleApplication struct defines the structure for this example.
// It includes configuration for application, models, video stream, and other operational parameters.
type larodExampleApplication struct {
	app             *acapapp.AcapApplication       // app represents the acap application
	DetectionModel  *axlarod.LarodModel            // DetectionModel is the model used for detecting objects in video frames.
	streamWidth     int                            // streamWidth specifies the width of the video stream.
	streamHeight    int                            // streamHeight specifies the height of the video stream.
	yoloInputWidth  int                            // yoloInputWidth specifies the width of the input tensor for the detection model.
	yoloInputHeight int                            // yoloInptHeight specifies the height of the input tensor for the detection model.
	fps             int                            // fps represents the frame rate of the video stream.
	sconfig         *axvdo.VideoSteamConfiguration // sconfig holds the configuration for the video stream.
	infer_result    *axlarod.JobResult             // infer_result holds the result of the detection model job.
	threshold       float32                        // threshold is the minimum score required for an object to be considered detected.
	overlayProvider *acapapp.OverlayProvider       // overlayProvider is used to draw overlay on the video stream.
	detections      []Detection                    // detections stores the detected objects.
	iouThreshold    float64                        // iouThreshold is the threshold for Intersection over Union (IoU) for non-maximum suppression.
}

// Initialize prepares and initializes all necessary components for the application.
// It sets up models, video streaming and processing configurations.
// Returns a configured instance of larodExampleApplication or an error if initialization fails.
func Initalize() (*larodExampleApplication, error) {

	lea := &larodExampleApplication{fps: 3, threshold: 0.6, yoloInputWidth: 640, yoloInputHeight: 640, detections: []Detection{}, iouThreshold: 0.5}

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

	// Initialize and start the video stream
	if err = lea.InitalizeAndStartVdo(); err != nil {
		return nil, err
	}

	// Initialize the preprocessing model
	if err = lea.app.FrameProvider.SetLarodPostProccessor("cpu-proc", axlarod.PreProccessOutputFormatRgbInterleaved, &axvdo.VdoResolution{Width: lea.yoloInputWidth, Height: lea.yoloInputHeight}); err != nil {
		return nil, err
	}

	// Initialize the detection model
	if err = lea.InitalizeDetectionModel("yolov5n.tflite", "axis-a8-dlpu-tflite"); err != nil {
		return nil, err
	}

	// Initialize the overlay provider
	if err = lea.InitOverlay(); err != nil {
		return nil, err
	}

	if err = lea.app.FrameProvider.Start(); err != nil {
		return nil, err
	}
	return lea, nil
}
