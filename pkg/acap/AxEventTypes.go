package acap

/*
#cgo pkg-config: axevent
#include <axsdk/axevent.h>
*/
import "C"

type AXEventValueType int

const (
	AXValueTypeInt     AXEventValueType = AXEventValueType(C.AX_VALUE_TYPE_INT)
	AXValueTypeBool    AXEventValueType = AXEventValueType(C.AX_VALUE_TYPE_BOOL)
	AXValueTypeDouble  AXEventValueType = AXEventValueType(C.AX_VALUE_TYPE_DOUBLE)
	AXValueTypeString  AXEventValueType = AXEventValueType(C.AX_VALUE_TYPE_STRING)
	AXValueTypeElement AXEventValueType = AXEventValueType(C.AX_VALUE_TYPE_ELEMENT)
)

var (
	OnfivNameSpaceTns1    string = "tns1"
	OnfivNameSpaceTnsAxis string = "tnsaxis"
)
