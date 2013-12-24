package gorgasm

import (
	"github.com/remogatto/egl/platform"
	"unsafe"
)

type WindowFocusChangedEvent struct {
	HasFocus bool
	Activity unsafe.Pointer
}

type ConfigurationChangedEvent struct {
	Activity unsafe.Pointer
}

type NativeWindowResizedEvent struct {
	Activity      unsafe.Pointer
	Width, Height int
}

type PauseEvent struct {
	Activity unsafe.Pointer
}

type ResumeEvent struct {
	Activity unsafe.Pointer
}

type CreateEvent struct {
	Activity       unsafe.Pointer
	SavedState     unsafe.Pointer
	SavedStateSize int
}

type DestroyEvent struct {
	Activity unsafe.Pointer
}

type NativeWindowCreatedEvent struct {
	EGLState platform.EGLState
	Activity unsafe.Pointer
	Window   unsafe.Pointer
}

type NativeWindowDestroyedEvent struct {
	Activity unsafe.Pointer
	Window   unsafe.Pointer
}

type InputQueueCreatedEvent struct {
	Activity   unsafe.Pointer
	InputQueue unsafe.Pointer
}

type InputQueueDestroyedEvent struct {
	Activity   unsafe.Pointer
	InputQueue unsafe.Pointer
}

type ActionUpDownEvent struct {
	Activity unsafe.Pointer
	Down     bool    // Is finger down on the screen?
	X, Y     float32 // Coordinates of the touched point
}

type ActionMoveEvent struct {
	Activity unsafe.Pointer
	X, Y     float32 // Coordinates of the touched point in movement
}
