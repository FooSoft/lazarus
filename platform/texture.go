package platform

import (
	"image/color"
	"unsafe"

	"github.com/FooSoft/lazarus/math"
	"github.com/go-gl/gl/v2.1/gl"
)

type Texture struct {
	size     math.Vec2i
	glHandle uint32
}

func newTextureFromRgba(colors []color.RGBA, width, height int) (*Texture, error) {
	var glHandleLast int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &glHandleLast)

	var glHandle uint32
	gl.GenTextures(1, &glHandle)
	gl.BindTexture(gl.TEXTURE_2D, glHandle)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&colors[0]))

	gl.BindTexture(gl.TEXTURE_2D, uint32(glHandleLast))
	return &Texture{size: math.Vec2i{X: width, Y: height}, glHandle: glHandle}, nil
}

func (t *Texture) Handle() Handle {
	return Handle(t.glHandle)
}

func (t *Texture) Size() math.Vec2i {
	return t.size
}

func (t *Texture) Destroy() error {
	if t.glHandle != 0 {
		gl.DeleteTextures(1, &t.glHandle)
		t.glHandle = 0
	}

	return nil
}