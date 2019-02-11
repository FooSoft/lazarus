package dcc

import (
	"io"
	"log"

	"github.com/FooSoft/lazarus/streaming"
)

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

func readDirectionHeader(bitReader *streaming.BitReader) (*directionHeader, error) {
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

	if err := bitReader.Error(); err != nil {
		return nil, err
	}

	return &dirHead, nil
}

func readDirection(reader io.ReadSeeker, fileHead fileHeader) (*direction, error) {
	bitReader := streaming.NewBitReader(reader)

	dirHead, err := readDirectionHeader(bitReader)
	if err != nil {
		return nil, err
	}

	frameHeads, dirBounds, err := readFrameHeaders(bitReader, fileHead, *dirHead)
	if err != nil {
		return nil, err
	}

	dirData := direction{bounds: dirBounds}
	for _, frameHead := range frameHeads {
		dirData.frames = append(dirData.frames, newFrame(frameHead, dirBounds))
	}

	var entries []pixelBufferEntry

	entries, err = dirData.decodeStage1(bitReader, entries)
	if err != nil {
		return nil, err
	}

	entries, err = dirData.decodeStage2(bitReader, entries)
	if err != nil {
		return nil, err
	}

	log.Println(dirData.bounds)
	return &dirData, nil
}

type direction struct {
	header directionHeader
	frames []frame
	bounds box
}

func (d *direction) decodeStage1(bitReader *streaming.BitReader, entries []pixelBufferEntry) ([]pixelBufferEntry, error) {
	return nil, nil
}

func (d *direction) decodeStage2(bitReader *streaming.BitReader, entries []pixelBufferEntry) ([]pixelBufferEntry, error) {
	return nil, nil
}
