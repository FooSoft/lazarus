package platform

import (
	"image/color"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type Texture struct {
	sdlTexture *sdl.Texture
}

func newTextureFromRgba(renderer *sdl.Renderer, colors []color.RGBA, width, height int) (*Texture, error) {
	surface, err := sdl.CreateRGBSurfaceFrom(
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

	if err != nil {
		return nil, nil
	}

	sdlTexture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		surface.Free()
	}

	return &Texture{sdlTexture}, nil
}

func (t *Texture) Destroy() error {
	if t == nil {
		return nil
	}

	if err := t.sdlTexture.Destroy(); err != nil {
		return err
	}

	t.sdlTexture = nil
	return nil
}
