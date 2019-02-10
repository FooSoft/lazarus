package dcc

import (
	"io"

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

type direction struct {
	header directionHeader
	frames []frame
	bounds bounds
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

	var entries []pixelBufferEntry

	entries, err = decodeDirectionStage1(bitReader, &dirData, entries)
	if err != nil {
		return nil, err
	}

	entries, err = decodeDirectionStage2(bitReader, &dirData, entries)
	if err != nil {
		return nil, err
	}

	return &dirData, nil
}

func decodeDirectionStage1(bitReader *streaming.BitReader, dirData *direction, entries []pixelBufferEntry) ([]pixelBufferEntry, error) {
	return nil, nil
}

func decodeDirectionStage2(bitReader *streaming.BitReader, dirData *direction, entries []pixelBufferEntry) ([]pixelBufferEntry, error) {
	return nil, nil
}
