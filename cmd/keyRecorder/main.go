package main

import (
	"bufio"
	"keyRecorder/internal/keyRecorder"
	"os"
)

func main() {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	keyRecorder.RecordKeys(0x1B, w)
}
