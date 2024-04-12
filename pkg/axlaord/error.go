package axlarod

/*
#cgo pkg-config: liblarod
#include "larod.h"
*/
import "C"

type LarodErrorCode int

const (
	LarodErrorNone            LarodErrorCode = 0
	LarodErrorJob             LarodErrorCode = -1
	LarodErrorLoadModel       LarodErrorCode = -2
	LarodErrorFD              LarodErrorCode = -3
	LarodErrorModelNotFound   LarodErrorCode = -4
	LarodErrorPermission      LarodErrorCode = -5
	LarodErrorConnection      LarodErrorCode = -6
	LarodErrorCreateSession   LarodErrorCode = -7
	LarodErrorKillSession     LarodErrorCode = -8
	LarodErrorInvalidChipID   LarodErrorCode = -9
	LarodErrorInvalidAccess   LarodErrorCode = -10
	LarodErrorDeleteModel     LarodErrorCode = -11
	LarodErrorTensorMismatch  LarodErrorCode = -12
	LarodErrorVersionMismatch LarodErrorCode = -13
	LarodErrorAlloc           LarodErrorCode = -14
	LarodErrorMaxErrno        LarodErrorCode = 1024
)

type LarodError struct {
	Code LarodErrorCode
	Msg  string
}

func (l *LarodError) Error() string {
	return l.Msg
}

func newLarodError(cError *C.larodError) *LarodError {
	if cError == nil {
		return nil
	}
	defer C.larodClearError(&cError)
	return &LarodError{
		Code: LarodErrorCode(cError.code),
		Msg:  C.GoString(cError.msg),
	}
}
