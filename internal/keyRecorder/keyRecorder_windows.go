package keyRecorder

import (
	"bufio"
	"fmt"
	"golang.org/x/sys/windows"
	"time"
)

var (
	user32 = windows.NewLazySystemDLL("user32.dll")
	// See: https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getasynckeystate
	getAsyncKeyStateProc = user32.NewProc("GetAsyncKeyState")
)

type keyStateChange struct {
	vKey      uint8
	down      bool
	timestamp int64
}

func isKeyDown(vKey uint8) bool {
	down, _, _ := getAsyncKeyStateProc.Call(uintptr(vKey))
	return down != 0
}

func getKeyStates() map[uint8]bool {
	states := map[uint8]bool{}

	// See: https://docs.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes
	// Not all of these are actually valid keys, but that's fine.
	for vKey := uint8(1); vKey < 255; vKey++ {
		states[vKey] = isKeyDown(vKey)
	}

	return states
}

func unixMilli() int64 {
	return time.Now().UnixNano() / 1_000_000
}

func watchStateChanges(changes chan keyStateChange, stopKey uint8) {
	var lastState map[uint8]bool
	for {
		// Get the new key states, compare them.
		newStates := getKeyStates()
		for vKey, down := range newStates {
			// Check if the stop key was pressed.
			if vKey == stopKey && down {
				close(changes)
				return
			}

			// Check if the key state changed.
			if down != lastState[vKey] {
				changes <- keyStateChange{
					vKey:      vKey,
					down:      down,
					timestamp: unixMilli(),
				}
			}
		}

		lastState = newStates
		time.Sleep(1 * time.Millisecond)
	}
}

func RecordKeys(stopKey uint8, w *bufio.Writer, ignore []int) {
	// Create a map from the ignore values for easy lookup.
	ignoreMap := make(map[uint8]bool, len(ignore))
	for _, v := range ignore {
		ignoreMap[uint8(v)] = true
	}

	// Make a first call to GetAsyncKeyState and discard the result,
	// so all next calls return clean results.
	getKeyStates()

	changes := make(chan keyStateChange)
	go watchStateChanges(changes, stopKey)

	w.WriteString("timestamps {\n")
	startMillis := int64(0)
	for change := range changes {
		// Check if the key should be ignored.
		if _, ok := ignoreMap[change.vKey]; ok {
			continue
		}

		// If this is the first event, set the start time.
		if startMillis == 0 {
			startMillis = change.timestamp
		}

		// Current timestamp.
		curMillis := change.timestamp - startMillis

		// Key operation.
		keyOp := "Up"
		if change.down {
			keyOp = "Down"
		}
		w.WriteString(fmt.Sprintf("\t%d vKey%s %d\n", curMillis, keyOp, change.vKey))
	}
	w.WriteString("}\n")
}
