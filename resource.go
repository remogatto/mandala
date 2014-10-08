// +build !android

package mandala

import (
	"io/ioutil"
	"path/filepath"
	"unsafe"
)

func loadResource(activity unsafe.Pointer, filename string) ([]byte, error) {
	// Open the file.
	buf, err := ioutil.ReadFile(filepath.Join(ResourcePath, filename))
	if err != nil {
		return nil, err
	}
	return buf, nil
}
