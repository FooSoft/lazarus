package platform

import (
	"errors"

	"github.com/FooSoft/imgui-go"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	platformIsInit       bool
	platformImguiContext *imgui.Context
	platformWindows      []Window
)

func Init() error {
	if platformIsInit {
		return errors.New("platform is already initialized")
	}

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return err
	}

	if err := gl.Init(); err != nil {
		return err
	}

	platformImguiContext = imgui.CreateContext(nil)
	return nil
}

func Shutdown() error {
	if !platformIsInit {
		return errors.New("platform is not yet initialized")
	}

	for _, w := range platformWindows {
		if err := w.Destroy(); err != nil {
			return err
		}
	}

	platformWindows = nil
	return nil
}

func CreateWindow(title string, width, height int) (Window, error) {
	window, err := newWindow(title, width, height)
	if err != nil {
		return nil, err
	}

	return window, err
}
