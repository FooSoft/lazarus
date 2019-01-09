package platform

import (
	"errors"
	"log"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	ErrPlatformNotInit = errors.New("platform is not initialized")
)

var platformState struct {
	isInit bool
}

func Initialize() error {
	log.Println("platform init")

	if platformState.isInit {
		return nil
	}

	runtime.LockOSThread()

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return err
	}

	platformState.isInit = true
	return nil
}

func Shutdown() error {
	log.Println("platform shutdown")

	if !platformState.isInit {
		return nil
	}

	if err := FileUnmountAll(); err != nil {
		return err
	}

	return nil
}

func Advance() (bool, error) {
	if !platformState.isInit {
		return false, ErrPlatformNotInit
	}

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			return false, nil
		default:
			if _, err := windowProcessEvent(event); err != nil {
				return false, err
			}
		}
	}

	run, err := windowAdvance()
	if !run {
		WindowDestroy()
	}

	return run, err
}
