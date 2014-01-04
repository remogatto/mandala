package testlib

import (
	"fmt"
)

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

	movements := 10

	// Begin counting move events from now
	t.resetActionMove <- movements

	if err := Move(10, 10, 20, 20); err != nil {
		panic(err)
	}

	close(t.testActionMove)

	count := 0
	for event := range t.testActionMove {
		t.Equal(float32(10.0)+float32(count), event.X)
		t.Equal(float32(10.0)+float32(count), event.Y)
		count++
	}

	t.testActionMove = nil
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
