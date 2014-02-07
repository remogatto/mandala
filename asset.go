// +build !android

package mandala

import (
	"io/ioutil"
	"path/filepath"
	"unsafe"
)

var (
	// The path in which the framework will search for resources.
	AssetPath string = "android/res"
)

func loadAsset(activity unsafe.Pointer, filename string) ([]byte, error) {
	// Open the file.
	buf, err := ioutil.ReadFile(filepath.Join(AssetPath, filename))
	if err != nil {
		return nil, err
	}
	return buf, nil
}
