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
)

func loadPalette(path string) (*dat.Palette, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return dat.NewFromReader(fp)
}

func loadSprite(path string) (*dc6.Sprite, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return dc6.NewFromReader(fp)
}

type scene struct {
	sprite  *dc6.Sprite
	palette *dat.Palette
	texture graphics.Texture

	directionIndex int
	frameIndex     int
}

func (s *scene) Name() string {
	return "Viewer"
}

func (s *scene) Destroy(window *platform.Window) error {
	return s.texture.Destroy()
}

func (s *scene) Advance(window *platform.Window) error {
	var (
		directionIndex = s.directionIndex
		frameIndex     = s.frameIndex
	)

	if s.texture == nil {
		if err := s.updateTexture(window); err != nil {
			return err
		}
	}

	imgui := window.Imgui()

	imgui.DialogBegin("DC6 Viewer")
	imgui.Image(s.texture)
	direction := s.sprite.Directions[directionIndex]
	if imgui.SliderInt("Direction", &directionIndex, 0, len(s.sprite.Directions)-1) {
		frameIndex = 0
	}
	frame := direction.Frames[frameIndex]
	imgui.SliderInt("Frame", &frameIndex, 0, len(direction.Frames)-1)
	imgui.Text(fmt.Sprintf("Size: %+v", frame.Size))
	imgui.Text(fmt.Sprintf("Offset: %+v", frame.Offset))
	if imgui.Button("Exit") {
		window.SetScene(nil)
	}
	imgui.DialogEnd()

	if directionIndex != s.directionIndex || frameIndex != s.frameIndex {
		s.directionIndex = directionIndex
		s.frameIndex = frameIndex
		s.updateTexture(window)
	}

	return nil
}

func (s *scene) updateTexture(window *platform.Window) error {
	frame := s.sprite.Directions[s.directionIndex].Frames[s.frameIndex]
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
	s.texture, err = window.CreateTextureRgb(colors, frame.Size)
	if err != nil {
		return err
	}

	return nil
}

func main() {
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

	sprite, err := loadSprite(flag.Arg(0))
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

	scene := &scene{sprite: sprite, palette: palette}
	window, err := platform.NewWindow("viewer", math.Vec2i{X: 1024, Y: 768}, scene)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer window.Destroy()

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
