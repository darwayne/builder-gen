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
		trace               = flag.Bool("trace", false, "trace enables raw AST tracing")
	)
	flag.Parse()

	fmt.Println(`
                                                                                                                  
88888888ba                88  88           88                               ,ad8888ba,                            
88      "8b               ""  88           88                              d8"'    '"8b                           
88      ,8P                   88           88                             d8'                                     
88aaaaaa8P'  88       88  88  88   ,adPPYb,88   ,adPPYba,  8b,dPPYba,     88              ,adPPYba,  8b,dPPYba,   
88""""""8b,  88       88  88  88  a8"    'Y88  a8P_____88  88P'   "Y8     88      88888  a8P_____88  88P'   '"8a  
88      '8b  88       88  88  88  8b       88  8PP"""""""  88             Y8,        88  8PP"""""""  88       88  
88      a8P  "8a,   ,a88  88  88  "8a,   ,d88  "8b,   ,aa  88              Y8a.    .a88  "8b,   ,aa  88       88  
88888888P"    '"YbbdP'Y8  88  88   '"8bbdP"Y8   '"Ybbd8"'  88               '"Y88888P"    '"Ybbd8"'  88       88  
                                                                                                                  
 `)

	var exclusions []string
	for _, ex := range strings.Split(*recursiveExclusions, ",") {
		if strings.TrimSpace(ex) != "" {
			exclusions = append(exclusions, ex)
		}
	}

	opts := generator.NewDirOptsBuilder().
		Recursive(*recursive).
		RecursiveExclusions(exclusions...).
		Trace(*trace).
		Build()
	if err := generator.Dir(*dir, opts...); err != nil {
		fmt.Printf("An error occurred:%+v\n", err)
		os.Exit(1)
	}
}
