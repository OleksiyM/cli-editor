package main

import (
	"flag"
	"fmt"
	"os"

	"cli-editor/internal/ui"  // Updated import path to match go.mod
)

func main() {
	// Set UTF-8 locale environment variables before doing anything else
	// This is essential for proper UTF-8 character handling in the terminal
	os.Setenv("LC_ALL", "en_US.UTF-8")
	os.Setenv("LANG", "en_US.UTF-8")
	os.Setenv("LC_CTYPE", "en_US.UTF-8")

	// Parse command line arguments
	flag.Parse()
	
	// Get filename from arguments
	var filename string
	if flag.NArg() > 0 {
		filename = flag.Arg(0)
	}
	
	// Set terminal to raw mode to capture all key events
	os.Setenv("TERM", "xterm-256color")
	
	// Create and run editor
	editor, err := ui.NewEditor(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing editor: %v\n", err)
		os.Exit(1)
	}
	
	// Start the editor
	if err := editor.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running editor: %v\n", err)
		os.Exit(1)
	}
}
