package platform

import (
	"github.com/FooSoft/lazarus/math"
	"github.com/FooSoft/lazarus/platform/imgui"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

type Window struct {
	sdlWindow    *sdl.Window
	sdlGlContext sdl.GLContext
	imguiContext *imgui_backend.Context
	scene        Scene
}

func newWindow(title string, size math.Vec2i, scene Scene) (*Window, error) {
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 2)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)

	sdlWindow, err := sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		int32(size.X),
		int32(size.Y),
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

	w := &Window{
		sdlWindow:    sdlWindow,
		sdlGlContext: sdlGlContext,
	}

	w.makeCurrent()

	w.imguiContext, err = imgui_backend.New(w.DisplaySize(), w.BufferSize())
	if err != nil {
		w.Destroy()
		return nil, err
	}

	if err := w.SetScene(scene); err != nil {
		w.Destroy()
		return nil, err
	}

	return w, nil
}

func (w *Window) SetScene(scene Scene) error {
	if w.scene == scene {
		return nil
	}

	if sceneDestroyer, ok := w.scene.(SceneDestroyer); ok {
		if err := sceneDestroyer.Destroy(w); err != nil {
			return err
		}
	}

	w.scene = scene

	if sceneCreator, ok := scene.(SceneCreator); ok {
		if err := sceneCreator.Create(w); err != nil {
			return err
		}
	}

	return nil
}

func (w *Window) Destroy() error {
	if w == nil || w.sdlWindow == nil {
		return nil
	}

	w.makeCurrent()

	if err := w.SetScene(nil); err != nil {
		return err
	}

	if err := w.imguiContext.Destroy(); err != nil {
		return err
	}
	w.imguiContext = nil

	sdl.GLDeleteContext(w.sdlGlContext)
	w.sdlGlContext = nil

	if err := w.sdlWindow.Destroy(); err != nil {
		return err
	}
	w.sdlWindow = nil

	removeWindow(w)
	return nil
}

func (w *Window) CreateTextureRgba(colors []math.Color4b, size math.Vec2i) (*Texture, error) {
	return newTextureFromRgba(colors, size)
}

func (w *Window) CreateTextureRgb(colors []math.Color3b, size math.Vec2i) (*Texture, error) {
	return newTextureFromRgb(colors, size)
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

func (w *Window) DisplaySize() math.Vec2i {
	width, height := w.sdlWindow.GetSize()
	return math.Vec2i{X: int(width), Y: int(height)}
}

func (w *Window) BufferSize() math.Vec2i {
	width, height := w.sdlWindow.GLGetDrawableSize()
	return math.Vec2i{X: int(width), Y: int(height)}
}

func (w *Window) advance() (bool, error) {
	w.makeCurrent()

	displaySize := w.DisplaySize()
	w.imguiContext.SetDisplaySize(displaySize)
	bufferSize := w.BufferSize()
	w.imguiContext.SetBufferSize(bufferSize)

	w.imguiContext.BeginFrame()

	gl.Viewport(0, 0, int32(displaySize.X), int32(displaySize.Y))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(displaySize.X), float64(displaySize.Y), 0, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	if sceneAdvancer, ok := w.scene.(SceneAdvancer); ok {
		if err := sceneAdvancer.Advance(w); err != nil {
			return false, err
		}
	}

	w.imguiContext.EndFrame()
	w.sdlWindow.GLSwap()

	return w.scene != nil, nil
}

func (w *Window) processEvent(event sdl.Event) (bool, error) {
	return w.imguiContext.ProcessEvent(event)
}

func (w *Window) makeCurrent() {
	w.sdlWindow.GLMakeCurrent(w.sdlGlContext)

}
