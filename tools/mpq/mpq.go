package main

import (
	"flag"
	"log"
	"os"

	"github.com/FooSoft/lazarus/filesystem"
)

func main() {
	// var (
	// 	wildcard = flag.String("wildcard", "*.*", "wildcard filter")
	// )
	flag.Parse()

	if flag.NArg() == 0 {
		flag.PrintDefaults()
		os.Exit(2)
	}

	fs := filesystem.New()
	if err := fs.Mount("", flag.Arg(0)); err != nil {
		log.Fatal(err)
	}

	paths, err := fs.List()
	if err != nil {
		log.Fatal(err)
	}

	log.Print(paths)
}
