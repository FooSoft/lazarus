package platform

import (
	"errors"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	platformIsInit  bool
	platformWindows []Window
)

func Init() error {
	if platformIsInit {
		return errors.New("platform is already initialized")
	}

	runtime.LockOSThread()

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return err
	}

	if err := gl.Init(); err != nil {
		return err
	}

	platformIsInit = true
	return nil
}

func ProcessEvents() error {
	if !platformIsInit {
		return errors.New("platform was not initialized")
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

func Shutdown() error {
	if !platformIsInit {
		return errors.New("platform was not initialized")
	}

	for _, w := range platformWindows {
		if err := w.Destroy(); err != nil {
			return err
		}
	}

	platformWindows = nil
	platformIsInit = false

	return nil
}

func CreateWindow(title string, width, height int) (Window, error) {
	if !platformIsInit {
		return nil, errors.New("platform was not initialized")
	}

	window, err := newWindow(title, width, height)
	if err != nil {
		return nil, err
	}

	platformWindows = append(platformWindows, window)

	return window, err
}
