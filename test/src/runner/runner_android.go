// +build android

package main

import (
	"runtime"

	"github.com/remogatto/gorgasm"
	"github.com/remogatto/gorgasm/test/src/testlib"
	"github.com/remogatto/prettytest"
)

type T struct{}

func (t T) Fail() {}

var t T

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	gorgasm.Verbose = true

	go prettytest.RunWithFormatter(
		t,
		testlib.NewTDDFormatter(),
		testlib.NewTestSuite(),
	)
}
