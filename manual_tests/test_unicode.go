package main

import (
	"fmt"
	"os"
	
	"github.com/gdamore/tcell/v2"
)

func main() {
	// Set UTF-8 environment variables
	os.Setenv("LC_ALL", "en_US.UTF-8")
	os.Setenv("LANG", "en_US.UTF-8")
	os.Setenv("LC_CTYPE", "en_US.UTF-8")
	
	// Initialize tcell screen
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating screen: %v\n", err)
		os.Exit(1)
	}

	screen.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))

	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing screen: %v\n", err)
		os.Exit(1)
	}
	
	defer screen.Fini()
	
	// Test Unicode characters
	testString := "hell привет"
	
	// Display the string
	for i, r := range testString {
		screen.SetContent(i, 0, r, nil, tcell.StyleDefault)
	}
	
	// Add instructions
	instructions := "Press any key to exit"
	for i, r := range instructions {
		screen.SetContent(i, 2, r, nil, tcell.StyleDefault)
	}
	
	screen.Show()
	
	// Wait for a key press
	screen.PollEvent()
	
	fmt.Println("Test string displayed:", testString)
}