package platform

import (
	"github.com/FooSoft/lazarus/math"
	"github.com/veandco/go-sdl2/sdl"
)

type Window interface {
	Destroy() error
}

type window struct {
	sdlWindow    *sdl.Window
	sdlGlContext sdl.GLContext
	sdlRenderer  *sdl.Renderer
	scene        Scene
}

func newWindow(title string, width, height int, scene Scene) (*window, error) {
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

	sdlRenderer, err := sdl.CreateRenderer(sdlWindow, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		sdlWindow.Destroy()
		return nil, err
	}

	return &window{sdlWindow, sdlGlContext, sdlRenderer, scene}, nil
}

func (w *window) Destroy() error {
	if w.sdlWindow != nil {
		if err := w.sdlWindow.Destroy(); err != nil {
			return err
		}
	}

	w.sdlGlContext = nil
	w.sdlWindow = nil

	return nil
}

func (w *window) advance() {
	w.scene.Advance()
}

func (w *window) render() {
	w.sdlWindow.GLMakeCurrent(w.sdlGlContext)
}

func (w *window) displaySize() math.Vec2i {
	width, height := w.sdlWindow.GetSize()
	return math.Vec2i{X: int(width), Y: int(height)}
}

func (w *window) bufferSize() math.Vec2i {
	width, height := w.sdlWindow.GLGetDrawableSize()
	return math.Vec2i{X: int(width), Y: int(height)}
}
