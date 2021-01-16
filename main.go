package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/darwayne/builder-gen/generator"
)

func main() {
	var (
		dir                 = flag.String("dir", "", "the directory to run builder-gen on. Will use working directory if not provided")
		recursive           = flag.Bool("recursive", false, "set to true to recursively iterate directories; auto excludes any directories starting with .")
		recursiveExclusions = flag.String("recursive-exclusions", "", "a comma separated list of directories to exclude when recursively iterating")
	)
	flag.Parse()

	var exclusions []string
	for _, ex := range strings.Split(*recursiveExclusions, ",") {
		if strings.TrimSpace(ex) != "" {
			exclusions = append(exclusions, ex)
		}
	}
	opts := generator.NewDirOptsBuilder().
		Recursive(*recursive).
		RecursiveExclusions(exclusions...).
		Build()
	if err := generator.Dir(*dir, opts...); err != nil {
		fmt.Printf("An error occurred:%+v\n", err)
		os.Exit(1)
	}
}
