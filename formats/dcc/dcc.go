package dcc

import (
	"encoding/binary"
	"io"
	"log"

	"github.com/FooSoft/lazarus/streaming"
)

type Sprite struct {
}

type extents struct {
	x1 int32
	y1 int32
	x2 int32
	y2 int32
}

type fileHeader struct {
	Signature    uint8
	Version      uint8
	DirCount     uint8
	FramesPerDir uint32
	Tag          uint32
	FinalDc6Size uint32
}

type directionHeader struct {
	CodedSize           uint32
	HasRawPixelEncoding bool
	CompressEqualCells  bool
	Variable0Bits       uint32
	WidthBits           uint32
	HeightBits          uint32
	OffsetXBits         int32
	OffsetYBits         int32
	OptionalBytesBits   uint32
	CodedBytesBits      uint32
}

type frameHeader struct {
	Variable0     uint32
	Width         uint32
	Height        uint32
	OffsetX       int32
	OffsetY       int32
	OptionalBytes uint32
	CodedBytes    uint32
	FrameBottomUp bool
	Extents       extents
}

func NewFromReader(reader io.ReadSeeker) (*Sprite, error) {
	var header fileHeader
	if err := binary.Read(reader, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	for i := 0; i < int(header.DirCount); i++ {
		var offsetDir uint32
		if err := binary.Read(reader, binary.LittleEndian, &offsetDir); err != nil {
			return nil, err
		}

		offset, err := reader.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, err
		}

		if _, err := reader.Seek(int64(offsetDir), io.SeekStart); err != nil {
			return nil, err
		}

		if err := readDirection(reader); err != nil {
			return nil, err
		}

		if _, err := reader.Seek(offset, io.SeekStart); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func readDirectionHeader(reader io.ReadSeeker) (*directionHeader, error) {
	r := streaming.NewBitReader(reader)

	codedSize, err := r.ReadBitsUnsigned(32)
	if err != nil {
		return nil, err
	}

	hasRawPixelEncoding, err := r.ReadBitsUnsigned(1)
	if err != nil {
		return nil, err
	}

	compressEqualCells, err := r.ReadBitsUnsigned(1)
	if err != nil {
		return nil, err
	}

	variable0Bits, err := r.ReadBitsUnsigned(4)
	if err != nil {
		return nil, err
	}

	widthBits, err := r.ReadBitsUnsigned(4)
	if err != nil {
		return nil, err
	}

	heightBits, err := r.ReadBitsUnsigned(4)
	if err != nil {
		return nil, err
	}

	offsetXBits, err := r.ReadBitsUnsigned(4)
	if err != nil {
		return nil, err
	}

	offsetYBits, err := r.ReadBitsUnsigned(4)
	if err != nil {
		return nil, err
	}

	optionalBytesBits, err := r.ReadBitsUnsigned(4)
	if err != nil {
		return nil, err
	}

	codedBytesBits, err := r.ReadBitsUnsigned(4)
	if err != nil {
		return nil, err
	}

	header := directionHeader{
		CodedSize:           uint32(codedSize),
		HasRawPixelEncoding: hasRawPixelEncoding == 1,
		CompressEqualCells:  compressEqualCells == 1,
		Variable0Bits:       uint32(variable0Bits),
		WidthBits:           uint32(widthBits),
		HeightBits:          uint32(heightBits),
		OffsetXBits:         int32(offsetXBits),
		OffsetYBits:         int32(offsetYBits),
		OptionalBytesBits:   uint32(optionalBytesBits),
		CodedBytesBits:      uint32(codedBytesBits),
	}

	return &header, nil
}

func readDirection(reader io.ReadSeeker) error {
	header, err := readDirectionHeader(reader)
	if err != nil {
		return err
	}

	log.Printf("%+v\n", header)
	return nil
}
