package platform

import (
	"runtime"

	"github.com/FooSoft/lazarus/math"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/veandco/go-sdl2/sdl"
)

var singleton struct {
	sdlIsInit bool
	windows   []*Window
}

type Handle uintptr

func Advance() (bool, error) {
	if err := advanceWindows(); err != nil {
		return false, err
	}

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

	return len(singleton.windows) > 0, nil
}

func NewWindow(title string, size math.Vec2i, scene Scene) (*Window, error) {
	if !singleton.sdlIsInit {
		runtime.LockOSThread()

		if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
			return nil, err
		}

		if err := gl.Init(); err != nil {
			return nil, err
		}

		singleton.sdlIsInit = true
	}

	w, err := newWindow(title, size, scene)
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

func advanceWindows() error {
	var windowsToRemove []*Window
	for _, window := range singleton.windows {
		run, err := window.advance()
		if err != nil {
			return err
		}
		if !run {
			windowsToRemove = append(windowsToRemove, window)
		}
	}

	for _, window := range windowsToRemove {
		removeWindow(window)
	}

	return nil
}

func processWindowEvents(event sdl.Event) error {
	for _, window := range singleton.windows {
		if _, err := window.processEvent(event); err != nil {
			return err
		}
	}

	return nil
}
