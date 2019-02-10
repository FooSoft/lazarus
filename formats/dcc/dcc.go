package dcc

import (
	"encoding/binary"
	"io"
)

type Sprite struct {
}

type box struct {
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

type pixelBufferEntry struct {
	values [4]byte
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
