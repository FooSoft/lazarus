package platform

import (
	"unsafe"

	"github.com/FooSoft/lazarus/graphics"
	"github.com/FooSoft/lazarus/math"
	"github.com/go-gl/gl/v2.1/gl"
)

type texture struct {
	glTexture uint32
	size      math.Vec2i
}

func NewTextureFromRgba(colors []math.Color4b, size math.Vec2i) (graphics.Texture, error) {
	if !WindowIsCreated() {
		return nil, ErrWindowNotExists
	}

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
	return &texture{glTexture, size}, nil
}

func NewTextureFromRgb(colors []math.Color3b, size math.Vec2i) (graphics.Texture, error) {
	if !WindowIsCreated() {
		return nil, ErrWindowNotExists
	}

	var glLastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &glLastTexture)

	var glTexture uint32
	gl.GenTextures(1, &glTexture)
	gl.BindTexture(gl.TEXTURE_2D, glTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(size.X), int32(size.Y), 0, gl.RGB, gl.UNSIGNED_BYTE, unsafe.Pointer(&colors[0]))

	gl.BindTexture(gl.TEXTURE_2D, uint32(glLastTexture))
	return &texture{glTexture, size}, nil
}

func (t *texture) Id() graphics.TextureId {
	return graphics.TextureId(t.glTexture)
}

func (t *texture) Size() math.Vec2i {
	return t.size
}

func (t *texture) Destroy() error {
	if t == nil || t.glTexture == 0 {
		return nil
	}

	gl.DeleteTextures(1, &t.glTexture)
	t.glTexture = 0

	return nil
}
