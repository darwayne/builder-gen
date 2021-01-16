package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/darwayne/builder-gen/generator"
)

func main() {
	var (
		dir = flag.String("dir", "", "the directory to run builder gen on will use working directory if not provided")
	)

	flag.Parse()

	if err := generator.Dir(*dir); err != nil {
		fmt.Printf("An error occurred:%+v\n", err)
		os.Exit(1)
	}
}
