package platform

import (
	"log"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

var platfromState struct {
	sdlIsInit bool
}

func Advance() (bool, error) {
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

func platformInit() error {
	if platfromState.sdlIsInit {
		return nil
	}

	runtime.LockOSThread()

	log.Println("sdl init")
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return err
	}

	platfromState.sdlIsInit = true
	return nil
}
