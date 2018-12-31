package platform

import (
	"image/color"
	"unsafe"

	"github.com/FooSoft/lazarus/math"
	"github.com/go-gl/gl/v2.1/gl"
)

type Texture struct {
	size      math.Vec2i
	glTexture uint32
}

func newTextureFromRgba(colors []color.RGBA, size math.Vec2i) (*Texture, error) {
	var glLastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &glLastTexture)

	var glTexture uint32
	gl.GenTextures(1, &glTexture)
	gl.BindTexture(gl.TEXTURE_2D, glTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(size.X), int32(size.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&colors[0]))

	gl.BindTexture(gl.TEXTURE_2D, uint32(glLastTexture))
	return &Texture{size, glTexture}, nil
}

func (t *Texture) Handle() Handle {
	return Handle(t.glTexture)
}

func (t *Texture) Size() math.Vec2i {
	return t.size
}

func (t *Texture) Destroy() error {
	if t == nil || t.glTexture == 0 {
		return nil
	}

	gl.DeleteTextures(1, &t.glTexture)
	t.glTexture = 0

	return nil
}
