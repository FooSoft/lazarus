package platform

import "github.com/veandco/go-sdl2/sdl"

type Window interface {
	Destroy() error
}

type window struct {
	sdlWindow    *sdl.Window
	sdlGlContext sdl.GLContext
}

func newWindow(title string, width, height int) (Window, error) {
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

	return &window{sdlWindow, sdlGlContext}, nil
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
