package gorgasm

// Window is an interface that abstracts native EGL surfaces.
type Window interface {
	SwapBuffers()
	MakeContextCurrent()
	GetSize() (int, int)
}
