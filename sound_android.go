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

type assetPlayer struct {
	fdPlayerObject C.SLObjectItf
	fdPlayerPlay   C.SLPlayItf
}

type AudioPlayer struct {
	filename    string
	assetPlayer *assetPlayer
	playCh      chan apPlayRequest
	stopCh      chan apStopRequest
}

func initSound() {
	C.createEngine(nil, nil)
}

func newAudioPlayer(activity *C.ANativeActivity, filename string) *AudioPlayer {
	ap := new(AudioPlayer)
	ap.filename = filename
	ap.assetPlayer = new(assetPlayer)

	cstring := C.CString(filename)
	defer C.free(unsafe.Pointer(cstring))

	C.createAssetAudioPlayer(
		activity,
		(*C.t_asset_ap)(ap.assetPlayer),
		cstring,
	)

	ap.playCh = make(chan apPlayRequest)
	ap.stopCh = make(chan apStopRequest)

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

	return ap
}

func (ap *AudioPlayer) requestLoopFunc() loop.LoopFunc {
	return func(l loop.Loop) error {
		for {
			select {
			case request := <-ap.playCh:
				ap.play()
				request.done <- true

			case <-ap.stopCh:
			}
		}
	}
}

func (ap *AudioPlayer) play() {
	C.setPlayingAssetAudioPlayer(ap.assetPlayer.fdPlayerPlay, C.JNI_TRUE)
}

func (ap *AudioPlayer) stop() {
	C.setPlayingAssetAudioPlayer(ap.assetPlayer.fdPlayerPlay, C.JNI_FALSE)
}

// The loop handles native sound events.
func androidSoundLoopFunc(activity *C.ANativeActivity, event chan interface{}) loop.LoopFunc {
	return func(l loop.Loop) error {

		// Initialize OpenSL
		initSound()

		for {
			select {
			case untypedEvent := <-event:
				switch event := untypedEvent.(type) {
				case apCreateRequest:
					event.apCh <- newAudioPlayer(activity, event.filename)
				}
			}
		}
	}
}
