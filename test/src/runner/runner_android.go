// +build android

package main

import (
	"runtime"

	"github.com/remogatto/mandala"
	"github.com/remogatto/mandala/test/src/testlib"
	"github.com/remogatto/prettytest"
)

type T struct{}

func (t T) Fail() {}

var t T

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	mandala.Verbose = true

	go prettytest.RunWithFormatter(
		t,
		testlib.NewTDDFormatter(),
		testlib.NewTestSuite(),
	)
}
