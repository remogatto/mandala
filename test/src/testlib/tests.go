package testlib

import (
	"fmt"
)

func (t *TestSuite) TestBasicCreationSequence() {
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

func (t *TestSuite) TestDraw() {
	t.True(<-t.testDraw)
}
