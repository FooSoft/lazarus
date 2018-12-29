package mpq

// #ifdef _MPQ_WINDOWS
// 		#include <windows.h>
// 		#include <stdlib.h>
// #endif
// #ifdef _MPQ_LINUX
//		#include <stdlib.h>
//		#define WINAPI
// 		DWORD GetLastError();
// #endif
//
// #define DWORD unsigned int
// #define LPDWORD unsigned int *
// #define LPOVERLAPPED void *
// #define TCHAR char
// #define HANDLE void *
// #define bool unsigned char
// #define LONG int
//
// bool  WINAPI SFileOpenArchive(const TCHAR * szMpqName, DWORD dwPriority, DWORD dwFlags, HANDLE * phMpq);
// bool  WINAPI SFileCloseArchive(HANDLE hMpq);
// bool  WINAPI SFileOpenFileEx(HANDLE hMpq, const char * szFileName, DWORD dwSearchScope, HANDLE * phFile);
// DWORD WINAPI SFileSetFilePointer(HANDLE hFile, LONG lFilePos, LONG * plFilePosHigh, DWORD dwMoveMethod);
// bool  WINAPI SFileReadFile(HANDLE hFile, void * lpBuffer, DWORD dwToRead, LPDWORD pdwRead, LPOVERLAPPED lpOverlapped);
// bool  WINAPI SFileCloseFile(HANDLE hFile);
//
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
	Seek(offset int64, whence int) (int64, error)
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
}

func (f *file) Read(data []byte) (int, error) {
	var bytesRead int
	if result := C.SFileReadFile(f.handle, unsafe.Pointer(&data[0]), C.uint(len(data)), (*C.uint)(unsafe.Pointer(&bytesRead)), nil); result == 0 {
		lastError := getLastError()
		if lastError == sysEOF { // ERROR_HANDLE_EOF
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
		method = 0 // FILE_BEGIN
	case io.SeekCurrent:
		method = 1 // FILE_CURRENT
	case io.SeekEnd:
		method = 2 // FILE_END
	}

	result := C.SFileSetFilePointer(f.handle, C.int(offset), nil, C.uint(method))
	if result == math.MaxUint32 { // SFILE_INVALID_SIZE
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
