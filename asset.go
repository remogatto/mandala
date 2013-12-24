// +build !android

package gorgasm

import (
	"io"
	"os"
	"path/filepath"
	"unsafe"
)

func loadAsset(activity unsafe.Pointer, filename string) (io.ReadCloser, error) {
	// Open the file.
	file, err := os.Open(filepath.Join(AssetsPath, filename))
	if err != nil {
		return nil, err
	}
	return file, nil
}
