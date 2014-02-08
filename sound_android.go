// +build android

package mandala

import (
	"fmt"
	"unsafe"
	"git.tideland.biz/goas/loop"
)

// #include <android/native_activity.h>
// #include <SLES/OpenSLES.h>
// #include <SLES/OpenSLES_Android.h>
// #include "sound_android.h"
//
// #cgo LDFLAGS: -landroid -lOpenSLES
import "C"

type apPlayRequest struct {
	buffer []byte
	doneCh chan bool
}

type bufferQueuePlayer struct {
	bqPlayerObject      C.SLObjectItf
	bqPlayerPlay        C.SLPlayItf
	bqPlayerBufferQueue C.SLAndroidSimpleBufferQueueItf
	bqPlayerVolume      C.SLVolumeItf
}

type audioPlayer struct {
	bqPlayer *bufferQueuePlayer
	playCh   chan apPlayRequest
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
	C.shutdownOpenSL()
}

func newAudioPlayer() (*audioPlayer, error) {
	ap := new(audioPlayer)
	ap.bqPlayer = new(bufferQueuePlayer)
	result := C.createBufferQueueAudioPlayer((*C.t_buffer_queue_ap)(ap.bqPlayer))

	if result != C.SL_RESULT_SUCCESS {
		return nil, fmt.Errorf("Error %d occured trying to create a buffer queue player", result)
	}

	ap.playCh = make(chan apPlayRequest)

	loop.GoRecoverable(
		ap.requestLoopFunc(),
		func(rs loop.Recoverings) (loop.Recoverings, error) {
			for _, r := range rs {
				Logf("%s", r.Reason)
				Logf("%s", Stacktrace())
			}
			return rs, fmt.Errorf("Unrecoverable loop\n")
		},
	)

	return ap, nil
}

func (ap *audioPlayer) requestLoopFunc() loop.LoopFunc {
	return func(l loop.Loop) error {
		for {
			select {
			case request := <-ap.playCh:
				ap.enqueue(request.buffer)
				if request.doneCh != nil {
					request.doneCh <- true
				}
			}
		}
	}
}

func (ap *audioPlayer) play(buffer []byte, doneCh chan bool) {
	ap.playCh <- apPlayRequest{buffer: buffer, doneCh: doneCh}
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
