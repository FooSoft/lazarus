package platform

import (
	"errors"
	"runtime"
	"time"

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

	window, err := newWindow(title, width, height, scene)
	if err != nil {
		return nil, err
	}

	singleton.windows = append(singleton.windows, window)

	return window, err
}

func ProcessEvents() error {
	if !singleton.isInit {
		return ErrWasNotInit
	}

	for running := true; running; {
		advanceWindows()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			default:
				processWindowEvents(event)
			}
		}

		<-time.After(time.Millisecond * 25)
	}

	return nil
}

func advanceWindows() {
	for _, window := range singleton.windows {
		window.advance()
	}
}

func processWindowEvents(event sdl.Event) {
	for _, window := range singleton.windows {
		window.processEvent(event)
	}
}
