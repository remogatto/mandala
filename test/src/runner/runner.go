// +build !android

package main

import (
	glfw "github.com/go-gl/glfw3"
	"github.com/remogatto/gorgasm"
	"github.com/remogatto/gorgasm/test/src/testlib"
	"github.com/remogatto/prettytest"
	"runtime"
	"testing"
)

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer glfw.Terminate()

	gorgasm.Verbose = true
	gorgasm.Debug = true

	if !glfw.Init() {
		panic("Can't init glfw!")
	}

	// Enable OpenGL ES 2.0.
	glfw.WindowHint(glfw.ClientApi, glfw.OpenglEsApi)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	window, err := glfw.CreateWindow(testlib.Width, testlib.Height, "Gorgasm Test", nil, nil)
	if err != nil {
		panic(err)
	}

	gorgasm.Init(window)

	prettytest.Run(new(testing.T), new(testlib.TestSuite))
}
