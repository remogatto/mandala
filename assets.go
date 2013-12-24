package gorgasm

import (
	"io"
	"unsafe"
)

var (
	AssetsPath string = "android"
)

type LoadAssetResponse struct {
	Buffer io.Reader
	Error  error
}

type LoadAssetRequest struct {
	Filename string
	Response chan LoadAssetResponse
}

// nativeEventsLoop handles lifecycle and input events of a NativeActivity.
type assetsLoop struct {
	pause, terminate chan int
	activity         unsafe.Pointer
	command          chan interface{}
}

// newNativeInputLoop returns a new nativeEventsLoop instance.
func newAssetsLoop() *assetsLoop {
	assetsLoop := &assetsLoop{
		pause:     make(chan int),
		terminate: make(chan int),
		command:   make(chan interface{}, 1),
	}
	return assetsLoop
}

func (l *assetsLoop) bindActivity(activity unsafe.Pointer) {
	l.activity = activity
}

// Pause returns the pause channel of the loop.
// If a value is sent to this channel, the loop will be paused.
func (l *assetsLoop) Pause() chan int {
	return l.pause
}

// Terminate returns the terminate channel of the loop.
// If a value is sent to this channel, the loop will be terminated.
func (l *assetsLoop) Terminate() chan int {
	return l.terminate
}

// Run runs nativeEventsLoop.
// The loop handles native input events.
func (l *assetsLoop) Run() {
	for {
		select {
		case <-l.pause:
			l.pause <- 0
		case <-l.terminate:
			l.terminate <- 0
		case untypedCommand := <-l.command:
			switch command := untypedCommand.(type) {
			case LoadAssetRequest:
				file, err := loadAsset(l.activity, command.Filename)
				command.Response <- LoadAssetResponse{file, err}
			}
		}
	}
}

// func LoadAsset(filename string) <-chan io.Reader {
// 	command := LoadAssetCommand{
// 		Filename: filename,
// 		Buffer:   make(chan io.Reader),
// 	}
// 	Assets <- command
// 	return command.Buffer
// }
