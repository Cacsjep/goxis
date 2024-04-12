package axlarod

/*
#cgo pkg-config: liblarod
#include "larod.h"
*/
import "C"

type LarodChip int

const (
	LarodChipInvalid           LarodChip = 0
	LarodChipDebug             LarodChip = 1
	LarodChipTFLiteCPU         LarodChip = 2
	LarodChipTPU               LarodChip = 4
	LarodChipCVFlowNN          LarodChip = 6
	LarodChipTFLiteGLGPU       LarodChip = 8
	LarodChipCVFlowProc        LarodChip = 9
	LarodChipACE               LarodChip = 10
	LarodChipLibYUV            LarodChip = 11
	LarodChipTFLiteArtpec8Dlpu LarodChip = 12
	LarodChipOpenCL            LarodChip = 13
)
