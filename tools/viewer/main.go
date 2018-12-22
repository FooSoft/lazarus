package main

import (
	"image/color"
	"log"
	"os"

	"github.com/FooSoft/lazarus/formats/dat"
	"github.com/FooSoft/lazarus/formats/dc6"
	"github.com/FooSoft/lazarus/graphics"
	"github.com/veandco/go-sdl2/sdl"
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

func loadSurface(spritePath, palettePath string) (*sdl.Surface, error) {
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

	return graphics.NewSurfaceFromRgba(colors, frame.Width, frame.Height)
}

func main() {
	window, err := sdl.CreateWindow("Viewer", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Destroy()

	sprite, err := loadSurface("/home/alex/loadingscreen.dc6", "/home/alex/pal.dat")
	if err != nil {
		log.Fatal(err)
	}

	texture, err := renderer.CreateTextureFromSurface(sprite)
	if err != nil {
		log.Fatal(err)
	}

	renderer.Clear()
	renderer.Copy(texture, &sdl.Rect{0, 0, 256, 256}, &sdl.Rect{0, 0, 256, 256})
	renderer.Present()

	for running := true; running; {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			}
			sdl.Delay(1)
		}
	}
}
