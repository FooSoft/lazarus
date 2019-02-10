package dcc

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/FooSoft/lazarus/streaming"
)

type Sprite struct {
}

type bounds struct {
	x1 int
	y1 int
	x2 int
	y2 int
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

type direction struct {
	header directionHeader
	frames []frame
	bounds bounds
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
}

type frame struct {
	header frameHeader
	bounds bounds
}

func NewFromReader(reader io.ReadSeeker) (*Sprite, error) {
	var fileHead fileHeader
	if err := binary.Read(reader, binary.LittleEndian, &fileHead); err != nil {
		return nil, err
	}

	var directions []direction
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

		dirData, err := readDirection(reader, fileHead)
		if err != nil {
			return nil, err
		}
		directions = append(directions, *dirData)

		if _, err := reader.Seek(offset, io.SeekStart); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func readDirectionHeader(bitReader *streaming.BitReader) directionHeader {
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

	return dirHead
}

func readFrameHeader(bitReader *streaming.BitReader, dirHead directionHeader) frameHeader {
	var frameHead frameHeader

	frameHead.Variable0 = uint32(bitReader.ReadUintPacked(int(dirHead.Variable0Bits)))
	frameHead.Width = uint32(bitReader.ReadUintPacked(int(dirHead.WidthBits)))
	frameHead.Height = uint32(bitReader.ReadUintPacked(int(dirHead.HeightBits)))
	frameHead.OffsetX = int32(bitReader.ReadIntPacked(int(dirHead.OffsetXBits)))
	frameHead.OffsetY = int32(bitReader.ReadIntPacked(int(dirHead.OffsetYBits)))
	frameHead.OptionalBytes = uint32(bitReader.ReadUintPacked(int(dirHead.OptionalBytesBits)))
	frameHead.CodedBytes = uint32(bitReader.ReadUintPacked(int(dirHead.CodedBytesBits)))
	frameHead.FrameBottomUp = bitReader.ReadBool()

	return frameHead
}

func readFrameHeaders(bitReader *streaming.BitReader, fileHead fileHeader, dirHead directionHeader) ([]frameHeader, error) {
	var frameHeads []frameHeader
	for i := 0; i < int(fileHead.FramesPerDir); i++ {
		frameHead := readFrameHeader(bitReader, dirHead)
		if err := bitReader.Error(); err != nil {
			return nil, err
		}

		if frameHead.OptionalBytes != 0 {
			return nil, errors.New("optional frame data not supported")
		}

		if frameHead.FrameBottomUp {
			return nil, errors.New("bottom-up frames are not supported")
		}

		frameHeads = append(frameHeads, frameHead)
	}

	return frameHeads, nil
}

func readDirection(reader io.ReadSeeker, fileHead fileHeader) (*direction, error) {
	bitReader := streaming.NewBitReader(reader)

	dirHead := readDirectionHeader(bitReader)
	if err := bitReader.Error(); err != nil {
		return nil, err
	}

	frameHeads, err := readFrameHeaders(bitReader, fileHead, dirHead)
	if err != nil {
		return nil, err
	}

	var dirData direction
	for i, frameHead := range frameHeads {
		frameData := frame{
			header: frameHead,
			bounds: bounds{
				x1: int(frameHead.OffsetX),
				y1: int(frameHead.OffsetY) - int(frameHead.Height) + 1,
				x2: int(frameHead.OffsetX) + int(frameHead.Width),
				y2: int(frameHead.OffsetY) + 1,
			},
		}

		dirData.frames = append(dirData.frames, frameData)

		if i == 0 {
			dirData.bounds = frameData.bounds
		} else {
			if dirData.bounds.x1 > frameData.bounds.x1 {
				dirData.bounds.x1 = frameData.bounds.x1
			}
			if dirData.bounds.y1 > frameData.bounds.y1 {
				dirData.bounds.y1 = frameData.bounds.y1
			}
			if dirData.bounds.x2 < frameData.bounds.x2 {
				dirData.bounds.x2 = frameData.bounds.x2
			}
			if dirData.bounds.y2 < frameData.bounds.y2 {
				dirData.bounds.y2 = frameData.bounds.y2
			}
		}
	}

	return &dirData, nil
}
