// +build android

package mandala

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"unsafe"
)

// #include <android/native_activity.h>
// #include "asset_android.h"
import "C"

func loadAsset(activity unsafe.Pointer, filename string) (io.Reader, error) {
	apkPath := C.GoString(C.getPackageName((*C.ANativeActivity)(activity)))

	// Open a zip archive for reading.
	r, err := zip.OpenReader(apkPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// Iterate through the files in the archive.
	for _, f := range r.File {
		if f.Name == filename {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			buffer := make([]byte, f.UncompressedSize64)
			_, err = io.ReadFull(rc, buffer)
			if err != nil {
				return nil, err
			}
			rc.Close()
			return bytes.NewBuffer(buffer), nil
		}
	}
	return nil, fmt.Errorf("Resource not found!")
}
