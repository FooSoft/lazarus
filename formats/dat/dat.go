package dat

import (
	"encoding/binary"
	"io"

	"github.com/FooSoft/lazarus/math"
)

type Palette struct {
	Colors [256]math.Color3b
}

func NewFromReader(r io.Reader) (*Palette, error) {
	p := new(Palette)
	if err := binary.Read(r, binary.LittleEndian, p); err != nil {
		return nil, err
	}

	return p, nil
}
