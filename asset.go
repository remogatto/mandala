// +build !android

package mandala

import (
	"io"
	"os"
	"path/filepath"
	"unsafe"
)

func loadAsset(activity unsafe.Pointer, filename string) (io.ReadCloser, error) {
	// Open the file.
	file, err := os.Open(filepath.Join(AssetPath, filename))
	if err != nil {
		return nil, err
	}
	return file, nil
}
