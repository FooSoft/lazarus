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

func NewBitReader(reader io.Reader) *BitReader {
	return &BitReader{reader: reader}
}

func (r *BitReader) ReadBitsSigned(count int) (int64, error) {
	value, err := r.readBits(count)
	if err != nil {
		return 0, err
	}

	if count > 0 {
		valueMasked := value &^ (1 << uint(count-1))
		if valueMasked != value {
			return -int64(valueMasked), nil
		}
	}

	return int64(value), nil
}

func (r *BitReader) ReadBitsUnsigned(count int) (uint64, error) {
	return r.readBits(count)
}

func (r *BitReader) ReadBitFlag() (bool, error) {
	value, err := r.readBits(1)
	if err != nil {
		return false, err
	}

	return value == 1, nil
}

func (r *BitReader) readBits(count int) (uint64, error) {
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

		r.offset += bitsRead
		count -= bitsRead
	}

	return value, nil
}
