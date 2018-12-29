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
	ErrAlreadyInit = errors.New("platform is already initialized")
	ErrWasNotInit  = errors.New("platform was not initialized")
)

var state struct {
	isInit  bool
	windows []*Window
}

type Scene interface {
	Init(window *Window) error
	Advance(window *Window) error
	Shutdown(window *Window) error
}

func Init() error {
	if state.isInit {
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

	state.isInit = true
	return nil
}

func Shutdown() error {
	if !state.isInit {
		return ErrWasNotInit
	}

	for _, w := range state.windows {
		if err := w.Destroy(); err != nil {
			return err
		}
	}

	if err := imgui_backend.Shutdown(); err != nil {
		return err
	}

	state.windows = nil
	state.isInit = false

	return nil
}

func CreateWindow(title string, width, height int, scene Scene) (*Window, error) {
	if !state.isInit {
		return nil, ErrWasNotInit
	}

	window, err := newWindow(title, width, height, scene)
	if err != nil {
		return nil, err
	}

	state.windows = append(state.windows, window)

	return window, err
}

func ProcessEvents() error {
	if !state.isInit {
		return ErrWasNotInit
	}

	for running := true; running; {
		advanceWindows()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			if !processEvent(event) {
				running = false
				break
			}
		}

		<-time.After(time.Millisecond * 25)
	}

	return nil
}

func advanceWindows() {
	for _, window := range state.windows {
		window.advance()
	}
}

func processEvent(event sdl.Event) bool {
	handled, _ := imgui_backend.ProcessEvent(event)
	if handled {
		return true
	}

	switch event.(type) {
	case *sdl.QuitEvent:
		return false
	}

	return true
}
