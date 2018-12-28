package platform

import (
	"errors"
	"runtime"
	"time"

	"github.com/FooSoft/lazarus/platform/imgui_backend"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	platformIsInit  bool
	platformWindows []Window
)

var (
	ErrAlreadyInit = errors.New("platform is already initialized")
	ErrWasNotInit  = errors.New("platform was not initialized")
)

func Init() error {
	if platformIsInit {
		return ErrAlreadyInit
	}

	runtime.LockOSThread()

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return err
	}

	if err := gl.Init(); err != nil {
		return err
	}

	if err := imgui_backend.Init(); err != nil {
		return err
	}

	platformIsInit = true
	return nil
}

func Shutdown() error {
	if !platformIsInit {
		return ErrWasNotInit
	}

	for _, w := range platformWindows {
		if err := w.Destroy(); err != nil {
			return err
		}
	}

	if err := imgui_backend.Shutdown(); err != nil {
		return err
	}

	platformWindows = nil
	platformIsInit = false

	return nil
}

func ProcessEvents() error {
	if !platformIsInit {
		return ErrWasNotInit
	}

	var terminate bool
	for !terminate {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				terminate = true
				break
			}
		}

		<-time.After(time.Millisecond * 25)
	}

	return nil
}

func CreateWindow(title string, width, height int) (Window, error) {
	if !platformIsInit {
		return nil, ErrWasNotInit
	}

	window, err := newWindow(title, width, height)
	if err != nil {
		return nil, err
	}

	platformWindows = append(platformWindows, window)

	return window, err
}
