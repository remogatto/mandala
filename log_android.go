// +build android

package mandala

// #include <android/log.h>
// #cgo LDFLAGS: -llog
import "C"

import (
	"bytes"
	"log"
	"unsafe"
)

const (
	// The default logcat tag for the framework
	MandalaLogcatTag = "Mandala"
)

func init() {
	log.SetOutput(AndroidWriter{Tag: MandalaLogcatTag})
}

type AndroidWriter struct {
	buf []byte
	Tag string
}

func (w AndroidWriter) Write(p []byte) (n int, err error) {
	ctag := C.CString(w.Tag)
	n = len(p)
	err = nil
	for nlidx := bytes.IndexByte(p, '\n'); nlidx != -1; nlidx = bytes.IndexByte(p, '\n') {
		w.buf = append(w.buf, p[:nlidx]...)
		p = p[nlidx+1:]
		w.buf = append(w.buf, 0)
		cstr := (*C.char)(unsafe.Pointer(&w.buf[0]))
		C.__android_log_write(C.ANDROID_LOG_INFO, ctag, cstr)
		w.buf = w.buf[:0]
	}
	w.buf = append(w.buf, p...)
	return
}
