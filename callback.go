// +build !android

package gorgasm

import (
	glfw "github.com/go-gl/glfw3"
)

func errorCallback(err glfw.ErrorCode, desc string) {
	Logf("%v: %v\n", err, desc)
}

func exitCallback(window *glfw.Window) {
	event <- DestroyEvent{}
}

func Init(window *glfw.Window) {

	glfw.SetErrorCallback(errorCallback)

	// Set callbacks associated with window events
	window.SetCloseCallback(exitCallback)

	event <- NativeWindowCreatedEvent{Window: window}
}
