// +build android

package main

import (
	"github.com/remogatto/gorgasm"
	"github.com/remogatto/gorgasm/integration_test/src/testlib"
	"github.com/remogatto/prettytest"
	"runtime"
)

type T struct{}

func (t T) Fail() {}

var t T

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	gorgasm.Verbose = true
	gorgasm.Debug = true

	go prettytest.RunWithFormatter(t, testlib.NewTDDFormatter(), new(testlib.TestSuite))
}
