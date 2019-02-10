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

func readFrameHeader(bitReader *streaming.BitReader, dirHead directionHeader) (*frameHeader, error) {
	var frameHead frameHeader

	frameHead.Variable0 = uint32(bitReader.ReadUintPacked(int(dirHead.Variable0Bits)))
	frameHead.Width = uint32(bitReader.ReadUintPacked(int(dirHead.WidthBits)))
	frameHead.Height = uint32(bitReader.ReadUintPacked(int(dirHead.HeightBits)))
	frameHead.OffsetX = int32(bitReader.ReadIntPacked(int(dirHead.OffsetXBits)))
	frameHead.OffsetY = int32(bitReader.ReadIntPacked(int(dirHead.OffsetYBits)))
	frameHead.OptionalBytes = uint32(bitReader.ReadUintPacked(int(dirHead.OptionalBytesBits)))
	frameHead.CodedBytes = uint32(bitReader.ReadUintPacked(int(dirHead.CodedBytesBits)))
	frameHead.FrameBottomUp = bitReader.ReadBool()

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

	return &frameHead, nil
}

func (h *frameHeader) bounds() box {
	return box{
		x1: int(h.OffsetX),
		y1: int(h.OffsetY) - int(h.Height) + 1,
		x2: int(h.OffsetX) + int(h.Width),
		y2: int(h.OffsetY) + 1,
	}
}

func readFrameHeaders(bitReader *streaming.BitReader, fileHead fileHeader, dirHead directionHeader) ([]frameHeader, box, error) {
	var (
		frameHeads []frameHeader
		boundsAll  box
	)

	for i := 0; i < int(fileHead.FramesPerDir); i++ {
		frameHead, err := readFrameHeader(bitReader, dirHead)
		if err != nil {
			return nil, box{}, err
		}

		bounds := frameHead.bounds()
		if i == 0 {
			boundsAll = bounds
		} else {
			if boundsAll.x1 > bounds.x1 {
				boundsAll.x1 = bounds.x1
			}
			if boundsAll.y1 > bounds.y1 {
				boundsAll.y1 = bounds.y1
			}
			if boundsAll.x2 < bounds.x2 {
				boundsAll.x2 = bounds.x2
			}
			if boundsAll.y2 < bounds.y2 {
				boundsAll.y2 = bounds.y2
			}
		}

		frameHeads = append(frameHeads, *frameHead)
	}

	return frameHeads, boundsAll, nil
}

type frame struct {
	header frameHeader
}

func newFrame(frameHead frameHeader, dirHead directionHeader) frame {
	return frame{frameHead}

}
