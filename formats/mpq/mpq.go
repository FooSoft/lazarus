package mpq

// #cgo windows CFLAGS: -D_MPQ_WINDOWS
// #cgo windows LDFLAGS: -Lstormlib -lstorm -lwininet -lz -lbz2 -lstdc++
// #cgo linux CFLAGS: -D_MPQ_LINUX
// #cgo linux LDFLAGS: -L./stormlib/ -lstorm -lz -lbz2 -lstdc++
// #ifdef _MPQ_WINDOWS
// #include "native_windows.h"
// #endif
// #ifdef _MPQ_LINUX
// #include "native_linux.h"
// #endif
import "C"
import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unsafe"
)

type File interface {
	Read(data []byte) (int, error)
	Seek(offset int64, whence int) (int64, error)
	Close() error
}

type MpqArchive interface {
	OpenFile(path string) (File, error)
	GetPaths() []string
	Close() error
}

func NewFromFile(path string) (MpqArchive, error) {
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))

	a := new(archive)
	if result := C.SFileOpenArchive(cs, 0, 0, &a.handle); result == 0 {
		return nil, fmt.Errorf("failed to open archive (%d)", getLastError())
	}

	if err := a.buildPathMap(); err != nil {
		a.Close()
		return nil, err
	}

	return a, nil
}

type file struct {
	handle unsafe.Pointer
}

func (f *file) Read(data []byte) (int, error) {
	var bytesRead int
	if result := C.SFileReadFile(f.handle, unsafe.Pointer(&data[0]), C.uint(len(data)), (*C.uint)(unsafe.Pointer(&bytesRead)), nil); result == 0 {
		lastError := getLastError()
		if lastError == C.ERROR_HANDLE_EOF {
			return bytesRead, io.EOF
		}

		return 0, fmt.Errorf("failed to read file (%d)", lastError)
	}

	return bytesRead, nil
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	var method uint
	switch whence {
	case io.SeekStart:
		method = C.FILE_BEGIN
	case io.SeekCurrent:
		method = C.FILE_CURRENT
	case io.SeekEnd:
		method = C.FILE_END
	}

	result := C.SFileSetFilePointer(f.handle, C.int(offset), nil, C.uint(method))
	if result == C.SFILE_INVALID_SIZE {
		return 0, fmt.Errorf("failed to set file pointer (%d)", getLastError())
	}

	return int64(result), nil
}

func (f *file) Close() error {
	if result := C.SFileCloseFile(f.handle); result == 0 {
		return fmt.Errorf("failed to close file (%d)", getLastError())
	}

	f.handle = nil
	return nil
}

type archive struct {
	handle unsafe.Pointer
	paths  map[string]string
}

func (a *archive) Close() error {
	if result := C.SFileCloseArchive(a.handle); result == 0 {
		return fmt.Errorf("failed to close archive (%d)", getLastError())
	}

	a.handle = nil
	a.paths = nil
	return nil
}

func (a *archive) OpenFile(path string) (File, error) {
	if pathInt, ok := a.paths[path]; ok {
		path = pathInt
	}

	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))

	file := new(file)
	if result := C.SFileOpenFileEx(a.handle, cs, 0, &file.handle); result == 0 {
		return nil, fmt.Errorf("failed to open file (%d)", getLastError())
	}

	return file, nil
}

func (a *archive) GetPaths() []string {
	var extPaths []string
	for extPath := range a.paths {
		extPaths = append(extPaths, extPath)
	}

	return extPaths
}

func (a *archive) buildPathMap() error {
	f, err := a.OpenFile("(listfile)")
	if err != nil {
		return err
	}
	defer f.Close()

	var buff bytes.Buffer
	if _, err := io.Copy(&buff, f); err != nil {
		return err
	}

	a.paths = make(map[string]string)

	lines := strings.Split(string(buff.Bytes()), "\r\n")
	for _, line := range lines {
		pathInt := strings.TrimSpace(line)
		if len(pathInt) > 0 {
			pathExt := sanitizePath(pathInt)
			a.paths[pathExt] = pathInt
		}
	}

	return nil
}

func sanitizePath(path string) string {
	return strings.ToLower(strings.Replace(path, "\\", string(os.PathSeparator), -1))
}

func getLastError() uint {
	return uint(C.GetLastError())
}
