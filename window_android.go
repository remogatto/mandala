// +build android

package mandala

import (
	"unsafe"
	"github.com/remogatto/egl"
	"github.com/remogatto/egl/platform"
)

// #include <android/native_window.h>
import "C"

type window struct {
	window   unsafe.Pointer
	eglState *platform.EGLState
}

// SwapBuffers swaps the surface with the display actually showing the
// result of rendering.
func (win *window) SwapBuffers() {
	egl.SwapBuffers(win.eglState.Display, win.eglState.Surface)
}

// MakeContextCurrent bind the OpenGL context to the current thread.
func (win *window) MakeContextCurrent() {
	if ok := egl.MakeCurrent(
		win.eglState.Display,
		win.eglState.Surface,
		win.eglState.Surface,
		win.eglState.Context); !ok {
		Fatalf("%s", egl.NewError(egl.GetError()))
	}
}

// GetSize returns the dimension of the rendering surface.
func (win *window) GetSize() (int, int) {
	return win.eglState.SurfaceWidth, win.eglState.SurfaceHeight
}

func (win *window) resize(androidWin unsafe.Pointer) {
	width := int(C.ANativeWindow_getWidth((*C.ANativeWindow)(androidWin)))
	height := int(C.ANativeWindow_getHeight((*C.ANativeWindow)(androidWin)))
	win.eglState.SurfaceWidth = width
	win.eglState.SurfaceHeight = height
}
