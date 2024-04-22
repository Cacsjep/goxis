package main

import (
	"math"
	"sort"

	"github.com/Cacsjep/goxis/pkg/axlarod"
)

const (
	BATCH_SIZE              = 1
	NUM_PREDICTIONS         = 25200
	ELEMENTS_PER_PREDICTION = 5 // x, y, w, h, confidence
	CLASSES                 = 80
	UINT_8_SIZE             = 1
	TENSOR_OUT_1_SIZE       = uint(BATCH_SIZE * NUM_PREDICTIONS * (ELEMENTS_PER_PREDICTION + CLASSES) * UINT_8_SIZE)
	BBOX_SIZE               = ELEMENTS_PER_PREDICTION + CLASSES
	V5N_QUANTIZATION        = 0.004190513398498297 // yolov5n.tflite
)

type YoloResult struct {
	Detections []Detection
}

type Detection struct {
	Box          axlarod.BoundingBox
	Confidence   float32
	BestClassIdx int
	BestScore    float32
}

// InitializeDetectionModel configures a detection model with the given model file and hardware chip.
// It sets up memory-mapped file configurations for input and output tensors.
// Returns an error if model initialization fails.
func (lea *larodExampleApplication) InitalizeDetectionModel(modelFilePath string, chipString string) error {
	model_defs := axlarod.MemMapConfiguration{
		InputTmpMapFiles: map[int]*axlarod.MemMapFile{
			0: lea.app.FrameProvider.PostProcessModel.Outputs[0].MemMapFile, // Using of ppmodel output as input for detection model
		},
		OutputTmpMapFiles: map[int]*axlarod.MemMapFile{
			0: {Size: TENSOR_OUT_1_SIZE}, // Output Tensor 1
		},
	}

	if lea.DetectionModel, err = lea.app.Larod.NewInferModel(modelFilePath, chipString, model_defs, nil); err != nil {
		return err
	}
	lea.app.AddModelCleaner(lea.DetectionModel)
	return nil
}

// getDResult retrieves the detection results from the model's output tensors.
func (lea *larodExampleApplication) getDResult() ([]Detection, error) {
	t1, err := lea.DetectionModel.Outputs[0].GetData(int(TENSOR_OUT_1_SIZE))
	if err != nil {
		return nil, err
	}

	detections := make([]Detection, 0, NUM_PREDICTIONS)
	for i := 0; i < NUM_PREDICTIONS; i++ {
		offset := i * BBOX_SIZE
		centerX := yolov5n_dequantize(t1[offset])
		centerY := yolov5n_dequantize(t1[offset+1])
		width := yolov5n_dequantize(t1[offset+2])
		height := yolov5n_dequantize(t1[offset+3])
		confidence := yolov5n_dequantize(t1[offset+4])

		if confidence > lea.threshold {
			top := centerY - height/2
			left := centerX - width/2
			bottom := centerY + height/2
			right := centerX + width/2

			det := Detection{
				Box: axlarod.BoundingBox{
					Top:    top,
					Left:   left,
					Bottom: bottom,
					Right:  right,
				},
				Confidence: confidence,
			}

			maxClassIdx := 0
			maxScore := yolov5n_dequantize(t1[offset+5])
			for j := 1; j < CLASSES; j++ {
				score := yolov5n_dequantize(t1[offset+5+j])
				if score > maxScore {
					maxScore = score
					maxClassIdx = j
				}
			}
			det.BestClassIdx = maxClassIdx
			det.BestScore = maxScore
			detections = append(detections, det)
		}
	}

	lea.detections = nonMaximumSuppression(detections)
	return lea.detections, nil
}

func nonMaximumSuppression(detections []Detection) []Detection {
	sort.Slice(detections, func(i, j int) bool {
		return detections[i].Confidence > detections[j].Confidence
	})

	var keep []Detection

	for i := 0; i < len(detections); i++ {
		suppressed := false
		for j := 0; j < len(keep); j++ {
			if computeIoU(detections[i].Box, keep[j].Box) > lea.iouThreshold {
				suppressed = true
				break
			}
		}
		if !suppressed {
			keep = append(keep, detections[i])
		}
	}
	return keep
}

// computeIoU calculates the Intersection over Union (IoU) of two bounding boxes.
func computeIoU(box1, box2 axlarod.BoundingBox) float64 {
	x1 := math.Max(float64(box1.Left), float64(box2.Left))
	y1 := math.Max(float64(box1.Top), float64(box2.Top))
	x2 := math.Min(float64(box1.Right), float64(box2.Right))
	y2 := math.Min(float64(box1.Bottom), float64(box2.Bottom))

	interArea := math.Max(0, x2-x1) * math.Max(0, y2-y1)
	if interArea == 0 {
		return 0
	}

	box1Area := (box1.Right - box1.Left) * (box1.Bottom - box1.Top)
	box2Area := (box2.Right - box2.Left) * (box2.Bottom - box2.Top)
	return interArea / (float64(box1Area+box2Area) - interArea)
}

// yolov5n_dequantize converts a quantized value to a float32 value.
func yolov5n_dequantize(q uint8) float32 {
	return V5N_QUANTIZATION * float32(q)
}

// Inference executes the model and retrieves the processed results.
// It ensures the model's file pointers are correctly positioned before execution.
// Returns a JobResult containing the inference results or an error if the inference process fails.
func (lea *larodExampleApplication) Inference() (*axlarod.JobResult, error) {

	// Rewind all output files position before each job.
	if err = lea.DetectionModel.RewindAllOutputsMemMapFiles(); err != nil {
		return nil, err
	}

	var result *axlarod.JobResult
	if result, err = lea.app.Larod.ExecuteJob(lea.DetectionModel, func() error {
		return nil // is feeded via memmap
	}, func() (any, error) {
		return lea.getDResult()
	}); err != nil {
		return nil, err
	}

	return result, nil
}
