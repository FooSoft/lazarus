package platform

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

type Platform interface {
	CreateWindow(title string, width, height int) (Window, error)
	Destroy() error
}

var globalPlatformInit bool

type platform struct {
	windows []Window
}

func New() (*Platform, error) {
	if !globalPlatformInit {
		if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
			return nil, err
		}

		if err := gl.Init(); err != nil {
			return nil, err
		}

		globalPlatformInit = true
	}

	return nil, nil
}

func (p *platform) CreateWindow(title string, width, height int) (Window, error) {
	window, err := newWindow(title, width, height)
	if err != nil {
		return nil, err
	}

	return window, err
}

func (p *platform) Destroy() error {
	for _, w := range p.windows {
		if err := w.Destroy(); err != nil {
			return err
		}
	}

	p.windows = nil
	return nil
}
