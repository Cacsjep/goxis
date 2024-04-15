package axlarod

import "fmt"

type PreProccessOutputFormat string

var (
	PreProccessOutputFormatRgbInterleaved PreProccessOutputFormat = "rgb-interleaved"
	PreProccessOutputFormatRgbPlanar      PreProccessOutputFormat = "rgb-planar"
)

type LarodResolution struct {
	Width  int
	Height int
}

func (lr *LarodResolution) ToArray() [2]int64 {
	return [2]int64{int64(lr.Width), int64(lr.Height)}
}

func (l *Larod) NewPreProccessModel(device string, inputSize LarodResolution, outputSize LarodResolution, outputFormat PreProccessOutputFormat) (*LarodModel, error) {
	var err error
	var ppmap *LarodMap
	var pp_model *LarodModel
	var pp_model_io *LarodModelIO

	if ppmap, err = NewLarodMapWithEntries([]*LarodMapEntries{
		{Key: "image.input.format", Value: "nv12", ValueType: LarodMapValueTypeStr},
		{Key: "image.input.size", Value: inputSize.ToArray(), ValueType: LarodMapValueTypeIntArr2},
		{Key: "image.output.format", Value: string(outputFormat), ValueType: LarodMapValueTypeStr},
		{Key: "image.output.size", Value: outputSize.ToArray(), ValueType: LarodMapValueTypeIntArr2},
	}); err != nil {
		return nil, err
	}

	if pp_model, err = l.LoadModelWithDeviceName(nil, device, LarodAccessPrivate, "", ppmap); err != nil {
		return nil, err
	}

	model_defs := ModelTmpMapDefiniton{
		InputTmpMapFiles: map[int]*TmpFile{
			0: {UsePitch0Size: true}, // Input Tensor 0
		},
		OutputTmpMapFiles: map[int]*TmpFile{
			0: {Size: 480 * 270 * 3}, // Output Tensor 0
		},
	}

	if pp_model_io, err = pp_model.CreateModelTensors(&model_defs); err != nil {
		return nil, err
	}

	rgb_buffer_size := pp_model_io.OutputPitches.Pitches[0]
	expectedSize := uint(outputSize.Width * outputSize.Height * 3)

	if rgb_buffer_size != expectedSize {
		return nil, fmt.Errorf("Expected size of RGB buffer is %d, but got %d", expectedSize, rgb_buffer_size)
	}

	_, err = pp_model.CreateJobRequest(pp_model_io.Inputs, pp_model_io.Outputs, nil)
	if err != nil {
		return nil, err
	}

	return pp_model, nil
}
