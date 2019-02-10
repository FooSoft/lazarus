package dcc

import (
	"encoding/binary"
	"fmt"
	"io"

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
	CodedSize           uint
	HasRawPixelEncoding bool
	CompressEqualCells  bool
	Variable0Bits       uint
	WidthBits           uint
	HeightBits          uint
	OffsetXBits         uint
	OffsetYBits         uint
	OptionalBytesBits   uint
	CodedBytesBits      uint
}

type frameHeader struct {
	Variable0     uint
	Width         uint
	Height        uint
	OffsetX       int
	OffsetY       int
	OptionalBytes uint
	CodedBytes    uint
	FrameBottomUp bool
	Extents       extents
}

func NewFromReader(reader io.ReadSeeker) (*Sprite, error) {
	var fileHead fileHeader
	if err := binary.Read(reader, binary.LittleEndian, &fileHead); err != nil {
		return nil, err
	}

	for i := 0; i < int(fileHead.DirCount); i++ {
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

		if err := readDirection(reader, fileHead); err != nil {
			return nil, err
		}

		if _, err := reader.Seek(offset, io.SeekStart); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func readDirectionHeader(bitReader *streaming.BitReader) *directionHeader {
	var dirHead directionHeader

	dirHead.CodedSize = uint(bitReader.ReadUint(32))
	dirHead.HasRawPixelEncoding = bitReader.ReadBool()
	dirHead.CompressEqualCells = bitReader.ReadBool()
	dirHead.Variable0Bits = uint(bitReader.ReadUint(4))
	dirHead.WidthBits = uint(bitReader.ReadUint(4))
	dirHead.HeightBits = uint(bitReader.ReadUint(4))
	dirHead.OffsetXBits = uint(bitReader.ReadInt(4))
	dirHead.OffsetYBits = uint(bitReader.ReadInt(4))
	dirHead.OptionalBytesBits = uint(bitReader.ReadUint(4))
	dirHead.CodedBytesBits = uint(bitReader.ReadUint(4))

	return &dirHead
}

func readFrameHeader(bitReader *streaming.BitReader, dirHead directionHeader) *frameHeader {
	var frameHead frameHeader

	frameHead.Variable0 = uint(bitReader.ReadUintPacked(int(dirHead.Variable0Bits)))
	frameHead.Width = uint(bitReader.ReadUintPacked(int(dirHead.WidthBits)))
	frameHead.Height = uint(bitReader.ReadUintPacked(int(dirHead.HeightBits)))
	frameHead.OffsetX = int(bitReader.ReadIntPacked(int(dirHead.OffsetXBits)))
	frameHead.OffsetY = int(bitReader.ReadIntPacked(int(dirHead.OffsetYBits)))
	frameHead.OptionalBytes = uint(bitReader.ReadUintPacked(int(dirHead.OptionalBytesBits)))
	frameHead.CodedBytes = uint(bitReader.ReadUintPacked(int(dirHead.CodedBytesBits)))
	frameHead.FrameBottomUp = bitReader.ReadBool()

	return &frameHead
}

func readDirection(reader io.ReadSeeker, fileHead fileHeader) error {
	bitReader := streaming.NewBitReader(reader)

	dirHead := readDirectionHeader(bitReader)
	frameHead := readFrameHeader(bitReader, *dirHead)

	fmt.Printf("%+v\n", dirHead)
	fmt.Printf("%+v\n", frameHead)

	return nil
}
