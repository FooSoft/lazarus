package streaming

import (
	"io"
)

type BitReader struct {
	reader    io.Reader
	bitOffset int
	tailByte  byte
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

	return twosComplement(value, count), nil
}

func (r *BitReader) ReadUint64(count int) (uint64, error) {
	buffer, bitOffset, err := r.readBytes(count)
	if err != nil {
		return 0, err
	}

	return readBits(buffer, bitOffset, count), nil
}

func (r *BitReader) readBytes(count int) ([]byte, int, error) {
	if count == 0 {
		return nil, 0, nil
	}

	var (
		bitOffsetInByte = r.bitOffset % 8
		bitsLeftInByte  = 8 - bitOffsetInByte
		bytesNeeded     = 1
	)

	if bitsLeftInByte < count {
		bitsOverrun := count - bitsLeftInByte
		bytesNeeded += bitsOverrun / 8
		if bitsOverrun%8 != 0 {
			bytesNeeded++
		}
	}

	buffer := make([]byte, bytesNeeded)
	bufferToRead := buffer
	if bitsLeftInByte < 8 {
		buffer[0] = r.tailByte
		bufferToRead = buffer[1:]
	}

	if _, err := io.ReadAtLeast(r.reader, bufferToRead, len(bufferToRead)); err != nil {
		return nil, 0, err
	}

	r.bitOffset += count
	r.tailByte = buffer[bytesNeeded-1]

	return buffer, bitOffsetInByte, nil
}

func readBits(buffer []byte, bitOffset, count int) uint64 {
	var result uint64

	remainder := count
	for byteOffset := 0; remainder > 0; byteOffset++ {
		bitsRead := 8 - bitOffset
		if bitsRead > remainder {
			bitsRead = remainder
		}

		bufferByte := buffer[byteOffset]
		bufferByte >>= uint(bitOffset)
		bufferByte &= ^(0xff << uint(bitsRead))

		result |= (uint64(bufferByte) << uint(count-remainder))

		remainder -= bitsRead
		bitOffset = 0
	}

	return result
}

func twosComplement(value uint64, bits int) int64 {
	signMask := uint64(1 << uint(bits-1))
	if value&signMask == 0 {
		return int64(value &^ signMask)
	} else {
		valueMask := ^(^uint64(0) << uint(bits-1))
		return -int64(valueMask & (^value + 1))
	}
}
