// +build android

package mandala

type apCreateResponse struct {
	ap  *AudioPlayer
	err error
}

type apCreateRequest struct {
	filename   string
	responseCh chan apCreateResponse
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
func CreateAudioPlayer(filename string) (*AudioPlayer, error) {
	responseCh := make(chan apCreateResponse)
	soundCh <- apCreateRequest{filename, responseCh}
	response := <-responseCh
	if response.err != nil {
		return nil, response.err
	}
	return response.ap, nil
}
