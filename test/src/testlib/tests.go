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

	t.Equal(float32(100.0), event.X)
	t.Equal(float32(100.0), event.Y)
}

// func (t *TestSuite) TestActionMove() {
// 	// Check only after at least the first frame has been rendered
// 	<-t.testDraw

// 	movements := 10

// 	// Begin counting move events from now
// 	t.resetActionMove <- movements

// 	for i := float32(0); i < float32(movements); i += 1.0 {
// 		if err := Move(110.0+i, 110.0+i); err != nil {
// 			panic(err)
// 		}
// 	}

// 	close(t.testActionMove)

// 	count := 0
// 	for event := range t.testActionMove {
// 		t.Equal(float32(110.0)+float32(count), event.X)
// 		t.Equal(float32(110.0)+float32(count), event.Y)
// 		count++
// 	}
// }

func (t *TestSuite) TestDraw() {
	t.True(<-t.testDraw)
}
