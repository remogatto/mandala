// +build !android

package gorgasm

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/remogatto/application"
	"github.com/remogatto/egl"
	"github.com/remogatto/egl/platform/xorg"
)

var X *xgbutil.XUtil

func newWindow(X *xgbutil.XUtil, width, height int) *xwindow.Window {
	var (
		err error
		win *xwindow.Window
	)
	win, err = xwindow.Generate(X)
	if err != nil {
		panic(err)
	}
	win.Create(X.RootWin(), 0, 0, width, height,
		xproto.CwBackPixel|xproto.CwEventMask,
		0, xproto.EventMaskButtonRelease)
	win.WMGracefulClose(
		func(w *xwindow.Window) {
			xevent.Detach(w.X, w.Id)
			mousebind.Detach(w.X, w.Id)
			// w.Destroy()
			xevent.Quit(X)
			application.Exit()
		})

	// In order to get ConfigureNotify events, we must listen to the window
	// using the 'StructureNotify' mask.
	win.Listen(xproto.EventMaskButtonPress |
		xproto.EventMaskButtonRelease |
		xproto.EventMaskKeyPress |
		xproto.EventMaskKeyRelease |
		xproto.EventMaskStructureNotify)

	win.Map()

	xevent.ConfigureNotifyFun(
		func(X *xgbutil.XUtil, ev xevent.ConfigureNotifyEvent) {
			event <- NativeWindowResizedEvent{
				Width:  int(ev.Width),
				Height: int(ev.Height),
			}
		}).Connect(X, win.Id)

	mousebind.Drag(X, win.Id, win.Id, "1", false,
		func(X *xgbutil.XUtil, rx, ry, ex, ey int) (bool, xproto.Cursor) {
			event <- ActionUpDownEvent{Down: true, X: float32(ex), Y: float32(ey)}
			return true, 0
		},
		func(X *xgbutil.XUtil, rx, ry, ex, ey int) {
			event <- ActionMoveEvent{X: float32(ex), Y: float32(ey)}
		},
		func(X *xgbutil.XUtil, rx, ry, ex, ey int) {
			event <- ActionUpDownEvent{Down: false, X: float32(ex), Y: float32(ey)}
		})

	xevent.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			// keybind.LookupString does the magic of implementing parts of
			// the X Keyboard Encoding to determine an english representation
			// of the modifiers/keycode tuple.
			// N.B. It's working for me, but probably isn't 100% correct in
			// all environments yet.
			modStr := keybind.ModifierString(e.State)
			keyStr := keybind.LookupString(X, e.State, e.Detail)
			if len(modStr) > 0 {
				application.Logf("Key: %s-%s\n", modStr, keyStr)
			} else {
				application.Logf("Key: %s", keyStr)
			}

			if keyStr == "p" || keyStr == "P" {
				event <- PauseEvent{}
			}
			if keyStr == "r" || keyStr == "R" {
				event <- ResumeEvent{}
			}

		}).Connect(X, win.Id)

	if err != nil {
		panic(err)
	}
	return win
}

func XorgInitialize(width, height int) {
	X, err := xgbutil.NewConn()
	if err != nil {
		panic(err)
	}

	mousebind.Initialize(X)
	keybind.Initialize(X)

	xWindow := newWindow(X, width, height)
	go xevent.Main(X)
	eglState := xorg.Initialize(egl.NativeWindowType(uintptr(xWindow.Id)), xorg.DefaultConfigAttributes, xorg.DefaultContextAttributes)
	go func() {
		initialize <- *eglState
		event <- NativeWindowCreatedEvent{EGLState: *eglState}
	}()
}
