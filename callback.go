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

func mouseButtonCallback(
	window *glfw.Window,
	button glfw.MouseButton,
	action glfw.Action,
	mod glfw.ModifierKey) {

	if button == glfw.MouseButton1 {
		down := action == glfw.Press
		x, y := window.GetCursorPosition()
		event <- ActionUpDownEvent{
			Down: down,
			X:    float32(x),
			Y:    float32(y),
		}
	}
}

func cursorPositionCallback(window *glfw.Window, x float64, y float64) {
	event <- ActionMoveEvent{
		X: float32(x),
		Y: float32(y),
	}
}

func Init(window *glfw.Window) {

	glfw.SetErrorCallback(errorCallback)

	// Set callbacks associated with window events
	window.SetCloseCallback(exitCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)
	window.SetCursorPositionCallback(cursorPositionCallback)

	// Begin sending events related to the creation process
	event <- CreateEvent{}
	event <- StartEvent{}
	event <- ResumeEvent{}
	event <- NativeWindowCreatedEvent{Window: window}
}
