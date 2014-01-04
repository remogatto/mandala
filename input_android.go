// +build android

package gorgasm

// #include <android/native_activity.h>
// #include <android/input.h>
import "C"

import (
	"runtime"
	"git.tideland.biz/goas/loop"
)

const LOOPER_ID_INPUT = 0

// The loop handles native input events.
func androidEventLoopFunc(event chan interface{}, looperCh chan *C.ALooper) loop.LoopFunc {
	return func(l loop.Loop) error {
		var inputQueue *C.AInputQueue
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		looper := C.ALooper_prepare(C.ALOOPER_PREPARE_ALLOW_NON_CALLBACKS)
		if looper == nil {
			Fatalf("ALooper_prepare returned nil")
		}

		looperCh <- looper

		for {
			select {
			case untypedEvent := <-event:
				switch event := untypedEvent.(type) {
				case inputQueueCreatedEvent:
					inputQueue = (*C.AInputQueue)(event.inputQueue)
					C.AInputQueue_attachLooper(inputQueue, looper, LOOPER_ID_INPUT, nil, nil)
				case inputQueueDestroyedEvent:
					inputQueue = (*C.AInputQueue)(event.inputQueue)
					C.AInputQueue_detachLooper(inputQueue)
				}
			default:
				if inputQueue != nil {
					ident := C.ALooper_pollAll(-1, nil, nil, nil)
					switch ident {
					case LOOPER_ID_INPUT:
						processInput(inputQueue)
					case C.ALOOPER_POLL_ERROR:
						Fatalf("ALooper_pollAll returned ALOOPER_POLL_ERROR\n")
					}
				}
			}
		}
	}
}

func dispatchEvent(nativeEvent *C.AInputEvent) bool {
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
			Debugf("Send event")
			event <- ActionMoveEvent{X: x, Y: y}
			Debugf("Event sent")
		}
	}
	return false
}

func processInput(inputQueue *C.AInputQueue) {
	var event *C.AInputEvent
	for {
		if ret := C.AInputQueue_getEvent(inputQueue, &event); ret < 0 {
			break
		}
		if C.AInputQueue_preDispatchEvent(inputQueue, event) != 0 {
			continue
		}
		handled := dispatchEvent(event)
		var handledC C.int
		if handled {
			handledC = 1
		}
		C.AInputQueue_finishEvent(inputQueue, event, handledC)
	}
}
