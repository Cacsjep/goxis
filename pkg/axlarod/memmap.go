package axlarod

/*
#include <sys/mman.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>
#include <string.h>
#include <stdlib.h>
#include <errno.h>

int createAndMapTmpFile(char* fileName, size_t fileSize, void** mappedAddr, int* fd) {
    *fd = mkstemp(fileName);
    if (*fd < 0) {
        return errno;
    }

    if (ftruncate(*fd, (off_t)fileSize) < 0) {
        close(*fd);
        return errno;
    }

	if (unlink(fileName)) {
        return errno;
    }

    *mappedAddr = mmap(NULL, fileSize, PROT_READ | PROT_WRITE, MAP_SHARED, *fd, 0);
    if (*mappedAddr == MAP_FAILED) {
        close(*fd);
        return errno;
    }

    return 0; // Success
}
*/
import "C"
import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"unsafe"
)

// MemMapFile represents a memory mapped file.
type MemMapFile struct {
	MemoryAddress unsafe.Pointer
	File          *os.File
	Size          uint
	FilePattern   string
	UsePitch0Size bool
}

// MemMapConfiguration represents the configuration for memory mapped files.
type MemMapConfiguration struct {
	InputTmpMapFiles  map[int]*MemMapFile
	OutputTmpMapFiles map[int]*MemMapFile
}

// configureMemMapFile configures the memory mapped file for the tensor.
func (model *LarodModel) configureMemMapFile(file_map map[int]*MemMapFile, tensors []*LarodTensor, pitches *LarodTensorPitches) error {
	var err error
	for i, f := range file_map {
		// Create file other wise reuse fd
		if f.File == nil {
			f.FilePattern = generateRandomMapFilePattern(fmt.Sprintf("%s-in", model.Name), i)
			if f.UsePitch0Size {
				f.Size = pitches.Pitches[0]
			}
			f.MemoryAddress, f.File, err = CreateAndMapTmpFile(f.FilePattern, f.Size)
			if err != nil {
				return err
			}
		}
		tensors[i].SetTensorFd(f.File.Fd())
		tensors[i].MemMapFile = f
	}
	return nil
}

// MapModelTmpFiles maps the temporary files for the model.
func (model *LarodModel) MapModelTmpFiles(m *MemMapConfiguration) error {
	if err := model.configureMemMapFile(m.InputTmpMapFiles, model.Inputs, model.InputPitches); err != nil {
		return err
	}
	if err := model.configureMemMapFile(m.OutputTmpMapFiles, model.Outputs, model.OutputPitches); err != nil {
		return err
	}
	return nil
}

// UnmapMemory unmaps the memory mapped file
func (t *MemMapFile) UnmapMemory() error {
	if ret := C.munmap(t.MemoryAddress, C.size_t(t.Size)); ret != 0 {
		return fmt.Errorf("failed to unmap memory: return code %d", int(ret))
	}
	return nil
}

// Seek the memory mapped file to the beginning
func (t *MemMapFile) Rewind() error {
	_, err := t.File.Seek(0, 0)
	return err
}

// CreateAndMapTmpFile creates a temporary file and maps it to memory
func CreateAndMapTmpFile(fileNamePattern string, fileSize uint) (unsafe.Pointer, *os.File, error) {
	cFileName := C.CString(fileNamePattern)
	defer C.free(unsafe.Pointer(cFileName))

	var mappedAddr unsafe.Pointer
	var fd C.int

	errCode := C.createAndMapTmpFile(cFileName, C.size_t(fileSize), &mappedAddr, &fd)
	if errCode != 0 {
		return nil, nil, fmt.Errorf("error creating and mapping file: %s", C.GoString(C.strerror(errCode)))
	}

	file := os.NewFile(uintptr(fd), C.GoString(cFileName))

	return mappedAddr, file, nil
}

// CopyDataToMappedMemory copies data to the mapped memory
func CopyDataToMappedMemory(mappedAddr unsafe.Pointer, data []byte) error {
	dataLen := len(data)
	if dataLen == 0 {
		return fmt.Errorf("data slice is empty")
	}
	C.memcpy(mappedAddr, unsafe.Pointer(&data[0]), C.size_t(dataLen))
	return nil
}

// CopyDataFromMappedMemory copies data from the mapped memory
func CopyDataFromMappedMemory(mappedAddr unsafe.Pointer, size int) ([]byte, error) {
	if size <= 0 {
		return nil, fmt.Errorf("size must be positive")
	}
	dataSlice := unsafe.Slice((*byte)(mappedAddr), size)
	copiedData := make([]byte, size)
	copy(copiedData, dataSlice)
	return copiedData, nil
}

// randomString generates a random string of a specified length using the provided character set.
func randomString(length int, charset string) string {
	var output strings.Builder
	charsetLength := len(charset)
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(charsetLength)
		output.WriteByte(charset[randomIndex])
	}
	return output.String()
}

// generatePath creates a path string with a random 10 character string.
func generateRandomMapFilePattern(prefix string, tensor_io_index int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomPart := randomString(10, charset)
	return fmt.Sprintf("/tmp/larod-%s-%d.%s-XXXXXX", prefix, tensor_io_index, randomPart)
}
