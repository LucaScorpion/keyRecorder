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
	vKey uint8
	down bool
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
	return time.Now().UnixNano() / 1000000
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
					vKey: vKey,
					down: down,
					timestamp: unixMilli(),
				}
			}
		}

		lastState = newStates
		time.Sleep(1 * time.Millisecond)
	}
}

func RecordKeys(stopKey uint8, w *bufio.Writer) {
	changes := make(chan keyStateChange)
	go watchStateChanges(changes, stopKey)

	prevMillis := unixMilli()
	for change := range changes {
		diffMillis := change.timestamp - prevMillis
		prevMillis = change.timestamp

		keyOp := "Up"
		if change.down {
			keyOp = "Down"
		}
		if _, err := w.WriteString(fmt.Sprintf("sleep %d\nvKey%s %d\n", diffMillis, keyOp, change.vKey)); err != nil {
			panic(err)
		}
	}
}
