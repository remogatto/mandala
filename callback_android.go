// +build android

package mandala

// #include <stdlib.h>
// #include <jni.h>
// #include <android/native_activity.h>
// #include <android/native_window.h>
// #include <android/input.h>
// #include "callback_android.h"
//
// #cgo LDFLAGS: -landroid
import "C"

import (
	"fmt"
	"unsafe"
	"git.tideland.biz/goas/loop"
	"github.com/remogatto/egl/platform/android"
)

var (
	// Internal channel used by the framework to handle events
	// like the creation/destruction of input queues.
	internalEvent chan interface{}
	looper        *C.ALooper
)

func handleCallbackError(act *C.ANativeActivity, err interface{}) {
	if err == nil {
		return
	}
	errStr := fmt.Sprintf("callback panic: %s stack: %s", err, Stacktrace())
	errStrC := C.CString(errStr)
	defer C.free(unsafe.Pointer(errStrC))
	if C.throwException(act, errStrC) == 0 {
		Fatalf("%v\n", errStr)
	}
}

//export onWindowFocusChanged
func onWindowFocusChanged(act *C.ANativeActivity, focusedC C.int) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	focused := int(focusedC) != 0
	Debugf("onWindowFocusChanged %v...\n", focused)
	event <- WindowFocusChangedEvent{Activity: unsafe.Pointer(act), HasFocus: focused}
	Debugf("onWindowFocusChanged done\n")
}

//export onConfigurationChanged
func onConfigurationChanged(act *C.ANativeActivity) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	event <- ConfigurationChangedEvent{}
	Debugf("onConfigurationChanged\n")
}

//export onNativeWindowResized
func onNativeWindowResized(act *C.ANativeActivity, win unsafe.Pointer) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	state := getState(act)
	state.window.resize(win)
	event <- NativeWindowResizedEvent{unsafe.Pointer(act), state.window}
	Debugf("onNativeWindowResized\n")
}

//export onNativeWindowRedrawNeeded
func onNativeWindowRedrawNeeded(act *C.ANativeActivity, win unsafe.Pointer) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	state := getState(act)
	state.window.resize(win)
	event <- NativeWindowRedrawNeededEvent{
		unsafe.Pointer(act),
		state.window,
	}
	Debugf("onNativeRedrawNeeded\n")
}

//export onPause
func onPause(act *C.ANativeActivity) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	Debugf("Pausing...\n")
	event <- PauseEvent{unsafe.Pointer(act)}
	Debugf("Paused...\n")
}

//export onResume
func onResume(act *C.ANativeActivity) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	Debugf("Resuming...\n")
	event <- ResumeEvent{unsafe.Pointer(act)}
	Debugf("Resumed...\n")
}

//export onStart
func onStart(act *C.ANativeActivity) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	event <- StartEvent{unsafe.Pointer(act)}
}

//export onCreate
func onCreate(act *C.ANativeActivity, savedState unsafe.Pointer, savedStateSize C.size_t) {
	defer func() {
		handleCallbackError(act, recover())
	}()

	internalEvent = make(chan interface{})
	looperCh := make(chan *C.ALooper)

	// Create a new state for the current activity and store it in
	// states global map.
	setState(act, &state{act, nil})

	// Initialize the native sound library
	err := initOpenSL()
	if err != nil {
		Logf(err.Error())
	} else {
		Debugf("OpenSL successfully initialized")
	}

	// Initialize the native event loop
	loop.GoRecoverable(
		androidEventLoopFunc(internalEvent, looperCh),
		func(rs loop.Recoverings) (loop.Recoverings, error) {
			for _, r := range rs {
				Logf("%s", r.Reason)
				Logf("%s", Stacktrace())
			}
			return rs, fmt.Errorf("Unrecoverable loop\n")
		},
	)

	activity <- unsafe.Pointer(act)
	looper = <-looperCh

	event <- CreateEvent{unsafe.Pointer(act), savedState, int(savedStateSize)}
}

//export onDestroy
func onDestroy(act *C.ANativeActivity) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	Debugf("onDestroy...\n")

	// Shutdown the sound engine
	Debugf("Shutdown OpenSL")
	shutdownOpenSL()

	Debugf("onDestroy done\n")
}

//export onNativeWindowCreated
func onNativeWindowCreated(act *C.ANativeActivity, win unsafe.Pointer) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	Debugf("onNativeWindowCreated...\n")

	eglState := android.Initialize(
		win,
		android.DefaultConfigAttributes,
		android.DefaultContextAttributes,
	)

	state := getState(act)
	state.window = &window{win, eglState}

	event <- NativeWindowCreatedEvent{
		Activity: unsafe.Pointer(act),
		Window:   state.window,
	}

	Debugf("onNativeWindowCreated done\n")
}

//export onNativeWindowDestroyed
func onNativeWindowDestroyed(act *C.ANativeActivity, win unsafe.Pointer) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	Debugf("onWindowDestroy...\n")

	state := getState(act)
	event <- NativeWindowDestroyedEvent{
		Activity: unsafe.Pointer(act),
		Window:   state.window,
	}
	Debugf("onWindowDestroy done\n")
}

//export onInputQueueCreated
func onInputQueueCreated(act *C.ANativeActivity, queue unsafe.Pointer) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	Debugf("onInputQueueCreated...\n")

	internalEvent <- inputQueueCreatedEvent{
		activity:   unsafe.Pointer(act),
		inputQueue: queue,
	}

	Debugf("onInputQueueCreated done\n")
}

//export onInputQueueDestroyed
func onInputQueueDestroyed(act *C.ANativeActivity, queue unsafe.Pointer) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	Debugf("onInputQueueDestroy...\n")

	C.ALooper_wake(looper)

	internalEvent <- inputQueueDestroyedEvent{
		unsafe.Pointer(act),
		queue,
	}

	Debugf("onInputQueueDestroy done\n")
}
