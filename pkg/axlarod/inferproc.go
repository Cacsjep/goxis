package axlarod

import (
	"github.com/Cacsjep/goxis/pkg/utils"
)

// NewInferModel creates a new model for inference, just a tiny wrapper around LoadModelWithDeviceName, CreateModelTensors and CreateJobRequest
func (l *Larod) NewInferModel(filename, device string, model_defs MemMapConfiguration, job_params *LarodMap) (*LarodModel, error) {
	var err error
	var infer_model *LarodModel

	if infer_model, err = l.LoadModelWithDeviceName(utils.StrPtr(filename), device, LarodAccessPrivate, "infer", nil); err != nil {
		return nil, err
	}

	if err = infer_model.CreateModelTensors(&model_defs); err != nil {
		return nil, err
	}

	_, err = infer_model.CreateJobRequest(infer_model.Inputs, infer_model.Outputs, job_params)
	if err != nil {
		return nil, err
	}

	return infer_model, nil
}
