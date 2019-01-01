package imgui

import (
	"log"

	imgui "github.com/FooSoft/imgui-go"
	"github.com/FooSoft/lazarus/math"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

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
		singleton.context = imgui.CreateContext(nil)

		keys := map[int]int{
			imgui.KeyTab:        sdl.SCANCODE_TAB,
			imgui.KeyLeftArrow:  sdl.SCANCODE_LEFT,
			imgui.KeyRightArrow: sdl.SCANCODE_RIGHT,
			imgui.KeyUpArrow:    sdl.SCANCODE_UP,
			imgui.KeyDownArrow:  sdl.SCANCODE_DOWN,
			imgui.KeyPageUp:     sdl.SCANCODE_PAGEUP,
			imgui.KeyPageDown:   sdl.SCANCODE_PAGEDOWN,
			imgui.KeyHome:       sdl.SCANCODE_HOME,
			imgui.KeyEnd:        sdl.SCANCODE_END,
			imgui.KeyInsert:     sdl.SCANCODE_INSERT,
			imgui.KeyDelete:     sdl.SCANCODE_DELETE,
			imgui.KeyBackspace:  sdl.SCANCODE_BACKSPACE,
			imgui.KeySpace:      sdl.SCANCODE_BACKSPACE,
			imgui.KeyEnter:      sdl.SCANCODE_RETURN,
			imgui.KeyEscape:     sdl.SCANCODE_ESCAPE,
			imgui.KeyA:          sdl.SCANCODE_A,
			imgui.KeyC:          sdl.SCANCODE_C,
			imgui.KeyV:          sdl.SCANCODE_V,
			imgui.KeyX:          sdl.SCANCODE_X,
			imgui.KeyY:          sdl.SCANCODE_Y,
			imgui.KeyZ:          sdl.SCANCODE_Z,
		}

		// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.
		io := imgui.CurrentIO()
		for imguiKey, nativeKey := range keys {
			io.KeyMap(imguiKey, nativeKey)
		}
	}

	log.Println("imgui context create")
	c := &Context{displaySize: displaySize, bufferSize: bufferSize}

	// Build texture atlas
	io := imgui.CurrentIO()
	image := io.Fonts().TextureDataRGBA32()

	// Store state
	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)

	// Create texture
	gl.GenTextures(1, &c.fontTexture)
	gl.BindTexture(gl.TEXTURE_2D, c.fontTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(image.Width),
		int32(image.Height),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		image.Pixels,
	)

	// Restore state
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))

	// Store texture identifier
	io.Fonts().SetTextureID(imgui.TextureID(c.fontTexture))

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
	imgui.CurrentIO().Fonts().SetTextureID(0)
	c.fontTexture = 0

	singleton.refCount--
	if singleton.refCount == 0 {
		log.Println("imgui global destroy")
		singleton.context.Destroy()
		singleton.context = nil
	}

	return nil
}
