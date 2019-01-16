package streaming

import (
	"errors"
	"io"
)

type BitReader struct {
	reader io.Reader
	offset int
	buffer [1]byte
}

func NewReader(reader io.Reader) *BitReader {
	return &BitReader{reader: reader}
}

func (r *BitReader) ReadBits(count int) (uint64, error) {
	if count > 64 {
		return 0, errors.New("cannot read more than 64 bits at a time")
	}

	var value uint64
	for count > 0 {
		bitOffset := r.offset % 8
		bitsLeft := 8 - bitOffset
		if bitsLeft == 8 {
			if _, err := r.reader.Read(r.buffer[:]); err != nil {
				return 0, err
			}
		}

		bitsRead := count
		if bitsRead > bitsLeft {
			bitsRead = bitsLeft
		}

		buffer := r.buffer[0]
		buffer <<= uint(bitOffset)
		buffer >>= (uint(bitOffset) + uint(bitsLeft-bitsRead))

		value <<= uint(bitsRead)
		value |= uint64(buffer)
	}

	return value, nil
}
