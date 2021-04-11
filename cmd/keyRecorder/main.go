package main

import (
	"bufio"
	"fmt"
	"github.com/alexflint/go-arg"
	"keyRecorder/internal/keyRecorder"
	"os"
)

var options struct {
	File     string `arg:"positional"`
	StopCode int    `arg:"-s" placeholder:"CODE" default:"27" help:"the virtual key code to listen for to stop recording (defaults to escape)"`
}

func main() {
	argParser := arg.MustParse(&options)

	if options.File == "" {
		fmt.Printf("Error: no output file specified\n\n")
		argParser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	// Open the file.
	f, err := os.Create(options.File)
	defer f.Close()
	if err != nil {
		if pathErr, ok := err.(*os.PathError); ok {
			fmt.Printf("An error occurred while trying to %s", pathErr.Error())
		} else {
			fmt.Printf("An error occurred while trying to open the output file: %s", err)
		}
		os.Exit(1)
	}

	// Create the file writer.
	w := bufio.NewWriter(f)
	defer w.Flush()

	// Start recording.
	keyRecorder.RecordKeys(uint8(options.StopCode), w)
}
