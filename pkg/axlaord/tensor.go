package axlarod

/*
#cgo pkg-config: liblarod
#include "larod.h"
*/
import "C"

type LarodTensorDataType int

const (
	LarodTensorDataTypeInvalid LarodTensorDataType = iota
	LarodTensorDataTypeUnspecified
	LarodTensorDataTypeBool
	LarodTensorDataTypeUint8
	LarodTensorDataTypeInt8
	LarodTensorDataTypeUint16
	LarodTensorDataTypeInt16
	LarodTensorDataTypeUint32
	LarodTensorDataTypeInt32
	LarodTensorDataTypeUint64
	LarodTensorDataTypeInt64
	LarodTensorDataTypeFloat16
	LarodTensorDataTypeFloat32
	LarodTensorDataTypeFloat64
	LarodTensorDataTypeMax
)

type LarodTensorLayout int

const (
	LarodTensorLayoutInvalid LarodTensorLayout = iota
	LarodTensorLayoutUnspecified
	LarodTensorLayoutNHWC
	LarodTensorLayoutNCHW
	LarodTensorLayout420SP
	LarodTensorLayoutMax
)
