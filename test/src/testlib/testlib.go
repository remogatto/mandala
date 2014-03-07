package testlib

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"path/filepath"
	"runtime"
	"time"

	"git.tideland.biz/goas/loop"
	"github.com/remogatto/imagetest"
	"github.com/remogatto/mandala"
	gl "github.com/remogatto/opengles2"
	"github.com/remogatto/prettytest"
	"github.com/remogatto/shaders"
)

const (
	FRAMES_PER_SECOND = 15
	GOPHER_PNG        = "gopher.png"
	TIMEOUT           = time.Second * 30
	expectedImgPath   = "drawable"
)

type TestSuite struct {
	prettytest.Suite

	rlControl        *renderLoopControl
	creationSequence []string
	exitSequence     []string
	moving           bool
	resetActionMove  chan int
	timeout          <-chan time.Time

	testDraw         chan image.Image
	testPause        chan bool
	testActionUpDown chan mandala.ActionUpDownEvent
	testActionMove   chan mandala.ActionMoveEvent
}

var (
	Width, Height                   int = 320, 480
	textureBuffer                   uint32
	unifTexture, attrPos, attrTexIn uint32
	currWidth, currHeight           int

	vertices = [24]float32{
		-1.0, -1.0, 0.0, 1.0, 0.0, 1.0,
		1.0, -1.0, 0.0, 1.0, 1.0, 1.0,
		1.0, 1.0, 0.0, 1.0, 1.0, 0.0,
		-1.0, 1.0, 0.0, 1.0, 0.0, 0.0,
	}
	vsh = shaders.VertexShader(`
        attribute vec4 pos;
        attribute vec2 texIn;
        varying vec2 texOut;
        void main() {
          gl_Position = pos;
          texOut = texIn;
        }`)
	fsh = shaders.FragmentShader(`
        precision mediump float;
        varying vec2 texOut;
        uniform sampler2D texture;
	void main() {
		gl_FragColor = texture2D(texture, texOut);
	}`)
)

type viewportSize struct {
	width, height int
}

type renderLoopControl struct {
	resizeViewport chan viewportSize
	pause          chan bool
	resume         chan bool
	window         chan mandala.Window
}

type renderState struct {
	window mandala.Window
}

// TestImage compares a saved expected image of a given filename with
// an actual image.Image that typically contains the result of a
// rendering. It returns the distance value and the two compared
// images.
func TestImage(filename string, act image.Image, adjust imagetest.Adjuster) (float64, image.Image, image.Image, error) {
	request := mandala.LoadResourceRequest{
		Filename: filepath.Join(expectedImgPath, filename),
		Response: make(chan mandala.LoadResourceResponse),
	}

	mandala.ResourceManager() <- request
	response := <-request.Response
	buffer := response.Buffer

	if response.Error != nil {
		return 1, nil, nil, response.Error
	}

	exp, err := png.Decode(bytes.NewBuffer(buffer))
	if err != nil {
		return 1, nil, nil, err
	}
	return imagetest.CompareDistance(exp, act, adjust), exp, act, nil
}

func (renderState *renderState) init(window mandala.Window) {
	window.MakeContextCurrent()

	renderState.window = window
	width, height := window.GetSize()

	// Set the viewport
	gl.Viewport(0, 0, gl.Sizei(width), gl.Sizei(height))
	check()

	// Compile the shaders
	program := shaders.NewProgram(fsh, vsh)
	program.Use()
	check()

	// Get attributes
	attrPos = program.GetAttribute("pos")
	attrTexIn = program.GetAttribute("texIn")
	unifTexture = program.GetUniform("texture")
	gl.EnableVertexAttribArray(attrPos)
	gl.EnableVertexAttribArray(attrTexIn)
	check()

	// Upload texture data
	img, err := loadImage(GOPHER_PNG)
	if err != nil {
		panic(err)
	}

	// Prepare the image to be placed on a texture.
	bounds := img.Bounds()
	imgWidth, imgHeight := bounds.Size().X, bounds.Size().Y
	buffer := make([]byte, imgWidth*imgHeight*4)
	index := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			buffer[index] = byte(r)
			buffer[index+1] = byte(g)
			buffer[index+2] = byte(b)
			buffer[index+3] = byte(a)
			index += 4
		}
	}

	gl.GenTextures(1, &textureBuffer)
	gl.BindTexture(gl.TEXTURE_2D, textureBuffer)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gl.Sizei(imgWidth), gl.Sizei(imgHeight), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Void(&buffer[0]))
	check()

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}

func (renderState *renderState) draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.VertexAttribPointer(attrPos, 4, gl.FLOAT, false, 6*4, &vertices[0])
	gl.VertexAttribPointer(attrTexIn, 2, gl.FLOAT, false, 6*4, &vertices[4])

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, textureBuffer)
	gl.Uniform1i(int32(unifTexture), 0)

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
	gl.Flush()
	gl.Finish()
}

func newRenderLoopControl() *renderLoopControl {
	return &renderLoopControl{
		make(chan viewportSize),
		make(chan bool),
		make(chan bool),
		make(chan mandala.Window),
	}
}

// Run runs renderLoop. The loop renders a frame and swaps the buffer
// at each tick received.
func (t *TestSuite) renderLoopFunc(control *renderLoopControl) loop.LoopFunc {
	return func(loop loop.Loop) error {

		// renderState stores rendering state variables such
		// as the EGL state
		renderState := new(renderState)

		// Lock/unlock the loop to the current OS thread. This is
		// necessary because OpenGL functions should be called from
		// the same thread.
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		// Create an instance of ticker and immediately stop
		// it because we don't want to swap buffers before
		// initializing a rendering state.
		ticker := time.NewTicker(time.Duration(1e9 / int(FRAMES_PER_SECOND)))
		ticker.Stop()

		for {
			select {
			case window := <-control.window:
				ticker.Stop()
				renderState.init(window)
				ticker = time.NewTicker(time.Duration(1e9 / int(FRAMES_PER_SECOND)))

			// At each tick render a frame and swap buffers.
			case <-ticker.C:
				renderState.draw()
				t.testDraw <- Screenshot(renderState.window)
				renderState.window.SwapBuffers()

			case <-control.resizeViewport:

			case <-control.pause:
				// store something
				ticker.Stop()

			case <-control.resume:
				// resume something

			case <-loop.ShallStop():
				ticker.Stop()
				return nil
			}
		}
	}
}

// eventLoopFunc is listening for events originating from the
// framwork.
func (t *TestSuite) eventLoopFunc(renderLoopControl *renderLoopControl) loop.LoopFunc {
	return func(loop loop.Loop) error {

		for {
			select {

			case c := <-t.resetActionMove:
				t.testActionMove = make(chan mandala.ActionMoveEvent, c)
				t.moving = true
				t.resetActionMove <- 0

			// Receive events from the framework.
			case untypedEvent := <-mandala.Events():

				switch event := untypedEvent.(type) {

				case mandala.CreateEvent:
					t.creationSequence = append(t.creationSequence, "onCreate")

				case mandala.StartEvent:
					t.creationSequence = append(t.creationSequence, "onStart")

				case mandala.NativeWindowCreatedEvent:
					renderLoopControl.window <- event.Window

					// Finger down/up on the screen.
				case mandala.ActionUpDownEvent:
					if !t.moving {
						t.testActionUpDown <- event
					}

					// Finger is moving on the screen.
				case mandala.ActionMoveEvent:
					if t.moving {
						t.testActionMove <- event
					}

				case mandala.NativeWindowDestroyedEvent:
					mandala.Debugf("Window destroyed")

				case mandala.DestroyEvent:
					// return nil

				case mandala.NativeWindowRedrawNeededEvent:
					mandala.Debugf("Redraw needed")

				case mandala.PauseEvent:

				case mandala.ResumeEvent:
					t.creationSequence = append(t.creationSequence, "onResume")

				}
			}
		}
	}
}

func (t *TestSuite) timeoutLoopFunc() loop.LoopFunc {
	return func(loop loop.Loop) error {
		time := <-t.timeout
		err := fmt.Errorf("Tests timed out after %v", time)
		mandala.Logf("%s %s", err.Error(), mandala.Stacktrace())
		t.Error(err)
		return nil
	}
}

func check() {
	error := gl.GetError()
	if error != 0 {
		mandala.Logf("An error occurred! Code: 0x%x", error)
	}
}

func loadImage(filename string) (image.Image, error) {
	// Request an resource to the resource manager. When the app runs on
	// an Android device, the apk will be unpacked and the file
	// will be read from it and copied to a byte buffer.
	request := mandala.LoadResourceRequest{
		Filename: "drawable/gopher.png",
		Response: make(chan mandala.LoadResourceResponse),
	}

	mandala.ResourceManager() <- request
	response := <-request.Response
	buffer := response.Buffer

	if response.Error != nil {
		return nil, response.Error
	}

	// Decode the image.
	img, err := png.Decode(bytes.NewBuffer(buffer))
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (t *TestSuite) BeforeAll() {
	// Create rendering loop control channels
	t.rlControl = newRenderLoopControl()
	// Start the rendering loop
	loop.GoRecoverable(
		t.renderLoopFunc(t.rlControl),
		func(rs loop.Recoverings) (loop.Recoverings, error) {
			for _, r := range rs {
				mandala.Logf("%s", r.Reason)
				mandala.Logf("%s", mandala.Stacktrace())
			}
			return rs, fmt.Errorf("Unrecoverable loop\n")
		},
	)
	// Start the event loop
	loop.GoRecoverable(
		t.eventLoopFunc(t.rlControl),
		func(rs loop.Recoverings) (loop.Recoverings, error) {
			for _, r := range rs {
				mandala.Logf("%s", r.Reason)
				mandala.Logf("%s", mandala.Stacktrace())
			}
			return rs, fmt.Errorf("Unrecoverable loop\n")
		},
	)

	// Start the timeout loop
	loop.GoRecoverable(
		t.timeoutLoopFunc(),
		func(rs loop.Recoverings) (loop.Recoverings, error) {
			for _, r := range rs {
				mandala.Logf("%s", r.Reason)
				mandala.Logf("%s", mandala.Stacktrace())
			}
			return rs, fmt.Errorf("Unrecoverable loop\n")
		},
	)

}

func (t *TestSuite) AfterAll() {
	// os.Exit(0)
}

func NewTestSuite() *TestSuite {
	return &TestSuite{
		rlControl:        newRenderLoopControl(),
		resetActionMove:  make(chan int),
		testDraw:         make(chan image.Image),
		testPause:        make(chan bool),
		testActionUpDown: make(chan mandala.ActionUpDownEvent),
		timeout:          time.After(TIMEOUT),
	}
}
