package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/FooSoft/lazarus/formats/dat"
	"github.com/FooSoft/lazarus/formats/dc6"
	"github.com/FooSoft/lazarus/graphics"
	"github.com/FooSoft/lazarus/math"
	"github.com/FooSoft/lazarus/platform"
	"github.com/FooSoft/lazarus/platform/imgui"
)

func loadPalette(path string) (*dat.Palette, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return dat.NewFromReader(fp)
}

func loadAnimation(path string) (*dc6.Animation, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return dc6.NewFromReader(fp)
}

type scene struct {
	animation *dc6.Animation
	palette   *dat.Palette
	texture   graphics.Texture

	directionIndex int
	frameIndex     int
}

func (s *scene) Name() string {
	return "viewer"
}

func (s *scene) Destroy() error {
	if s.texture != nil {
		return s.texture.Destroy()
	}

	return nil
}

func (s *scene) Advance() error {
	var (
		directionIndex = s.directionIndex
		frameIndex     = s.frameIndex
	)

	if s.texture == nil {
		if err := s.updateTexture(); err != nil {
			return err
		}
	}

	imgui.Begin("DC6 Viewer")
	imgui.Image(s.texture)
	direction := s.animation.Directions[directionIndex]
	if imgui.SliderInt("Direction", &directionIndex, 0, len(s.animation.Directions)-1) {
		frameIndex = 0
	}
	frame := direction.Frames[frameIndex]
	imgui.SliderInt("Frame", &frameIndex, 0, len(direction.Frames)-1)
	imgui.Columns(2)
	imgui.Text("Size")
	imgui.NextColumn()
	imgui.Text("%+v", frame.Size)
	imgui.NextColumn()
	imgui.Text("Offset")
	imgui.NextColumn()
	imgui.Text("%+v", frame.Offset)
	imgui.Columns(1)
	if imgui.Button("Exit") {
		platform.WindowSetScene(nil)
	}
	imgui.End()

	if directionIndex != s.directionIndex || frameIndex != s.frameIndex {
		s.directionIndex = directionIndex
		s.frameIndex = frameIndex
		s.updateTexture()
	}

	return nil
}

func (s *scene) updateTexture() error {
	frame := s.animation.Directions[s.directionIndex].Frames[s.frameIndex]
	colors := make([]math.Color3b, frame.Size.X*frame.Size.Y)
	for y := 0; y < frame.Size.Y; y++ {
		for x := 0; x < frame.Size.X; x++ {
			colors[y*frame.Size.X+x] = s.palette.Colors[frame.Data[y*frame.Size.X+x]]
		}
	}

	if s.texture != nil {
		if err := s.texture.Destroy(); err != nil {
			return err
		}
	}

	var err error
	s.texture, err = platform.NewTextureFromRgb(colors, frame.Size)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	platform.Initialize()
	defer platform.Shutdown()

	var (
		palettePath = flag.String("palette", "", "path to palette file")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] file\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "Parameters:\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
	}

	animation, err := loadAnimation(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var palette *dat.Palette
	if len(*palettePath) > 0 {
		palette, err = loadPalette(*palettePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		palette = dat.NewFromGrayscale()
	}

	scene := &scene{animation: animation, palette: palette}
	if err := platform.WindowCreate("Viewer", math.Vec2i{X: 1024, Y: 768}, scene); err != nil {
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
