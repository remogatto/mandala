package testlib

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"path/filepath"

	"github.com/remogatto/imagetest"
	"github.com/remogatto/mandala"
)

const (
	distanceThreshold = 0.02
)

// Compare the result of rendering against the saved expected image.
func testImage(filename string, act image.Image) (float64, image.Image, image.Image, error) {
	request := mandala.LoadResourceRequest{
		Filename: filepath.Join(expectedImgPath, filename),
		Response: make(chan mandala.LoadResourceResponse),
	}

	mandala.ResourceManager() <- request
	response := <-request.Response
	buffer := response.Buffer

	if response.Error != nil {
		return 1, nil, nil, response.Error
	}

	exp, err := png.Decode(bytes.NewBuffer(buffer))
	if err != nil {
		return 1, nil, nil, err
	}
	return imagetest.CompareDistance(exp, act, imagetest.Scale), exp, act, nil
}

func (t *TestSuite) TestResourceManager() {

	request := mandala.LoadResourceRequest{
		Filename: "drawable/gopher.png",
		Response: make(chan mandala.LoadResourceResponse),
	}

	mandala.ResourceManager() <- request

	response := <-request.Response
	buffer := response.Buffer

	t.True(response.Error == nil, "An error occured during resource opening")

	_, err := png.Decode(bytes.NewBuffer(buffer))
	t.True(err == nil, "An error occured during png decoding")

	// Load a non existent resource
	request = mandala.LoadResourceRequest{
		Filename: "res/doesntexist",
		Response: make(chan mandala.LoadResourceResponse),
	}

	mandala.ResourceManager() <- request

	response = <-request.Response
	buffer = response.Buffer

	t.True(buffer == nil)
	t.True(response.Error != nil)

	// Use the helper API for loading resources
	responseCh := make(chan mandala.LoadResourceResponse)
	mandala.ReadResource("drawable/gopher.png", responseCh)
	response = <-responseCh

	buffer = response.Buffer
	t.True(buffer != nil)
	t.True(response.Error == nil, "An error occured during resource opening")

	_, err = png.Decode(bytes.NewBuffer(buffer))
	t.True(err == nil, "An error occured during png decoding")
}

func (t *TestSuite) TestBasicCreationSequence() {
	// Check only after at least the first frame has been rendered
	<-t.testDraw

	exp := []string{"onCreate", "onStart", "onResume"}

	if a := t.Equal(
		len(exp),
		len(t.creationSequence),
		fmt.Sprintf("Triggered/Catched events were %v", t.creationSequence),
	); a.Passed {
		for i, exp := range []string{"onCreate", "onStart", "onResume"} {
			t.Equal(exp, t.creationSequence[i])
		}
	}
}

func (t *TestSuite) TestActionUpDown() {
	// Check only after at least the first frame has been rendered
	<-t.testDraw

	if err := Tap(100.0, 100.0); err != nil {
		panic(err)
	}

	event := <-t.testActionUpDown

	t.True(event.Down)
	t.Equal(float32(100.0), event.X)
	t.Equal(float32(100.0), event.Y)
}

func (t *TestSuite) TestActionMove() {
	// Check only after at least the first frame has been rendered
	<-t.testDraw

	// Move the cursor on the initial position (this has no effect
	// on the android-side but it's necessary during the xorg
	// test)
	if err := Move(11, 11, 11, 11); err != nil {
		panic(err)
	}

	movements := 10

	// Begin counting move events from now
	t.resetActionMove <- movements
	<-t.resetActionMove

	if err := Move(10, 10, 20, 20); err != nil {
		panic(err)
	}

	close(t.testActionMove)

	for event := range t.testActionMove {
		t.True(event.X > 0.0)
		t.True(event.Y > 0.0)
	}

	t.moving = false
}

func (t *TestSuite) TestDraw() {
	filename := GOPHER_PNG
	distance, _, _, err := testImage(filename, <-t.testDraw)
	if err != nil {
		mandala.Fatalf(err.Error())
	}
	t.True(distance < distanceThreshold, fmt.Sprintf("Image differs by distance %f", distance))
}

func (t *TestSuite) TestAudio() {
	// Open an audio file from resources
	responseCh := make(chan mandala.LoadResourceResponse)
	mandala.ReadResource("raw/android.raw", responseCh)
	response := <-responseCh
	buffer := response.Buffer
	t.True(response.Error == nil, "An error occured during resource opening")

	// Create the audio player
	player, err := mandala.NewAudioPlayer()
	t.True(err == nil)

	// Hum...this seems to fail for now
	max, err := player.GetMaxVolumeLevel()
	t.True(err == nil, "Error in getting the max volume level")
	t.True(max > 0, fmt.Sprintf("Max volume level is %d", max))

	if player != nil {
		player.Play(buffer, nil)
	}
}

func (t *TestSuite) TestBasicExitSequence() {
	t.Pending()
}
