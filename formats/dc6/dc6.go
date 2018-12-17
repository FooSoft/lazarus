package dc6

import (
	"bytes"
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
	Length    uint32
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

		buff := bytes.NewBuffer(make([]byte, frameHeader.Width*frameHeader.Height))
		// if err := extractFrame(reader, buff, frameHeader); err != nil {
		if err := extractFrame(reader, nil, frameHeader); err != nil {
			return nil, err
		}

		frame := Frame{
			int(frameHeader.Width),
			int(frameHeader.Height),
			int(frameHeader.OffsetX),
			int(frameHeader.OffsetY),
			buff.Bytes(),
		}

		sprite.Frames = append(sprite.Frames, frame)
	}

	return sprite, nil
}

func extractFrame(reader io.ReadSeeker, writer io.WriteSeeker, header FrameHeader) error {
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
