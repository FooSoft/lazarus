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

	nbCellsX   int
	nbCellsY   int
	dirOffsetX int
	dirOffsetY int

	cellSameAsPrevious []bool
	cellWidths         []int
	cellHeights        []int

	data []byte
}

func newFrame(frameHead frameHeader, dirBounds box) frame {
	frameBounds := frameHead.bounds()

	frameData := frame{
		header:     frameHead,
		dirOffsetX: frameBounds.x1 - dirBounds.x1,
		dirOffsetY: frameBounds.y1 - dirBounds.y1,
	}

	widthFirstColumn := 4 - frameData.dirOffsetX%4
	frameWidth := frameBounds.x2 - frameBounds.x1
	if frameWidth-widthFirstColumn <= 1 {
		frameData.nbCellsX = 1
	} else {
		temp := frameWidth - widthFirstColumn - 1
		frameData.nbCellsX = 2 + temp/4
		if temp%4 == 0 {
			frameData.nbCellsX--
		}
	}

	heightFirstRow := 4 - frameData.dirOffsetY%4
	frameHeight := frameBounds.y2 - frameBounds.y1
	if frameHeight-heightFirstRow <= 1 {
		frameData.nbCellsY = 1
	} else {
		temp := frameHeight - heightFirstRow - 1
		frameData.nbCellsY = 2 + temp/4
		if temp%4 == 0 {
			frameData.nbCellsY--
		}
	}

	frameData.cellWidths = make([]int, frameData.nbCellsX)
	for i := range frameData.cellWidths {
		frameData.cellWidths[i] = 4
	}

	if frameData.nbCellsX == 1 {
		frameData.cellWidths[0] = frameWidth
	} else {
		frameData.cellWidths[0] = widthFirstColumn
		nbColumnsExcludingFirstAndLast := frameData.nbCellsX - 2
		widthExcludingFirstAndLastColumns := 4 * nbColumnsExcludingFirstAndLast
		frameData.cellWidths[frameData.nbCellsX-1] = frameWidth - (widthFirstColumn + widthExcludingFirstAndLastColumns)
	}

	frameData.cellHeights = make([]int, frameData.nbCellsY)
	for i := range frameData.cellHeights {
		frameData.cellHeights[i] = 4
	}

	if frameData.nbCellsY == 1 {
		frameData.cellHeights[0] = frameHeight
	} else {
		frameData.cellHeights[0] = heightFirstRow
		nbRowsExcludingFirstAndLast := frameData.nbCellsY - 2
		heightExcludingFirstAndLastRows := 4 * nbRowsExcludingFirstAndLast
		frameData.cellHeights[frameData.nbCellsY-1] = frameHeight - (heightFirstRow + heightExcludingFirstAndLastRows)
	}

	frameData.cellSameAsPrevious = make([]bool, frameData.nbCellsX*frameData.nbCellsY)
	frameData.data = make([]byte, frameWidth*frameHeight)

	return frameData

}
