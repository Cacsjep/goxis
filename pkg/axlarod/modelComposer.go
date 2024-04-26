package axlarod

import (
	"math"
	"sort"
)

type ModelComposer struct {
	Labels              []string
	DequantizeFunc      func(int8) float32
	OutputParser        func(rawModelOuput []float32, mc *ModelComposer) []Detection
	larodModel          *LarodModel
	Threshold           float32
	larod               *Larod
	IouThreshold        float64
	OutputTensorPitches *LarodTensorPitches
}

type Detection struct {
	Box        BoundingBox
	Confidence float32
	ClassIndex int
	ClassLabel string
}

func InizalizeModelComposer(larod *Larod, modelFilePath string, chipString string, modelInput *MemMapFile, modelComposer *ModelComposer) error {
	var err error
	model_defs := MemMapConfiguration{
		InputTmpMapFiles: map[int]*MemMapFile{
			0: modelInput, // Using of ppmodel output as input for detection model
		},
		OutputTmpMapFiles: map[int]*MemMapFile{
			0: {UsePitch0Size: true},
		},
	}
	if modelComposer.larodModel, err = larod.NewInferModel(modelFilePath, chipString, model_defs, nil); err != nil {
		return err
	}
	modelComposer.larod = larod
	pitches, err := modelComposer.larodModel.Outputs[0].GetTensorPitches()
	if err != nil {
		return err
	}
	modelComposer.OutputTensorPitches = pitches
	return nil
}

func (mc *ModelComposer) getDResult() ([]Detection, error) {
	quant_output, err := mc.larodModel.Outputs[0].GetData(int(mc.OutputTensorPitches.Pitches[0]))
	if err != nil {
		return nil, err
	}
	output := make([]float32, len(quant_output))
	for i, byteVal := range quant_output {
		output[i] = mc.DequantizeFunc(int8(byteVal))
	}
	return mc.nonMaximumSuppression(mc.OutputParser(output, mc)), nil
}

func (mc *ModelComposer) Inference() (*JobResult, error) {
	var err error
	if err = mc.larodModel.RewindAllOutputsMemMapFiles(); err != nil {
		return nil, err
	}
	var result *JobResult
	if result, err = mc.larod.ExecuteJob(mc.larodModel, func() error {
		return nil
	}, func() (any, error) {
		return mc.getDResult()
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (mc *ModelComposer) Clean() error {
	return mc.larod.DestroyModel(mc.larodModel)
}

func (mc *ModelComposer) computeIoU(box1, box2 BoundingBox) float64 {
	x1, y1 := math.Max(float64(box1.Left), float64(box2.Left)), math.Max(float64(box1.Top), float64(box2.Top))
	x2, y2 := math.Min(float64(box1.Right), float64(box2.Right)), math.Min(float64(box1.Bottom), float64(box2.Bottom))
	interArea := math.Max(0, x2-x1) * math.Max(0, y2-y1)
	if interArea == 0 {
		return 0
	}
	box1Area, box2Area := (box1.Right-box1.Left)*(box1.Bottom-box1.Top), (box2.Right-box2.Left)*(box2.Bottom-box2.Top)
	return interArea / (float64(box1Area+box2Area) - interArea)
}

func (mc *ModelComposer) nonMaximumSuppression(detections []Detection) []Detection {
	sort.Slice(detections, func(i, j int) bool {
		return detections[i].Confidence > detections[j].Confidence
	})

	keep := make([]Detection, 0, len(detections))
	suppressed := make([]bool, len(detections))

	for i := 0; i < len(detections); i++ {
		if suppressed[i] {
			continue
		}
		for j := i + 1; j < len(detections); j++ {
			if !suppressed[j] && mc.computeIoU(detections[i].Box, detections[j].Box) > mc.IouThreshold {
				suppressed[j] = true
			}
		}
		keep = append(keep, detections[i])
	}
	return keep
}
