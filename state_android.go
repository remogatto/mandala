// +build android

package gorgasm

import (
	"sync"
)

// #include <android/native_activity.h>
import "C"

var (
	states map[*C.ANativeActivity]*state
	// Global mutexes.
	rwMutex sync.RWMutex
	mutex   sync.Mutex
)

type state struct {
	nativeActivity *C.ANativeActivity
	window         *window
}

func getState(act *C.ANativeActivity) *state {
	mutex.Lock()
	defer mutex.Unlock()
	state, ok := states[act]
	if !ok {
		Fatalf("%s", "Invalid activity reference")
	}
	return state
}

func setState(act *C.ANativeActivity, state *state) *state {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	states[act] = state
	return state
}

func deleteState(act *C.ANativeActivity, state *state) {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	delete(states, act)
}

func init() {
	states = make(map[*C.ANativeActivity]*state)
}
