package axlarodv2

import (
	"os"
	"unsafe"
)

type TmpFile struct {
	MemoryAddress unsafe.Pointer
	File          *os.File
	Size          uint
	FilePattern   string
	UsePitch0Size bool
}

type TempFileConfiguration struct {
	InputTmpMapFiles  map[int]*TmpFile
	OutputTmpMapFiles map[int]*TmpFile
}

type Tensor struct {
	ptr     *C.larodTensor
	fd      *os.File
	Map     unsafe.Pointer
	TmpFile TmpFile
}

type LarodModel struct {
	Connection *LarodConnection
	Model      *C.larodModel
	Inputs     []*Tensor
	Outputs    []*Tensor
	Access     C.larodAccess
	Chip       C.larodChip
	Name       string
	Params     *LarodMap
	Error      *LarodError
}
