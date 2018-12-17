package dc6

import (
	"encoding/binary"
	"io"

	"github.com/FooSoft/lazarus/streaming"
)

type fileHeader struct {
	Version         uint32
	UnusedFlags     uint32
	UnusedFormat    uint32
	UnusedSkipColor uint32
	DirCount        uint32
	FramesPerDir    uint32
}

type frameHeader struct {
	Flip            uint32
	Width           uint32
	Height          uint32
	OffsetX         uint32
	OffsetY         uint32
	UnusedAllocSize uint32
	UnusedNextBlock uint32
	Length          uint32
}

type Direction struct {
	Frames []Frame
}

type Frame struct {
	Width   int
	Height  int
	OffsetX int
	OffsetY int
	Data    []byte
}

type Sprite struct {
	Directions []Direction
}

func New(reader io.ReadSeeker) (*Sprite, error) {
	sprite := new(Sprite)

	var fileHead fileHeader
	if err := binary.Read(reader, binary.LittleEndian, &fileHead); err != nil {
		return nil, err
	}

	frameCount := int(fileHead.DirCount * fileHead.FramesPerDir)

	var frameOffsets []uint32
	for i := 0; i < frameCount; i++ {
		var frameOffset uint32
		if err := binary.Read(reader, binary.LittleEndian, &frameOffset); err != nil {
			return nil, err
		}

		frameOffsets = append(frameOffsets, frameOffset)
	}

	sprite.Directions = make([]Direction, fileHead.FramesPerDir)

	for i := 0; i < frameCount; i++ {
		var frameHead frameHeader
		if _, err := reader.Seek(int64(frameOffsets[i]), io.SeekStart); err != nil {
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

		frame := Frame{
			int(frameHead.Width),
			int(frameHead.Height),
			int(frameHead.OffsetX),
			int(frameHead.OffsetY),
			data,
		}

		direction := &sprite.Directions[i/int(fileHead.FramesPerDir)]
		direction.Frames = append(direction.Frames, frame)
	}

	return sprite, nil
}

func extractFrame(reader io.ReadSeeker, writer io.WriteSeeker, header frameHeader) error {
	var (
		x uint32
		y = header.Height - 1
	)

	var offset uint32
	for offset < header.Length {
		var chunkSize byte
		if err := binary.Read(reader, binary.LittleEndian, &chunkSize); err != nil {
			return err
		}

		if chunkSize == 0x80 {
			x = 0
			y--
		} else if (chunkSize & 0x80) != 0 {
			x += uint32(chunkSize & 0x7f)
		} else {
			if _, err := writer.Seek(int64(header.Width*y+x), io.SeekStart); err != nil {
				return err
			}
			if _, err := io.CopyN(writer, reader, int64(chunkSize)); err != nil {
				return err
			}

			offset += uint32(chunkSize)
			x += uint32(chunkSize)
		}
	}

	return nil
}
