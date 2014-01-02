// +build android

package gorgasm

// #include <android/log.h>
// #cgo LDFLAGS: -llog
import "C"

import (
	"bytes"
	"log"
	"unsafe"
)

var ctag *C.char = C.CString("Gorgasm")

func init() {
	log.SetOutput(AndroidWriter{})
}

type AndroidWriter []byte

func (buf AndroidWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	err = nil
	for nlidx := bytes.IndexByte(p, '\n'); nlidx != -1; nlidx = bytes.IndexByte(p, '\n') {
		buf = append(buf, p[:nlidx]...)
		p = p[nlidx+1:]
		buf = append(buf, 0)
		cstr := (*C.char)(unsafe.Pointer(&buf[0]))
		C.__android_log_write(C.ANDROID_LOG_INFO, ctag, cstr)
		buf = buf[:0]
	}
	buf = append(buf, p...)
	return
}
