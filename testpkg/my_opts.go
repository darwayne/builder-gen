package testpkg

import (
	"fmt"
	"io"
	"time"
)

type MyOpts struct {
	// ::builder-gen

	Interesting     *bool
	Yo              string
	Total           int
	Duration        time.Duration
	Complex         *func(string) func(reader io.Reader) time.Duration
	Hmm             io.Reader
	SimpleChan      chan struct{}
	ReceiveOnlyChan <-chan string
	ComplexChan     chan func(string) io.Reader
	Map             map[string]string
	MultiArray      [][]string
	FixedSizeArray  [3]string
}

func Yo() {
	//force fmt import to verify generate removes unused imports
	fmt.Sprintf("%+v", 1)
}
