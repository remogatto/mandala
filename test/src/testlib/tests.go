package testlib

import (
	"fmt"
	"image/png"

	"github.com/remogatto/mandala"
)

func (t *TestSuite) TestAssetManagerLoadResponse() {

	request := mandala.LoadAssetRequest{
		Filename: "res/drawable/gopher.png",
		Response: make(chan mandala.LoadAssetResponse),
	}

	mandala.AssetManager() <- request

	response := <-request.Response
	buffer := response.Buffer

	t.True(response.Error == nil, "An error occured during resource opening")

	_, err := png.Decode(buffer)
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
	t.True(<-t.testDraw)
}

// func (t *TestSuite) TestBasicExitSequence() {
// 	if err := Back(); err != nil {
// 		panic(err)
// 	}
// 	<-t.testPause
// 	exp := []string{"onPause"}

// 	if a := t.Equal(
// 		len(exp),
// 		len(t.creationSequence),
// 		fmt.Sprintf("Triggered/Catched events were %v", t.exitSequence),
// 	); a.Passed {
// 		for i, exp := range []string{"onPause"} {
// 			t.Equal(exp, t.creationSequence[i])
// 		}
// 	}
// }
