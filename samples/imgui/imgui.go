package main

import (
	"fmt"
	"os"
	"time"

	"github.com/FooSoft/lazarus/math"
	"github.com/FooSoft/lazarus/platform"
	"github.com/FooSoft/lazarus/platform/imgui"
)

type scene struct{}

func (s *scene) Name() string {
	return "imgui"
}

func (s *scene) Advance() error {
	imgui.ShowDemoWindow()
	return nil
}

func main() {
	if err := platform.Initialize(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer platform.Shutdown()

	if err := platform.WindowCreate("ImGui", math.Vec2i{X: 1024, Y: 768}, new(scene)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer platform.WindowDestroy()

	for {
		run, err := platform.Advance()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if !run {
			break
		}

		<-time.After(time.Millisecond * 25)
	}
}
