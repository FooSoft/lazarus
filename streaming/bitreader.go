package streaming

import (
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

func (r *BitReader) ReadBool() (bool, error) {
	value, err := r.ReadUint64(1)
	return value == 1, err
}

func (r *BitReader) ReadInt8(count int) (int8, error) {
	value, err := r.ReadInt64(count)
	return int8(value), err
}

func (r *BitReader) ReadUint8(count int) (uint8, error) {
	value, err := r.ReadUint64(count)
	return uint8(value), err
}

func (r *BitReader) ReadInt16(count int) (int16, error) {
	value, err := r.ReadInt64(count)
	return int16(value), err
}

func (r *BitReader) ReadUint16(count int) (uint16, error) {
	value, err := r.ReadUint64(count)
	return uint16(value), err
}

func (r *BitReader) ReadInt32(count int) (int32, error) {
	value, err := r.ReadInt64(count)
	return int32(value), err
}

func (r *BitReader) ReadUint32(count int) (uint32, error) {
	value, err := r.ReadUint64(count)
	return uint32(value), err
}

func (r *BitReader) ReadInt64(count int) (int64, error) {
	value, err := r.ReadUint64(count)
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

func (r *BitReader) ReadUint64(count int) (uint64, error) {
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
