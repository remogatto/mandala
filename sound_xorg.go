// +build !android

package mandala

type apPlayRequest struct {
	buffer []byte
	doneCh chan bool
}

type audioPlayer struct{}

func newAudioPlayer() (*audioPlayer, error) {
	return &audioPlayer{}, nil
}

func (ap *audioPlayer) play(buffer []byte, doneCh chan bool) {
	// do nothing
}

func (ap *audioPlayer) setVolumeLevel(value int) error {
	return nil
}

func (ap *audioPlayer) getMaxVolumeLevel() (int, error) {
	return 0, nil
}
