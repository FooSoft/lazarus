package dat

import (
	"encoding/binary"
	imageColor "image/color"
	"io"
)

type Palette struct {
	Colors [256]imageColor.NRGBA
}

type color struct {
	B byte
	G byte
	R byte
}

func NewFromReader(reader io.Reader) (*Palette, error) {
	var colors [256]color
	if err := binary.Read(reader, binary.LittleEndian, &colors); err != nil {
		return nil, err
	}

	palette := new(Palette)
	for i, color := range colors {
		palette.Colors[i] = imageColor.NRGBA{color.R, color.G, color.B, 0xff}
	}

	return palette, nil
}
