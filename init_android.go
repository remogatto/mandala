// +build android

package gorgasm

// #include <stdlib.h>
// #include <jni.h>
// #include <android/native_activity.h>
// #include <android/input.h>
// #include "init_android.h"
//
// #cgo LDFLAGS: -landroid
import "C"

import (
	"fmt"
	"github.com/remogatto/application"
	"github.com/remogatto/egl"
	"github.com/remogatto/egl/platform"
	"log"
	"runtime/debug"
	"unsafe"
)

var (
	INITIAL_WINDOW_WIDTH, INITIAL_WINDOW_HEIGHT int
	DefaultContextAttributes                    = []int32{
		egl.CONTEXT_CLIENT_VERSION, 2,
		egl.NONE,
	}
	// DefaultConfigAttributes = []int32{
	// 	egl.RED_SIZE, 8,
	// 	egl.GREEN_SIZE, 8,
	// 	egl.BLUE_SIZE, 8,
	// 	egl.DEPTH_SIZE, 8,
	// 	egl.RENDERABLE_TYPE, egl.OPENGL_ES2_BIT,
	// 	egl.SURFACE_TYPE, egl.WINDOW_BIT,
	// 	egl.NONE,
	// }
	activityNativeEventsLoop *nativeEventsLoop
)

func assert(cond bool) {
	if !cond {
		debug.PrintStack()
		panic("Assertion failed!")
	}
}

func handleCallbackError(act *C.ANativeActivity, err interface{}) {
	if err == nil {
		return
	}
	errStr := fmt.Sprintf("callback panic: %s stack: %s", err, debug.Stack())
	errStrC := C.CString(errStr)
	defer C.free(unsafe.Pointer(errStrC))
	if C.throwException(act, errStrC) == 0 {
		log.Fatalf("%v\n", errStr)
	}
}

//export onWindowFocusChanged
func onWindowFocusChanged(act *C.ANativeActivity, focusedC C.int) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	focused := int(focusedC) != 0
	application.Debugf("onWindowFocusChanged %v...\n", focused)
	event <- WindowFocusChangedEvent{Activity: unsafe.Pointer(act), HasFocus: focused}
	application.Debugf("onWindowFocusChanged done\n")
}

//export onConfigurationChanged
func onConfigurationChanged(act *C.ANativeActivity) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	application.Logf("onConfigurationChanged\n")
}

//export onNativeWindowResized
func onNativeWindowResized(act *C.ANativeActivity) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	application.Logf("onNativeWindowResized\n")
}

//export onPause
func onPause(act *C.ANativeActivity) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	application.Logf("Pausing...\n")
	event <- PauseEvent{}
	application.Logf("Paused...\n")
}

//export onResume
func onResume(act *C.ANativeActivity) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	application.Logf("Resuming...\n")

	select {
	case event <- ResumeEvent{}:
	default:
	}

	application.Logf("Resumed...\n")
}

//export onCreate
func onCreate(act *C.ANativeActivity, savedState unsafe.Pointer, savedStateSize C.size_t) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	application.Debugf("onCreate...\n")

	activityNativeEventsLoop = newNativeEventsLoop(act)
	application.Register("nativeEventsLoop", activityNativeEventsLoop)

	activityAssetsLoop.bindActivity(unsafe.Pointer(act))

	// When the Android application is created we launch our
	// internal loops. Application's specific loops will be
	// launched in the main() function on the client side.
	go application.Run()

	application.Debugf("onCreate done\n")
}

//export onDestroy
func onDestroy(act *C.ANativeActivity) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	application.Debugf("onDestroy...\n")
	application.Debugf("onDestroy done\n")
}

//export onNativeWindowCreated
func onNativeWindowCreated(act *C.ANativeActivity, win unsafe.Pointer) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	application.Debugf("onNativeWindowCreated...\n")

	eglState := eglInitialize(win)

	// Send the EGL state to the client code to signal that the
	// rendering can start.
	initialize <- eglState

	event <- NativeWindowCreatedEvent{
		Activity: unsafe.Pointer(act),
		Window:   win,
		EGLState: eglState,
	}

	application.Debugf("onNativeWindowCreated done\n")
}

//export onNativeWindowDestroyed
func onNativeWindowDestroyed(act *C.ANativeActivity, win unsafe.Pointer) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	application.Debugf("onWindowDestroy...\n")
	// egl.DestroySurface(.Display, platform.Surface)
	application.Debugf("onWindowDestroy done\n")
	application.Exit()
}

//export onInputQueueCreated
func onInputQueueCreated(act *C.ANativeActivity, queue unsafe.Pointer) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	log.Printf("onInputQueueCreated...\n")

	activityNativeEventsLoop.event <- InputQueueCreatedEvent{
		Activity:   unsafe.Pointer(act),
		InputQueue: queue,
	}

	log.Printf("onInputQueueCreated done\n")
}

//export onInputQueueDestroyed
func onInputQueueDestroyed(act *C.ANativeActivity, queue unsafe.Pointer) {
	defer func() {
		handleCallbackError(act, recover())
	}()
	application.Logf("onInputQueueDestroy...\n")
	activityNativeEventsLoop.event <- InputQueueDestroyedEvent{
		unsafe.Pointer(act),
		queue,
	}
	application.Logf("onInputQueueDestroy done\n")
}

func getEGLDisp(disp egl.NativeDisplayType) egl.Display {
	if !egl.BindAPI(egl.OPENGL_ES_API) {
		panic("Error: eglBindAPI() failed")
	}

	egl_dpy := egl.GetDisplay(egl.DEFAULT_DISPLAY)
	if egl_dpy == egl.NO_DISPLAY {
		panic("Error: eglGetDisplay() failed\n")
	}

	var egl_major, egl_minor int32
	if !egl.Initialize(egl_dpy, &egl_major, &egl_minor) {
		panic("Error: eglInitialize() failed\n")
	}
	return egl_dpy
}

func EGLCreateWindowSurface(eglDisp egl.Display, config egl.Config, win egl.NativeWindowType) egl.Surface {
	eglSurf := egl.CreateWindowSurface(eglDisp, config, win, nil)
	if eglSurf == egl.NO_SURFACE {
		panic("Error: eglCreateWindowSurface failed\n")
	}
	return eglSurf
}

func getEGLNativeVisualId(eglDisp egl.Display, config egl.Config) int32 {
	var vid int32
	if !egl.GetConfigAttrib(eglDisp, config, egl.NATIVE_VISUAL_ID, &vid) {
		panic("Error: eglGetConfigAttrib() failed\n")
	}
	return vid
}

func chooseEGLConfig(eglDisp egl.Display) egl.Config {
	eglAttribs := []int32{
		egl.RED_SIZE, 8,
		egl.GREEN_SIZE, 8,
		egl.BLUE_SIZE, 8,
		egl.DEPTH_SIZE, 8,
		egl.RENDERABLE_TYPE, egl.OPENGL_ES2_BIT,
		egl.SURFACE_TYPE, egl.WINDOW_BIT,
		egl.NONE,
	}

	var config egl.Config
	var num_configs int32
	if !egl.ChooseConfig(eglDisp, eglAttribs, &config, 1, &num_configs) {
		panic("Error: couldn't get an EGL visual config\n")
	}

	return config
}

func eglInitialize(win unsafe.Pointer) platform.EGLState {
	var (
		width, height int32
		eglState      platform.EGLState
	)
	eglState.Display = getEGLDisp(egl.DEFAULT_DISPLAY)
	eglState.Config = chooseEGLConfig(eglState.Display)
	eglState.VisualId = getEGLNativeVisualId(eglState.Display, eglState.Config)
	C.ANativeWindow_setBuffersGeometry((*[0]byte)(win), 0, 0, C.int32_t(eglState.VisualId))
	eglState.Surface = EGLCreateWindowSurface(eglState.Display, eglState.Config, egl.NativeWindowType(win))
	egl.QuerySurface(eglState.Display, eglState.Surface, egl.WIDTH, &width)
	egl.QuerySurface(eglState.Display, eglState.Surface, egl.HEIGHT, &height)
	egl.BindAPI(egl.OPENGL_ES_API)
	eglState.Context = egl.CreateContext(eglState.Display, eglState.Config, egl.NO_CONTEXT, &DefaultContextAttributes[0])
	eglState.SurfaceWidth = int(width)
	eglState.SurfaceHeight = int(height)
	return eglState
}
