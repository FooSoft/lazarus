package dcc

import (
	"encoding/binary"
	"fmt"
	"io"
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
	var fileHead fileHeader
	if err := binary.Read(reader, binary.LittleEndian, &fileHead); err != nil {
		return nil, err
	}

	dirOffsets := make([]uint32, fileHead.DirCount)
	for i := 0; i < int(fileHead.DirCount); i++ {
		if err := binary.Read(reader, binary.LittleEndian, &dirOffsets[i]); err != nil {
			return nil, err
		}
	}

	fmt.Printf("%+v\n", fileHead)
	fmt.Printf("%+v\n", dirOffsets)
	return nil, nil
}
