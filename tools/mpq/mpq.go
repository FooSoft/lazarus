package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/FooSoft/lazarus/formats/mpq"
	"github.com/bmatcuk/doublestar"
)

func list(mpqPath, filter string) error {
	arch, err := mpq.NewFromFile(mpqPath)
	if err != nil {
		return err
	}
	defer arch.Close()

	for _, resPath := range arch.GetPaths() {
		match, err := doublestar.Match(filter, resPath)
		if err != nil {
			return err
		}

		if match {
			fmt.Println(resPath)
		}
	}

	return nil
}

func extract(mpqPath, filter, targetDir string) error {
	arch, err := mpq.NewFromFile(mpqPath)
	if err != nil {
		return err
	}
	defer arch.Close()

	for _, resPath := range arch.GetPaths() {
		match, err := doublestar.Match(filter, resPath)
		if err != nil {
			return err
		}

		if !match {
			continue
		}

		fmt.Println(resPath)

		resFile, err := arch.OpenFile(resPath)
		if err != nil {
			return err
		}
		defer resFile.Close()

		sysPath := path.Join(targetDir, resPath)
		if err := os.MkdirAll(path.Dir(sysPath), 0777); err != nil {
			return err
		}

		sysFile, err := os.Create(sysPath)
		if err != nil {
			return err
		}

		if _, err := io.Copy(sysFile, resFile); err != nil {
			sysFile.Close()
			return err
		}

		sysFile.Close()
	}

	return nil
}

func main() {
	var (
		filter    = flag.String("filter", "**", "wildcard file filter")
		targetDir = flag.String("target", ".", "target directory")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] command [files]\n", path.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "Parameters:\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
	}

	switch flag.Arg(0) {
	case "list":
		for i := 1; i < flag.NArg(); i++ {
			if err := list(flag.Arg(i), *filter); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	case "extract":
		for i := 1; i < flag.NArg(); i++ {
			if err := extract(flag.Arg(i), *filter, *targetDir); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	default:
		flag.Usage()
		os.Exit(2)
	}
}
