package mandala

type AudioPlayer struct {
	// The native instance of the audio player
	ap *audioPlayer
}

// CreateAudioPlayer instantiates a player for the given filename.
func NewAudioPlayer() (*AudioPlayer, error) {
	ap, err := newAudioPlayer()
	if err != nil {
		return nil, err
	}
	return &AudioPlayer{ap}, nil
}

// Play tells the player to play the named track and send a value to
// doneCh when done. The channel can be nil, in that case nothing is
// sent to it.
func (ap *AudioPlayer) Play(buffer []byte, doneCh chan bool) {
	ap.ap.play(buffer, doneCh)
}

// GetVolumeScale returns the [min,max] values for volume. If the
// device doesn't support volume controls, it returns an error.
func (ap *AudioPlayer) GetMaxVolumeLevel() (int, error) {
	return ap.ap.getMaxVolumeLevel()
}

// SetVolume sets the volume for the player.
func (ap *AudioPlayer) SetVolumeLevel(value int) error {
	return ap.ap.setVolumeLevel(value)
}

func (ap *AudioPlayer) Destroy() {
	ap.ap.destroy()
}
