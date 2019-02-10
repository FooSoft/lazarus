package streaming

import (
	"errors"
	"io"
)

type BitReader struct {
	reader    io.Reader
	bitOffset int
	tailByte  byte
	err       error
}

func NewBitReader(reader io.Reader) *BitReader {
	return &BitReader{reader: reader}
}

func (r *BitReader) ReadBool() bool {
	return r.ReadUint(1) != 0
}

func (r *BitReader) ReadIntPacked(countPacked int) int64 {
	if r.err != nil {
		return 0
	}

	count, err := unpackSize(countPacked)
	if err != nil {
		r.err = err
		return 0
	}

	return r.ReadInt(count)
}

func (r *BitReader) ReadInt(count int) int64 {
	return twosComplement(r.ReadUint(count), count)
}

func (r *BitReader) ReadUintPacked(countPacked int) uint64 {
	if r.err != nil {
		return 0
	}

	count, err := unpackSize(countPacked)
	if err != nil {
		r.err = err
		return 0
	}

	return r.ReadUint(count)
}

func (r *BitReader) ReadUint(count int) uint64 {
	if r.err != nil {
		return 0
	}

	buffer, bitOffset, err := r.readBytes(count)
	if err != nil {
		r.err = err
		return 0
	}

	return readBits(buffer, bitOffset, count)
}

func (r *BitReader) Error() error {
	return r.err
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

func unpackSize(sizePacked int) (int, error) {
	sizes := []int{0, 1, 2, 4, 6, 8, 10, 12, 14, 16, 20, 24, 26, 28, 30, 32}
	if sizePacked >= len(sizes) {
		return 0, errors.New("invalid packed size")
	}

	return sizes[sizePacked], nil
}
