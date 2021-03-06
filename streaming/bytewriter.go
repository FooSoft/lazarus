package streaming

import (
	"errors"
	"io"
)

type Writer struct {
	data   []byte
	offset int
}

func NewWriter(data []byte) *Writer {
	return &Writer{data, 0}
}

func (w *Writer) Write(data []byte) (int, error) {
	length := len(data)
	if w.offset+length > len(w.data) {
		length = len(w.data) - w.offset
	}

	if length == 0 {
		return 0, errors.New("cannot write past end of buffer")
	}

	copy(w.data[w.offset:], data[:length])
	w.offset += length

	return length, nil
}

func (w *Writer) Seek(offset int64, whence int) (int64, error) {
	result := w.offset
	switch whence {
	case io.SeekStart:
		result = int(offset)
	case io.SeekCurrent:
		result = w.offset + int(offset)
	case io.SeekEnd:
		result = len(w.data) - int(offset)
	}

	if result < 0 {
		return int64(w.offset), errors.New("cannot seek before beginning of buffer")
	}
	if result >= len(w.data) {
		return int64(w.offset), errors.New("cannot seek past end of buffer")
	}

	w.offset = result
	return int64(w.offset), nil
}
