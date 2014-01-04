package gorgasm

import (
	"unsafe"
)

type NativeWindowCreatedEvent struct {
	Activity unsafe.Pointer
	Window   Window
}

type NativeWindowDestroyedEvent struct {
	Activity unsafe.Pointer
	Window   Window
}

type WindowFocusChangedEvent struct {
	Activity unsafe.Pointer
	HasFocus bool
}

type ConfigurationChangedEvent struct {
	Activity unsafe.Pointer
}

type NativeWindowResizedEvent struct {
	Activity unsafe.Pointer
	Window   Window
}

type NativeWindowRedrawNeededEvent struct {
	Activity unsafe.Pointer
	Window   Window
}

type PauseEvent struct {
	Activity unsafe.Pointer
}

type ResumeEvent struct {
	Activity unsafe.Pointer
}

type StartEvent struct {
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

type ActionUpDownEvent struct {
	Activity unsafe.Pointer
	Down     bool
	X, Y     float32 // Coordinates of the touched point
}

type ActionMoveEvent struct {
	Activity unsafe.Pointer
	X, Y     float32 // Coordinates of the touched point in movement
}

// Internal events

type inputQueueCreatedEvent struct {
	activity   unsafe.Pointer
	inputQueue unsafe.Pointer
}

type inputQueueDestroyedEvent struct {
	activity   unsafe.Pointer
	inputQueue unsafe.Pointer
}
