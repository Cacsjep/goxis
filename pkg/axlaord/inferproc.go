package axlarod

import (
	"github.com/Cacsjep/goxis/pkg/utils"
)

type TmpFileSize int

func (l *Larod) NewInferModel(filename, device string, model_defs ModelTmpMapDefiniton) (*LarodModel, error) {
	var err error
	var infer_model *LarodModel
	var pp_model_io *LarodModelIO

	if infer_model, err = l.LoadModelWithDeviceName(utils.StrPtr(filename), device, LarodAccessPrivate, "", nil); err != nil {
		return nil, err
	}

	if pp_model_io, err = infer_model.CreateModelTensors(&model_defs); err != nil {
		return nil, err
	}

	_, err = infer_model.CreateJobRequest(pp_model_io.Inputs, pp_model_io.Outputs, nil)
	if err != nil {
		return nil, err
	}

	return infer_model, nil
}
