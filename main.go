package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/darwayne/builder-gen/generator"
)

func main() {
	var (
		dir       = flag.String("dir", "", "the directory to run builder-gen on. Will use working directory if not provided")
		recursive = flag.Bool("recursive", false, "set to true to recursively iterate directories")
	)
	flag.Parse()

	opts := generator.NewDirOptsBuilder().Recursive(*recursive).Build()
	if err := generator.Dir(*dir, opts...); err != nil {
		fmt.Printf("An error occurred:%+v\n", err)
		os.Exit(1)
	}
}
