package dcc

import (
	"encoding/binary"
	"errors"
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

func readDirectionHeader(reader io.ReadSeeker) (*directionHeader, error) {
	r := streaming.NewBitReader(reader)

	var (
		dirHead directionHeader
		err     error
	)

	dirHead.CodedSize, err = r.ReadUint32(32)
	if err != nil {
		return nil, err
	}

	dirHead.HasRawPixelEncoding, err = r.ReadBool()
	if err != nil {
		return nil, err
	}

	dirHead.CompressEqualCells, err = r.ReadBool()
	if err != nil {
		return nil, err
	}

	dirHead.Variable0Bits, err = r.ReadUint32(4)
	if err != nil {
		return nil, err
	}

	dirHead.WidthBits, err = r.ReadUint32(4)
	if err != nil {
		return nil, err
	}

	dirHead.HeightBits, err = r.ReadUint32(4)
	if err != nil {
		return nil, err
	}

	dirHead.OffsetXBits, err = r.ReadInt32(4)
	if err != nil {
		return nil, err
	}

	dirHead.OffsetYBits, err = r.ReadInt32(4)
	if err != nil {
		return nil, err
	}

	dirHead.OptionalBytesBits, err = r.ReadUint32(4)
	if err != nil {
		return nil, err
	}

	dirHead.CodedBytesBits, err = r.ReadUint32(4)
	if err != nil {
		return nil, err
	}

	return &dirHead, nil
}

func readFrameHeader(reader io.ReadSeeker, dirHead directionHeader) (*frameHeader, error) {
	r := streaming.NewBitReader(reader)

	var (
		frameHead frameHeader
		err       error
	)

	frameHead.Variable0, err = readPackedUint32(r, int(dirHead.Variable0Bits))
	if err != nil {
		return nil, err
	}

	frameHead.Width, err = readPackedUint32(r, int(dirHead.WidthBits))
	if err != nil {
		return nil, err
	}

	frameHead.Height, err = readPackedUint32(r, int(dirHead.HeightBits))
	if err != nil {
		return nil, err
	}

	frameHead.OffsetX, err = readPackedInt32(r, int(dirHead.OffsetXBits))
	if err != nil {
		return nil, err
	}

	frameHead.OffsetY, err = readPackedInt32(r, int(dirHead.OffsetYBits))
	if err != nil {
		return nil, err
	}

	frameHead.OptionalBytes, err = readPackedUint32(r, int(dirHead.OptionalBytesBits))
	if err != nil {
		return nil, err
	}

	frameHead.CodedBytes, err = readPackedUint32(r, int(dirHead.CodedBytesBits))
	if err != nil {
		return nil, err
	}

	frameHead.FrameBottomUp, err = r.ReadBool()
	if err != nil {
		return nil, err
	}

	return &frameHead, nil
}

func readDirection(reader io.ReadSeeker, fileHead fileHeader) error {
	dirHead, err := readDirectionHeader(reader)
	if err != nil {
		return err
	}

	frameHead, err := readFrameHeader(reader, *dirHead)
	if err != nil {
		return err
	}

	// fmt.Printf("%+v\n", dirHead)
	fmt.Printf("%+v\n", frameHead)

	return nil
}

func readPackedInt32(reader *streaming.BitReader, packedSize int) (int32, error) {
	width, err := unpackSize(packedSize)
	if err != nil {
		return 0, err
	}

	return reader.ReadInt32(width)
}

func readPackedUint32(reader *streaming.BitReader, packedSize int) (uint32, error) {
	width, err := unpackSize(packedSize)
	if err != nil {
		return 0, err
	}

	return reader.ReadUint32(width)
}

func unpackSize(packedSize int) (int, error) {
	sizes := []int{0, 1, 2, 4, 6, 8, 10, 12, 14, 16, 20, 24, 26, 28, 30, 32}
	if packedSize >= len(sizes) {
		return 0, errors.New("invalid packed size")
	}

	return sizes[packedSize], nil
}
