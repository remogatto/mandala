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
	request := mandala.LoadAssetRequest{
		Filename: filepath.Join(expectedImgPath, filename),
		Response: make(chan mandala.LoadAssetResponse),
	}

	mandala.AssetManager() <- request
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

func (t *TestSuite) TestAssetManagerLoadResponse() {

	request := mandala.LoadAssetRequest{
		Filename: "res/drawable/gopher.png",
		Response: make(chan mandala.LoadAssetResponse),
	}

	mandala.AssetManager() <- request

	response := <-request.Response
	buffer := response.Buffer

	t.True(response.Error == nil, "An error occured during resource opening")

	_, err := png.Decode(bytes.NewBuffer(buffer))
	t.True(err == nil, "An error occured during png decoding")

	// Load a non existent resource
	request = mandala.LoadAssetRequest{
		Filename: "res/doesntexist",
		Response: make(chan mandala.LoadAssetResponse),
	}

	mandala.AssetManager() <- request

	response = <-request.Response
	buffer = response.Buffer

	t.True(buffer == nil)
	t.True(response.Error != nil)
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
	player := mandala.CreateAudioPlayer("background.mp3")
	t.True(player != nil)

	done := make(chan bool)
	player.Play(done)
	<-done
	player.Stop()
}

func (t *TestSuite) TestBasicExitSequence() {
	t.Pending()
}
