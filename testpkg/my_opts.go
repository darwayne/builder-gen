package testpkg

import (
	"fmt"
	"time"
)

type MyOpts struct {
	// ::builder-gen

	Interesting *bool
	Yo          string
	Total       int
	Duration    time.Duration
}

func Yo() {
	//force fmt import to verify generate removes unused imports
	fmt.Sprintf("%+v", 1)
}
