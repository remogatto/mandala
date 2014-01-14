package mandala

// Window is an interface that abstracts native EGL surfaces.
type Window interface {

	// Swap display <-> surface buffers.
	SwapBuffers()

	// Bind context to the current rendering thread.
	MakeContextCurrent()

	// Get the size of the window as width,height values.
	GetSize() (int, int)
}
