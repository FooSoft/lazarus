package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path"
	"path/filepath"

	"github.com/FooSoft/lazarus/formats/dat"
	"github.com/FooSoft/lazarus/formats/dc6"
)

func loadPalette(path string) (*dat.DatPalette, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return dat.NewFromReader(fp)
}

func loadSprite(path string) (*dc6.Dc6Animation, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return dc6.NewFromReader(fp)
}

func extractSprite(spritePath string, palette *dat.DatPalette, targetDir string) error {
	sprite, err := loadSprite(spritePath)
	if err != nil {
		return err
	}

	for di, direction := range sprite.Directions {
		for fi, frame := range direction.Frames {
			img := image.NewRGBA(image.Rect(0, 0, frame.Size.X, frame.Size.Y))
			for y := 0; y < frame.Size.Y; y++ {
				for x := 0; x < frame.Size.X; x++ {
					colorSrc := palette.Colors[frame.Data[y*frame.Size.X+x]]
					colorDst := color.RGBA{colorSrc.R, colorSrc.G, colorSrc.B, 0xff}
					img.Set(x, y, colorDst)
				}
			}

			basePath := filepath.Base(spritePath)
			targetPath := fmt.Sprintf("%s_%d_%d.png", filepath.Join(targetDir, basePath), di, fi)

			fp, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			defer fp.Close()

			if err := png.Encode(fp, img); err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	targetDir := flag.String("target", ".", "target directory")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] palette_file [dc6_files]\n", path.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "Parameters:\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
	}

	palette, err := loadPalette(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for i := 1; i < flag.NArg(); i++ {
		if err := extractSprite(flag.Arg(1), palette, *targetDir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
