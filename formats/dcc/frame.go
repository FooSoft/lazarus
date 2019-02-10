package dcc

import (
	"errors"

	"github.com/FooSoft/lazarus/streaming"
)

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

		if frameHead.Width == 0 || frameHead.Height == 0 {
			return nil, errors.New("invalid frame dimensions")
		}

		frameHeads = append(frameHeads, frameHead)
	}

	return frameHeads, nil
}
