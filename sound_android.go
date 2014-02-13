// +build android

package mandala

import (
	"fmt"
	"sync"
	"unsafe"
)

// #include <android/native_activity.h>
// #include <SLES/OpenSLES.h>
// #include <SLES/OpenSLES_Android.h>
// #include "sound_android.h"
//
// #cgo LDFLAGS: -landroid -lOpenSLES
import "C"

type bufferQueuePlayer struct {
	bqPlayerObject      C.SLObjectItf
	bqPlayerPlay        C.SLPlayItf
	bqPlayerBufferQueue C.SLAndroidSimpleBufferQueueItf
	bqPlayerVolume      C.SLVolumeItf
}

type audioPlayer struct {
	rwMutex  sync.RWMutex
	bqPlayer *bufferQueuePlayer
}

func jniBool(value C.jboolean) bool {
	if value == C.JNI_TRUE {
		return true
	}
	return false
}

func initOpenSL() error {
	result := C.initOpenSL()
	if result != C.SL_RESULT_SUCCESS {
		return fmt.Errorf("Unable to initialize the native audio library. Error code is %x", int(result))
	}
	return nil
}

func shutdownOpenSL() {
	Debugf("Shutting down OpenSL ES")
	C.shutdownOpenSL()
}

func newAudioPlayer() (*audioPlayer, error) {
	ap := new(audioPlayer)
	ap.bqPlayer = new(bufferQueuePlayer)
	result := C.createBufferQueueAudioPlayer((*C.t_buffer_queue_ap)(ap.bqPlayer))
	if result != C.SL_RESULT_SUCCESS {
		return nil, fmt.Errorf("Error %d occured trying to create a buffer queue player", result)
	}
	return ap, nil
}

func (ap *audioPlayer) play(buffer []byte, doneCh chan bool) {
	ap.enqueue(buffer)
}

func (ap *audioPlayer) destroy() {
	Debugf("Destroying audio player at 0x%x", ap)
	ap.rwMutex.Lock()
	C.destroyBufferQueueAudioPlayer((*C.t_buffer_queue_ap)(ap.bqPlayer))
	ap.rwMutex.Unlock()
}

func (ap *audioPlayer) enqueue(buffer []byte) {
	C.enqueueBuffer((*C.t_buffer_queue_ap)(ap.bqPlayer), unsafe.Pointer(&buffer[0]), C.SLuint32(len(buffer)))
}

func (ap *audioPlayer) setVolumeLevel(value int) error {
	result := C.setVolumeLevel((*C.t_buffer_queue_ap)(ap.bqPlayer), C.SLmillibel(value))
	if result != C.SL_RESULT_SUCCESS {
		return fmt.Errorf("Unable to set volume level. Error code is %x", int(result))
	}
	return nil
}

func (ap *audioPlayer) getMaxVolumeLevel() (int, error) {
	var maxLevel C.SLmillibel
	result := C.getMaxVolumeLevel((*C.t_buffer_queue_ap)(ap.bqPlayer), &maxLevel)
	if result != C.SL_RESULT_SUCCESS {
		return 0, fmt.Errorf("Unable to get max volume level. Error code is %x", int(result))
	}
	return int(maxLevel), nil
}

//export playerCallback
func playerCallback() {
	Logf("Player done!\n")
}
