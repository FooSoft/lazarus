package platform

import (
	"image/color"

	imgui "github.com/FooSoft/imgui-go"
	"github.com/FooSoft/lazarus/math"
	"github.com/FooSoft/lazarus/platform/imgui_backend"
	"github.com/veandco/go-sdl2/sdl"
)

type Window struct {
	sdlWindow    *sdl.Window
	sdlGlContext sdl.GLContext
	sdlRenderer  *sdl.Renderer
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

	sdlRenderer, err := sdl.CreateRenderer(sdlWindow, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		sdlWindow.Destroy()
		return nil, err
	}

	sdlGlContext, err := sdlWindow.GLCreateContext()
	if err != nil {
		sdlWindow.Destroy()
		return nil, err
	}

	window := &Window{sdlWindow, sdlGlContext, sdlRenderer, scene}
	if err := scene.Init(window); err != nil {
		return nil, err
	}

	return window, nil
}

func (w *Window) Destroy() error {
	if w.sdlWindow == nil {
		return nil
	}

	if err := w.scene.Shutdown(w); err != nil {
		return err
	}

	if err := w.sdlWindow.Destroy(); err != nil {
		return err
	}

	w.sdlGlContext = nil
	w.sdlWindow = nil

	return nil
}

func (w *Window) CreateTextureRgba(colors []color.RGBA, width, height int) (*Texture, error) {
	return newTextureFromRgba(w.sdlRenderer, colors, width, height)
}

func (w *Window) RenderTexture(texture *Texture, srcRect, dstRect math.Rect4i) {
	w.sdlRenderer.Copy(
		texture.sdlTexture,
		&sdl.Rect{X: int32(srcRect.X), Y: int32(srcRect.Y), W: int32(srcRect.W), H: int32(srcRect.H)},
		&sdl.Rect{X: int32(dstRect.X), Y: int32(dstRect.Y), W: int32(dstRect.W), H: int32(dstRect.H)},
	)
}

func (w *Window) advance() {
	imgui_backend.NewFrame(w.displaySize())
	w.scene.Advance(w)

	imgui.Render()

	w.sdlWindow.GLMakeCurrent(w.sdlGlContext)
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
