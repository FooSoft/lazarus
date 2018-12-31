package platform

import (
	"image/color"

	imgui "github.com/FooSoft/imgui-go"
	"github.com/FooSoft/lazarus/math"
	"github.com/FooSoft/lazarus/platform/imgui_backend"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

type Window struct {
	sdlWindow    *sdl.Window
	sdlGlContext sdl.GLContext
	scene        Scene
}

func newWindow(title string, width, height int, scene Scene) (*Window, error) {
	sdlWindow, err := sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		int32(width),
		int32(height),
		sdl.WINDOW_OPENGL,
	)
	if err != nil {
		return nil, err
	}

	sdlGlContext, err := sdlWindow.GLCreateContext()
	if err != nil {
		sdlWindow.Destroy()
		return nil, err
	}

	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 2)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)

	w := &Window{sdlWindow, sdlGlContext, scene}
	if err := scene.Init(w); err != nil {
		w.Destroy()
		return nil, err
	}

	return w, nil
}

func (w *Window) Destroy() error {
	if w.sdlWindow == nil {
		return nil
	}

	if err := w.scene.Shutdown(w); err != nil {
		return err
	}
	w.scene = nil

	sdl.GLDeleteContext(w.sdlGlContext)
	w.sdlGlContext = nil

	if err := w.sdlWindow.Destroy(); err != nil {
		return err
	}
	w.sdlWindow = nil

	return nil
}

func (w *Window) CreateTextureRgba(colors []color.RGBA, width, height int) (*Texture, error) {
	return newTextureFromRgba(colors, width, height)
}

func (w *Window) RenderTexture(texture *Texture, position math.Vec2i) {
	size := texture.Size()

	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, uint32(texture.Handle()))

	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(0, 0)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(0, float32(size.Y))
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(float32(size.X), float32(size.Y))
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(float32(size.X), 0)
	gl.End()
}

func (w *Window) advance() {
	size := w.displaySize()
	imgui_backend.NewFrame(size)

	gl.Viewport(0, 0, int32(size.X), int32(size.Y))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(size.X), float64(size.Y), 0, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	w.scene.Advance(w)

	imgui.Render()
	imgui_backend.Render(w.displaySize(), w.bufferSize(), imgui.RenderedDrawData())

	w.sdlWindow.GLSwap()
}

func (w *Window) displaySize() math.Vec2i {
	width, height := w.sdlWindow.GetSize()
	return math.Vec2i{X: int(width), Y: int(height)}
}

func (w *Window) bufferSize() math.Vec2i {
	width, height := w.sdlWindow.GLGetDrawableSize()
	return math.Vec2i{X: int(width), Y: int(height)}
}
