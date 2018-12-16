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
	"strings"
	"unsafe"
)

type File interface {
	Read(data []byte) (int, error)
	GetSize() (int, error)
	Close() error
}

type Archive interface {
	OpenFile(path string) (File, error)
	GetPaths() ([]string, error)
	Close() error
}

func New(path string) (Archive, error) {
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))

	a := new(archive)
	if result := C.SFileOpenArchive(cs, 0, 0, &a.handle); result == 0 {
		return nil, fmt.Errorf("failed to open archive (%d)", getLastError())
	}

	return a, nil
}

type file struct {
	handle unsafe.Pointer
	offset int
	size   int
}

func (f *file) Read(data []byte) (int, error) {
	size, err := f.GetSize()
	if err != nil {
		return 0, err
	}

	bytesRemaining := size - f.offset
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

func (f *file) GetSize() (int, error) {
	if f.size != math.MaxUint32 {
		return f.size, nil
	}

	size := int(C.SFileGetFileSize(f.handle, nil))
	if size == -1 {
		return 0, fmt.Errorf("failed to get file size (%d)", getLastError())
	}

	f.size = size
	return size, nil
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

type archive struct {
	handle unsafe.Pointer
	paths  []string
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
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))

	file := &file{size: math.MaxUint32}
	if result := C.SFileOpenFileEx(a.handle, cs, 0, &file.handle); result == 0 {
		return nil, fmt.Errorf("failed to open file (%d)", getLastError())
	}

	return file, nil
}

func (a *archive) GetPaths() ([]string, error) {
	if len(a.paths) > 0 {
		return a.paths, nil
	}

	f, err := a.OpenFile("(listfile)")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buff bytes.Buffer
	if _, err := io.Copy(&buff, f); err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(buff.Bytes()), "\r\n") {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			a.paths = append(a.paths, line)
		}
	}

	return a.paths, nil
}

func getLastError() uint {
	return uint(C.GetLastError())
}
