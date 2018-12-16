package mpq

// #cgo LDFLAGS: -L./stormlib/ -lstorm -lz -lbz2 -lstdc++
// #include <stdlib.h>
// #define WINAPI
// #define DWORD   		unsigned int
// #define HANDLE  		void *
// #define LPDWORD 		unsigned int *
// #define LPOVERLAPPED void *
// #define TCHAR   		char
// #define bool    		unsigned char
// bool  WINAPI SFileOpenArchive(const TCHAR * szMpqName, DWORD dwPriority, DWORD dwFlags, HANDLE * phMpq);
// bool  WINAPI SFileCloseArchive(HANDLE hMpq);
// bool  WINAPI SFileOpenFileEx(HANDLE hMpq, const char * szFileName, DWORD dwSearchScope, HANDLE * phFile);
// DWORD WINAPI SFileGetFileSize(HANDLE hFile, LPDWORD pdwFileSizeHigh);
// bool  WINAPI SFileReadFile(HANDLE hFile, void * lpBuffer, DWORD dwToRead, LPDWORD pdwRead, LPOVERLAPPED lpOverlapped);
// bool  WINAPI SFileCloseFile(HANDLE hFile);
// DWORD GetLastError();
import "C"
import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"unsafe"
)

type File interface {
	Read(data []byte) (int, error)
	GetSize() int
	Close() error
}

type Archive interface {
	OpenFile(path string) (File, error)
	GetPaths() []string
	Close() error
}

func NewFromFile(path string) (Archive, error) {
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
	offset int
	size   int
}

func (f *file) Read(data []byte) (int, error) {
	bytesRemaining := f.size - f.offset
	if bytesRemaining == 0 {
		return 0, io.EOF
	}

	bytesRequested := len(data)
	if bytesRequested > bytesRemaining {
		bytesRequested = bytesRemaining
	}

	var bytesRead int
	if result := C.SFileReadFile(f.handle, unsafe.Pointer(&data[0]), C.unsigned(bytesRequested), (*C.unsigned)(unsafe.Pointer(&bytesRead)), nil); result == 0 {
		return 0, fmt.Errorf("failed to read file (%d)", getLastError())
	}

	f.offset += bytesRead
	return bytesRead, nil
}

func (f *file) GetSize() int {
	return f.size
}

func (f *file) Close() error {
	if result := C.SFileCloseFile(f.handle); result == 0 {
		return fmt.Errorf("failed to close file (%d)", getLastError())
	}

	f.handle = nil
	f.offset = 0
	f.size = 0

	return nil
}

func (f *file) buildSize() error {
	size := int(C.SFileGetFileSize(f.handle, nil))
	if size == -1 {
		return fmt.Errorf("failed to get file size (%d)", getLastError())
	}

	f.size = size
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

	file := &file{size: math.MaxUint32}
	if result := C.SFileOpenFileEx(a.handle, cs, 0, &file.handle); result == 0 {
		return nil, fmt.Errorf("failed to open file (%d)", getLastError())
	}

	if err := file.buildSize(); err != nil {
		file.Close()
		return nil, err
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
			pathExt := santizePath(pathInt)
			a.paths[pathExt] = pathInt
		}
	}

	return nil
}

func santizePath(path string) string {
	return strings.ToLower(strings.Replace(path, "\\", string(os.PathSeparator), -1))
}

func getLastError() uint {
	return uint(C.GetLastError())
}
