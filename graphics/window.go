package graphics

import "github.com/veandco/go-sdl2/sdl"

type Window struct {
	window   *sdl.Window
	renderer *sdl.Renderer
}

func NewWindow(title string, width, height int) (*Window, error) {
	window, err := sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(width),
		int32(height),
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return nil, err
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		window.Destroy()
		return nil, err
	}

	return &Window{window, renderer}, nil
}

func (w *Window) Destroy() {

}
