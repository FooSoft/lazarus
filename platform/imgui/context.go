package imgui

import (
	"log"

	"github.com/FooSoft/lazarus/math"
	"github.com/go-gl/gl/v2.1/gl"
)

var singleton struct {
	context  *imGuiContext
	refCount int
}

type Context struct {
	buttonsDown [3]bool
	lastTime    uint64
	fontTexture uint32
	displaySize math.Vec2i
	bufferSize  math.Vec2i
}

func New(displaySize, bufferSize math.Vec2i) (*Context, error) {
	singleton.refCount++
	if singleton.refCount == 1 {
		log.Println("imgui global create")
		singleton.context = CreateContext()
	}

	log.Println("imgui context create")
	c := &Context{displaySize: displaySize, bufferSize: bufferSize}

	pixels, width, height := FontImage()

	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	var fontTexture uint32
	gl.GenTextures(1, &c.fontTexture)
	gl.BindTexture(gl.TEXTURE_2D, c.fontTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, pixels)
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))

	SetFontTexture(uintptr(fontTexture))
	return c, nil
}

func (c *Context) SetDisplaySize(displaySize math.Vec2i) {
	c.displaySize = displaySize
}

func (c *Context) SetBufferSize(bufferSize math.Vec2i) {
	c.bufferSize = bufferSize
}

func (c *Context) Destroy() error {
	if c == nil || c.fontTexture == 0 {
		return nil
	}

	log.Println("imgui context destroy")
	gl.DeleteTextures(1, &c.fontTexture)
	SetFontTexture(uintptr(0))
	c.fontTexture = 0

	singleton.refCount--
	if singleton.refCount == 0 {
		log.Println("imgui global destroy")
		singleton.context.Destroy()
		singleton.context = nil
	}

	return nil
}
