package gorgasm

import (
	"github.com/remogatto/application"
	"github.com/remogatto/egl/platform"
)

var (
	// This the channel from which client code receive the EGL
	// state. Then rendering loops could start.
	Init <-chan platform.EGLState

	// This is the public global event channel. Client code listen
	// to this channel in order to receive system events.
	Events <-chan interface{}

	// Send commands to this channel in order to manage
	// application's assets.
	Assets chan<- interface{}

	// This is the internal global event channel. All system
	// events should be sent to this channel.
	event chan interface{}

	// This is the internal global init channel. Platform-specific
	// code sends an EGL state to it.
	initialize chan platform.EGLState

	activityAssetsLoop *assetsLoop
)

func init() {
	event = make(chan interface{})
	initialize = make(chan platform.EGLState, 1)

	activityAssetsLoop = newAssetsLoop()
	application.Register("gorgasm.assetsLoop", activityAssetsLoop)

	Events = event
	Assets = activityAssetsLoop.command
	Init = initialize
}
