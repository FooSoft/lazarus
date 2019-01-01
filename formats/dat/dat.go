package dat

import (
	"encoding/binary"
	"io"

	"github.com/FooSoft/lazarus/math"
)

type DatPalette struct {
	Colors [256]math.Color3b
}

type color struct {
	B byte
	G byte
	R byte
}

func NewFromReader(reader io.Reader) (*DatPalette, error) {
	var colors [256]color
	if err := binary.Read(reader, binary.LittleEndian, &colors); err != nil {
		return nil, err
	}

	palette := new(DatPalette)
	for i, color := range colors {
		palette.Colors[i] = math.Color3b{R: color.R, G: color.G, B: color.B}
	}

	return palette, nil
}

func NewFromGrayscale() *DatPalette {
	palette := new(DatPalette)
	for i := 0; i < 256; i++ {
		value := uint8(i)
		palette.Colors[i] = math.Color3b{R: value, G: value, B: value}
	}

	return palette
}
