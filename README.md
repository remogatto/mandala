# What's that?

Gorgasm is a framework for writing Android native, hardware
accelerated, concurrent applications in Go. You can develop, test and
run your application on your (xorg) desktop and then compile and
deploy it to an Android device.

# How does it work?

Gorgasm use [Goandroid](https://github.com/) toolchain to compile Go
applications for Android. The graphics abstraction is obtained using
EGL.

The framework provides an Events channel from which client code listen
for events happening during program execution. Examples of events are
the interaction with the screen, pausing/resuming of the application,
etc. The framework abstracts the Android native events providing a way
to build, run and test the application on the desktop with the promise
that it will behave the same on the device once deployed. Oh well,
this is the aim of the framework, at least!

A typical application has two loops: one for listening to events, one
for rendering.

~~~go
// ...

func (l *renderLoop) Run() {
	for {
		select {
			// Pause the loop
		case <-l.pause:
			l.ticker.Stop()
			l.pause <- 0

			// Terminate the loop
		case <-l.terminate:
			l.terminate <- 0

		case <-l.resume:
			// resume application state

		case l.eglState = <-gorgasm.Init:
			if ok := egl.MakeCurrent(
				l.eglState.Display,
				l.eglState.Surface,
				l.eglState.Surface,
				l.eglState.Context,
			); !ok {
			       // Handle errors
 			}

 			// setup viewport, textures, etc.
			initGL()

			// Start the ticker for rendering at a given frame rate.
			l.ticker = time.NewTicker(time.Duration(1e9 / int(FRAMES_PER_SECOND)))

			// At each tick render the frame.
		case <-l.ticker.C:
			draw()
		}
	}
}

func (l *eventsLoop) Run() {
	for {
		select {
		case <-l.pause:
			l.pause <- 0
		case <-l.terminate:
			l.terminate <- 0

			// Receive events from the framework.
		case untypedEvent := <-gorgasm.Events:
			switch event := untypedEvent.(type) {
				// Finger down/up on the screen.
			case gorgasm.ActionUpDownEvent:
				// do something

				// Finger is moving on the screen.
			case gorgasm.ActionMoveEvent:
				// do something

				// Application was paused.
			case gorgasm.PauseEvent:
				l.renderLoop.pause <- 1
				
				// Application was resumed.
			case gorgasm.ResumeEvent:
				l.renderLoop.resume <- 1
			}
		}
	}
}

// ...
~~~

In order to dealing with application resources (images, sounds,
configuration files, etc.), the framework provides an Assets
channel. Client code sends request to Assets in order to obtain
resources as <pre>io.Reader</pre> instances. In the desktop
application this simply means opening the file at the given path. In
the Android application the framework will unpack the apk archive on
the fly getting the requested resources from it. However, is the
framework responsibility to deal with the right native method for
opening file. From the client-code point of view the the request will
be the same, for example:

~~~go
// ...

// Request an asset to the asset manager. When the app runs on
// an Android device, the apk will be unpacked and the file
// will be read from it and copied to a byte buffer.
request := gorgasm.LoadAssetRequest{
	"foo.png",
	make(chan gorgasm.LoadAssetResponse),
}

gorgasm.Assets <- request
response := <-request.Response

if response.Error != nil {
	return nil, response.Error
}

// Decode the image.
img, err := png.Decode(response.Buffer)
if err != nil {
	return nil, err
}

// ...
~~~

# Prerequisites

* Android NDK
* Goandroid
* EGL (on a debian-like system: <pre></pre>)

# Install

<pre>
go get https://github.com/remogatto/gorgasm
</pre>

# Quick start

Once you have a working environment, in order to create a basic
application,

<pre>
$ go install github.com/remogatto/gotask
$ git clone https://github.com/remogatto/gorgasm-environment
$ cd gorgasm-environment
$ gotask create HelloAndroid
</pre>

See gorgasm-environment for furher info.

# License

See LICENSE.md
