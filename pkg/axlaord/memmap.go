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

type TmpFile struct {
	MemoryAddress unsafe.Pointer
	File          *os.File
	Size          uint
	FilePattern   string
	UsePitch0Size bool
}

func (t *TmpFile) UnmapMemory() error {
	if ret := C.munmap(t.MemoryAddress, C.size_t(t.Size)); ret != 0 {
		return fmt.Errorf("failed to unmap memory: return code %d", int(ret))
	}
	return nil
}

type ModelTmpMapDefiniton struct {
	InputTmpMapFiles  map[int]*TmpFile
	OutputTmpMapFiles map[int]*TmpFile
}

func (l *LarodModel) MapModelTmpFiles(m *ModelTmpMapDefiniton) error {
	var err error
	for i, tmp_file := range m.InputTmpMapFiles {
		tmp_file.FilePattern = generateRandomMapFilePattern()
		if tmp_file.UsePitch0Size {
			tmp_file.Size = l.LarodModelIO.InputPitches.Pitches[0]
		}
		tmp_file.MemoryAddress, tmp_file.File, err = CreateAndMapTmpFile(tmp_file.FilePattern, tmp_file.Size)
		if err != nil {
			return err
		}
		l.LarodModelIO.Inputs[i].SetTensorFd(tmp_file.File.Fd())
		l.LarodModelIO.Inputs[i].TmpFile = tmp_file
		fmt.Println("Input Tensor", i, "mapped to", tmp_file.FilePattern, "with size", tmp_file.Size, "address", tmp_file.MemoryAddress)
	}
	for i, tmp_file := range m.OutputTmpMapFiles {
		tmp_file.FilePattern = generateRandomMapFilePattern()
		tmp_file.MemoryAddress, tmp_file.File, err = CreateAndMapTmpFile(tmp_file.FilePattern, tmp_file.Size)
		if err != nil {
			return err
		}
		l.LarodModelIO.Outputs[i].SetTensorFd(tmp_file.File.Fd())
		l.LarodModelIO.Outputs[i].TmpFile = tmp_file
		fmt.Println("Output Tensor", i, "mapped to", tmp_file.FilePattern, "with size", tmp_file.Size, "address", tmp_file.MemoryAddress)
	}
	return nil
}

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

func CopyDataToMappedMemory(mappedAddr unsafe.Pointer, data []byte) error {
	// Calculate the length of the data to copy
	dataLen := len(data)
	if dataLen == 0 {
		return fmt.Errorf("data slice is empty")
	}

	// Perform the memory copy
	C.memcpy(mappedAddr, unsafe.Pointer(&data[0]), C.size_t(dataLen))
	return nil
}

func CopyDataFromMappedMemory(mappedAddr unsafe.Pointer, size int) ([]byte, error) {
	if size <= 0 {
		return nil, fmt.Errorf("size must be positive")
	}
	data := make([]byte, size)
	C.memcpy(unsafe.Pointer(&data[0]), mappedAddr, C.size_t(size))
	return data, nil
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
func generateRandomMapFilePattern() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomPart := randomString(10, charset)
	return "/tmp/larod." + randomPart + "-XXXXXX"
}
