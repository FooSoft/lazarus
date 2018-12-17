package dc6

import (
	"encoding/binary"
	"io"
)

const (
	FlagIsSerialized = 1 << iota
	FlagIsLoadedInHw
	FlagIs24Bits
)

type FileHeader struct {
	Version      uint32
	Flags        uint32
	Format       uint32
	SkipColor    uint32
	DirCount     uint32
	FramesPerDir uint32
}

type FrameHeader struct {
	Flip      uint32
	Width     uint32
	Height    uint32
	OffsetX   uint32
	OffsetY   uint32
	AllocSize uint32
	NextBlock uint32
	Block     uint32
}

type Frame struct {
	Width   int
	Height  int
	OffsetX int
	OffsetY int
	Data    []byte
}

type Dc6 struct {
	Frames []Frame
}

func New(reader io.ReadSeeker) (*Dc6, error) {
	sprite := new(Dc6)

	var fileHeader FileHeader
	if err := binary.Read(reader, binary.LittleEndian, &fileHeader); err != nil {
		return nil, err
	}

	var frameOffsets []uint32
	for i := 0; i < int(fileHeader.DirCount*fileHeader.FramesPerDir); i++ {
		var frameOffset uint32
		if err := binary.Read(reader, binary.LittleEndian, &frameOffset); err != nil {
			return nil, err
		}

		frameOffsets = append(frameOffsets, frameOffset)
	}

	for _, frameOffset := range frameOffsets {
		var frameHeader FrameHeader
		if _, err := reader.Seek(int64(frameOffset), io.SeekStart); err != nil {
			return nil, err
		}

		if err := binary.Read(reader, binary.LittleEndian, &frameHeader); err != nil {
			return nil, err
		}

		frame := Frame{
			int(frameHeader.Width),
			int(frameHeader.Height),
			int(frameHeader.OffsetX),
			int(frameHeader.OffsetY),
			make([]byte, frameHeader.Width*frameHeader.Height),
		}

		if _, err := io.ReadFull(reader, frame.Data); err != nil {
			return nil, err
		}

		if err := extractFrame(reader, frameHeader); err != nil {
			return nil, err
		}

		sprite.Frames = append(sprite.Frames, frame)
	}

	return sprite, nil
}

func extractFrame(reader io.ReadSeeker, header FrameHeader) error {
	return nil
}
