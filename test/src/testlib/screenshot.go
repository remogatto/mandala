package testlib

import (
	"image"
	"image/color"

	"github.com/remogatto/mandala"
	gl "github.com/remogatto/opengles2"
)

// Screenshot takes a screenshot of the current window and return a
// RGBA image.Image.
func Screenshot(window mandala.Window) image.Image {
	width, height := window.GetSize()

	// Allocate the pixel buffer
	pixels := make([]byte, width*height*4)
	gl.PixelStorei(gl.PACK_ALIGNMENT, 1)

	// Read the framebuffer
	gl.ReadPixels(0, 0, gl.Sizei(width), gl.Sizei(height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Void(&pixels[0]))

	// Create a RGBA image
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)

	index := 0
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			color := color.RGBA{
				pixels[index],
				pixels[index+1],
				pixels[index+2],
				pixels[index+3],
			}
			img.Set(x, rect.Max.Y-y, color)
			index += 4
		}
	}

	return img
}
