package testlib

import (
	"fmt"
	"git.tideland.biz/goas/loop"
	"github.com/remogatto/gorgasm"
	gl "github.com/remogatto/opengles2"
	"github.com/remogatto/prettytest"
	"image"
	"image/png"
	"log"
	"runtime"
	"time"
)

const (
	FRAMES_PER_SECOND = 15
	GOPHER_PNG        = "res/drawable/gopher.png"
)

type TestSuite struct {
	prettytest.Suite
}

var (
	Width, Height                   int = 320, 480
	rlControl                       *renderLoopControl
	verticesArrayBuffer             uint32
	textureBuffer                   uint32
	unifTexture, attrPos, attrTexIn uint32
	currWidth, currHeight           int

	vertices = [24]float32{
		-1.0, -1.0, 0.0, 1.0, 0.0, 1.0,
		1.0, -1.0, 0.0, 1.0, 1.0, 1.0,
		1.0, 1.0, 0.0, 1.0, 1.0, 0.0,
		-1.0, 1.0, 0.0, 1.0, 0.0, 0.0,
	}
	vsh = `
        attribute vec4 pos;
        attribute vec2 texIn;
        varying vec2 texOut;
        void main() {
          gl_Position = pos;
          texOut = texIn;
        }
`
	fsh = `
        precision mediump float;
        varying vec2 texOut;
        uniform sampler2D texture;
	void main() {
		gl_FragColor = texture2D(texture, texOut);
	}
`
)

type viewportSize struct {
	width, height int
}

type renderLoopControl struct {
	resizeViewport chan viewportSize
	pause          chan bool
	resume         chan bool
	window         chan gorgasm.Window
	testDraw       chan bool
}

type renderState struct {
	window gorgasm.Window
}

func checkShaderCompileStatus(shader uint32) {
	var stat int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &stat)
	if stat == 0 {
		var length int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)
		infoLog := gl.GetShaderInfoLog(shader, gl.Sizei(length), nil)
		log.Fatalf("Compile error in shader %d: \"%s\"\n", shader, infoLog[:len(infoLog)-1])
	}
}

func checkProgramLinkStatus(pid uint32) {
	var stat int32
	gl.GetProgramiv(pid, gl.LINK_STATUS, &stat)
	if stat == 0 {
		var length int32
		gl.GetProgramiv(pid, gl.INFO_LOG_LENGTH, &length)
		infoLog := gl.GetProgramInfoLog(pid, gl.Sizei(length), nil)
		log.Fatalf("Link error in program %d: \"%s\"\n", pid, infoLog[:len(infoLog)-1])
	}
}

// Create a fragment shader from a string and return its reference.
func FragmentShader(s string) uint32 {
	shader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(shader, 1, &s, nil)
	gl.CompileShader(shader)
	checkShaderCompileStatus(shader)
	return shader
}

// Create a vertex shader from a string and return its reference.
func VertexShader(s string) uint32 {
	shader := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(shader, 1, &s, nil)
	gl.CompileShader(shader)
	checkShaderCompileStatus(shader)
	return shader
}

// Create a program from vertex and fragment shaders.
func Program(fsh, vsh uint32) uint32 {
	p := gl.CreateProgram()
	gl.AttachShader(p, fsh)
	gl.AttachShader(p, vsh)
	gl.LinkProgram(p)
	checkProgramLinkStatus(p)
	return p
}

func (renderState *renderState) init(window gorgasm.Window) {
	window.MakeContextCurrent()

	renderState.window = window
	width, height := window.GetSize()

	// Set the viewport
	gl.Viewport(0, 0, gl.Sizei(width), gl.Sizei(height))
	check()

	// Compile the shaders
	program := Program(FragmentShader(fsh), VertexShader(vsh))
	gl.UseProgram(program)
	check()

	// Get attributes
	attrPos = uint32(gl.GetAttribLocation(program, "pos"))
	attrTexIn = uint32(gl.GetAttribLocation(program, "texIn"))
	unifTexture = gl.GetUniformLocation(program, "texture")
	gl.EnableVertexAttribArray(attrPos)
	gl.EnableVertexAttribArray(attrTexIn)
	check()

	// Upload vertices data
	gl.GenBuffers(1, &verticesArrayBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, verticesArrayBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.SizeiPtr(len(vertices))*4, gl.Void(&vertices[0]), gl.STATIC_DRAW)
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
	gl.BindBuffer(gl.ARRAY_BUFFER, verticesArrayBuffer)
	gl.VertexAttribPointer(attrPos, 4, gl.FLOAT, false, 6*4, 0)

	// bind texture - FIX size of vertex

	gl.VertexAttribPointer(attrTexIn, 2, gl.FLOAT, false, 6*4, 4*4)

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
		make(chan gorgasm.Window),
		make(chan bool),
	}
}

// Run runs renderLoop. The loop renders a frame and swaps the buffer
// at each tick received.
func renderLoopFunc(control *renderLoopControl) loop.LoopFunc {
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
				renderState.window.SwapBuffers()
				control.testDraw <- true

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
func eventLoopFunc(renderLoopControl *renderLoopControl) loop.LoopFunc {
	return func(loop loop.Loop) error {

		for {
			select {
			// Receive events from the framework.
			case untypedEvent := <-gorgasm.Events():
				switch event := untypedEvent.(type) {

				case gorgasm.NativeWindowCreatedEvent:
					renderLoopControl.window <- event.Window

				// Finger down/up on the screen.
				case gorgasm.ActionUpDownEvent:
					if event.Down {
						gorgasm.Logf("Finger is DOWN at %f %f", event.X, event.Y)
					} else {
						gorgasm.Logf("Finger is now UP")
					}

					// Finger is moving on the screen.
				case gorgasm.ActionMoveEvent:
					gorgasm.Logf("Finger is moving at coord %f %f", event.X, event.Y)

				case gorgasm.DestroyEvent:
					gorgasm.Logf("Stop rendering...\n")
					gorgasm.Logf("Quitting from application...\n")
					return nil

				case gorgasm.NativeWindowRedrawNeededEvent:

				case gorgasm.PauseEvent:
					gorgasm.Logf("Application was paused. Stopping rendering ticker.")
					renderLoopControl.pause <- true

				case gorgasm.ResumeEvent:
					gorgasm.Logf("Application was resumed. Reactivating rendering ticker.")
					renderLoopControl.resume <- true

				}
			}
		}
	}
}

func check() {
	error := gl.GetError()
	if error != 0 {
		gorgasm.Logf("An error occurred! Code: 0x%x", error)
	}
}

func loadImage(filename string) (image.Image, error) {
	// Request an asset to the asset manager. When the app runs on
	// an Android device, the apk will be unpacked and the file
	// will be read from it and copied to a byte buffer.
	request := gorgasm.LoadAssetRequest{
		filename,
		make(chan gorgasm.LoadAssetResponse),
	}
	gorgasm.AssetManager() <- request
	response := <-request.Response

	if response.Error != nil {
		return nil, response.Error
	}

	// Decode the image.
	img, err := png.Decode(response.Buffer)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (t *TestSuite) BeforeAll() {
	gorgasm.Verbose = true
	gorgasm.Debug = true

	// Create rendering loop control channels
	rlControl = newRenderLoopControl()
	// Start the rendering loop
	loop.GoRecoverable(
		renderLoopFunc(rlControl),
		func(rs loop.Recoverings) (loop.Recoverings, error) {
			for _, r := range rs {
				gorgasm.Logf("%s", r.Reason)
				gorgasm.Logf("%s", gorgasm.Stacktrace())
			}
			return rs, fmt.Errorf("Unrecoverable loop\n")
		},
	)
	// Start the event loop
	loop.GoRecoverable(
		eventLoopFunc(rlControl),
		func(rs loop.Recoverings) (loop.Recoverings, error) {
			for _, r := range rs {
				gorgasm.Logf("%s", r.Reason)
				gorgasm.Logf("%s", gorgasm.Stacktrace())
			}
			return rs, fmt.Errorf("Unrecoverable loop\n")
		},
	)

}

func (t *TestSuite) TestDraw() {
	t.True(<-rlControl.testDraw)
}
