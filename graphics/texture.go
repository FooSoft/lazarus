package graphics

import (
	"image/color"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

func NewSurfaceFromRgba(colors []color.RGBA, width, height int) (*sdl.Surface, error) {
	return sdl.CreateRGBSurfaceFrom(
		unsafe.Pointer(&colors[0]),
		int32(width),
		int32(height),
		32,
		width*4,
		0x000000ff,
		0x0000ff00,
		0x00ff0000,
		0xff000000,
	)
}
