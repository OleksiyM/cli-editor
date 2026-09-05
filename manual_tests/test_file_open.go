package main

import (
	"fmt"
	"os"

	"cli-editor/internal/ui"
)

func main() {
	// Test 1: Try to open a non-existent file
	fmt.Println("Test 1: Opening non-existent file")
	editor, err := ui.NewEditor("nonexistent.txt")
	if err != nil {
		fmt.Printf("Error creating editor: %v\n", err)
		os.Exit(1)
	}

	// Check if error dialog was created
	if editor != nil && editor.ActiveDialog() != nil {
		fmt.Println("Error dialog was created successfully")
	} else {
		fmt.Println("No error dialog was created")
	}

	// Test 2: Try to open an existing file
	fmt.Println("\nTest 2: Opening existing file")
	editor2, err := ui.NewEditor("l.txt")
	if err != nil {
		fmt.Printf("Error creating editor: %v\n", err)
		os.Exit(1)
	}

	// Check if there's no error dialog
	if editor2 != nil && editor2.ActiveDialog() == nil {
		fmt.Println("No error dialog for existing file - correct behavior")
	} else {
		fmt.Println("Unexpected error dialog for existing file")
	}
}