package dc6

import (
	"encoding/binary"
	"io"

	"github.com/FooSoft/lazarus/math"
	"github.com/FooSoft/lazarus/streaming"
)

type fileHeader struct {
	Version      uint32
	_            uint32 // unused: flags
	_            uint32 // unused: format
	_            uint32 // unused: skipColor
	DirCount     uint32
	FramesPerDir uint32
}

type frameHeader struct {
	Flip    uint32
	Width   uint32
	Height  uint32
	OffsetX int32
	OffsetY int32
	_       uint32 // unused: allocSize
	_       uint32 // unused: nextBlock
	Length  uint32
}

type Direction struct {
	Frames []Frame
}

type Frame struct {
	Size   math.Vec2i
	Offset math.Vec2i
	Data   []byte
}

type Sprite struct {
	Directions []Direction
}

func NewFromReader(reader io.ReadSeeker) (*Sprite, error) {
	sprite := new(Sprite)

	var fileHead fileHeader
	if err := binary.Read(reader, binary.LittleEndian, &fileHead); err != nil {
		return nil, err
	}

	var frameOffsets []uint32
	for i := uint32(0); i < fileHead.DirCount*fileHead.FramesPerDir; i++ {
		var frameOffset uint32
		if err := binary.Read(reader, binary.LittleEndian, &frameOffset); err != nil {
			return nil, err
		}

		frameOffsets = append(frameOffsets, frameOffset)
	}

	sprite.Directions = make([]Direction, fileHead.DirCount)

	for i, frameOffset := range frameOffsets {
		var frameHead frameHeader
		if _, err := reader.Seek(int64(frameOffset), io.SeekStart); err != nil {
			return nil, err
		}

		if err := binary.Read(reader, binary.LittleEndian, &frameHead); err != nil {
			return nil, err
		}

		data := make([]byte, frameHead.Width*frameHead.Height)
		writer := streaming.NewWriter(data)
		if err := extractFrame(reader, writer, frameHead); err != nil {
			return nil, err
		}

		var (
			size      = math.Vec2i{X: int(frameHead.Width), Y: int(frameHead.Height)}
			offset    = math.Vec2i{X: int(frameHead.OffsetX), Y: int(frameHead.OffsetY)}
			frame     = Frame{size, offset, data}
			direction = &sprite.Directions[i/int(fileHead.FramesPerDir)]
		)

		direction.Frames = append(direction.Frames, frame)
	}

	return sprite, nil
}

func extractFrame(reader io.ReadSeeker, writer io.WriteSeeker, header frameHeader) error {
	var x, y uint32
	for readOffset := uint32(0); readOffset < header.Length; readOffset++ {
		var chunk byte
		if err := binary.Read(reader, binary.LittleEndian, &chunk); err != nil {
			return err
		}

		if chunk&0x80 > 0 {
			if skipLength := uint32(chunk & 0x7f); skipLength > 0 {
				x += skipLength
			} else {
				x = 0
				y++
			}

		} else {
			writeOffset := int64(header.Width*(header.Height-y-1) + x)
			if header.Flip != 0 {
				writeOffset = int64(header.Width*y + x)
			}

			if _, err := writer.Seek(writeOffset, io.SeekStart); err != nil {
				return err
			}
			if _, err := io.CopyN(writer, reader, int64(chunk)); err != nil {
				return err
			}

			readOffset += uint32(chunk)
			x += uint32(chunk)
		}
	}

	return nil
}
