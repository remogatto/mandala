package mandala

import (
	"unsafe"
)

// A NativeWindowCreatedEvent is triggered when the native window is
// created. It indicates that the native surface is ready for
// rendering.
type NativeWindowCreatedEvent struct {
	Activity unsafe.Pointer
	Window   Window
}

// A NativeWindowDestroyedEvent is triggered when the native window is
// destroyed and nothing can be rendered on it anymore.
type NativeWindowDestroyedEvent struct {
	Activity unsafe.Pointer
	Window   Window
}

// A WindowFocusChangedEvent is triggered when the native window gain
// or lost its focus.
type WindowFocusChangedEvent struct {
	Activity unsafe.Pointer
	HasFocus bool
}

// A ConfigurationChangedEvent is triggered when the application
// changes its configuration.
type ConfigurationChangedEvent struct {
	Activity unsafe.Pointer
}

// A NativeWindowResizedEvent is triggered when the native window is
// resized. It happens, for example, when the device is rotated.
type NativeWindowResizedEvent struct {
	Activity unsafe.Pointer
	Window   Window
}

// A NativeWindowRedrawNeededEvent is triggered when the native window
// needs to be redrawn.
type NativeWindowRedrawNeededEvent struct {
	Activity unsafe.Pointer
	Window   Window
}

// A PauseEvent is triggered when the application is paused. It
// happens, for example, when the back button is pressed and the
// application goes in background. Please note that the framework will
// wait for a value from Paused channel before actually pause the
// application.
type PauseEvent struct {
	Activity unsafe.Pointer
	Paused   chan bool
}

// A ResumeEvent is triggered when the application goes
// foreground. This doesn't mean that a native surface is ready for
// rendering.
type ResumeEvent struct {
	Activity unsafe.Pointer
}

// CreateEvent is the first event triggered by the application.
type CreateEvent struct {
	Activity       unsafe.Pointer
	SavedState     unsafe.Pointer
	SavedStateSize int
}

// A StartEvent is triggered after a CreateEvent. It initiates the
// "visible" lifespan of the application.
type StartEvent struct {
	Activity unsafe.Pointer
}

// DestroyEvent is the last event triggered by the application before
// terminate.
type DestroyEvent struct {
	Activity unsafe.Pointer
}

// ActionUpDownEvent is triggered when the user has the finger down/up
// the device's surface.
type ActionUpDownEvent struct {
	Activity unsafe.Pointer

	// The finger is down on the surface
	Down bool

	// Coordinates of the touched point
	X, Y float32
}

// ActionMoveEvent is triggered when the user moves the finger on the
// device surface.
type ActionMoveEvent struct {
	Activity unsafe.Pointer

	// Coordinates of the touched point in movement
	X, Y float32
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
