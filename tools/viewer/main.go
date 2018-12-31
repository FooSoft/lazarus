package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"

	imgui "github.com/FooSoft/imgui-go"
	"github.com/FooSoft/lazarus/formats/dat"
	"github.com/FooSoft/lazarus/formats/dc6"
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
	texture *platform.Texture

	directionIndex int32
	frameIndex     int32
}

func (s *scene) Init(window *platform.Window) error {
	return nil
}

func (s *scene) Shutdown(window *platform.Window) error {
	if s.texture == nil {
		return nil
	}

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

	imgui.Begin("DC6 Viewer")
	size := s.texture.Size()
	imgui.Image(imgui.TextureID(s.texture.Handle()), imgui.Vec2{X: float32(size.X), Y: float32(size.Y)})
	direction := s.sprite.Directions[directionIndex]
	if imgui.SliderInt("Direction", &directionIndex, 0, int32(len(s.sprite.Directions))-1) {
		frameIndex = 0
	}
	frame := direction.Frames[frameIndex]
	imgui.SliderInt("Frame", &frameIndex, 0, int32(len(direction.Frames))-1)
	imgui.Text(fmt.Sprintf("Size: %dx%d", frame.Width, frame.Height))
	imgui.Text(fmt.Sprintf("Offset: %dx%d", frame.OffsetX, frame.OffsetY))

	if directionIndex != s.directionIndex || frameIndex != s.frameIndex {
		s.directionIndex = directionIndex
		s.frameIndex = frameIndex
		s.updateTexture(window)
	}

	imgui.End()
	return nil
}

func (s *scene) updateTexture(window *platform.Window) error {
	frame := s.sprite.Directions[s.directionIndex].Frames[s.frameIndex]
	colors := make([]color.RGBA, frame.Width*frame.Height)
	for y := 0; y < frame.Height; y++ {
		for x := 0; x < frame.Width; x++ {
			colors[y*frame.Width+x] = s.palette.Colors[frame.Data[y*frame.Width+x]]
		}
	}

	if s.texture != nil {
		if err := s.texture.Destroy(); err != nil {
			return err
		}
	}

	var err error
	s.texture, err = window.CreateTextureRgba(colors, frame.Width, frame.Height)
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

	platform.Init()
	defer platform.Shutdown()

	scene := &scene{sprite: sprite, palette: palette}
	window, err := platform.CreateWindow("Viewer", 1280, 720, scene)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()

	if err := platform.ProcessEvents(); err != nil {
		log.Fatal(err)
	}
}
