package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/FooSoft/lazarus/formats/dat"
	"github.com/FooSoft/lazarus/formats/dc6"
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

func extractSprite(palettePath, spritePath, targetDir string) error {
	_, err := loadPalette(palettePath)
	if err != nil {
		return err
	}

	_, err = loadSprite(spritePath)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var (
		targetDir = flag.String("target", ".", "target directory")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] palette_file dc6_file\n", path.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "Parameters:\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(2)
	}

	if err := extractSprite(flag.Arg(0), flag.Arg(1), *targetDir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
