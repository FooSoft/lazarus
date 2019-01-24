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

	var (
		header directionHeader
		err    error
	)

	header.CodedSize, err = r.ReadUint32(32)
	if err != nil {
		return nil, err
	}

	header.HasRawPixelEncoding, err = r.ReadBool()
	if err != nil {
		return nil, err
	}

	header.CompressEqualCells, err = r.ReadBool()
	if err != nil {
		return nil, err
	}

	header.Variable0Bits, err = r.ReadUint32(4)
	if err != nil {
		return nil, err
	}

	header.WidthBits, err = r.ReadUint32(4)
	if err != nil {
		return nil, err
	}

	header.HeightBits, err = r.ReadUint32(4)
	if err != nil {
		return nil, err
	}

	header.OffsetXBits, err = r.ReadInt32(4)
	if err != nil {
		return nil, err
	}

	header.OffsetYBits, err = r.ReadInt32(4)
	if err != nil {
		return nil, err
	}

	header.OptionalBytesBits, err = r.ReadUint32(4)
	if err != nil {
		return nil, err
	}

	header.CodedBytesBits, err = r.ReadUint32(4)
	if err != nil {
		return nil, err
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
