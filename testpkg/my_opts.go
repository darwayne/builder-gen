package testpkg

import "fmt"

type MyOpts struct {
	// ::builder-gen

	Interesting *bool
	Yo          string
	Total       int
}

func Yo() {
	//force fmt import to verify generate removes unused imports
	fmt.Sprintf("%+v", 1)
}
