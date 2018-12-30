package main

import (
	"image/color"
	"log"
	"os"

	imgui "github.com/FooSoft/imgui-go"
	"github.com/FooSoft/lazarus/formats/dat"
	"github.com/FooSoft/lazarus/formats/dc6"
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

func loadTexture(window *platform.Window, spritePath, palettePath string) (*platform.Texture, error) {
	sprite, err := loadSprite(spritePath)
	if err != nil {
		return nil, err
	}

	palette, err := loadPalette(palettePath)
	if err != nil {
		return nil, err
	}

	frame := sprite.Directions[0].Frames[0]
	colors := make([]color.RGBA, frame.Width*frame.Height)
	for y := 0; y < frame.Height; y++ {
		for x := 0; x < frame.Width; x++ {
			colors[y*frame.Width+x] = palette.Colors[frame.Data[y*frame.Width+x]]
		}
	}

	return window.CreateTextureRgba(colors, frame.Width, frame.Height)
}

type scene struct {
	texture *platform.Texture
}

func (s *scene) Init(window *platform.Window) error {
	var err error
	s.texture, err = loadTexture(window, "/home/alex/loadingscreen.dc6", "/home/alex/pal.dat")
	return err
}

func (s *scene) Advance(window *platform.Window) error {
	imgui.Text("Hello")

	window.RenderTexture(
		s.texture,
		math.Rect4i{X: 0, Y: 0, W: 256, H: 256},
		math.Rect4i{X: 0, Y: 0, W: 256, H: 256},
	)

	return nil
}

func (s *scene) Shutdown(window *platform.Window) error {
	return s.texture.Destroy()
}

func main() {
	platform.Init()
	defer platform.Shutdown()

	window, err := platform.CreateWindow("Viewer", 1280, 720, new(scene))
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()

	if err := platform.ProcessEvents(); err != nil {
		log.Fatal(err)
	}
}
