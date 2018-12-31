package platform

import (
	"errors"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	ErrAlreadyInit = errors.New("platform is already initialized")
	ErrWasNotInit  = errors.New("platform was not initialized")
)

var singleton struct {
	isInit  bool
	windows []*Window
}

type Handle uintptr

type Scene interface {
	Init(window *Window) error
	Advance(window *Window) error
	Shutdown(window *Window) error
}

func Init() error {
	if singleton.isInit {
		return ErrAlreadyInit
	}

	runtime.LockOSThread()

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return err
	}

	if err := gl.Init(); err != nil {
		return err
	}

	singleton.isInit = true
	return nil
}

func Advance() (bool, error) {
	if !singleton.isInit {
		return false, ErrWasNotInit
	}

	advanceWindows()

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			return false, nil
		default:
			if err := processWindowEvents(event); err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

func Shutdown() error {
	if !singleton.isInit {
		return ErrWasNotInit
	}

	for _, w := range singleton.windows {
		if err := w.Destroy(); err != nil {
			return err
		}
	}

	singleton.windows = nil
	singleton.isInit = false

	return nil
}

func CreateWindow(title string, width, height int, scene Scene) (*Window, error) {
	if !singleton.isInit {
		return nil, ErrWasNotInit
	}

	w, err := newWindow(title, width, height, scene)
	if err != nil {
		return nil, err
	}

	appendWindow(w)
	return w, err
}

func appendWindow(window *Window) {
	singleton.windows = append(singleton.windows, window)
}

func removeWindow(window *Window) bool {
	for i, w := range singleton.windows {
		if w == window {
			singleton.windows = append(singleton.windows[:i], singleton.windows[i+1:]...)
			return true
		}
	}

	return false
}

func advanceWindows() {
	for _, window := range singleton.windows {
		window.advance()
	}
}

func processWindowEvents(event sdl.Event) error {
	for _, window := range singleton.windows {
		if _, err := window.processEvent(event); err != nil {
			return err
		}
	}

	return nil
}
