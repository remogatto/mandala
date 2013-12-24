// +build android

package gorgasm

// #include <android/native_activity.h>
// #include <android/input.h>
// #include "init_android.h"
import "C"

import (
	"log"
	"runtime"
)

const LOOPER_ID_INPUT = 0

// nativeEventsLoop handles lifecycle and input events of a NativeActivity.
type nativeEventsLoop struct {
	pause, terminate chan int

	event chan interface{}

	activity     *C.ANativeActivity
	inputQueue   *C.AInputQueue
	nativeLooper *C.ALooper
}

// newNativeInputLoop returns a new nativeEventsLoop instance.
func newNativeEventsLoop(activity *C.ANativeActivity) *nativeEventsLoop {
	nativeEventsLoop := &nativeEventsLoop{
		pause:     make(chan int),
		terminate: make(chan int),
		event:     make(chan interface{}, 1),
		activity:  activity,
	}
	return nativeEventsLoop
}

// Pause returns the pause channel of the loop.
// If a value is sent to this channel, the loop will be paused.
func (l *nativeEventsLoop) Pause() chan int {
	return l.pause
}

// Terminate returns the terminate channel of the loop.
// If a value is sent to this channel, the loop will be terminated.
func (l *nativeEventsLoop) Terminate() chan int {
	return l.terminate
}

// Run runs nativeEventsLoop.
// The loop handles native input events.
func (l *nativeEventsLoop) Run() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	looper := C.ALooper_prepare(C.ALOOPER_PREPARE_ALLOW_NON_CALLBACKS)
	if looper == nil {
		panic("ALooper_prepare returned nil")
	}

	l.nativeLooper = looper

	for {
		select {
		case <-l.pause:
			l.pause <- 0
		case <-l.terminate:
			l.terminate <- 0
		case untypedEvent := <-l.event:
			switch event := untypedEvent.(type) {
			case InputQueueCreatedEvent:
				l.inputQueue = (*C.AInputQueue)(event.InputQueue)
				C.AInputQueue_attachLooper(l.inputQueue, l.nativeLooper, LOOPER_ID_INPUT, nil, nil)
			case InputQueueDestroyedEvent:
				C.AInputQueue_detachLooper(l.inputQueue)
			}
		default:
			if l.inputQueue != nil {
				ident := C.ALooper_pollAll(-1, nil, nil, nil)
				switch ident {
				case LOOPER_ID_INPUT:
					l.processInput(l.inputQueue)
				case C.ALOOPER_POLL_ERROR:
					log.Fatalf("ALooper_pollAll returned ALOOPER_POLL_ERROR\n")
				}
			}

		}
	}
}

func (l *nativeEventsLoop) dispatchEvent(nativeEvent *C.AInputEvent) bool {
	switch C.AInputEvent_getType(nativeEvent) {
	case C.AINPUT_EVENT_TYPE_MOTION:
		action := C.AMotionEvent_getAction(nativeEvent) & C.AMOTION_EVENT_ACTION_MASK
		switch action {
		case C.AMOTION_EVENT_ACTION_UP:
			down := false
			x := float32(C.AMotionEvent_getX(nativeEvent, 0))
			y := float32(C.AMotionEvent_getY(nativeEvent, 0))
			event <- ActionUpDownEvent{Down: down, X: x, Y: y}
		case C.AMOTION_EVENT_ACTION_DOWN:
			down := true
			x := float32(C.AMotionEvent_getX(nativeEvent, 0))
			y := float32(C.AMotionEvent_getY(nativeEvent, 0))
			event <- ActionUpDownEvent{Down: down, X: x, Y: y}
		case C.AMOTION_EVENT_ACTION_MOVE:
			x := float32(C.AMotionEvent_getX(nativeEvent, 0))
			y := float32(C.AMotionEvent_getY(nativeEvent, 0))
			event <- ActionMoveEvent{X: x, Y: y}
		}
	}
	return false
}

func (l *nativeEventsLoop) processInput(inputQueue *C.AInputQueue) {
	var event *C.AInputEvent
	for {
		if ret := C.AInputQueue_getEvent(inputQueue, &event); ret < 0 {
			break
		}
		if C.AInputQueue_preDispatchEvent(inputQueue, event) != 0 {
			continue
		}
		handled := l.dispatchEvent(event)
		var handledC C.int
		if handled {
			handledC = 1
		}
		C.AInputQueue_finishEvent(inputQueue, event, handledC)
	}
}
