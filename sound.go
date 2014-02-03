package mandala

type apCreateRequest struct {
	filename string
	apCh     chan *AudioPlayer
}

type apPlayRequest struct {
	done chan bool
}

type apStopRequest struct {
}

// GetFilename returns the filename used to instantiate the audio
// player.
func (ap *AudioPlayer) GetFilename() string {
	return ap.filename
}

// Play starts playing the sound.
func (ap *AudioPlayer) Play(done chan bool) {
	ap.playCh <- apPlayRequest{done}
}

// Stop stops the sound.
func (ap *AudioPlayer) Stop() {
	ap.stopCh <- apStopRequest{}
}

// CreateAudioPlayer instantiates a player for the given filename.
func CreateAudioPlayer(filename string) *AudioPlayer {
	audioPlayerCh := make(chan *AudioPlayer)
	soundCh <- apCreateRequest{filename, audioPlayerCh}
	return <-audioPlayerCh
}
