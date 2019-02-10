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
	CodedSize           uint32
	HasRawPixelEncoding bool
	CompressEqualCells  bool
	Variable0Bits       uint8
	WidthBits           uint8
	HeightBits          uint8
	OffsetXBits         uint8
	OffsetYBits         uint8
	OptionalBytesBits   uint8
	CodedBytesBits      uint8
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

	dirHead.CodedSize = uint32(bitReader.ReadUint(32))
	dirHead.HasRawPixelEncoding = bitReader.ReadBool()
	dirHead.CompressEqualCells = bitReader.ReadBool()
	dirHead.Variable0Bits = uint8(bitReader.ReadUint(4))
	dirHead.WidthBits = uint8(bitReader.ReadUint(4))
	dirHead.HeightBits = uint8(bitReader.ReadUint(4))
	dirHead.OffsetXBits = uint8(bitReader.ReadInt(4))
	dirHead.OffsetYBits = uint8(bitReader.ReadInt(4))
	dirHead.OptionalBytesBits = uint8(bitReader.ReadUint(4))
	dirHead.CodedBytesBits = uint8(bitReader.ReadUint(4))

	return &dirHead
}

func readFrameHeader(bitReader *streaming.BitReader, dirHead directionHeader) *frameHeader {
	var frameHead frameHeader

	frameHead.Variable0 = uint32(bitReader.ReadUintPacked(int(dirHead.Variable0Bits)))
	frameHead.Width = uint32(bitReader.ReadUintPacked(int(dirHead.WidthBits)))
	frameHead.Height = uint32(bitReader.ReadUintPacked(int(dirHead.HeightBits)))
	frameHead.OffsetX = int32(bitReader.ReadIntPacked(int(dirHead.OffsetXBits)))
	frameHead.OffsetY = int32(bitReader.ReadIntPacked(int(dirHead.OffsetYBits)))
	frameHead.OptionalBytes = uint32(bitReader.ReadUintPacked(int(dirHead.OptionalBytesBits)))
	frameHead.CodedBytes = uint32(bitReader.ReadUintPacked(int(dirHead.CodedBytesBits)))
	frameHead.FrameBottomUp = bitReader.ReadBool()

	return &frameHead
}

func readDirection(reader io.ReadSeeker, fileHead fileHeader) error {
	bitReader := streaming.NewBitReader(reader)

	dirHead := readDirectionHeader(bitReader)
	if err := bitReader.Error(); err != nil {
		return err
	}

	frameHead := readFrameHeader(bitReader, *dirHead)
	if err := bitReader.Error(); err != nil {
		return err
	}

	fmt.Printf("%+v\n", dirHead)
	fmt.Printf("%+v\n", frameHead)

	return nil
}
