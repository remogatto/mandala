package gorgasm

import (
	"github.com/remogatto/application"
	"github.com/remogatto/egl/platform"
	"github.com/remogatto/prettytest"
	"testing"
	"time"
)

type testSuite struct {
	prettytest.Suite
}

type eventsLoop struct {
	pause, terminate chan int
	pauseTest        chan bool
}

type renderLoop struct {
	pause, terminate chan int
	ticker           *time.Ticker
	eglState         platform.EGLState
	testInit         chan platform.EGLState
}

func newEventsLoop() *eventsLoop {
	return &eventsLoop{
		pause:     make(chan int),
		terminate: make(chan int),
		pauseTest: make(chan bool),
	}
}

func (loop *eventsLoop) Pause() chan int {
	return loop.pause
}

func (loop *eventsLoop) Terminate() chan int {
	return loop.terminate
}

func (loop *eventsLoop) Run() {
	for {
		select {
		case <-loop.pause:
			loop.pause <- 0
		case <-loop.terminate:
			loop.terminate <- 0
		case untypedEvent := <-Events:
			switch untypedEvent.(type) {
			case PauseEvent:
				loop.pauseTest <- true
			}
		}
	}
}

func newRenderLoop(eglState platform.EGLState) *renderLoop {
	return &renderLoop{
		pause:     make(chan int),
		terminate: make(chan int),
		ticker:    time.NewTicker(time.Duration(time.Second)),
		eglState:  eglState,
		testInit:  make(chan platform.EGLState),
	}
}

func (loop *renderLoop) Pause() chan int {
	return loop.pause
}

func (loop *renderLoop) Terminate() chan int {
	return loop.terminate
}

func (loop *renderLoop) Run() {
	for {
		select {
		case <-loop.pause:
			loop.pause <- 0
		case <-loop.terminate:
			loop.terminate <- 0
		case <-loop.ticker.C:
			loop.draw()

		case <-loop.testInit:
			loop.testInit <- loop.eglState

		}
	}
}

func (loop *renderLoop) draw() {
	// do nothing
}

func (t *testSuite) BeforeAll() {
	application.Register("eventsLoop", newEventsLoop())

	// Initialize the EGL surface, should be non-blocking.
	initialize <- platform.EGLState{
		Display:       1,
		Context:       0xdeadbeef,
		Surface:       0xdeadbeef,
		SurfaceWidth:  640,
		SurfaceHeight: 480,
	}

	go application.Run()

	eglState := <-Init

	application.Register("renderLoop", newRenderLoop(eglState))
	application.Register("eventsLoop", newEventsLoop())
	application.Start("renderLoop")
	application.Start("eventsLoop")

	go func() {
		for {
			select {
			case <-application.ExitCh:
				return
			case err := <-application.ErrorCh:
				application.Logf(err.(application.Error).Error())
			}
		}
	}()

}

func (t *testSuite) TestInit() {
	loop, err := application.Loop("renderLoop")
	renderLoop, ok := loop.(*renderLoop)

	t.Nil(err)
	t.True(ok)

	renderLoop.testInit <- platform.EGLState{}
	eglState := <-renderLoop.testInit
	t.Equal(640, eglState.SurfaceWidth)
	t.Equal(480, eglState.SurfaceHeight)
}

func (t *testSuite) TestEvents() {
	loop, err := application.Loop("eventsLoop")
	eventsLoop, ok := loop.(*eventsLoop)

	t.Nil(err)
	t.True(ok)

	event <- PauseEvent{}
	t.True(<-eventsLoop.pauseTest)
}

func (t *testSuite) TestLoadAssetRequest() {
	AssetsPath = "testdata"

	request := LoadAssetRequest{
		"doesntexist",
		make(chan LoadAssetResponse),
	}
	Assets <- request
	response := <-request.Response
	t.True(response.Error != nil)

	request = LoadAssetRequest{
		"foo.txt",
		make(chan LoadAssetResponse),
	}
	Assets <- request
	response = <-request.Response
	t.True(response.Error == nil)
	t.True(response.Buffer != nil)
}

func TestGorgasm(t *testing.T) {
	prettytest.Run(t, new(testSuite))
}
