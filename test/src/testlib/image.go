package testlib

import (
	"image"
	"image/png"
	"os"
)

// Save the given image to file.
func SaveImage(filename string, img image.Image) error {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return err
	}
	err = png.Encode(f, img)
	if err != nil {
		return err
	}
	return nil
}
