// +build !android

package main

import (
	"runtime"
	"testing"

	glfw "github.com/go-gl/glfw3"
	"github.com/remogatto/mandala"
	"github.com/remogatto/mandala/test/src/testlib"
	"github.com/remogatto/prettytest"
)

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer glfw.Terminate()

	mandala.Verbose = true

	if !glfw.Init() {
		panic("Can't init glfw!")
	}

	// Enable OpenGL ES 2.0.
	glfw.WindowHint(glfw.ClientApi, glfw.OpenglEsApi)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	window, err := glfw.CreateWindow(testlib.Width, testlib.Height, "Mandala Test", nil, nil)
	if err != nil {
		panic(err)
	}

	mandala.Init(window)

	go prettytest.Run(new(testing.T), testlib.NewTestSuite())

	for !window.ShouldClose() {
		glfw.WaitEvents()
	}
}
